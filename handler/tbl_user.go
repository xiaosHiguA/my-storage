package handler

import (
	"MyStorage/meta"
	"MyStorage/model"
	"net/http"
)

func SignupHandler(writer http.ResponseWriter, request *http.Request) {
	var tblUser = new(model.TblUser)
	request.ParseForm()
	tblUser.UserName = request.Form.Get("tbl_user")
	tblUser.UserPwd = request.Form.Get("tbl_pass")
	if len(tblUser.UserName) < 3 || len(tblUser.UserPwd) < 5 {
		writer.Write([]byte("lnvalid paramter"))
		return
	}
	//查询用户是否存在
	tblName := meta.GetTbUser(tblUser.UserName)
	if tblName == meta.UserNAMENULL {
		if ok := meta.CreateUser(tblUser); ok {
			writer.Write([]byte("用户注册成功"))
		}
	}
	writer.Write([]byte("用户已存在"))
}
