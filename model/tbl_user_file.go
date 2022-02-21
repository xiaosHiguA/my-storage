package model

type TblUserFile struct {
	Id        int    `json:"id"`
	UserName  string `json:"user_name"`
	UserToken string `json:"user_token"`
}

func (t *TblUserFile) TableName() string {
	return "tbl_user_file"
}
