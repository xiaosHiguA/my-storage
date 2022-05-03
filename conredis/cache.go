package conredis

import (
	"fmt"
	"github.com/garyburd/redigo/redis"
	"github.com/spf13/viper"
	"log"
	"os"
	"time"
)

type poolRedis struct {
	pool          *redis.Pool
	hostRedisAddr string
	pass          string
}

var pool *redis.Pool

func getConfig() *poolRedis {
	var pool poolRedis
	path, _ := os.Getwd()

	viper.SetConfigName("cache_config")
	viper.SetConfigType("yml")
	viper.AddConfigPath(path + "/config/")

	err := viper.ReadInConfig()
	if err != nil {
		log.Println("获取配置信息错误: ", err)
		return nil
	}

	pool.hostRedisAddr = viper.GetString("caches.redisHost")
	pool.pass = viper.GetString("caches.redisPass")
	return &pool
}

func newRedisConfig() *redis.Pool {
	pl := getConfig()
	return &redis.Pool{
		MaxIdle:     50, //最大连接数
		MaxActive:   30,
		IdleTimeout: 300 * time.Second,
		Dial: func() (redis.Conn, error) {
			// 1 连接
			conn, err := redis.Dial("tcp", pl.hostRedisAddr)
			if err != nil {
				fmt.Println("conn redis ", err)
				return nil, err
			}
			// 2. 访问认证
			if _, err = conn.Do("AUTH", pl.pass); err != nil {
				fmt.Println("AUTH: redis 密码认证失败...", err)
				conn.Close()
				return nil, err
			}
			return conn, nil
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error { //检测redis连接健康状态
			if time.Since(t) < time.Minute {
				return nil
			}
			_, err := c.Do("PING")
			return err
		},
	}
}

func InitRedis() {
	pool = newRedisConfig()
}

func GetRedisPool() *redis.Pool {
	return pool
}
