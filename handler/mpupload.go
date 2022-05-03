package handler

import (
	"MyStorage/conredis"
	"MyStorage/meta"
	"MyStorage/model"
	"MyStorage/util"
	"fmt"
	"github.com/garyburd/redigo/redis"
	"log"
	"math"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"
	"time"
)

//MultipartUploadInfo 分块上传的文件
type MultipartUploadInfo struct {
	FileHash   string
	FileSize   int
	UploadID   string //每次个分块的文件标识
	ChunkSize  int    //分块上传的大小
	ChunkCount int    //一个文件分块的数量
	//已经完成上传的分块索引列表
	ChunkExists []int
}

const (
	// ChunkDir : 上传的分块所在目录
	ChunkDir = "/data/chunks/"
	// MergeDir : 合并后的文件所在目录
	MergeDir = "/data/merge/"
	// ChunkKeyPrefix : 分块信息对应的redis键前缀
	ChunkKeyPrefix = "MP_"
	// HashUpIDKeyPrefix : 文件hash映射uploadid对应的redis键前缀
	HashUpIDKeyPrefix = "HASH_UPID_"
)

//InitialMultipartUploadHandler : 初始化分块上传
func InitialMultipartUploadHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	userName := r.Form.Get("username")
	fileHash := r.Form.Get("filehash")
	fileSize, err := strconv.Atoi(r.Form.Get("filesize"))
	fmt.Println("file size", fileSize)
	if err != nil {
		w.Write(util.NewRespMsg(-1, "params invalid", nil).JsonByte())
		return
	}

	//判断文件是否上传过用户表
	if meta.IsUserFileUploaded(userName, fileHash) {
		w.Write(util.NewRespMsg(int(util.FileAlreadExists), "文件已存在", nil).JsonByte())
		return
	}

	//获取redis的连接
	rConn := conredis.GetRedisPool().Get()
	defer rConn.Close()
	//通过的文件 hash 判断 是否 断点续传,并获取uploadID
	var uploadID string
	keyExist, _ := redis.Bool(rConn.Do("EXISTS", HashUpIDKeyPrefix+fileHash))
	if keyExist { //如果存在,将上传过文件的uploadID取出
		uploadID, err = redis.String(rConn.Do("GET", HashUpIDKeyPrefix+fileHash))
		if err != nil {
			w.Write(util.NewRespMsg(-1, "Upload part failed", err.Error()).JsonByte())
			return
		}
	}

	var chunksExist []int
	//如果文件没有上传过,就时第一次上传,创建一个uploadID
	if uploadID == "" {
		uploadID = userName + fmt.Sprintf("%x", time.Now().UnixNano())
	} else { // 断点续传则根据uploadID获取已上传过的文件分块列表
		chunks, vaErr := redis.Values(rConn.Do("HGETALL", ChunkKeyPrefix+uploadID))
		if vaErr != nil {
			w.Write(util.NewRespMsg(-2, "Upload part failed", vaErr.Error()).JsonByte())
		}
		for i := 0; i < len(chunks); i += 2 {
			k := string(chunks[i].([]byte))
			v := string(chunks[i].([]byte))
			if strings.HasPrefix(k, "chkidx_") && v == "1" { //校验key 是否以 chkidx_开头
				chunkIdx, _ := strconv.Atoi(k[7:len(k)])
				chunksExist = append(chunksExist, chunkIdx)
			}
		}
	}

	//组装 分块上传的信息
	upload := MultipartUploadInfo{
		FileHash:   fileHash,
		FileSize:   fileSize,
		UploadID:   userName + fmt.Sprintf("%x", time.Now().UnixNano()),   //用户名+当前上传的时间
		ChunkSize:  5 * 1024 * 1024,                                       //每个分块的大小（5M）
		ChunkCount: int(math.Ceil(float64(fileSize) / (5 * 1024 * 1024))), //文件的容量除以每个分块大小(5m) 然后在向上取整=整的文件分块的数量
	}

	if len(upload.ChunkExists) <= 0 {
		hKey := ChunkKeyPrefix + upload.UploadID
		rConn.Do("HSET", ChunkKeyPrefix+upload.UploadID, "chunkCount", upload.ChunkCount)
		rConn.Do("HSET", ChunkKeyPrefix+upload.UploadID, "fileHash", upload.FileHash)
		rConn.Do("HSET", ChunkKeyPrefix+upload.UploadID, "fileSize", upload.FileSize)
		rConn.Do("EXPIRE", hKey, 43200)
		//设置每个hKey 的有效期
		rConn.Do("EXPIRE", hKey, 43200)

		//设置首次上传文件的UploadID key值,并且有效期为半天
		rConn.Do("SET", HashUpIDKeyPrefix+fileHash, upload.UploadID, "EX", 43200)
	}

	// 7. 将响应初始化数据返回到客户端
	w.Write(util.NewRespMsg(0, "OK", upload).JsonByte())
}

