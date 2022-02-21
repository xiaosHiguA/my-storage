package meta

import (
	"MyStorage/gormdb"
	"MyStorage/model"
	"MyStorage/util"
	"log"
	"time"
)

const (
	UserNAMENULL = ""
)

//CreateUser 用户注册
func CreateUser(user *model.TblUser) bool {
	db := gormdb.GetDb()
	//md5加密密码
	user.UserPwd = util.TblUserMD5(user.UserPwd)
	user.SignupAt = time.Now()
	user.LastActive = time.Now()
	if err := db.Create(user).Error; err != nil {
		return false
	}
	return true
}

//GetTbUser 根据用户名登录
func GetTbUser(userName string) string {
	var tblUser = &model.TblUser{}
	db := gormdb.GetDb()
	db.Take(&tblUser, "user_name= ?", userName)
	if len(tblUser.UserName) <= 1 {
		return UserNAMENULL
	}
	return tblUser.UserName
}

// GetUser 根据用户名 查询用户信息
func GetUser(u string) *model.TblUser {
	var tblUser = &model.TblUser{}
	db := gormdb.GetDb()
	if err := db.Take(tblUser, "user_name=?", u).Error; err != nil {
		log.Println("select user err: ", err)
		return nil
	}
	return tblUser
}

func SaveToken(tblUserFile model.TblUserFile) bool {
	db := gormdb.GetDb()
	if err := db.Create(&tblUserFile).Error; err != nil {
		return false
	}
	return true
}
