package meta

import (
	"MyStorage/gormdb"
	"MyStorage/model"
	"log"
)

//OnUserFileUploadFinished 更新用户文件表
func OnUserFileUploadFinished(tblUserFile *model.TblUserFile) bool {
	db := gormdb.GetDb()
	if err := db.Model(&model.TblUserFile{}).Create(tblUserFile).Error; err != nil {
		log.Println("update user_file err: ", err)
		return false
	}
	return true
}

//IsUserFileUploaded  根据用户名与文件查询对应的文件名
func IsUserFileUploaded(username, fileHash string) bool {
	var userFile model.TblUserFile
	db := gormdb.GetDb()
	err := db.Where("user_name=? and file_sha1=? and status= 1", username, fileHash).Limit(1).First(&userFile).Error
	if err != nil {
		log.Println("select tblUserFile: ", err)
		return false
	}
	return true
}

// QueryUserFileMetas ：批量获取用户文件的信息
func QueryUserFileMetas(userName string, limit int) ([]model.TblUserFile, error) {
	var tblUserFile []model.TblUserFile
	db := gormdb.GetDb()
	err := db.Model(&model.TblUserFile{}).Where("user_name=?", userName).Limit(limit).First(&tblUserFile).Error
	if err != nil {
		log.Println("select limit tbl_user_file err: ", err)
	}
	return tblUserFile, err
}
