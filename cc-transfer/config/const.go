package config

const (
	PwdSalt = "#jzh666"			//注册密码加密盐值
	TokenPre = "token_"			//token在redis中存储的key前缀，TokenPre + username组成token在redis中存储的key前缀
	TokenMaxAge = 36000			//token在cookie和redis中存在的最大时长，为10h
)
