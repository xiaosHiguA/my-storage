package handler

import (
	"MyStorage/persistentlayer"
	"MyStorage/util"

	"MyStorage/meta"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"time"
)

//UpFileLoaHandler 上传文件
func UpFileLoaHandler(writer http.ResponseWriter, request *http.Request) {
	//判断请求是Get 请求 渲染 上传页面
	if request.Method == "GET" {
		data, err := ioutil.ReadFile("./static/view/index.html")
		if err != nil {
			io.WriteString(writer, "interNel server error")
			return
		}
		io.WriteString(writer, string(data))
	} else if request.Method == "POST" {
		file, head, err := request.FormFile("file")
		if err != nil {
			fmt.Printf("file to get data ,err %s\n", err.Error())
			return
		}
		defer file.Close()

		// 在相对路径 创建 文件位置
		_ = os.Mkdir("./fileData", os.ModePerm)
		var fileMeta *meta.FileMeta
		//为true 是文件不存在
		//
		fileMeta = &meta.FileMeta{
			FileName: head.Filename,
			Location: "./fileData" + head.Filename,             //文件存放的位置
			UploadAt: time.Now().Format("2006-01-02 15:04:05"), //创建时间
		}

		//在本地创建新的一个文件
		newFile, err := os.Create("./fileData" + fileMeta.Location)
		if err != nil {
			fmt.Printf("failed to create file err:%s\n", err.Error())
			return
		}
		defer newFile.Close()
		//完成文件拷贝
		fileMeta.FileSize, err = io.Copy(newFile, file)
		if err != nil {
			fmt.Printf("file to save data into file ,err: %s\n", err.Error())
			return
		}
		//Seek将下一次在文件上读取或写入的偏移量设置为偏移量
		//newFile.Seek(0, 0) //光标默认在文件开头，设置光标的位置在：距离文件开头
		//计算sha1值
		fileMeta.FileShl = util.FileSha1(newFile)
		meta.UpdateFileMeta(fileMeta)
	}
}

//UploadSucHandler 上传文件 保存到
func UploadSucHandler(writer http.ResponseWriter, request *http.Request) {
	io.WriteString(writer, "Upload finished")
}

// GetFileMetaHandler 获取单个文件元信息
func GetFileMetaHandler(writer http.ResponseWriter, request *http.Request) {
	request.ParseForm()
	fileHash := request.Form["filehash"][0]
	//获取单条信息
	tblFile := persistentlayer.GetFileData(fileHash)

	data, err := json.Marshal(tblFile)
	if err != nil {
		return
	}
	writer.Write(data)
}

// FileQueryHandler 获取文件列表
func FileQueryHandler(writer http.ResponseWriter, request *http.Request) {
	request.ParseForm()
	limit, err := strconv.Atoi(request.Form.Get("limit"))
	if err != nil {
		return
	}
	fileMeta := meta.GetFileMetaList(limit)
	fileByte, err := json.Marshal(fileMeta)
	if err != nil {
		return
	}
	writer.Write(fileByte)
}

// RemoveListFileMetaHandler 删除文件
func RemoveListFileMetaHandler(writer http.ResponseWriter, request *http.Request) {
	request.ParseForm()
	fileHash := request.Form.Get("filehash")
	//查询队列是否存在该文件
	fileData := meta.SelectFileMeta(fileHash)
	ok := meta.RemoveFileMetaList(fileData.FileShl)
	if ok {
		writer.WriteHeader(http.StatusOK)
	}
	writer.WriteHeader(http.StatusInternalServerError)
}

//DownloadHandler 下载元文件
func DownloadHandler(writer http.ResponseWriter, request *http.Request) {
	request.ParseForm()
	fSha1 := request.Form.Get("filehash")
	//根据文件的Key的获取文件
	fm := meta.SelectFileMeta(fSha1)
	f, err := os.Open(fm.Location)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer f.Close()
	data, err := ioutil.ReadAll(f)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	//让浏览器可以是识别返回去的信息
	writer.Header().Set("Content-Type", "application/octect-stream")
	writer.Header().Set("content-disposition", "attachment; filename=\""+fm.FileName+"\"")
	writer.Write(data)
}

//UpFileMetaHandler 更新文件
func UpFileMetaHandler(writer http.ResponseWriter, request *http.Request) {
	request.ParseForm()
	opType := request.Form.Get("op")
	fileHash1 := request.Form.Get("filehash")
	fileName := request.Form.Get("filename")
	if len(opType) == 0 {
		return
	}

	if request.Method != "POST" {
		writer.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	//查询 这个元文件是否存在缓存队列中
	fileMeta := meta.SelectFileMeta(fileHash1)
	//验证结构体的值不为空
	if fileMeta == nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	fileMeta.FileName = fileName
	//
	meta.UpdateFileMeta(fileMeta)
	fileData, err := json.Marshal(fileMeta)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	writer.WriteHeader(http.StatusOK)
	writer.Header().Set("Content-type", "application/json")
	writer.Write(fileData)
}
