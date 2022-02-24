package model

import "time"

//TblUserFile 用户文件列表
type TblUserFile struct {
	UserName   string    `json:"user_name" gorm:"column:"user_name"`
	FileSha1   string    `json:"file_sha_1"  gorm:"column:"file_sha_1"`
	FileSize   int64     `json:"file_size"  gorm:"column:"file_size"`
	FileName   string    `json:"file_name"  gorm:"column:"file_name"`
	UploadAt   time.Time `json:"upload_at"  gorm:"column:"upload_at"`
	LastUpdate time.Time `json:"last_update"  gorm:"column:"last_update"`
	Status     int       `json:"status"  gorm:"column:"status"`
}

func (TblUserFile) TableName() string {
	return "tbl_user_file"
}
