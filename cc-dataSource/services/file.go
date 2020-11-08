package services

import (
	errs "cc-dataSource/error"
	"cc-dataSource/models/dao"
	"errors"
	"io/ioutil"
	"os"
)

//TODO 文件上传逻辑 加上对于ceph的支持
func UploadFile(username, fileHash, filename, curServerAddr string, file []byte) error {
	//判断该用户有没有上传过该文件
	userFile, err := dao.GetUserFileByUsername(username)
	if err != nil {
		return err
	}
	for _, item := range userFile {
		if item.Filename == filename {
			return errors.New("A file with the same name has been uploaded before ")
		}
	}

	//判断此服务器中存不存在此hash值的文件，也即触不触发秒传逻辑
	isExecuted, err := FastUpload(fileHash, filename, username, curServerAddr)
	if err != nil {
		return err
	}
	if isExecuted == true {
		return nil
	}

	//不触发秒传，则走普通上传逻辑
	_, err = CommonUpload(file, username, fileHash, filename, curServerAddr)
	if err != nil {
		return err
	}

	return nil
}

//TODO 文件下载逻辑 加上对于ceph的支持
func DownloadFile(username, filename, curServerAddr string) ([]byte, error) {
	//根据用户和文件名获得该文件hash，从而获得该文件元信息
	userFile, err := dao.GetUserFileByUsernameAndFilename(username, filename)
	if err != nil {
		return nil, err
	}
	fileMeta, err := dao.GetFileMetaByHash([]string{userFile.FileHash})
	if err != nil {
		return nil, err
	}
	if fileMeta == nil {
		return nil, errors.New("No such file ")
	}

	//判断本地服务器有没有存储该文件
	fileServer, err := dao.GetFileServerByHash(userFile.FileHash)
	if err != nil {
		return nil, err
	}
	if !ContainCurServer(fileServer, curServerAddr) {
		return nil, errors.New("Local server does not contain this file! ")
	}

	//若本地存了，就在本地直接返回数据量给客户端
	var data []byte
	file, err := os.Open(fileMeta[0].FileLocalAddr)
	if err != nil {
		return nil, err
	} else {
		defer file.Close()
		data, err = ioutil.ReadAll(file)
		if err != nil {
			return nil, err
		}
		return data, nil
	}
}

//文件删除逻辑（最好写成事务）
//这里删的是文件元信息，也就是文件的根信息
func DeleteFile(fileHash []string) error {
	var err error
	for _, item := range fileHash {
		//删本机文件
		fileList, err := dao.GetFileMetaByHash([]string{item})
		if err != nil {
			if err.Error() == "sql: no rows in result set" {
				return errors.New(errs.NoFile)
			}
			return err
		}
		file := fileList[0]
		f, err := os.Open(file.FileLocalAddr)
		if err == nil {
			_ = f.Close()
			err = os.Remove(file.FileLocalAddr)
			if err != nil {
				return err
			}
		}
	}

	//删用户文件信息
	err = dao.DeleteUserFileByFileHash(fileHash)
	if err != nil && err.Error() != "sql: no rows in result set" {
		return err
	}

	//删文件元信息
	err = dao.DeleteFileMetaByHash(fileHash)
	if err != nil && err.Error() != "sql: no rows in result set" {
		return err
	}

	//删文件服务器信息
	err = dao.DeleteFileServerByHash(fileHash)
	if err != nil && err.Error() != "sql: no rows in result set" {
		return err
	}
	return nil
}

//获取文件元信息逻辑
func GetFileMeta(fileHash []string) ([]dao.FileMeta, error) {
	fileMeta, err := dao.GetFileMetaByHash(fileHash)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return nil, errors.New(errs.NoFile)
		}
		return nil, err
	}
	return fileMeta, nil
}
