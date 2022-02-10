package meta

import (
	"MyStorage/db"
	"MyStorage/model"
	"MyStorage/util"
	"time"
)

const (
	UserNAMENULL = ""
)

func CreateUser(user *model.TblUser) bool {
	db := db.GetDb()
	//md5加密密码
	user.UserPwd = util.TblUser(user.UserPwd)
	user.SignupAt = time.Now()
	user.LastActive = time.Now()
	if err := db.Create(user).Error; err != nil {
		return false
	}
	return true
}

func GetTbUser(userName string) string {
	var tblUser = &model.TblUser{}
	db := db.GetDb()
	db.Take(&tblUser, "user_name= ?", userName)
	if len(tblUser.UserName) <= 1 {
		return UserNAMENULL
	}
	return tblUser.UserName
}
