package middleware

import (
	rPool "cc-transfer/cache"
	"cc-transfer/config"
	"cc-transfer/models/dao"
	"cc-transfer/models/dto"
	"github.com/garyburd/redigo/redis"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Auth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		username, err := ctx.Cookie("username")
		if err != nil {
			ctx.Abort()
			ctx.JSON(http.StatusOK, gin.H{"code":-1,"msg":"Cookie Parsing Error"})
			return
		}
		token, err := ctx.Cookie("token")
		if err != nil {
			ctx.Abort()
			ctx.JSON(http.StatusOK, gin.H{"code":-1,"msg":"Cookie Parsing Error"})
			return
		}

		//验证token是否合法（合法包括两方面，是否正确，以及是否过期）
		valid, err := isTokenValid(username, token)
		if err != nil {
			ctx.Abort()
			ctx.JSON(http.StatusOK, gin.H{"code":-1,"msg":"Token Parsing Error"})
			return
		}
		if !valid {
			ctx.Abort()
			ctx.JSON(http.StatusOK, gin.H{"code":-1,"msg":"Please re-login"})
			return
		}

		//确保当前用户使用的server目前存在
		exists, serverAddr, err := checkServer(username)
		if err != nil {
			ctx.Abort()
			ctx.JSON(http.StatusOK, gin.H{"code":-1,"msg":"Token Parsing Error"})
			return
		}
		//如果不存在，更新当前cookie
		if !exists {
			ctx.Abort()
			ctx.SetCookie("remote_addr", serverAddr, config.TokenMaxAge, "/", "127.0.0.1", false, true)
			ctx.JSON(http.StatusOK, gin.H{"code":-1,"msg":"Remote Server Address Parsing Error"})
			return
		}
		ctx.Next()
	}
}

func isTokenValid(username, token string) (bool, error) {
	parsingData, err := Decrypt(token)
	if err != nil {
		return false, err
	}
	parsingName := string(parsingData[:])
	if username != parsingName {
		return false, nil
	}

	rConn := rPool.RedisPool().Get()
	defer rConn.Close()

	exists, err := redis.Bool(rConn.Do("EXISTS", config.TokenPre + username))
	if err != nil {
		return false, err
	}

	return exists, nil
}

//如果发现这次的serverAddr和MySQL中记录的不一致，则更新MySQL中记录
func checkServer(username string) (bool, string, error) {
	user, err := dao.GetUserByName(username)
	if err != nil {
		return false, "", err
	}
	serverAddr, err := dto.GlobalHashCircle.Get(username)
	if err != nil {
		return false, "", err
	}

	if user.ServerAddr != serverAddr {
		err = dao.DeleteUserByName(username)
		if err != nil {
			return false, "", err
		}
		err = dao.InsertUser(username, serverAddr, user.Pwd, user.Email, user.Phone)
		return false, serverAddr, nil
	}

	return true, serverAddr, nil
}
