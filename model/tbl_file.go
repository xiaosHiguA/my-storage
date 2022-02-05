package model

import (
	"time"
)

// TblFile 文件存储
type TblFile struct {
	FileSha1 string    `json:"file_sha_1" gorm:"column:"file_sha1"` //文件哈希后值
	FileName string    `json:"file_name" gorm:"column:"file_name"`  //文件昵称
	FileSize int64     `json:"file_size"  gorm:"column:"file_size"` //文件大小
	FileAddr string    `json:"file_addr"  gorm:"column:"file_addr"` //存放位置
	CreateAt time.Time `json:"create_at"  gorm:"column:"create_at"` //创建时间
	UpdateAt time.Time `json:"update_at"  gorm:"column:"update_at"` //修改时间
	Status   int       `json:"status"  gorm:"column:"status"`       //文件状态
	Ext1     int       `json:"ext_1"  gorm:"column:"ext_1"`
	Ext2     string    `json:"ext_2" gorm:"column:"ext_2"`
}

func (TblFile) TableName() string {
	return "tbl_file"
}
