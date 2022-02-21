package util

import (
	"crypto/md5"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"strconv"
	"time"
)

const (
	Salt     = "_*Salt"
	TokenKey = "_*#102token"
)

//FileSha1 计算hash
func FileSha1(file *os.File) string {
	_sha1 := sha1.New()
	io.Copy(_sha1, file)
	return hex.EncodeToString(_sha1.Sum(nil))
}

func TblUserMD5(str string) string {
	w := md5.New()
	io.WriteString(w, str+Salt)
	//将str写入到w中
	return fmt.Sprintf("%x", w.Sum(nil))
}

// Hex2Dec : 十六进制转十进制
func Hex2Dec(val string) int64 {
	n, err := strconv.ParseInt(val, 16, 0)
	if err != nil {
		fmt.Println(err)
	}
	return n
}

func Sha1(data []byte) string {
	_sha1 := sha1.New()
	_sha1.Write(data)
	return hex.EncodeToString(_sha1.Sum([]byte("")))
}

func TokenMD5(data []byte) string {
	_md5 := md5.New()
	_md5.Write(data)
	return hex.EncodeToString(_md5.Sum([]byte("")))
}

//GenToken 生成token
func GenToken(userName string) string {
	// 40位字符:md5(username+timestamp+token_salt)+timestamp[:8]
	ts := fmt.Sprintf("%x", time.Now().Unix())
	token := TokenMD5([]byte(userName + ts + TokenKey))
	return token + ts[:8]
}
