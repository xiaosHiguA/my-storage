package handler

import (
	"MyStorage/util"
	"net/http"
	"time"
)

//HttpInterceptor http请求拦截器
func HttpInterceptor(h http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			r.ParseForm()
			username := r.Form.Get("username")
			token := r.Form.Get("token")
			//验证token
			if len(username) < 3 || !IsTokenValid(token) {
				resp := util.NewRespMsg(
					int(util.StatusInvalidToken),
					"token无效",
					nil,
				)
				w.Write(resp.JsonByte())
				return
			}
			h(w, r)
		})
}

//IsTokenValid 验证token
func IsTokenValid(token string) bool {
	if len(token) != 40 {
		return false
	}
	//验证token是否有效
	tokenString := token[len(token)-8:]
	if util.Hex2Dec(tokenString) < time.Now().Unix()-86400 {
		return false
	}
	return true
}
