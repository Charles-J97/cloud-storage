package route

import (
	"cc-dataSource/controller"
	"github.com/gin-gonic/gin"
)

func Router() *gin.Engine {
	router := gin.Default()

	router.GET("/test", controller.Test)

	//文件元信息接口
	file := router.Group("/file")
	file.POST("/upload", controller.UploadFileHandler)
	file.GET("/download", controller.DownloadFileHandler)
	file.DELETE("/single", controller.DeleteSingleFileHandler)
	file.DELETE("/batch", controller.DeleteBatchFileHandler)
	file.GET("/single", controller.GetSingleFileMetaHandler)
	file.GET("/batch", controller.GetBatchFileMetaHandler)

	//用户文件接口
	userFile := router.Group("/user_file")
	userFile.GET("/batch", controller.GetUserFileInfoHandler)
	userFile.PUT("/rename", controller.UpdateUserFilenameHandler)
	userFile.DELETE("/single", controller.DeleteSingleUserFileHandler)
	userFile.DELETE("/batch", controller.DeleteBatchUserFileHandler)
	userFile.POST("/share", controller.ShareUserFileHandler)

	//分块上传接口
	mpUpload := router.Group("/mp_upload")
	mpUpload.GET("/initiate", controller.InitiateMultipartUploadHandler)
	mpUpload.POST("/upload", controller.MultipartUploadHandler)
	mpUpload.POST("/complete", controller.CompleteMultipartUploadHandler)
	mpUpload.GET("/progress", controller.MultipartUploadProgressHandler)
	mpUpload.DELETE("/cancel", controller.MultipartUploadCancelHandler)

	return router
}
