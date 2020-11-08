package services

import (
	rPool "cc-dataSource/cache"
	"cc-dataSource/models/dao"
	"github.com/garyburd/redigo/redis"
	"os"
	"path"
)

//秒传
func FastUpload(fileHash, filename, username, curServerAddr string) (bool, error) {
	fileServer, err := dao.GetFileServerByHash(fileHash)
	if err != nil {
		//若是文件服务器表中没有此fileHash的文件，此秒传不执行
		if err.Error() == "sql: no rows in result set" {
			return false, nil
		}
		return false, err
	} else if !ContainCurServer(fileServer, curServerAddr) {
		//若有文件之前被传过，但是传的不是现在这个服务器，则秒传也不执行
		return false, nil
	} else {
		//否则，秒传执行，此时仅需更新用户文件表即可
		tmpUserFile := dao.UserFile{
			Username: username,
			Filename: filename,
			FileHash: fileHash,
		}
		err = dao.InsertUserFile(tmpUserFile)
		if err != nil {
			return false, err
		}
		return true, nil
	}
}

//普通上传
//最后此函数会更新文件表，用户文件表，以及文件服务器表记录
func CommonUpload(file []byte, username, fileHash, filename, curServerAddr string) (*dao.FileMeta, error) {
	fileMeta := dao.FileMeta{
		FileHash:       fileHash,
		Filename:       filename,
		FileLocalAddr:  LocalPath + "/" + username + "/" + filename,
		FileRemoteAddr: RemotePath + "/" + filename,
	}
	//创建一个本地文件来存储接收到的文件流，Create参数为路径
	//第一个参数是目录，第二个参数是权限，0744表明当前用户权限是7，其他用户权限是4(只读)。MkDirAll如果遇到已经存在的目录，就会不创建，但是也不报错
	err := os.MkdirAll(path.Dir(fileMeta.FileLocalAddr), 0744)
	if err != nil {
		return nil, err
	}
	newFile, err := os.Create(fileMeta.FileLocalAddr)
	if err != nil {
		return nil, err
	}
	defer newFile.Close()

	//把接收到的文件流写入到新创立的文件夹中
	fileMeta.FileSize, err = newFile.Write(file)
	if err != nil {
		return nil, err
	}

	//更新文件表记录
	err = dao.InsertFileMeta(fileMeta)
	if err != nil {
		return nil, err
	}

	//更新用户文件表记录
	tmpUserFile := dao.UserFile{
		Username: username,
		Filename: filename,
		FileHash: fileHash,
	}
	err = dao.InsertUserFile(tmpUserFile)
	if err != nil {
		return nil, err
	}

	//更新文件服务器表记录
	tmpFileServer := dao.FileServer{
		FileHash: fileHash,
		ServerAddr: curServerAddr,
	}
	err = dao.InsertFileServer(tmpFileServer)
	if err != nil {
		return nil, err
	}

	return &fileMeta, nil
}

//分块上传（接收其中一个chunkIndex分块）
func MpUpload(file []byte, fileHash, chunkIndex string) error {
	rConn := rPool.RedisPool().Get()
	defer rConn.Close()

	//先通过fileHash取得uploadId
	uploadId, err := redis.String(rConn.Do("GET", HashUpIdPrefix+fileHash))
	if err != nil {
		return err
	}

	filePath := MpPath + "/" + uploadId + "/" + chunkIndex
	//第一个参数是目录，第二个参数是权限，0744表明当前用户权限是7，其他用户权限是4(只读)。MkDirAll如果遇到已经存在的目录，就会不创建，但是也不报错
	err = os.MkdirAll(path.Dir(filePath), 0744)
	if err != nil {
		return err
	}

	newFile, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer newFile.Close()

	//把接收到的文件流copy到新创立的文件夹中
	_, err = newFile.Write(file)
	if err != nil {
		return err
	}

	//分块上传完成后，更新redis缓存，把本次操作的chunkIndex置为1，表明此块已完成传输
	_, _ = rConn.Do("HSET", MpInfoPrefix + uploadId, "chunkIndex_" + chunkIndex, 1)
	return nil
}

func ContainCurServer(fileServers []dao.FileServer, curServerAddr string) bool {
	for _, item := range fileServers {
		if item.ServerAddr == curServerAddr {
			return true
		}
	}
	return false
}