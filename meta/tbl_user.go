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

//UpdateToken 更新用户token
func UpdateToken(tblUserFile model.TblUserToken) bool {
	db := gormdb.GetDb()
	if err := db.Model(&tblUserFile).Where("user_name=?", tblUserFile.UserName).Update("user_token", tblUserFile.UserToken).Error; err != nil {
		log.Println("UpdateToken user err: ", err)
		return false
	}
	return true
}

func GetUserToken(userName string) *model.TblUserToken {
	var userToken = model.TblUserToken{}
	db := gormdb.GetDb()
	if err := db.Where("user_name=?", userName).First(&userToken).Error; err != nil {
		log.Println("get token err: ", err)
		return nil
	}
	return &userToken
}

func SaveToken(tblUserFile *model.TblUserToken) bool {
	db := gormdb.GetDb()
	if err := db.Create(tblUserFile).Error; err != nil {
		log.Println("SaveToken token err: ", err)
		return false
	}
	return true
}
