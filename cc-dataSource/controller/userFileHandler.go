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

//查询用户文件信息接口 GET
func GetUserFileInfoHandler(ctx *gin.Context) {
	count, _ := strconv.Atoi(ctx.Query("count"))
	username := ctx.Query("username")

	userFiles, err := services.GetUserFileInfo(username, int64(count))
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"code": -1, "msg": err.Error()})
		return
	}
	res := util.RespMsg{
		Code: 0,
		Msg:  "OK",
		Data: userFiles,
	}
	ctx.Data(http.StatusOK, "octet-stream", res.JSONBytes())
}

//修改（重命名）用户文件信息接口 PUT
func UpdateUserFilenameHandler(ctx *gin.Context) {
	var userFile dto.UpdateUserFilename
	err := ctx.ShouldBindBodyWith(&userFile, binding.JSON)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"code": -1, "msg": err.Error()})
		return
	}

	err = services.UpdateUserFilename(userFile.Username, userFile.Filename, userFile.FileNewName)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"code": -1, "msg": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"code": 0, "msg": "OK"})
}

//删除多个用户文件信息接口 DELETE
func DeleteSingleUserFileHandler(ctx *gin.Context) {
	username := ctx.Query("username")
	filename := ctx.Query("filename")

	err := services.DeleteUserFile(username, []string{filename})
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"code": -1, "msg": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"code": 0, "msg": "OK"})
}

//删除单个用户文件信息接口 DELETE
func DeleteBatchUserFileHandler(ctx *gin.Context) {
	username := ctx.Query("username")
	filename := ctx.QueryArray("filename")

	err := services.DeleteUserFile(username, filename)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"code": -1, "msg": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"code": 0, "msg": "OK"})
}

//将该username中指定名字的文件共享给另一username POST
func ShareUserFileHandler(ctx *gin.Context) {
	var userFile dto.ShareUserFile
	err := ctx.ShouldBindBodyWith(&userFile, binding.JSON)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"code": -1, "msg": err.Error()})
		return
	}

	err = services.ShareUserFile(userFile.Username, userFile.Filename, userFile.NewUsername, userFile.ServerAddr)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"code": -1, "msg": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"code": 0, "msg": "OK"})
}
