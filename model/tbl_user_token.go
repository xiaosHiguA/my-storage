package model

//TblUserToken 列表
type TblUserToken struct {
	Id        int    `json:"id"`
	UserName  string `json:"user_name"`
	UserToken string `json:"user_token"`
}

func (t *TblUserToken) TableName() string {
	return "tbl_user_token"
}
