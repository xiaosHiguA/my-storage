package main

import (
	"MyStorage/handler"

	"fmt"
	"net/http"
)

func main() {
	// 静态资源处理
	http.Handle("/static/",
		http.StripPrefix("/static/",
			http.FileServer(http.Dir("./static"))))
	//文件存储模块
	http.HandleFunc("/file/upload", handler.UpFileLoaHandler)
	http.HandleFunc("/file/meta", handler.GetFileMetaHandler)
	http.HandleFunc("/file/query", handler.HttpInterceptor(handler.FileQueryHandler))
	http.HandleFunc("/file/downloadFile", handler.DownloadHandler)
	http.HandleFunc("/file/upload/upFile", handler.UpFileMetaHandler)
	http.HandleFunc("/file/delete", handler.RemoveListFileMetaHandler)

	//用户模块
	http.HandleFunc("/user/signup", handler.RegisterHandler)
	http.HandleFunc("/user/login", handler.TblUserLoginHandle)
	http.HandleFunc("/user/info", handler.HttpInterceptor(handler.UserInfoHandler))
	err := http.ListenAndServe(":8083", nil)
	if err != nil {
		fmt.Printf("Failed to start server ,err: %s", err.Error())
	}
	fmt.Println("启动监听：...")
}
