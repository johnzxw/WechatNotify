package Controller

import (
	"strconv"
	"fmt"
	"net/http"
	"crypto/md5"
	"github.com/labstack/echo"
	"../Tool"
	"encoding/json"
	"github.com/garyburd/redigo/redis"
	"github.com/bitly/go-simplejson"
	"sort"
	"strings"
	"crypto/sha1"
	"gopkg.in/mgo.v2"
)

var (
	wxAppId           string
	wxAppSecret       string
	oauth2RedirectURI string
	oauth2Scope       string
	imgUrl            string
	token             string
	tableName         string
)

var (
	RedisHost string
	RedisPass string
	RedisDb   int
)

type GetAccessStruct struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
}

type GetWechatUserStruct struct {
	Total      int            `json:"total"`
	Count      int            `json:"count"`
	Data       WechatUserData `json:"data"`
	NextOpenid string         `json:"next_openid"`
}
type WechatUserData struct {
	OpenidArray []WechatUserOpenid `json:"openid"`
}

type WechatUserOpenid struct {
	Openid string
}

type WechatSendResultStruct struct {
	Errcode int    `json:"errcode"`
	Errmsg  string `json:"errmsg"`
}

type WechatUserInfoStruct struct {
	Subscribe     int    `json:"subscribe"`
	Openid        string `json:"openid"`
	Nickname      string `json:"nickname"`
	Sex           int    `json:"sex"`
	Language      string `json:"language"`
	City          string `json:"city"`
	Province      string `json:"province"`
	Country       string `json:"country"`
	Headimgurl    string `json:"headimgurl"`
	SubscribeTime int64  `json:"subscribe_time"`
	//Unionid       string `json:"unionid"`
	Remark        string `json:"remark"`
	Groupid       int    `json:"groupid"`
	//TagidList     string `json:"tagid_list"`
}

func init() {
	wechatConfig := Tool.GetConfig("wechat")
	wxAppId = wechatConfig["wxAppId"]
	wxAppSecret = wechatConfig["wxAppSecret"]
	oauth2RedirectURI = wechatConfig["oauth2RedirectURI"]
	oauth2Scope = wechatConfig["oauth2Scope"]
	imgUrl = wechatConfig["imgUrl"]
	token = wechatConfig["token"]
	tableName = wechatConfig["tableName"]

	RedisTmp := Tool.GetConfig("redisConfig")
	RedisHost = RedisTmp["host"]
	RedisPass = RedisTmp["password"]
	RedisDbTmp, _ := strconv.Atoi(RedisTmp["db"])
	RedisDb = RedisDbTmp
}

//授权
func WechatAuth(c echo.Context) error {
	echoStr := c.QueryParam("echostr")
	//接口校验
	if echoStr != "" {
		if validateWechatSign(c) {
			return c.HTML(http.StatusOK, c.QueryParam("echostr"))
		} else {
			return c.HTML(http.StatusOK, "fail")
		}
	}
	//关注校验
	openid := c.QueryParam("openid")
	if openid != "" {
		if validateWechatSign(c) {
			saveWechatForOpenid(openid, c)
			xmlData := "<xml><ToUserName>< ![CDATA[toUser] ]></ToUserName><FromUserName>< ![CDATA[FromUser] ]></FromUserName><CreateTime>123456789</CreateTime><MsgType>< ![CDATA[event] ]></MsgType><Event>< ![CDATA[subscribe] ]></Event></xml>"
			return c.XML(http.StatusOK, xmlData)
		}
	}

	return c.Redirect(http.StatusFound, imgUrl)
}

func SendMessage(c echo.Context) error {
	Url := "https://api.weixin.qq.com/cgi-bin/message/template/send?access_token=" + getAccessToken(false)
	AllWechatUser := GetWechatUserOpenid()
	for _, data := range AllWechatUser {
		PostField := map[string]string{
			"touser":      data.Openid,
			"template_id": "OOlaOfGYVlTb2RTeFYMaNrccYfmfvtrP7WN93bk3ljU",
			"data": "{'ip':{'value':" + c.RealIP() + "','color':'#173177'}," +
				"'datetime':{'value':'" + c.RealIP() + "','color':'#173177'}," +
				"'where':{'value':'" + c.RealIP() + "','color':'#173177'}}",
		}
		Tool.Post(Url, PostField)
		//if SendResult == "" {
		//	return c.JSON(http.StatusOK, map[string]string{
		//		"status":  "500",
		//		"message": "send message error",
		//	})
		//}
		//resultArray := &WechatSendResultStruct{}
		//errs := json.Unmarshal([]byte(SendResult), resultArray)
		//if errs != nil || resultArray.Errmsg != "ok" {
		//	return c.JSON(http.StatusOK, map[string]string{
		//		"status":  "500",
		//		"message": resultArray.Errmsg,
		//	})
		//}
	}
	return c.JSON(http.StatusOK, map[string]string{
		"status":  "200",
		"message": "success",
	})
}

/**
 *拉取所有关注的用户
 */
