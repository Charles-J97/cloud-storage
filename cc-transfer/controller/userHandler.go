package controller

import (
	"cc-transfer/config"
	util "cc-transfer/middleware"
	"cc-transfer/models/dto"
	"cc-transfer/service"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"net/http"
)

//用户注册接口 POST
func UserSignUpHandler(ctx *gin.Context) {
	var user dto.UserSignUp
	err := ctx.ShouldBindBodyWith(&user, binding.JSON)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"code":-1,"msg":err.Error()})
		return
	}

	//通过分布式一致性Hash算法求得该用户对应的服务器地址是多少
	serverAddr, err := dto.GlobalHashCircle.Get(user.Username)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"code":-1,"msg":err.Error()})
		return
	}
	err = service.UserSignUp(user.Username, serverAddr, user.Pwd, user.Email, user.Phone)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"code":-1,"msg":err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"code":0,"msg":"OK"})
}

//用户登录接口 GET
func UserLogInHandler(ctx *gin.Context) {
	username := ctx.Query("username")
	pwd := ctx.Query("pwd")

	token, serverInfo, err := service.UserLogin(username, pwd)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"code":-1,"msg":err.Error()})
		return
	}
	ctx.SetCookie("username", username, config.TokenMaxAge, "/", "127.0.0.1", false, true)
	ctx.SetCookie("token", token, config.TokenMaxAge, "/", "127.0.0.1", false, true)
	ctx.SetCookie("remote_addr", serverInfo, config.TokenMaxAge, "/", "127.0.0.1", false, true)
	ctx.JSON(http.StatusOK, gin.H{"code":0,"msg":"OK"})
}

//用户信息查询接口 GET
func UserInfoQueryHandler(ctx *gin.Context) {
	username := ctx.Query("username")
	userInfo, err := service.UserInfoQuery(username)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{"code":-1,"msg":err.Error()})
		return
	}
	res := util.RespMsg{
		Code: 0,
		Msg: "OK",
		Data: userInfo,
	}
	ctx.Data(http.StatusOK, "octet-stream", res.JSONBytes())
}