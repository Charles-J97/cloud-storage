package service

import (
	rPool "cc-transfer/cache"
	"cc-transfer/config"
	util "cc-transfer/middleware"
	"cc-transfer/models/dao"
	"cc-transfer/models/dto"
	"errors"
)

//用户注册逻辑
func UserSignUp(username, serverAddr, pwd, email, phone string) error {
	//给用户名和密码制定某些规则
	if len(username) < 6 || len(username) > 16 || len(pwd) < 6 || len(pwd) > 20 {
		return errors.New("Username or password is invalid. ")
	}
	//给密码加盐加密，然后存储到DB中
	encodePwd := util.Sha1([]byte(pwd + config.PwdSalt))
	err := dao.InsertUser(username, serverAddr, encodePwd, email, phone)
	if err != nil {
		return err
	}
	return nil
}

//用户登录逻辑
func UserLogin(username, pwd string) (string, string, error) {
	encodedPwd := util.Sha1([]byte(pwd + config.PwdSalt))
	//若是查询用户不存在，会直接返回在error里
	user, err := dao.GetUserByName(username)
	if err != nil {
		return "", "", err
	}
	if user.Pwd != encodedPwd {
		return "", "", errors.New("Password is wrong ")
	}

	//若用户存在，且密码正确，则把用户名通过AES双向加密后，返回其token
	token, err := util.Encrypt([]byte(username))
	if err != nil {
		return "", "", err
	}

	//查找该用户对应的服务器地址
	serverInfo, err := dto.GlobalHashCircle.Get(username)
	if err != nil {
		return "", "", err
	}

	//将用户名，登录时间，所用token和该用户对应的服务器地址记录在mysql中
	err = dao.InsertUserLoginTime(username, token, serverInfo)
	if err != nil {
		return "", "", err
	}

	//将token存在redis中，来保证其时效性
	rConn := rPool.RedisPool().Get()
	defer rConn.Close()
	_, err = rConn.Do("SET", config.TokenPre + username, token)
	if err != nil {
		return "", "", err
	}
	_, err = rConn.Do("EXPIRE", config.TokenPre + username, config.TokenMaxAge)
	if err != nil {
		return "", "", err
	}

	return token, serverInfo, nil
}

//用户信息查询逻辑
func UserInfoQuery(username string) (*dao.User, error) {
	userInfo, err := dao.GetUserByName(username)
	if err != nil {
		return nil, err
	}
	return userInfo, nil
}