package dto

type FileMpUpload struct {
	FileHash   string
	ChunkIndex string
	FileBytes  []byte
}

type CompleteMultipartUpload struct {
	Username string `json:"username"`
	FileHash string `json:"fileHash"`
	Filename string `json:"filename"`
	ServerAddr string `json:"serverAddr"`
}
