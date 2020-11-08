package controller

import (
	util "cc-dataSource/middleware"
	"cc-dataSource/models/dto"
	"cc-dataSource/services"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"net/http"
	"strconv"
)

//分块上传初始化接口 GET
//chunkSize单位是B，5M就是5*1024*1024
func InitiateMultipartUploadHandler(ctx *gin.Context) {
	username := ctx.Query("username")
	fileHash := ctx.Query("fileHash")
	filename := ctx.Query("filename")
	curServerAddr := ctx.Query("serverAddr")
	fileSize, _ := strconv.Atoi(ctx.Query("fileSize"))
	chunkSize, _ := strconv.Atoi(ctx.Query("chunkSize"))

	upInfo, err := services.InitiateMultipartUpload(username, fileHash, filename, curServerAddr, fileSize, chunkSize)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"code": -1, "msg": err.Error()})
		return
	}
	//若是upInfo为空，则代表触发了秒传
	if upInfo == nil {
		ctx.JSON(http.StatusOK, gin.H{"code": 0, "msg": "OK! Uploaded fast"})
	}
	res := util.RespMsg{
		Code: 0,
		Msg:  "OK",
		Data: upInfo,
	}
	ctx.Data(0, "octet-stream", res.JSONBytes())
}

//分块上传接口 POST
func MultipartUploadHandler(ctx *gin.Context) {
	fileUpload := dto.FileMpUpload{}
	err := ctx.ShouldBindBodyWith(&fileUpload, binding.JSON)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"code": -1, "msg": err.Error()})
		return
	}

	err = services.MpUpload(fileUpload.FileBytes, fileUpload.FileHash, fileUpload.ChunkIndex)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"code": -1, "msg": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"code": 0, "msg": "OK"})
}

//通知分块上传完成接口 POST
func CompleteMultipartUploadHandler(ctx *gin.Context) {
	var mpUpload dto.CompleteMultipartUpload
	err := ctx.ShouldBindBodyWith(&mpUpload, binding.JSON)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"code": -1, "msg": err.Error()})
		return
	}
	err = services.CompleteMultipartUpload(mpUpload.Username, mpUpload.FileHash, mpUpload.Filename, mpUpload.ServerAddr)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"code": -1, "msg": err.Error()})
	}
	ctx.JSON(http.StatusOK, gin.H{"code": 0, "msg": "OK"})
}

//显示分块上传进度接口 GET
func MultipartUploadProgressHandler(ctx *gin.Context) {
	fileHash := ctx.Query("fileHash")

	progress, err := services.MultipartUploadProgress(fileHash)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"code": -1, "msg": err.Error()})
		return
	}
	res := util.RespMsg{
		Code: 0,
		Msg:  "OK",
		Data: progress,
	}
	ctx.Data(0, "octet-stream", res.JSONBytes())
}

//取消分块上传接口 DEL
func MultipartUploadCancelHandler(ctx *gin.Context) {
	fileHash := ctx.Query("fileHash")

	err := services.MultipartUploadCancel(fileHash)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"code": -1, "msg": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"code": 0, "msg": "OK"})
}
