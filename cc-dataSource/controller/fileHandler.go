package controller

import (
	util "cc-dataSource/middleware"
	"cc-dataSource/models/dto"
	"cc-dataSource/services"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"net/http"
)

//文件上传接口 POST
//http.ResponseWriter对象用于向用户返回数据，*http.Request用于接受用户请求
func UploadFileHandler(ctx *gin.Context) {
	fileUpload := dto.FileUpload{}
	err := ctx.ShouldBindBodyWith(&fileUpload, binding.JSON)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"code": -1, "msg": err.Error()})
		return
	}

	err = services.UploadFile(fileUpload.Username, fileUpload.FileHash, fileUpload.Filename, fileUpload.ServerAddr, fileUpload.FileBytes)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"code": -1, "msg": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"code": 0, "msg": "OK"})
}

//文件下载接口 GET
func DownloadFileHandler(ctx *gin.Context) {
	username := ctx.Query("username")
	filename := ctx.Query("filename")
	curServerAddr := ctx.Query("curServerAddr")

	data, err := services.DownloadFile(username, filename, curServerAddr)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"code": -1, "msg": err.Error()})
		return
	}
	ctx.Header("content-disposition", "attachment; filename=\""+filename+"\"")
	ctx.Data(http.StatusOK, "application/octect-stream", data)
}

//单个文件删除接口 DEL
func DeleteSingleFileHandler(ctx *gin.Context) {
	fileHash := ctx.Query("fileHash")

	err := services.DeleteFile([]string{fileHash})
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"code": -1, "msg": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"code": 0, "msg": "OK"})
}

//批量文件删除接口 DEL
func DeleteBatchFileHandler(ctx *gin.Context) {
	fileHash := ctx.QueryArray("fileHash")

	err := services.DeleteFile(fileHash)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"code": -1, "msg": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"code": 0, "msg": "OK"})
}

//获取单个文件元信息接口 GET
func GetSingleFileMetaHandler(ctx *gin.Context) {
	fileHash := ctx.Query("fileHash")

	fileMeta, err := services.GetFileMeta([]string{fileHash})
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"code": -1, "msg": err.Error()})
		return
	}
	res := util.RespMsg{
		Code: 0,
		Msg:  "OK",
		Data: fileMeta[0],
	}
	ctx.Data(0, "octet-stream", res.JSONBytes())
}

//获取批量文件元信息接口 GET
func GetBatchFileMetaHandler(ctx *gin.Context) {
	fileHash := ctx.QueryArray("fileHash")

	fileMeta, err := services.GetFileMeta(fileHash)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"code": -1, "msg": err.Error()})
		return
	}
	res := util.RespMsg{
		Code: 0,
		Msg:  "OK",
		Data: fileMeta,
	}
	ctx.Data(0, "octet-stream", res.JSONBytes())
}
