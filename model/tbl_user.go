package model

type TblUser struct {
	Id             int    `json:"id"`
	UserName       string `json:"user_name"`
	UserPwd        string `json:"user_pwd"`
	Email          string `json:"email"`
	Phone          string `json:"phone"`
	EmailValidated string `json:"email_validated"`
	PhoneValidated string `json:"phone_validated"`
	SignupAt       string `json:"signup_at"`
	LastActive     string `json:"last_active"`
	Profile        string `json:"profile"`
	Status         string `json:"status"`
}
