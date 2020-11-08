package services

import (
	rPool "cc-dataSource/cache"
	util "cc-dataSource/middleware"
	"cc-dataSource/models/dao"
	mp "cc-dataSource/models/dto"
	"errors"
	"fmt"
	"github.com/garyburd/redigo/redis"
	_ "github.com/garyburd/redigo/redis"
	"math"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	LocalPath      = "/Users/charles/data/cloud_storage/local" // "/username/fileHash"
	RemotePath     = "/"
	MpPath         = "/Users/charles/data/cloud_storage/mpUpload" // "/uploadId/chunkIndex"
	MpInfoPrefix   = "MP_INFO_"                                   // + uploadId，存的是被分块的文件信息
	HashUpIdPrefix = "HASH_UPLOAD_ID_"                            // + fileHash，存的是uploadId
)

//用于判断分块文件存放目录和本地文件存放目录是否可以被创建
func init() {
	if err := os.MkdirAll(MpPath, 0744); err != nil {
		fmt.Println("This directory cannot be created: " + MpPath)
		os.Exit(1)
	}
	if err := os.MkdirAll(LocalPath, 0744); err != nil {
		fmt.Println("This directory cannot be created: " + LocalPath)
		os.Exit(1)
	}
}

//初始化分块上传信息逻辑
//返回当前分块上传信息：当前文件hash，当前文件size，当前uploadId，当前文件被分成多少块，每块多大
//支持断点续传功能，也即用户发起一个新的上传请求后，服务端会首先判断该hash对应的文件存不存在，也即支不支持秒传；再会判断该文件有没有之前上传过但是没有上传完
//若是符合第二种情况，是不支持秒传的，但是可以启动断点续传
func InitiateMultipartUpload(username, fileHash, filename, curServerAddr string, fileSize, chunkSize int) (*mp.MultipartUploadInfo, error) {
	rConn := rPool.RedisPool().Get()
	defer rConn.Close()

	//判断该文件有没有被别人上传过，也即触不触发秒传逻辑
	isExecuted, err := FastUpload(fileHash, filename, username, curServerAddr)
	if err != nil {
		return nil, err
	}
	if isExecuted == true {
		return nil, nil
	}

	//判断该文件是否之前已经进行过分块上传，但是没传完，也即触不触发断点续传逻辑
	//具体思路：判断redis中是否存在和此fileHash对应的uploadId，若是不存在，说明是首次上传；若是存在，说明该文件之前已经被分块上传过但没有传完（因为传完的话会删掉这条记录）
	//若是断点续传，则跟据uploadId获取已上传过的文件分块列表，就是先HGETAll，然后再获取所有前缀是"chunkIndex_"并且为"1"的chunk，将其填入ChunkExists []int列表中
	var uploadId string
	var chunksExist []int
	exists, err := redis.Bool(rConn.Do("EXISTS", HashUpIdPrefix+fileHash))
	//若是首次上传，则新建uploadId
	if !exists {
		uploadId = username + fmt.Sprintf("%x", time.Now().UnixNano())
	} else {
		uploadId, err = redis.String(rConn.Do("GET", HashUpIdPrefix+fileHash))
		if err != nil {
			return nil, err
		}
		data, err := redis.Values(rConn.Do("HGETALL", MpInfoPrefix+uploadId))
		if err != nil {
			return nil, err
		}
		for i := 0; i < len(data); i += 2 {
			k := string(data[i].([]byte))
			v := string(data[i+1].([]byte))
			if strings.HasPrefix(k, "chunkIndex_") && v == "1" {
				chunkIdx, _ := strconv.Atoi(k[11:])
				chunksExist = append(chunksExist, chunkIdx)
			}
		}
	}

	//生成分块初始化信息
	upInfo := mp.MultipartUploadInfo{
		FileHash:     fileHash,
		FileSize:     fileSize,
		UploadId:     uploadId,
		ChunkSize:    chunkSize,
		ChunkCount:   int(math.Ceil(float64(fileSize) / float64(chunkSize))), //fileSize/ChunkSize，然后向上取整，再转为int
		ChunksExists: chunksExist,
	}

	//将初始化信息写入redis缓存中
	//Redis数据只能用于缓存的原因是：不设置超时时间的话，会很占内存
	var args []interface{}
	fieldKey := MpInfoPrefix + uploadId
	keyAndValues := map[string]interface{}{"fileHash": upInfo.FileHash, "fileSize": upInfo.FileSize, "chunkSize": upInfo.ChunkSize, "chunkCount": upInfo.ChunkCount}
	args = append(args, fieldKey)
	for k, v := range keyAndValues {
		args = append(args, k, v)
	}
	_, _ = rConn.Do("HSET", args...)
	//设置HSET的超时时间
	_, _ = rConn.Do("EXPIRE", fieldKey, 43200)
	//将文件Hash和uploadID做映射，并设置超时时间
	_, _ = rConn.Do("SET", HashUpIdPrefix+fileHash, upInfo.UploadId, "EX", 43200)

	return &upInfo, nil
}

