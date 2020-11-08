package dao

import "cc-dataSource/models"

type FileMeta struct {
	FileHash       string `json:"fileHash" db:"file_hash"`
	Filename       string `json:"filename" db:"filename"`
	FileSize       int    `json:"fileSize" db:"file_size"`
	FileLocalAddr  string `json:"fileLocalAddr" db:"file_local_addr"`
	FileRemoteAddr string `json:"fileRemoteAddr" db:"file_remote_addr"`
	CreateTime     string `json:"createTime" db:"create_time"`
	UpdateTime     string `json:"updateTime" db:"update_time"`
}

//保存meta到MySQL数据库中
func InsertFileMeta(fileMeta FileMeta) error {
	fileHash := fileMeta.FileHash
	filename := fileMeta.Filename
	fileSize := fileMeta.FileSize
	fileLocalAddr := fileMeta.FileLocalAddr
	fileRemoteAddr := fileMeta.FileRemoteAddr

	_, err := GetFileMetaByHash([]string{fileHash})
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			_, err = models.DBConn().Exec("INSERT ignore INTO file (`file_hash`, `filename`, `file_size`, `file_local_addr`, `file_remote_addr`, `status`) values (?,?,?,?,?,1)", fileHash, filename, fileSize, fileLocalAddr, fileRemoteAddr)
			if err != nil {
				return err
			}
			return nil
		}
		return err
	}

	return nil
}

//根据Hash获取文件的元信息对象
func GetFileMetaByHash(fileHash []string) ([]FileMeta, error) {
	fileMeta := make([]FileMeta, len(fileHash))
	for i, item := range fileHash {
		err := models.DBConn().Get(&fileMeta[i], "SELECT file_hash, filename, file_size, file_local_addr, file_remote_addr, create_time, update_time FROM file WHERE file_hash = ? AND status = 1 limit 1", item)
		if err != nil {
			return nil, err
		}
	}
	return fileMeta, nil
}

//删除指定Hash的文件元信息(软删除)
func DeleteFileMetaByHash(fileHash []string) error {
	for _, item := range fileHash {
		_, err := models.DBConn().Exec("UPDATE file SET status = 0 WHERE file_hash = ?", item)
		if err != nil {
			return err
		}
	}
	return nil
}
