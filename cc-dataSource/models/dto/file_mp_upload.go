package dto

type CompleteMultipartUpload struct {
	Username string `json:"username"`
	FileHash string `json:"fileHash"`
	Filename string `json:"filename"`
	ServerAddr string `json:"serverAddr"`
}

//分块初始化信息
type MultipartUploadInfo struct {
	FileHash     string
	FileSize     int
	UploadId     string
	ChunkSize    int   //每个分块的大小
	ChunkCount   int   //分块的数量
	ChunksExists []int //已经上传完成的分块索引列表
}

//分块上传req结构体
type FileMpUpload struct {
	FileHash   string
	ChunkIndex string
	FileBytes  []byte
}
