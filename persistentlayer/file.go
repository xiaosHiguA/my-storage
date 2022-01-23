package persistentlayer

import (
	"MyStorage/db"
	"MyStorage/model"
	"log"
)

func OnFileUploadFinished(file *model.TblFile) bool {
	db := db.GetDb()
	result := db.Create(file).Omit("ext_1", "ext_2")
	if result.Error != nil {
		log.Println("存储文件失败", result.Error.Error())
		return false
	}
	return true
}

// GetFileData 取文件信息
func GetFileData(fileHash string) *model.TblFile {
	var tblFile *model.TblFile
	db := db.GetDb()
	if rest := db.Find(tblFile, "file_sha_1= ?", fileHash); rest.Error != nil {
		return nil
	}
	return tblFile
}
