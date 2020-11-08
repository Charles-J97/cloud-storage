package dto

type FileUpload struct {
	Username   string
	FileHash   string
	Filename   string
	ServerAddr string
	FileBytes  []byte
}
