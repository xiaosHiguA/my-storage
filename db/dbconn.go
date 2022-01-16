package db

import (
	"MyStorage/config"
	"encoding/json"
	"fmt"
	"github.com/jinzhu/gorm"
	"log"
	"os"
)

var dbSql *gorm.DB

func init() {
	confData := getSqlConnConfig()
	connString := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=true&loc=%s",
		confData.Username,
		confData.Password,
		confData.Host,
		confData.Port,
		confData.Database,
		confData,
		confData.Loc)
	db, err := gorm.Open(confData.DriverName, connString)
	if err != nil {
		log.Println("连接mysql err: ", err)
		return
	}
	dbSql = db
}

// 获取sql配置文件信息
func getSqlConnConfig() config.SqlConnDate {
	var conf = config.SqlConnDate{}
	file, err := os.Open("./sql.json")
	if err != nil {
		log.Println("打开json文件出错: ", err)
		return conf
	}
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(conf); err != nil {
		log.Println("解析json err: ", err)
		return conf
	}
	return conf
}

func GetDb() *gorm.DB {
	return dbSql
}
