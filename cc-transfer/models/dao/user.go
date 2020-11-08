package dao

import (
	"cc-transfer/models"
	"fmt"
)

type User struct {
	Username string `json:"username" db:"username"`
	ServerAddr string `json:"serverAddr" db:"server_addr"`
	Pwd string `json:"pwd" db:"user_pwd"`
	Email string `json:"email" db:"user_email"`
	Phone string `json:"phone" db:"user_phone"`
	EmailValidated bool `json:"emailValidated" db:"email_validated"`
	PhoneValidated bool `json:"phoneValidated" db:"phone_validated"`
	SignUpTime string `json:"signUpTime" db:"sign_up_time"`
	LatestUpdateTime string `json:"latestUpdateTime" db:"latest_update_time"`
}

func InsertUser(username, serverAddr, pwd, email, phone string) error {
	_, err := models.DBConn().Exec("INSERT INTO user (`username`, `server_addr`, `user_pwd`, `user_email`, `user_phone`, `status`) values (?,?,?,?,?,1)", username, serverAddr, pwd, email, phone)
	if err != nil {
		fmt.Printf("Data inserts failed, error:[%v]", err.Error())
		return err
	} else {
		fmt.Println("Data inserts successfully")
	}
	return nil
}

func GetUserByName(username string) (*User, error) {
	var user User
	err := models.DBConn().Get(&user, "SELECT username, server_addr, user_pwd, user_email, user_phone, email_validated, phone_validated, sign_up_time, latest_update_time FROM user WHERE username = ? AND status = 1 limit 1", username)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func DeleteUserByName(username string) error {
	_, err := models.DBConn().Exec("UPDATE user SET status = 0 WHERE username = ? AND status = 1", username)
	if err != nil {
		return err
	}
	return nil
}