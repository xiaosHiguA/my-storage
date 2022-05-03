package meta

import (
	"MyStorage/gormdb"
	"MyStorage/model"
	"log"
)

//OnFileUploadFinished 新增/更新文件元信息到mysql中
func OnFileUploadFinished(file *model.TblFile) bool {
	db := gormdb.GetDb()
	result := db.Model(&model.TblFile{}).Create(file).Omit("ext_1", "ext_2")
	if result.Error != nil {
		log.Println("存储文件失败", result.Error.Error())
		return false
	}
	return true
}

// GetFileData 取单个文件信息
func GetFileData(fileSha1 string) (*model.TblFile, error) {
	var tblFile = &model.TblFile{}
	db := gormdb.GetDb()
	if rest := db.Model(&model.TblFile{}).First(&tblFile, "file_sha1=?", fileSha1); rest.Error != nil {
		log.Println("获取单个文件信息错误: ", rest.Error)
		return nil, rest.Error
	}
	return tblFile, nil
}

func OnUserFileUploadFile() *model.TblFile {
	tab := &model.TblFile{}
	db := gormdb.GetDb()
	db.Where("")
	return tab
}
