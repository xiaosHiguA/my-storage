package model

import "database/sql"

type TblFile struct {
	FileSha1 string       `json:"file_sha_1"`
	FileName string       `json:"file_name"`
	FileSize string       `json:"file_size"`
	FileAddr string       `json:"file_addr"`
	CreateAt sql.NullTime `json:"create_at"`
	UpdateAt sql.NullTime `json:"update_at"`
	Status   int          `json:"status"`
	Ext1     int          `json:"ext_1"`
	Ext2     []string     `json:"ext_2"`
}

func (TblFile) TableName() string {
	return "TblFile"
}
