package meta

import (
	"MyStorage/gormdb"
	"MyStorage/model"
	"log"
)

//OnUserFileUploadFinished 更新用户文件表
func OnUserFileUploadFinished(tblUserFile *model.TblUserFile) bool {
	db := gormdb.GetDb()
	if err := db.Create(tblUserFile).Error; err != nil {
		log.Println("update user_file err: ", err)
		return false
	}
	return true
}
