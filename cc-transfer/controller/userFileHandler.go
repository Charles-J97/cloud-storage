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

//查询用户文件信息接口 GET
func GetUserFileInfoHandler(ctx *gin.Context) {
	count := ctx.Query("count")
	username := ctx.Query("username")
	serverAddr, err := ctx.Cookie("remote_addr")
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"code":-1,"msg":err.Error()})
		return
	}

	url := "http://" + serverAddr + config.GetUserFileInfoUrl + "?username=" + username + "&count=" + count + "&serverAddr=" + serverAddr

	respMsg, err := util.HandleGETAndDELRequest(url, "GET")
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"code":-1,"msg":err.Error()})
		return
	}

	ctx.Data(http.StatusOK, "octet-stream", respMsg.JSONBytes())
}

//修改（重命名）用户文件信息接口 PUT
func UpdateUserFileNameHandler(ctx *gin.Context) {
	serverAddr, err := ctx.Cookie("remote_addr")
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"code":-1,"msg":err.Error()})
		return
	}
	jsonBody, err := ioutil.ReadAll(ctx.Request.Body)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"code":-1,"msg":err.Error()})
		return
	}
	var updateUserFileName dto.UpdateUserFilename
	err = json.Unmarshal(jsonBody, &updateUserFileName)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"code":-1,"msg":err.Error()})
		return
	}
	updateUserFileName.ServerAddr = serverAddr
	jsonBody, err = json.Marshal(updateUserFileName)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"code":-1,"msg":err.Error()})
		return
	}

	url := "http://" + serverAddr + config.UpdateUserFilenameUrl

	respMsg, err := util.HandlePOSTAndPUTRequest(jsonBody, url, "PUT")
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"code":-1,"msg":err.Error()})
		return
	}

	ctx.Data(http.StatusOK, "octet-stream", respMsg.JSONBytes())
}

//删除单个用户文件信息接口 DELETE
func DeleteSingleUserFileHandler(ctx *gin.Context) {
	username := ctx.Query("username")
	filename := ctx.Query("filename")
	serverInfo, err := ctx.Cookie("remote_addr")
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"code":-1,"msg":err.Error()})
		return
	}

	url := "http://" + serverInfo + config.DeleteSingleUserFileUrl + "?username=" + username + "&filename=" + filename

	respMsg, err := util.HandleGETAndDELRequest(url, "DELETE")
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"code":-1,"msg":err.Error()})
		return
	}

	ctx.Data(http.StatusOK, "octet-stream", respMsg.JSONBytes())
}

//删除多个用户文件信息接口 DELETE
func DeleteBatchUserFileHandler(ctx *gin.Context) {
	username := ctx.Query("username")
	filename := ctx.QueryArray("filename")
	serverInfo, err := ctx.Cookie("remote_addr")
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"code":-1,"msg":err.Error()})
		return
	}

	url := "http://" + serverInfo + config.DeleteBatchUserFileUrl + "?username=" + username
	for _, item := range filename {
		url = url + "&filename=" + item
	}

	respMsg, err := util.HandleGETAndDELRequest(url, "GET")
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"code":-1,"msg":err.Error()})
		return
	}

	ctx.Data(http.StatusOK, "octet-stream", respMsg.JSONBytes())
}

//将该username中指定名字的文件共享给另一username POST
func ShareUserFileHandler(ctx *gin.Context) {
	serverInfo, err := ctx.Cookie("remote_addr")
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"code":-1,"msg":err.Error()})
		return
	}

	var userFile dto.ShareUserFile
	err = ctx.ShouldBindBodyWith(&userFile, binding.JSON)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"code": -1, "msg": err.Error()})
		return
	}
	userFile.ServerAddr = serverInfo

	jsonBody, err := json.Marshal(userFile)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"code":-1,"msg":err.Error()})
		return
	}

	url := "http://" + serverInfo + config.ShareUserFileUrl

	respMsg, err := util.HandlePOSTAndPUTRequest(jsonBody, url, "POST")
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"code":-1,"msg":err.Error()})
		return
	}

	ctx.Data(http.StatusOK, "octet-stream", respMsg.JSONBytes())
}
