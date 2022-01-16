package meta

import "sort"

//FileMeta 源文件信息
type FileMeta struct {
	FileShl  string //
	FileName string // 文件昵称
	FileSize int64  // 文件大小
	Location string // 存放位置
	UploadAt string // 更新信息
}

//用一个map可以查询文件的信息
var fileMetas map[string]*FileMeta

func init() {
	fileMetas = make(map[string]*FileMeta, 10000)
}

// UpdateFileMeta 更新文件
func UpdateFileMeta(fileMeta *FileMeta) {
	//通过fileMeta.FileSha1作为key,每个文件为value
	fileMetas[fileMeta.FileName] = fileMeta
}

//GetFileMeta : 新增/更新文件源信息
func GetFileMeta(fileSha1 string) *FileMeta {
	if value, ok := fileMetas[fileSha1]; ok {
		return value
	}
	//不存在返回空
	return fileMetas[fileSha1]
}

// GetFileMetaList 获取批量文件
func GetFileMetaList(cont int) []FileMeta {
	var ListFileMeta []FileMeta
	//获取元文件
	for _, v := range fileMetas {
		d := v
		ListFileMeta = append(ListFileMeta, *d)
	}
	//给元素进行排序
	sort.Stable(ByUploadTime(ListFileMeta))
	if cont > len(ListFileMeta) { //
		return ListFileMeta
	}
	return ListFileMeta[0:cont]
}

//RemoveFileMetaList 根据元文件名称删除元文件对应
func RemoveFileMetaList(sha1 string) bool {
	//查看当前元文件队列中是否存
	if v, ok := fileMetas[sha1]; ok {
		delete(fileMetas, v.FileShl)
	}
	return false
}

func Doc() {

}
