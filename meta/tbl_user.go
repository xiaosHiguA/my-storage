package meta

import (
	"MyStorage/db"
	"MyStorage/model"
	"MyStorage/util"
)

const (
	UserNAMENULL = ""
)

func CreateUser(user *model.TblUser) bool {
	db := db.GetDb()
	//md5加密码
	user.UserPwd = util.TblUser(user.UserPwd)
	if err := db.Create(user).Error; err != nil {
		return false
	}
	return true
}

func GetTbUser(userName string) string {
	var user string
	db := db.GetDb()
	if rest := db.First(user, "tbl_user=?", userName); rest.Error != nil {
		return UserNAMENULL
	}
	return user
}
