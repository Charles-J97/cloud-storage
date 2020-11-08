package services

import (
	"bytes"
	"cc-dataSource/error"
	"cc-dataSource/models/dao"
	"cc-dataSource/models/dto"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"os"
)

//查询用户文件信息逻辑
//若count为0（也就是没传值），则表明查询该username对应的所有文件信息，否则就查询最近上传的的k个文件信息
func GetUserFileInfo(username string, k int64) ([]dao.UserFile, error) {
	var userFiles []dao.UserFile
	var err error
	//判断是否存在此用户
	_, err = dao.GetUserByName(username)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return nil, errors.New(errs.NoUser)
		}
		return nil, err
	}

	//查询用户文件
	if k == 0 {
		userFiles, err = dao.GetUserFileByUsername(username)
		if err != nil {
			if err.Error() == "sql: no rows in result set" {
				return nil, errors.New(errs.NoFile)
			}
			return nil, err
		}
	} else {
		userFiles, err = dao.GetUserFileByUsernameAndCount(username, k)
		if err != nil {
			if err.Error() == "sql: no rows in result set" {
				return nil, errors.New(errs.NoFile)
			}
			return nil, err
		}
	}
	return userFiles, nil
}

//修改（重命名）用户文件信息逻辑
//根据username和filename唯一确定用户文件表中的一项，再将那项的filename改成新名字
func UpdateUserFilename(username, filename, fileNewName string) error {
	//先判断是否存在此用户
	_, err := dao.GetUserByName(username)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return errors.New(errs.NoUser)
		}
		return err
	}

	//判断该用户是否存在该文件
	_, err = dao.GetUserFileByUsernameAndFilename(username, filename)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return errors.New(username + errs.UserHasNoFile + filename)
		}
		return err
	}

	//更改文件名
	err = dao.UpdateUserFilenameByUsernameAndFileHash(username, filename, fileNewName)
	if err != nil {
		return err
	}
	return nil
}

//删除用户文件信息逻辑
//这里的删除不删除用户元信息，仅删除用户文件表中的映射关系
func DeleteUserFile(username string, filename []string) error {
	//判断是否存在此用户
	_, err := dao.GetUserByName(username)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return errors.New(errs.NoUser)
		}
		return err
	}

	//判断该用户是否存在这些文件
	for _, item := range filename {
		_, err = dao.GetUserFileByUsernameAndFilename(username, item)
		if err != nil {
			if err.Error() == "sql: no rows in result set" {
				return errors.New(username + errs.UserHasNoFile + item)
			}
			return err
		}
	}

	//删除该用户下的该文件
	err = dao.DeleteUserFileByUsernameAndFilename(username, filename)
	if err != nil {
		return err
	}
	return nil
}

func ShareUserFile(username, filename, newUsername, curServerAddr string) error {
	//判断两个用户名是否均合法
	user, err := dao.GetUserByName(username)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return errors.New(errs.NoUser)
		}
		return err
	}
	newUser, err := dao.GetUserByName(newUsername)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return errors.New(errs.NoUser)
		}
		return err
	}

	//判断本用户有没有该文件
	//要共享的文件就是这个userFile
	userFile, err := dao.GetUserFileByUsernameAndFilename(username, filename)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return errors.New(username + errs.UserHasNoFile + filename)
		}
		return err
	}

	//再判断新用户有没有该fileHash的文件
	//根据新用户名字找到其所有的文件
	newUserFiles, err := dao.GetUserFileByUsername(newUsername)
	if err != nil {
		return err
	}

	//若新用户存在该hash值的文件时，直接用新文件名配上文件hash值更新其用户文件表即可
	for _, item := range newUserFiles {
		if userFile.FileHash == item.FileHash {
			newUserFile := dao.UserFile{
				Username: newUsername,
				Filename: filename,
				FileHash: userFile.FileHash,
			}
			err = dao.InsertUserFile(newUserFile)
			if err != nil {
				return err
			}
			return nil
		}
	}

	//若新用户不存在该hash值的文件，则判断新用户和老用户处不处于同一server下
	//若处于同一server，直接用新文件名配上文件hash值更新其用户文件表即可
	if user.ServerAddr == newUser.ServerAddr {
		newUserFile := dao.UserFile{
			Username: newUsername,
			Filename: filename,
			FileHash: userFile.FileHash,
		}
		err = dao.InsertUserFile(newUserFile)
		if err != nil {
			return err
		}
	} else {
		//若处于不同server，则先对新用户server发起上传文件请求
		fileMeta, err := dao.GetFileMetaByHash([]string{userFile.FileHash})
		if err != nil {
			return err
		}

		err = sendFileUploadRequest(fileMeta[0].FileHash, fileMeta[0].Filename, fileMeta[0].FileLocalAddr, newUsername, newUser.ServerAddr)
		if err != nil {
			return err
		}
	}

	return nil
}

func sendFileUploadRequest(fileHash, filename, fileAddr, newUserName, newUserServerAddr string) error {
	url := "http://" + newUserServerAddr + "/file/upload"
	file, err := os.Open(fileAddr)
	if err != nil {
		return err
	}
	defer file.Close()
	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}

	fileUpload := dto.FileUpload{
		Username: newUserName,
		FileHash: fileHash,
		Filename: filename,
		ServerAddr: newUserServerAddr,
		FileBytes: fileBytes,
	}

	fileJson, err := json.Marshal(fileUpload)
	if err != nil {
		return err
	}

	reader := bytes.NewReader(fileJson)
	req, err := http.NewRequest("POST", url, reader)
	if err != nil {
		return err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	respMsgBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	respMsg := make(map[string]interface{})
	err = json.Unmarshal(respMsgBytes, &respMsg)
	if err != nil {
		return err
	}
	if int(respMsg["code"].(float64)) != 0 {
		return errors.New(respMsg["msg"].(string))
	}
	return nil
}
