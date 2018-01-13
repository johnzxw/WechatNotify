package service

import (
	"fmt"
	"strconv"

	"../Tool"
	"github.com/labstack/echo"
	"gopkg.in/boj/redistore.v1"
)

var (
	host             string
	password         string
	db               string
	size             int
	secretKey        string
	domain           string
	MaxAge           int
	sessionName      string
	sessionKeyPreFix string
)

func init() {
	configData := Tool.GetConfig("redisConfig")
	host, _ = configData["host"]
	password, _ = configData["password"]
	db, _ = configData["db"]
	configDataSession := Tool.GetConfig("sessionConfig")
	sizeTmp, _ := configDataSession["size"]
	size, _ = strconv.Atoi(sizeTmp)
	secretKey = configDataSession["secret-key"]
	domain = configDataSession["Domain"]
	maxAgeTmp := configDataSession["maxAge"]
	MaxAge, _ = strconv.Atoi(maxAgeTmp)
	sessionName = configDataSession["sessionName"]
	sessionKeyPreFix = configDataSession["sessionKeyPreFix"]
}

/**
 * redis for mongo
 *
 */
func GetSessionByName(name string, c echo.Context) map[interface{}]interface{} {
	store, err := redistore.NewRediStoreWithDB(size, "tcp", host, password, db, []byte(secretKey))
	if err != nil {
		panic(err.Error())
	}
	defer store.Close()

	session, err := store.Get(c.Request(), name)
	if err != nil {
		fmt.Println(err.Error())
	}
	return session.Values
}

/**
 *获取sessionid
 */
func GetSessionId(name string, c echo.Context) string {
	store, err := redistore.NewRediStoreWithDB(size, "tcp", host, password, db, []byte(secretKey))
	if err != nil {
		panic(err.Error())
	}
	defer store.Close()

	session, err := store.Get(c.Request(), name)
	if err != nil {
		fmt.Println(err.Error())
	}
	return sessionKeyPreFix + session.ID
}

/**
 *
 */
func DeleteSessionByName(name string, c echo.Context) bool {
	store, err := redistore.NewRediStoreWithDB(size, "tcp", host, password, db, []byte(secretKey))
	if err != nil {
		panic(err.Error())
	}
	defer store.Close()

	session, err := store.Get(c.Request(), name)
	if err != nil {
		fmt.Println(session)
	}
	session.Options.MaxAge = -1
	if err = session.Save(c.Request(), c.Response()); err != nil {
		fmt.Println(err)
		return false
	}
	return true
}

func SetSessionByName(name string, data map[string]string, durTime int, c echo.Context) bool {
	store, err := redistore.NewRediStoreWithDB(size, "tcp", host, password, db, []byte(secretKey))
	if err != nil {
		panic(err)
	}
	defer store.Close()
	if len(name) == 0 {
		name = sessionName
	}
	store.SetMaxAge(durTime)
	store.SetKeyPrefix(sessionKeyPreFix)

	store.Options.MaxAge = durTime
	store.Options.Domain = domain
	//js not found
	store.Options.HttpOnly = true
	//https send
	store.Options.Secure = true

	session, err := store.Get(c.Request(), name)
	if err != nil {
		fmt.Println(session)
	}
	for k, v := range data {
		session.Values[k] = v
	}

	session.Options.Domain = domain
	session.Options.MaxAge = durTime
	//js not found
	session.Options.Secure = true
	//https send
	session.Options.HttpOnly = true
	if err = session.Save(c.Request(), c.Response()); err != nil {
		fmt.Println(err)
		return false
	}
	return true
}
