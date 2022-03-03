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
func GetFileData(fileHash string) *model.TblFile {
	var tblFile = &model.TblFile{}
	db := gormdb.GetDb()
	if rest := db.First(tblFile, "file_sha1=?", fileHash); rest.Error != nil {
		log.Println("获取单个文件信息错误: ", rest.Error)
		return nil
	}
	return tblFile
}

func GetArticleList() []*model.TblFile {
	tab := make([]*model.TblFile, 0)
	db := gormdb.GetDb()
	db.Where("")
	return tab
}
