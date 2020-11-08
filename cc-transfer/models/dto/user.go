package dto

type UserSignUp struct {
	Username string `json:"username"`
	Pwd string `json:"pwd"`
	Email string `json:"email"`
	Phone string `json:"phone"`
}