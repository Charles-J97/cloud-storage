package controller

import (
	"cc-transfer/config"
	util "cc-transfer/middleware"
	"cc-transfer/models/dto"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
)

//文件上传接口 POST
//http.ResponseWriter对象用于向用户返回数据，*http.Request用于接受用户请求
func UploadFileHandler(ctx *gin.Context) {
	//接收文件流及存储到本地目录。客户端默认会以FormFile形式传输文件
	serverAddr, err := ctx.Cookie("remote_addr")
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"code":-1,"msg":err.Error()})
		return
	}
	file, err := ctx.FormFile("file")
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"code":-1,"msg":err.Error()})
		return
	}
	username := ctx.Request.FormValue("username")
	fileHash := ctx.Request.FormValue("fileHash")
	filename := file.Filename
	f, _ := file.Open()
	defer f.Close()
	fileBytes, _ := ioutil.ReadAll(f)

	fileUpload := dto.FileUpload{
		Username: username,
		FileHash: fileHash,
		Filename: filename,
		ServerAddr: serverAddr,
		FileBytes: fileBytes,
	}

	jsonFileUpload, err := json.Marshal(fileUpload)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"code":-1,"msg":err.Error()})
		return
	}

	url := "http://" + serverAddr + config.UploadFileUrl

	respMsg, err := util.HandlePOSTAndPUTRequest(jsonFileUpload, url, "POST")
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"code":-1,"msg":err.Error()})
		return
	}

	ctx.Data(http.StatusOK, "octet-stream", respMsg.JSONBytes())
}

//文件下载接口 GET
func DownloadFileHandler(ctx *gin.Context) {
	username := ctx.Query("username")
	filename := ctx.Query("filename")
	serverInfo, err := ctx.Cookie("remote_addr")
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"code":-1,"msg":err.Error()})
		return
	}

	url := "http://" + serverInfo + config.DownloadFileUrl + "?username=" + username + "&filename=" + filename + "&curServerAddr=" + serverInfo
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"code":-1,"msg":err.Error()})
		return
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"code":-1,"msg":err.Error()})
		return
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"code":-1,"msg":err.Error()})
		return
	}

	ctx.Header("content-disposition", "attachment; filename=\"" + filename + "\"")
	ctx.Data(http.StatusOK, "application/octect-stream", data)
}

//单个文件删除接口 DEL
func DeleteSingleFileHandler(ctx *gin.Context) {
	fileHash := ctx.Query("fileHash")
	serverInfo, err := ctx.Cookie("remote_addr")
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"code":-1,"msg":err.Error()})
		return
	}

	url := "http://" + serverInfo + config.DeleteSingleFileUrl + "?fileHash=" + fileHash

	respMsg, err := util.HandleGETAndDELRequest(url, "DELETE")
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"code":-1,"msg":err.Error()})
		return
	}
	ctx.Data(http.StatusOK, "octet-stream", respMsg.JSONBytes())
}

//多个文件删除接口 DEL
func DeleteBatchFileHandler(ctx *gin.Context) {
	fileHash := ctx.QueryArray("fileHash")
	serverInfo, err := ctx.Cookie("remote_addr")
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"code":-1,"msg":err.Error()})
		return
	}

	url := "http://" + serverInfo + config.DeleteBatchFileUrl + "?fileHash=" + fileHash[0]
	for i:=1; i<len(fileHash); i++ {
		url = url + "&fileHash=" + fileHash[i]
	}

	respMsg, err := util.HandleGETAndDELRequest(url, "DELETE")
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"code":-1,"msg":err.Error()})
		return
	}
	ctx.Data(http.StatusOK, "octet-stream", respMsg.JSONBytes())
}

//获取单个文件元信息接口 GET
func GetSingleFileMetaHandler(ctx *gin.Context) {
	fileHash := ctx.Query("fileHash")
	serverInfo, err := ctx.Cookie("remote_addr")
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"code":-1,"msg":err.Error()})
		return
	}

	url := "http://" + serverInfo + config.GetSingleFileMetaUrl + "?fileHash=" + fileHash

	respMsg, err := util.HandleGETAndDELRequest(url, "GET")
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"code":-1,"msg":err.Error()})
		return
	}

	ctx.Data(http.StatusOK, "octet-stream", respMsg.JSONBytes())
}

//获取批量文件元信息接口 GET
func GetBatchFileMetaHandler(ctx *gin.Context) {
	fileHash := ctx.QueryArray("fileHash")
	serverInfo, err := ctx.Cookie("remote_addr")
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"code":-1,"msg":err.Error()})
		return
	}

	url := "http://" + serverInfo + config.GetBatchFileMetaUrl + "?fileHash=" + fileHash[0]

	respMsg, err := util.HandleGETAndDELRequest(url, "GET")
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"code":-1,"msg":err.Error()})
		return
	}

	ctx.Data(http.StatusOK, "octet-stream", respMsg.JSONBytes())
}