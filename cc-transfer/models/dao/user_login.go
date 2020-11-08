package dao

import (
	"cc-transfer/models"
	"fmt"
)

type UserLoginTIme struct {
	Username string `json:"username" db:"username"`
	LoginTime string `json:"loginTime" db:"login_time"`
	Token string `json:"token" db:"token"`
	ServerAddr string `json:"serverAddr" db:"server_addr"`
}

func InsertUserLoginTime(username, token, serverAddr string) error {
	_, err := models.DBConn().Exec("INSERT INTO user_login (`username`, `token`, `server_addr`) values (?,?,?)", username, token, serverAddr)
	if err != nil {
		fmt.Printf("Data inserts failed, error:[%v]", err.Error())
		return err
	} else {
		fmt.Println("Data inserts successfully")
	}
	return nil
}
