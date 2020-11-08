package dao

import "cc-dataSource/models"

type FileServer struct {
	FileHash   string `json:"fileHash" db:"file_hash"`
	ServerAddr string `json:"serverAddr" db:"server_addr"`
}

func InsertFileServer(fileServer FileServer) error {
	fileHash := fileServer.FileHash
	serverAddr := fileServer.ServerAddr
	_, err := models.DBConn().Exec("INSERT IGNORE INTO file_server (`file_hash`, `server_addr`, `status`) values (?,?,1)", fileHash, serverAddr)
	if err != nil {
		return err
	}
	return nil
}

func GetFileServerByHash(fileHash string) ([]FileServer, error) {
	var fileServer []FileServer
	err := models.DBConn().Select(&fileServer, "SELECT file_hash, server_addr FROM file_server WHERE file_hash = ? and status = 1", fileHash)
	if err != nil {
		return nil, err
	}
	return fileServer, nil
}

func DeleteFileServerByHash(fileHash []string) error {
	for _, item := range fileHash {
		_, err := models.DBConn().Exec("UPDATE file_server SET status = 0 WHERE file_hash = ?", item)
		if err != nil {
			return err
		}
	}
	return nil
}