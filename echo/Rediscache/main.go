package Rediscache

import (
	"fmt"
	"strconv"

	"../Tool"
	"github.com/garyburd/redigo/redis"
)

func GetConn() redis.Conn {
	configData := Tool.GetConfig("redisConfig")
	host, _ := configData["host"]
	c, err := redis.Dial("tcp", host)
	if err != nil {
		fmt.Println("Connect to redis error", err)
		panic(err)
	}
	password, _ := configData["password"]
	if password != "" {
		c.Do("AUTH", password)
	}
	db, _ := configData["db"]
	if db != "" {
		dbNum, _ := strconv.Atoi(db)
		c.Do("SELECT", dbNum)
	}
	return c
}

//判断key存不存在
func KeyExists(name string) bool {
	redisCoonect := GetConn()
	defer redisCoonect.Close()
	value, err := redis.Int64(redisCoonect.Do("EXISTS", name))
	if err != nil {
		return false
	}
	if value != 1 {
		return false
	}
	return true
}

//得到key的值
func GetValue(name string) string {
	redisCoonect := GetConn()
	defer redisCoonect.Close()
	value, err := redis.String(redisCoonect.Do("GET", name))
	if err != nil {
		return ""
	}
	return value
}

//设置key-value 有过期时间
func SetValueTime(name, value string, ex int) bool {
	redisCoonect := GetConn()
	defer redisCoonect.Close()
	_, err := redisCoonect.Do("SET", name, value, "EX", ex)
	if err != nil {
		fmt.Println("redis set failed:", err)
		return false
	} else {
		return true
	}
}

/**
 *删除key
 */
func DeleteValue(name string) bool {
	redisCoonect := GetConn()
	defer redisCoonect.Close()
	_, err := redisCoonect.Do("DEL", name)
	if err != nil {
		fmt.Println("redis set failed:", err)
		return false
	} else {
		return true
	}
}

//设置key-value
func SetValue(name, value string) bool {
	redisCoonect := GetConn()
	defer redisCoonect.Close()
	_, err := redisCoonect.Do("SET", name, value)
	if err != nil {
		fmt.Println("redis set failed:", err)
		return false
	} else {
		return true
	}

}

//得到一个自增的数字
func Incr(name string) int64 {
	redisCoonect := GetConn()
	defer redisCoonect.Close()
	Num, _ := redis.Int64(redisCoonect.Do("INCR", name))
	return Num
}
