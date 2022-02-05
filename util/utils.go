package util

import (
	"crypto/md5"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io"
	"os"
)

//
func FileSha1(file *os.File) string {
	_sha1 := sha1.New()
	io.Copy(_sha1, file)
	return hex.EncodeToString(_sha1.Sum(nil))
}

func TblUser(str string) string {
	w := md5.New()
	io.WriteString(w, str)
	//将str写入到w中
	return fmt.Sprintf("%x", w.Sum(nil))
}
