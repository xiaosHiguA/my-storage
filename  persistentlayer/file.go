package persistentlayer

import (
	"MyStorage/db"
	"MyStorage/model"
	"log"
)

func OnFileUploadFinished(file *model.TblFile) bool {
	db := db.GetDb()
	defer db.Close()
	result := db.Create(file)
	if result.Error != nil {
		log.Println("存储文件失败",result.Error.Error())
		return false
	}
	return true
}
