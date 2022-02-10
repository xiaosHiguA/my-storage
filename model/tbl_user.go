package model

import "time"

type TblUser struct {
	Id             int       `json:"id"`
	UserName       string    `json:"user_name"`
	UserPwd        string    `json:"user_pwd"`
	Email          string    `json:"email"`
	Phone          string    `json:"phone"`
	EmailValidated string    `json:"email_validated"`
	PhoneValidated string    `json:"phone_validated"`
	SignupAt       time.Time `json:"signup_at"`
	LastActive     time.Time `json:"last_active"`
	Profile        string    `json:"profile"`
	Status         string    `json:"status"`
}

func (TblUser) TableName() string {
	return "tbl_user"
}
