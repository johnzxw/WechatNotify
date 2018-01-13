package service

import (
	"crypto/md5"
	"fmt"

	"../Define"
	"../Rediscache"
	"../Tool"
	"../mongoCache"
	"github.com/labstack/echo"
	"gopkg.in/mgo.v2/bson"
)

var UserLoginSessionName = "userlogin"

/**
 *用户登录
 */
func AdminUserLogin(name, password string, c echo.Context) bool {
	Validata, data := ValidateAdminUserNamePass(name, password)

	if Validata == true {
		mapData := Tool.Struct2Map(data)
		mapData["Id"] = data.Id.Hex()
		if SetSessionByName(UserLoginSessionName, mapData, MaxAge, c) == true {
			SetAdminUserLastLoginTime(data.Id.Hex())
			SetAdminUserLoginLog(mapData["Id"], c)
			SetAdminUserSessionId(mapData["Id"], c)
			return true
		} else {
			return false
		}
	}
	return false
}

/**
 *验证用户名 密码是否正确
 */
func ValidateAdminUserNamePass(name, password string) (bool, *Define.AdminUserDefine) {
	data := GetAdminDataByName(name)
	if data.Name == "" {
		return false, data
	}
	pass := GetAdminPass(password, data.Salt)
	if pass != data.Pass {
		return false, data
	}
	return true, data
}

/**
 *验证用户是否登录
 */
func ValidateAdminIsLogin(c echo.Context) bool {
	data := GetUserDataBySession(c)
	if len(data) == 0 {
		return false
	}
	return true
}

/**
 *添加管理员
 */
func AdminUserAdd(data Define.AdminUserDefine) string {
	return mongoCache.SaveAdminUserData(data)
}

/**
 * 获取用户信息 map[string]string
 */
func GetUserDataBySession(c echo.Context) map[interface{}]interface{} {
	return GetSessionByName(UserLoginSessionName, c)
}

/**
 *返回用户加密之后的密码
 */
func GetAdminPass(password, salt string) string {
	pass1 := md5.Sum([]byte(password))
	pass2 := fmt.Sprintf("%x", pass1)
	pass3 := md5.Sum([]byte(pass2 + salt))
	return fmt.Sprintf("%x", pass3)
}

/**
 *
 */
func GetAdminDataByName(name string) *Define.AdminUserDefine {
	return mongoCache.GetAdminUserDataByUsername(name)
}

func GetAdminDataById(id string) *Define.AdminUserDefine {
	return mongoCache.GetAdminUserDataByObjectId(id)
}

/**
 * 获取用户sessionid 保存到数据库 用户单点登录
 */
func SetAdminUserSessionId(name string, c echo.Context) bool {
	sessionid := GetSessionId(UserLoginSessionName, c)
	return Rediscache.SetValueTime(name, sessionid, MaxAge)
}

/**
 *通过cookie获取sessionid
 */
func GetAdminUserSessionId(c echo.Context) string {
	return GetSessionId(UserLoginSessionName, c)
}

/**
 *更新最后一次登录时间
 */
func SetAdminUserLastLoginTime(Id string) bool {
	where := bson.M{"_id": bson.ObjectIdHex(Id)}
	change := bson.M{"$set": bson.M{"LastLoginTime": Tool.GetDateString()}}
	return mongoCache.UpdateAdminUserDataForWhere(where, change)
}

/**
 *记录登录日志
 */
func SetAdminUserLoginLog(name string, c echo.Context) bool {
	data := &Define.AdminUserLoginLogDefine{}
	data.Id = bson.NewObjectId()
	data.Name = name
	data.LoginIp = c.RealIP()
	data.LoginTime = Tool.GetDateString()
	data.LoginTimeStamp = Tool.GetTimeStamp()
	return mongoCache.SetAdminUserLoginLog(*data)
}

func AddWechatUserGroup(name, ForName string) bool {

	wechatUserGroup := &Define.WechatUserGroupStruct{}
	wechatUserGroup.Name = name
	wechatUserGroup.ForName = ForName

	return mongoCache.SaveWechatUserGroup(*wechatUserGroup)
}

//获取所有的person数据
func GetAllWechatGroupData() map[int]map[string]string {
	tmpData := mongoCache.GetAllWechatGroupData()
	result := map[int]map[string]string{}
	i := 1
	tmpAdminName := map[string]string{}
	for _, data := range tmpData {
		result[i] = map[string]string{}
		result[i]["Id"] = data.Id.Hex()
		result[i]["Name"] = data.Name
		result[i]["AddTime"] = Tool.GetDateStringByStamp(data.AddTime)
		if AdminName, ok := tmpAdminName[data.ForName]; ok {
			result[i]["ForName"] = AdminName
		} else {
			tmp := GetAdminDataById(data.ForName)
			result[i]["ForName"] = tmp.Name
			tmpAdminName[data.ForName] = tmp.Name
		}
		i++
	}
	return result
}

/**
 *根据groupid 删除改组下所有的人
 */
func DeleteAdminUserGroupBelong(groupId string) bool {
	return mongoCache.DeleteAdminUserGroupBelong(groupId)
}

func GetAllAdminUserGroupBelong(groupId string) map[int]map[string]string {
	tmpData := mongoCache.GetAllAdminUserGroupBelong(groupId)
	result := map[int]map[string]string{}
	i := 1
	for _, data := range tmpData {
		result[i] = map[string]string{}
		result[i]["Uid"] = data.Uid
		i++
	}
	return result
}

/**
 *往组里面添加人
 */
func AddAdminUserGroupBelong(groupId string, uid map[string]string) bool {
	DeleteAdminUserGroupBelong(groupId)
	AddTime := Tool.GetTimeStamp()
	for _, u := range uid {
		WechatUserGroupBelongStruct := &Define.WechatUserGroupBelongStruct{}
		WechatUserGroupBelongStruct.Id = bson.NewObjectId()
		WechatUserGroupBelongStruct.Uid = u
		WechatUserGroupBelongStruct.GroupId = groupId
		WechatUserGroupBelongStruct.AddTime = AddTime
		mongoCache.SaveWechatUserGroupBelong(*WechatUserGroupBelongStruct)
	}
	return true
}