//TODO 通知分块上传完成逻辑 增加对于Ceph的支持
//在前端发完最后一个chunk包后发送此通知，此时要进行块的合并，分块文件以及redis记录的删除，以及用户文件表和文件表的更新
func CompleteMultipartUpload(username, fileHash, filename, serverAddr string) error {
	rConn := rPool.RedisPool().Get()
	defer rConn.Close()

	uploadId, err := redis.String(rConn.Do("GET", HashUpIdPrefix + fileHash))
	if err != nil {
		return err
	}

	//查询redis，通过uploadId查询是否所有的分块都已上传，以及此文件的一些基本信息
	fieldKey := MpInfoPrefix + uploadId
	data, err := redis.Values(rConn.Do("HGETALL", fieldKey))
	if err != nil {
		return err
	}

	var fileSize int
	totalCount := 0 //已上传块数
	chunkCount := 0 //应该上传的块数
	for i := 0; i < len(data); i += 2 {
		key := string(data[i].([]byte))
		value := string(data[i+1].([]byte))
		if key == "chunkCount" {
			chunkCount, _ = strconv.Atoi(value)
		} else if strings.HasPrefix(key, "chunkIndex_") && value == "1" {
			totalCount++
		} else if key == "fileSize" {
			fileSize, _ = strconv.Atoi(value)
		}
	}
	if totalCount != chunkCount {
		return errors.New("Multipart uploading goes wrong ")
	}

	//合并分块，并将合并后的文件存入本地中
	fileAddr := LocalPath + "/" + username + "/" + filename
	chunkAddr := MpPath + "/" + uploadId
	if mergeSuc := util.MergeChuncksByShell(chunkAddr, fileAddr, fileHash); !mergeSuc {
		return errors.New("Multi-part uploading failed when merging parts ")
	}

	//删除分块文件数据（本地数据和redis数据都删）
	err = deleteMpUploadDataByFileHash(fileHash)
	if err != nil {
		return err
	}

	//更新文件表，用户文件表，文件服务器表
	fileMeta := dao.FileMeta {
		FileHash:      fileHash,
		FileSize:      fileSize,
		FileLocalAddr: fileAddr,
	}
	err = dao.InsertFileMeta(fileMeta)
	if err != nil {
		return err
	}

	tmpUserFile := dao.UserFile {
		Username: username,
		Filename: filename,
		FileHash: fileHash,
	}
	err = dao.InsertUserFile(tmpUserFile)
	if err != nil {
		return err
	}

	tmpFileServer := dao.FileServer {
		FileHash: fileHash,
		ServerAddr: serverAddr,
	}
	err = dao.InsertFileServer(tmpFileServer)
	if err != nil {
		return err
	}

	return nil
}

//显示分块上传进度逻辑
func MultipartUploadProgress(fileHash string) (int, error) {
	rConn := rPool.RedisPool().Get()
	defer rConn.Close()

	//根据fileHash找到上传的uploadId
	//先判断该fileHash对应的文件在不在上传，也就是redis中存不存在他的数据
	exists, err := redis.Bool(rConn.Do("EXISTS", HashUpIdPrefix+fileHash))
	if err != nil {
		return 0, err
	}
	if !exists {
		return 0, errors.New("This file has not been uploaded multiply before ")
	}
	//若存在，则找uploadId
	uploadId, err := redis.String(rConn.Do("GET", HashUpIdPrefix+fileHash))
	if err != nil {
		return 0, err
	}

	fieldKey := MpInfoPrefix + uploadId
	data, err := redis.Values(rConn.Do("HGETALL", fieldKey))
	if err != nil {
		return 0, err
	}
	totalCount := 0
	chunkCount := 0
	for i := 0; i < len(data); i += 2 {
		key := string(data[i].([]byte))
		value := string(data[i+1].([]byte))
		if key == "chunkCount" {
			chunkCount, _ = strconv.Atoi(value)
		} else if strings.HasPrefix(key, "chunkIndex_") && value == "1" {
			totalCount++
		}
	}
	progress := float64(totalCount) / float64(chunkCount) * 100.0
	return int(math.Floor(progress)), nil
}

//取消分块上传逻辑
func MultipartUploadCancel(fileHash string) error {
	rConn := rPool.RedisPool().Get()
	defer rConn.Close()
	err := deleteMpUploadDataByFileHash(fileHash)
	if err != nil {
		return err
	}
	return nil
}

//根据文件hash删除redis中此次upload操作的记录，以及本地分块文件
func deleteMpUploadDataByFileHash(fileHash string) error {
	rConn := rPool.RedisPool().Get()
	defer rConn.Close()

	exists, err := redis.Bool(rConn.Do("EXISTS", HashUpIdPrefix+fileHash))
	if err != nil {
		return err
	}
	if !exists {
		return errors.New("This file has not been uploaded multiply before ")
	}

	uploadId, err := redis.String(rConn.Do("GET", HashUpIdPrefix+fileHash))
	if err != nil {
		return err
	}

	//先删本机数据
	err = os.RemoveAll(MpPath + "/" + uploadId)
	if err != nil {
		return err
	}

	//再删redis中记录
	_, err = rConn.Do("DEL", MpInfoPrefix+uploadId)
	if err != nil {
		return err
	}
	_, err = rConn.Do("DEL", HashUpIdPrefix+fileHash)
	if err != nil {
		return err
	}
	return nil
}
