package dao

import "cc-dataSource/models"

//用户文件表结构体
type UserFile struct {
	Username         string `json:"username" db:"username"`
	Filename         string `json:"filename" db:"filename"`
	FileHash         string `json:"fileHash" db:"file_hash"`
	UploadTime       string `json:"uploadTime" db:"upload_time"`
	LatestUpdateTime string `json:"latestUpdateTime" db:"latest_update_time"`
}

//插入用户文件记录
func InsertUserFile(userFile UserFile) error {
	username := userFile.Username
	filename := userFile.Filename
	fileHash := userFile.FileHash
	var tmpUserFile UserFile
	err := models.DBConn().Get(&tmpUserFile, "SELECT username, file_hash, filename, upload_time, latest_update_time FROM user_file WHERE username = ? AND filename = ? AND status = 1", username, filename)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			_, err = models.DBConn().Exec("INSERT INTO user_file (`username`, `filename`, `file_hash`, `status`) VALUES (?,?,?,1)", username, filename, fileHash)
			if err != nil {
				return err
			}
			return nil
		}
		return err
	}
	return nil
}

//获取该用户所有文件信息
func GetUserFileByUsername(username string) ([]UserFile, error) {
	var userFile []UserFile
	err := models.DBConn().Select(&userFile, "SELECT username, file_hash, filename, upload_time, latest_update_time FROM user_file WHERE username = ? AND status = 1 ORDER BY upload_time DESC", username)
	if err != nil {
		return nil, err
	}
	return userFile, nil
}

//根据用户名获取最近k个文件信息
func GetUserFileByUsernameAndCount(username string, k int64) ([]UserFile, error) {
	var userFile []UserFile
	err := models.DBConn().Select(&userFile, "SELECT username, file_hash, filename, upload_time, latest_update_time FROM user_file WHERE username = ? AND status = 1 ORDER BY upload_time DESC limit ?", username, k)
	if err != nil {
		return nil, err
	}
	return userFile, nil
}

//根据用户名和文件名获取单个文件信息
func GetUserFileByUsernameAndFilename(username, filename string) (*UserFile, error) {
	var userFile UserFile
	err := models.DBConn().Get(&userFile, "SELECT username, file_hash, filename, upload_time, latest_update_time FROM user_file WHERE username = ? AND filename = ? AND status = 1 limit 1", username, filename)
	if err != nil {
		return nil, err
	}
	return &userFile, nil
}

//修改该用户特定文件名
func UpdateUserFilenameByUsernameAndFileHash(username, filename, fileNewName string) error {
	_, err := models.DBConn().Exec("UPDATE user_file SET filename = ? WHERE username = ? AND filename = ? AND status = 1", fileNewName, username, filename)
	if err != nil {
		return err
	}
	return nil
}

//根据用户名和文件名删除用户文件记录（软删除）
func DeleteUserFileByUsernameAndFilename(username string, filename []string) error {
	for _, item := range filename {
		_, err := models.DBConn().Exec("UPDATE user_file SET status = 0 WHERE username = ? AND filename = ?", username, item)
		if err != nil {
			return err
		}
	}
	return nil
}

//根据文件hash值删除用户文件记录（软删除）
func DeleteUserFileByFileHash(fileHash []string) error {
	for _, item := range fileHash {
		_, err := models.DBConn().Exec("UPDATE user_file SET status = 0 WHERE file_hash = ?", item)
		if err != nil {
			return err
		}
	}
	return nil
}
