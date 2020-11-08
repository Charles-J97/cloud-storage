package route

import (
	"cc-transfer/controller"
	util "cc-transfer/middleware"
	"github.com/gin-gonic/gin"
)

func Router() *gin.Engine {
	router := gin.Default()

	//用户接口
	router.POST("/user/signUp", controller.UserSignUpHandler)
	router.GET("/user/logIn", controller.UserLogInHandler)

	//加入中间件，用于校验token
	router.Use(util.Auth())

	router.GET("/user/queryInfo", controller.UserInfoQueryHandler)

	//文件元信息接口
	router.POST("/file/upload", controller.UploadFileHandler)
	router.GET("/file/download", controller.DownloadFileHandler)
	router.DELETE("/file/single", controller.DeleteSingleFileHandler)
	router.DELETE("/file/batch", controller.DeleteBatchFileHandler)
	router.GET("/file/single", controller.GetSingleFileMetaHandler)
	router.GET("/file/batch", controller.GetBatchFileMetaHandler)

	//用户文件接口
	router.GET("/user_file/batch", controller.GetUserFileInfoHandler)
	router.PUT("/user_file/rename", controller.UpdateUserFileNameHandler)
	router.DELETE("/user_file/single", controller.DeleteSingleUserFileHandler)
	router.DELETE("/user_file/batch", controller.DeleteBatchUserFileHandler)
	router.POST("/user_file/share", controller.ShareUserFileHandler)

	//分块上传接口
	router.GET("/mp_upload/initiate", controller.InitiateMultipartUploadHandler)
	router.POST("/mp_upload/upload", controller.MultipartUploadHandler)
	router.POST("/mp_upload/complete", controller.CompleteMultipartUploadHandler)
	router.GET("/mp_upload/progress", controller.MultipartUploadProgressHandler)
	router.DELETE("/mp_upload/cancel", controller.MultipartUploadCancelHandler)

	return router
}

