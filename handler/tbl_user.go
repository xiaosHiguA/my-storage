package handler

import (
	"MyStorage/meta"
	"MyStorage/model"
	"MyStorage/util"
	"net/http"
)

//RegisterHandler  用户注册
func RegisterHandler(writer http.ResponseWriter, request *http.Request) {
	if request.Method == http.MethodGet {
		http.Redirect(writer, request, "/static/view/signup.html", http.StatusFound)
		return
	}
	var tblUser = new(model.TblUser)
	request.ParseForm()
	tblUser.UserName = request.Form.Get("username")
	tblUser.UserPwd = request.Form.Get("password")
	if len(tblUser.UserName) < 3 || len(tblUser.UserPwd) < 5 {
		writer.Write([]byte("Invalid parameter"))
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

//TblUserLoginHandle : 用户登录
func TblUserLoginHandle(writer http.ResponseWriter, request *http.Request) {
	if request.Method == http.MethodGet {
		http.Redirect(writer, request, "/static/view/signin.html", http.StatusFound)
	}
	request.ParseForm()
	userName := request.Form.Get("username")
	userPwd := request.Form.Get("password")
	tblUser := meta.GetUser(userName)
	if tblUser.UserName == "" {
		resp := util.RespMsg{
			Code: 500,
			Msg:  "用户不存在",
			Data: nil,
		}
		writer.Write(resp.JsonByte())
		return
	}
	//验证用户密码
	pwd := util.TblUserMD5(userPwd)
	if len(pwd) != len(tblUser.UserPwd) {
		resp := util.RespMsg{
			Code: 500,
			Msg:  "密码错误",
		}
		writer.Write(resp.JsonByte())
		return
	}
	//生成token
	token := util.GenToken(tblUser.UserName)
	var tblUserFile = model.TblUserToken{
		UserName:  userName,
		UserToken: token,
	}
	//查询token 表是否有此用户
	userToken := meta.GetUserToken(userName)
	if userToken.UserName == meta.UserNAMENULL {
		if !meta.SaveToken(&tblUserFile) {
			resp := util.RespMsg{
				Code: 500,
				Msg:  "密码错误",
			}
			writer.Write(resp.JsonByte())
			return
		}
	} else {
		//保存用户token
		if !meta.UpdateToken(tblUserFile) {
			resp := util.RespMsg{
				Code: 500,
				Msg:  "内部错误",
			}
			writer.Write(resp.JsonByte())
			return
		}
	}
	resp := util.RespMsg{
		Code: 0,
		Msg:  "OK",
		Data: struct {
			Location string
			Username string
			Token    string
		}{
			Location: "http://" + request.Host + "/static/view/home.html",
			Username: userName,
			Token:    token,
		},
	}
	writer.Write(resp.JsonByte())
}

// UserInfoHandler 获取用户信息
func UserInfoHandler(writer http.ResponseWriter, request *http.Request) {
	request.ParseForm()
	userName := request.Form.Get("username")
	tblUser := meta.GetUserToken(userName)
	if tblUser != nil {
		resp := util.RespMsg{
			Code: 0,
			Msg:  "ok",
			Data: tblUser,
		}
		writer.Write(resp.JsonByte())
	}
	writer.WriteHeader(http.StatusForbidden)
}
