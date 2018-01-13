package Controller

import (
	"../Rediscache"
	"../service"
	"fmt"
	"github.com/labstack/echo"
	"net/http"
)

var UserLoginRedirectUrl = "/main"
var UserNotLoginRedirectUrl = "/login"
/**
 *登录
 */
func AdminLogin(c echo.Context) error {
	//logined redirect
	if service.ValidateAdminIsLogin(c) {
		c.Redirect(http.StatusMovedPermanently, UserLoginRedirectUrl)
	}

	account := c.FormValue("account")
	password := c.FormValue("password")
	code := c.FormValue("code")
	if len(account) > 0 && len(password) > 0 && len(code) > 0 {
		returnJson := map[string]string{}
		returnJson["status"] = "500"
		returnJson["message"] = "验证码错误"
		//validate captcha
		if service.ValidateCaptcha(code, c) == false {
			return c.JSON(http.StatusOK, returnJson)
		}

		//Validate username password
		if service.AdminUserLogin(account, password, c) == false {
			returnJson["message"] = "用户名密码错误"
			return c.JSON(http.StatusOK, returnJson)
		}
		returnJson["status"] = "200"
		returnJson["message"] = "success"
		returnJson["url"] = UserLoginRedirectUrl
		return c.JSON(http.StatusOK, returnJson)
	}
	data := map[string]string{}
	return c.Render(http.StatusOK, "adminLogin.html", data)
}

/**
 * 主界面
 */
func AdminMain(c echo.Context) error {
	if service.ValidateAdminIsLogin(c) == false {
		c.Redirect(http.StatusMovedPermanently, UserNotLoginRedirectUrl)
	}
	UserData := service.GetUserDataBySession(c)
	if _, ok := UserData["Id"]; ok {
		//判断session是否多处登录 ，之前的登录session删除
		redisSessionId := Rediscache.GetValue(fmt.Sprintf("%s", UserData["Id"]))
		sessionId := service.GetAdminUserSessionId(c)
		//如果当前的是sessionid和redis存储的不一样， 删除当前的session
		if redisSessionId != sessionId && redisSessionId != "" {
			Rediscache.DeleteValue(sessionId)
			c.Redirect(http.StatusOK, UserNotLoginRedirectUrl)
		}
	}
	return c.Render(http.StatusOK, "adminMain.html", UserData)
}

/**
 *管理用户组
 */
func WechatUserGroup(c echo.Context) error {
	if service.ValidateAdminIsLogin(c) == false {
		c.Redirect(http.StatusMovedPermanently, UserNotLoginRedirectUrl)
	}
	groupName := c.FormValue("groupName")
	if len(groupName) > 0 {
		returnJson := map[string]string{}
		returnJson["status"] = "500"
		CurrentUser := service.GetUserDataBySession(c)
		CurrentUserId := fmt.Sprintf("%s", CurrentUser["Id"])
		if service.AddWechatUserGroup(groupName, CurrentUserId) == false {
			returnJson["message"] = "用户组添加失败，请检查名称"
			return c.JSON(http.StatusOK, returnJson)
		}
		returnJson["status"] = "200"
		returnJson["message"] = "success"
		return c.JSON(http.StatusOK, returnJson)
	}
	data := map[string]interface{}{}
	data["data"] = service.GetAllWechatGroupData()
	return c.Render(http.StatusOK, "wechatUserGroup.html", data)
}

/**
 *退出界面
 */
func AdminLogout(c echo.Context) error {
	if service.DeleteSessionByName(service.UserLoginSessionName, c) {
		c.Redirect(http.StatusMovedPermanently, UserNotLoginRedirectUrl)
	} else {
		c.Redirect(http.StatusMovedPermanently, UserLoginRedirectUrl)
	}
	return c.String(http.StatusOK, "")
}

func UserGroupEdit(c echo.Context) error {
	if service.ValidateAdminIsLogin(c) == false {
		c.Redirect(http.StatusMovedPermanently, UserNotLoginRedirectUrl)
	}
	groupId := c.QueryParam("id")

	uidArray := c.FormValue("uidArray")
	if len(uidArray) > 0 {
		dataResult := map[string]string{
			"status":  "200",
			"message": "success",
		}
		uidArray := map[string]string{
			"kdf":  "sdfef",
			"sdfe": "sdfe",
		}
		service.AddAdminUserGroupBelong(groupId, uidArray)
		return c.JSON(http.StatusOK, dataResult)
	}
	//查询出所有的用户
	AllUserListTmp := service.GetAllWechatGroupData()
	var AllUserList map[string]string
	for _, value := range AllUserListTmp {
		AllUserList[value["Id"]] = value["Name"]
	}
	//查询出所有的groupid下的用户
	AllGroupUserListTmp := service.GetAllAdminUserGroupBelong(groupId)
	var AllGroupUserList map[string]string
	for _, value := range AllGroupUserListTmp {
		AllGroupUserList[value["Uid"]] = value["Uid"]
	}
	var data map[string]map[string]string
	data["AllUserList"] = AllUserList
	data["AllGroupUserList"] = AllGroupUserList
	return c.Render(http.StatusOK, "wechatuserGroupEdit.html", data)
}
