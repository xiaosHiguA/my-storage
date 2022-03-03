package util

import (
	"encoding/json"
	"log"
)

const (
	TokenSalt = "_tokenSalt"
)

type RespMsg struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func NewRespMsg(code int, msg string, data interface{}) *RespMsg {
	return &RespMsg{
		Code: code,
		Msg:  msg,
		Data: data,
	}
}

// JsonByte : 对象转json格式的二进制数组
func (resp *RespMsg) JsonByte() []byte {
	r, err := json.Marshal(resp)
	if err != nil {
		return nil
	}
	return r
}

// JSONString : 对象转json格式的string
func (resp *RespMsg) JSONString() string {
	r, err := json.Marshal(resp)
	if err != nil {
		log.Println(err)
	}
	return string(r)
}
