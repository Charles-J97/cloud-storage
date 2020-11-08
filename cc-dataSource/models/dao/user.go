package dao

import "cc-dataSource/models"

type User struct {
	Username         string `json:"username" db:"username"`
	ServerAddr       string `json:"serverAddr" db:"server_addr"`
	Pwd              string `json:"pwd" db:"user_pwd"`
	Email            string `json:"email" db:"user_email"`
	Phone            string `json:"phone" db:"user_phone"`
	EmailValidated   bool   `json:"emailValidated" db:"email_validated"`
	PhoneValidated   bool   `json:"phoneValidated" db:"phone_validated"`
	SignUpTime       string `json:"signUpTime" db:"sign_up_time"`
	LatestUpdateTime string `json:"latestUpdateTime" db:"latest_update_time"`
}

func GetUserByName(username string) (*User, error) {
	var user User
	err := models.DBConn().Get(&user, "SELECT username, server_addr, user_pwd, user_email, user_phone, email_validated, phone_validated, sign_up_time, latest_update_time FROM user WHERE username = ? AND status = 1 limit 1", username)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
