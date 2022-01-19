package main

import (
	"MyStorage/handler"
	"fmt"
	"net/http"
)

func main() {
	http.Handle("/static/",
		http.StripPrefix("/static/",
			http.FileServer(http.Dir("./static"))))

	http.HandleFunc("/file/upload", handler.UpFileLoaHandler)
	http.HandleFunc("/file/upload/suc", handler.UploadSucHandler)
	http.HandleFunc("/file/meta", handler.GetFileMetaHandler)
	http.HandleFunc("/file/query", handler.FileQueryHandler)
	http.HandleFunc("/file/downloadFile", handler.DownloadHandler)
	http.HandleFunc("/file/upload/upFile", handler.UpFileMetaHandler)
	http.HandleFunc("/file/delete", handler.RemoveListFileMetaHandler)
	err := http.ListenAndServe(":8083", nil)
	if err != nil {
		fmt.Printf("Failed to start server ,err: %s", err.Error())
	}
}
