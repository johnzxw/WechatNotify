package service

import (
	"fmt"
	"strconv"
	"time"

	"../Tool"
	"github.com/labstack/echo"
)

/**
 *验证码是否正确
 */
func ValidateCaptcha(captchaValue string, c echo.Context) bool {
	captchaConfig := Tool.GetConfig("captchaConfig")
	captchaName := captchaConfig["captchaSessionName"]
	data := GetSessionByName(captchaName, c)
	if len(data) == 0 {
		return false
	}
	durTimeConfigTmp := captchaConfig["captchaDurtime"]
	durTimeConfig, _ := strconv.ParseInt(durTimeConfigTmp, 10, 64)
	currentTime := time.Now().Unix()
	//校验验证码时间是否过期
	if _, ok := data["time"]; ok == false {
		return false
	}
	durTime := fmt.Sprintf("%s", data["time"])
	durTimeInt64, _ := strconv.ParseInt(durTime, 10, 64)
	if durTimeInt64+durTimeConfig <= currentTime {
		return false
	}

	if captcha, ok := data["captcha"]; ok {
		captchaData := fmt.Sprintf("%s", captcha)
		if captchaData == captchaValue {
			return true
		}
	}
	return false
}