//UploadPartHandler 上传将文件分块的接口
func UploadPartHandler(w http.ResponseWriter, r *http.Request) {
	//解析参数
	r.ParseForm()
	uploadID := r.Form.Get("uploadid")
	chunkSha1 := r.Form.Get("chkhash")
	chunkIndex := r.Form.Get("index")

	rConn := conredis.GetRedisPool().Get()
	defer rConn.Close()

	//创建文件路径,用于文件存储
	filePath := ChunkDir + uploadID + "/" + chunkIndex
	err := os.MkdirAll(path.Dir(filePath), 0744)
	log.Println("mkdir filePath err: ", err)
	fd, err := os.Create(filePath)
	if err != nil {
		log.Println("create: ", err)
		w.Write(util.NewRespMsg(-1, "Upload part faild", nil).JsonByte())

		return
	}
	buf := make([]byte, 1024*1024)
	for {
		n, err := r.Body.Read(buf) //获取请求里的文件
		fd.Write(buf[:n])          //将文件写入创建好的文件中
		if err != nil {
			break
		}
	}

	// 校验分块hash
	cmdSH1, err := util.ComputeSha1ByShell(filePath)
	if err != nil || cmdSH1 != chunkSha1 {
		fmt.Printf("Verify chunk sha1 failed, compare OK: %t, err:%+v\n", cmdSH1 == chunkSha1, err)
		w.Write(util.NewRespMsg(-2, "Verify hash failed, chkIdx:"+chunkIndex, nil).JsonByte())
		return
	}

	//更新redis缓存状态
	rConn.Do("HSET", ChunkKeyPrefix+filePath, uploadID, "chkidx_"+chunkIndex, 1)

	w.Write(util.NewRespMsg(0, "ok", nil).JsonByte())
	//数据响应
}

// CompleteUploadHandler : 通知上传合并
func CompleteUploadHandler(w http.ResponseWriter, r *http.Request) {
	// 1. 解析请求参数
	r.ParseForm()
	uploadID := r.Form.Get("uploadid")
	username := r.Form.Get("username")
	filehash := r.Form.Get("filehash")
	filesize := r.Form.Get("filesize")
	filename := r.Form.Get("filename")

	rConn := conredis.GetRedisPool().Get()
	defer rConn.Close()
	//根据uploadID上传分块是否上传完毕
	data, err := redis.Values(rConn.Do("HGETALL", ChunkKeyPrefix+uploadID))
	if err != nil {
		w.Write(util.NewRespMsg(-1, "complete upload failed", nil).JsonByte())
		return
	}
	totalCount := 0 //整个文件的容量
	chunkCount := 0 //分块的数量
	for i := 0; i < len(data); i += 2 {
		k := string(data[i].([]byte))
		v := string(data[i].([]byte))
		if k == "chunkCount" {
			totalCount, _ = strconv.Atoi(v)
		} else if strings.HasPrefix(k, "chkidx_") && v == "1" { //判断是否以chkidx_开头
			chunkCount++ //统计分块上传的数量
		}
	}
	if totalCount != chunkCount {
		w.Write(util.NewRespMsg(-2, "Invalid request", nil).JsonByte())
		return
	}

	//更新唯一用户文件表以及文件表
	fileSize, _ := strconv.Atoi(filesize)
	var fileData = &model.TblFile{
		FileSha1: filehash,
		FileName: filename,
		FileSize: int64(fileSize),
	}
	if !meta.OnFileUploadFinished(fileData) {
		log.Println("更新文件表失败")
		w.Write(util.NewRespMsg(-2, "合并上传分块异常", nil).JsonByte())
		return
	}

	var tblUserFile = &model.TblUserFile{
		UserName: username,
		FileSha1: filehash,
		FileName: filename,
		FileSize: int64(fileSize),
	}
	if !meta.OnUserFileUploadFinished(tblUserFile) {
		log.Println("更新唯一用户文件表失败")
		w.Write(util.NewRespMsg(-2, "合并上传分块异常", nil).JsonByte())
		return
	}

	// 更新于2020-04: 删除已上传的分块文件及redis分块信息
	_, delHashErr := rConn.Do("DEL", HashUpIDKeyPrefix+filehash)
	delUploadID, delUploadInfoErr := redis.Int64(rConn.Do("DEL", ChunkKeyPrefix+uploadID))
	if delUploadID != 1 || delUploadInfoErr != nil || delHashErr != nil {
		w.Write(util.NewRespMsg(-4, "Complete upload part failed", nil).JsonByte())
		return
	}

	delRes := util.RemovePathByShell(ChunkDir + uploadID)

	if !delRes {
		fmt.Printf("Failed to delete chuncks as upload comoleted, uploadID: %s\n", uploadID)
	}

	w.Write(util.NewRespMsg(200, "合并成功", nil).JsonByte())
}

// CancelUploadHandler : 文件取消上传接口
func CancelUploadHandler(w http.ResponseWriter, r *http.Request) {
	// 1. 解析用户请求参数
	r.ParseForm()
	filehash := r.Form.Get("filehash")
	// 2. 获得redis的一个连接
	rPool := conredis.GetRedisPool().Get()
	defer rPool.Close()

	// 3. 检查uploadID是否存在，如果存在则删除
	uploadID, err := redis.String(rPool.Do("GET", HashUpIDKeyPrefix+filehash))
	if err != nil || uploadID == "" {
		w.Write(util.NewRespMsg(-1, "Cancel upload part failed", nil).JsonByte())
		return
	}
	//删除分块
	_, delHashErr := rPool.Do("DEL", HashUpIDKeyPrefix+filehash)
	_, delUploadInfoErr := rPool.Do("DEL", ChunkKeyPrefix+uploadID)
	if delHashErr != nil || delUploadInfoErr != nil {
		w.Write(util.NewRespMsg(-2, "Cancel upload part failed", nil).JsonByte())
		return
	}
	// 4. 删除已上传的分块文件
	delChkRes := util.RemovePathByShell(ChunkDir + uploadID)
	if !delChkRes {
		fmt.Printf("Failed to delete chunks as upload canceled, uploadID:%s\n", uploadID)
	}

	// 5. 响应客户端
	w.Write(util.NewRespMsg(0, "OK", nil).JsonByte())
}