func GetWechatUser(c echo.Context) error {
	url := "https://api.weixin.qq.com/cgi-bin/user/get?access_token=" + getAccessToken(false) + "&next_openid="
	wechatUser := Tool.Get(url)
	Arr, err := simplejson.NewJson([]byte(wechatUser))
	if err != nil {
		fmt.Println("json格式化失败")
	}
	total, _ := Arr.Get("total").Int64()
	count, _ := Arr.Get("count").Int64()
	Number := (total / count) + 1

	arr, _ := Arr.Get("data").Get("openid").Array()
	saveWechatUser(arr, c)
	//多次拉取
	if Number > 1 {
		var i int64
		for i = 2; i <= Number; i++ {
			nextOpenid, _ := Arr.Get("next_openid").String()
			urlTmp := url + nextOpenid
			wechatUser := Tool.Get(urlTmp)
			Arr, err := simplejson.NewJson([]byte(wechatUser))
			if err != nil {
				fmt.Println("json格式化失败")
			}
			arr, _ := Arr.Get("data").Get("openid").Array()
			saveWechatUser(arr, c)
		}
	}
	return c.HTML(http.StatusOK, "")
}

/**
 *批次保存用户信息
 */
func saveWechatUser(data []interface{}, c echo.Context) bool {
	for _, openidTmp := range data {
		openid := fmt.Sprintf("%s", openidTmp)
		saveWechatForOpenid(openid, c)
	}
	return true
}

/**
 *根据 openid 保存用户信息
 */
func saveWechatForOpenid(openid string, c echo.Context) bool {
	url := "https://api.weixin.qq.com/cgi-bin/user/info?access_token=" + getAccessToken(false) + "&openid=" + openid + "&lang=zh_CN"
	WechatUser := Tool.Get(url)
	c.Logger().Warn("get wechat user:" + openid + " info!" + WechatUser)
	WechatUserStruct := &WechatUserInfoStruct{}
	errs := json.Unmarshal([]byte(WechatUser), WechatUserStruct)
	if errs == nil && WechatUserStruct.SubscribeTime != 0 {
		saveResult := SaveWechatUserInfoData(*WechatUserStruct)
		if saveResult == true {
			return true
		}
	}

	return false
}

/**
 *获取accesstoken
 */
func getAccessToken(isrefer bool) string {
	RedisNameTmp := md5.Sum([]byte(wxAppId + wxAppSecret))
	RedisName := fmt.Sprintf("%x", RedisNameTmp)
	data := GetValue(RedisName)
	//请求接口
	if data == "" || isrefer {
		requestUrl := "https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential&appid=" + wxAppId + "&secret=" + wxAppSecret
		fmt.Println("通过微信接口获取access_token")
		wechatData := Tool.Get(requestUrl)
		resultArray := &GetAccessStruct{}
		errs := json.Unmarshal([]byte(wechatData), resultArray)
		if errs != nil {
			fmt.Println("json 解析失败！ url：")
		}
		//saveRedis
		SetValueTime(RedisName, resultArray.AccessToken, resultArray.ExpiresIn-50)

		return resultArray.AccessToken

	}
	return data
}

/**
 *校验signature
 */
func validateWechatSign(c echo.Context) bool {
	timestamp := c.QueryParam("timestamp")
	nonce := c.QueryParam("nonce")
	tmpSli := []string{
		token,
		timestamp,
		nonce,
	}
	sort.Strings(tmpSli)

	joinString := strings.Join(tmpSli, "")
	tmpSha1 := sha1.New()
	tmpSha1.Write([]byte(joinString))
	sign := fmt.Sprintf("%x", tmpSha1.Sum(nil))
	signWecht := c.QueryParam("signature")
	if sign == signWecht {
		return true
	}

	return false
}

/**
 * redis 链接
 */
func GetConn() redis.Conn {
	c, err := redis.Dial("tcp", RedisHost)
	if err != nil {
		fmt.Println("Connect to redis error", err)
		panic(err)
	}
	if RedisPass != "" {
		c.Do("AUTH", RedisPass)
	}

	if RedisDb != 0 {
		c.Do("SELECT", RedisDb)
	}
	return c
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

var mgoSession *mgo.Session

/**
 * 公共方法，获取session，如果存在则拷贝一份
 */
func GetSession() *mgo.Session {
	mongoUrl := Tool.GetConfig("mongoUrl")
	URL := mongoUrl["name"]
	if mgoSession == nil {
		var err error
		mgoSession, err = mgo.Dial(URL)
		if err != nil {
			panic("mongodb connect faild! " + err.Error()) //直接终止程序运行
		}
	}
	//mgoSession.SetMode(mgo.Monotonic, true)
	//最大连接池默认为4096
	return mgoSession.Clone()
}

//公共方法，获取collection对象
func witchCollection(collection string, s func(*mgo.Collection) error) error {
	session := GetSession()
	defer session.Close()
	dataBaseTmp := Tool.GetConfig("mongoUrl")
	dataBase := dataBaseTmp["dataBase"]
	c := session.DB(dataBase).C(collection)
	return s(c)
}

/**
 *保存用户信息
 */
func SaveWechatUserInfoData(p WechatUserInfoStruct) bool {
	query := func(c *mgo.Collection) error {
		return c.Insert(p)
	}
	err := witchCollection(tableName, query)
	if err == nil {
		return true
	}
	return false
}

/**
 *获取所有的关注者的信息
 */
func GetWechatUserOpenid() []WechatUserInfoStruct {
	var persons []WechatUserInfoStruct
	query := func(c *mgo.Collection) error {
		return c.Find(nil).All(&persons)
	}
	witchCollection(tableName, query)
	return persons
}
