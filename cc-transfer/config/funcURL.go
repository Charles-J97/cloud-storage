package config

const (
	UploadFileUrl = "/file/upload"
	DownloadFileUrl = "/file/download"
	DeleteSingleFileUrl = "/file/single"
	DeleteBatchFileUrl = "/file/batch"
	GetSingleFileMetaUrl = "/file/single"
	GetBatchFileMetaUrl = "/file/batch"

	GetUserFileInfoUrl = "/user_file/batch"
	UpdateUserFilenameUrl = "/user_file/rename"
	DeleteSingleUserFileUrl = "/user_file/single"
	DeleteBatchUserFileUrl = "/user_file/batch"
	ShareUserFileUrl = "/user_file/share"

	InitiateMultipartUploadUrl = "/mp_upload/initiate"
	MultipartUploadUrl = "/mp_upload/upload"
	CompleteMultipartUploadUrl = "/mp_upload/complete"
	MultipartUploadProgressUrl = "/mp_upload/progress"
	MultipartUploadCancelUrl = "/mp_upload/cancel"
)
