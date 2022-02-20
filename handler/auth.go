package handler

import (
	"MyStorage/meta"
	"net/http"
)

//HttpInterceptor http请求拦截器
func HttpInterceptor(h http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			r.ParseForm()
			username := r.Form.Get("username")
			token := r.Form.Get("token")
			//验证token
			if len(username) < 3 || !meta.AnalysisToken(token) {
				w.WriteHeader(http.StatusForbidden) //返回403
				return
			}
			h(w, r)
		})
}
