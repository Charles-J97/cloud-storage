package dto

type UpdateUserFilename struct {
	Username    string `json:"username"`
	Filename    string `json:"filename"`
	FileNewName string `json:"fileNewName"`
	ServerAddr  string `json:"serverAddr"`
}

type ShareUserFile struct {
	Username    string `json:"username"`
	Filename    string `json:"filename"`
	NewUsername string `json:"newUsername"`
	ServerAddr  string `json:"serverAddr"`
}
