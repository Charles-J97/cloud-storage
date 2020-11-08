package controller

import (
	util "cc-dataSource/middleware"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Test(ctx *gin.Context) {
	username, _ := ctx.Cookie("username")
	token, _ := ctx.Cookie("token")
	server, _ := ctx.Cookie("remote_addr")
	res := util.RespMsg{
		Code: 0,
		Msg:  "OK",
		Data: struct {
			Username string
			Token    string
			Server   string
		}{
			Username: username,
			Token:    token,
			Server:   server,
		},
	}
	ctx.Data(http.StatusOK, "octet-stream", res.JSONBytes())
}
