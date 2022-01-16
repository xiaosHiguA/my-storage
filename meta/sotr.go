package meta

import "time"

const baseFormat = "2006-01-02 15:04:05"

/*
	go 实现基于对象的排序， 需要实现 "interface"的接口
*/
type ByUploadTime []FileMeta

// Len 长度
func (a ByUploadTime) Len() int {

	return len(a)
}

//Swap 将的列表元素进行对换
func (a ByUploadTime) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

//Less 按照上传的时间进行排序
func (a ByUploadTime) Less(i, j int) bool {
	iTime, _ := time.Parse(baseFormat, a[i].UploadAt) //按照指定格式解析
	jTime, _ := time.Parse(baseFormat, a[j].UploadAt) //按照指定格式解析
	return iTime.UnixNano() > jTime.UnixNano()
}
