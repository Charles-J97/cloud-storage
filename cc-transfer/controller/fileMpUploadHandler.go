package controller

import (
	"cc-transfer/config"
	util "cc-transfer/middleware"
	"cc-transfer/models/dto"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"io/ioutil"
	"net/http"
)

//分块上传初始化接口 GET
func InitiateMultipartUploadHandler(ctx *gin.Context) {
	username := ctx.Query("username")
	fileHash := ctx.Query("fileHash")
	filename := ctx.Query("filename")
	curServerAddr := ctx.Query("serverAddr")
	fileSize := ctx.Query("fileSize")
	chunkSize := ctx.Query("chunkSize")
	serverInfo, err := ctx.Cookie("remote_addr")
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"code":-1,"msg":err.Error()})
		return
	}

	url := "http://" + serverInfo + config.InitiateMultipartUploadUrl + "?username=" + username + "&fileHash=" + fileHash + "&filename=" + filename + "&serverAddr=" + curServerAddr + "&fileSize=" + fileSize + "&chunkSize=" + chunkSize

	respMsg, err := util.HandleGETAndDELRequest(url, "GET")
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"code":-1,"msg":err.Error()})
		return
	}

	ctx.Data(http.StatusOK, "octet-stream", respMsg.JSONBytes())
}

//上传分块接口 POST
func MultipartUploadHandler(ctx *gin.Context) {
	serverInfo, err := ctx.Cookie("remote_addr")
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"code":-1,"msg":err.Error()})
		return
	}

	file, err := ctx.FormFile("file")
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"code":-1,"msg":err.Error()})
		return
	}

	fileHash := ctx.Request.FormValue("fileHash")
	chunkIndex := ctx.Request.FormValue("chunkIndex")
	f, _ := file.Open()
	defer f.Close()
	fileBytes, _ := ioutil.ReadAll(f)

	fileUpload := dto.FileMpUpload{
		FileHash: fileHash,
		ChunkIndex: chunkIndex,
		FileBytes: fileBytes,
	}
	jsonFileMpUpload, err := json.Marshal(fileUpload)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"code":-1,"msg":err.Error()})
		return
	}

	url := "http://" + serverInfo + config.MultipartUploadUrl

	respMsg, err := util.HandlePOSTAndPUTRequest(jsonFileMpUpload, url, "POST")
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"code":-1,"msg":err.Error()})
		return
	}

	ctx.Data(http.StatusOK, "octet-stream", respMsg.JSONBytes())
}

//通知分块上传完成接口 POST
func CompleteMultipartUploadHandler(ctx *gin.Context) {
	serverAddr, err := ctx.Cookie("remote_addr")
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"code":-1,"msg":err.Error()})
		return
	}

	var mpUpload dto.CompleteMultipartUpload
	err = ctx.ShouldBindBodyWith(&mpUpload, binding.JSON)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"code": -1, "msg": err.Error()})
		return
	}
	mpUpload.ServerAddr = serverAddr

	jsonBody, err := json.Marshal(mpUpload)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"code":-1,"msg":err.Error()})
		return
	}

	url := "http://" + serverAddr + config.CompleteMultipartUploadUrl

	respMsg, err := util.HandlePOSTAndPUTRequest(jsonBody, url, "POST")
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"code":-1,"msg":err.Error()})
		return
	}

	ctx.Data(http.StatusOK, "octet-stream", respMsg.JSONBytes())
}

//显示分块上传进度接口 GET
func MultipartUploadProgressHandler(ctx *gin.Context) {
	fileHash := ctx.Query("fileHash")
	serverInfo, err := ctx.Cookie("remote_addr")
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"code":-1,"msg":err.Error()})
		return
	}

	url := "http://" + serverInfo + config.MultipartUploadProgressUrl + "?fileHash=" + fileHash

	respMsg, err := util.HandleGETAndDELRequest(url, "GET")
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"code":-1,"msg":err.Error()})
		return
	}

	ctx.Data(http.StatusOK, "octet-stream", respMsg.JSONBytes())
}

//取消分块上传接口 DEL
func MultipartUploadCancelHandler(ctx *gin.Context) {
	fileHash := ctx.Query("fileHash")
	serverInfo, err := ctx.Cookie("remote_addr")
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"code":-1,"msg":err.Error()})
		return
	}

	url := "http://" + serverInfo + config.MultipartUploadCancelUrl + "?fileHash=" + fileHash

	respMsg, err := util.HandleGETAndDELRequest(url, "DELETE")
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"code":-1,"msg":err.Error()})
		return
	}

	ctx.Data(http.StatusOK, "octet-stream", respMsg.JSONBytes())
}