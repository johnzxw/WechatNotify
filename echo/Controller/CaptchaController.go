package Controller

import (
	"strconv"

	"../Tool"
	"../service"
	"github.com/labstack/echo"
)

func Imgcaptcha(c echo.Context) error {
	//图形验证码宽度和高度
	var captchaWidth int
	var captchaHeight int
	captchaConfig := Tool.GetConfig("captchaConfig")
	captchaWidth, _ = strconv.Atoi(captchaConfig["captchaWidth"])
	captchaHeight, _ = strconv.Atoi(captchaConfig["captchaHeight"])
	d := make([]byte, 4)
	s := Tool.NewLen(4)
	ss := ""
	d = []byte(s)
	for v := range d {
		d[v] %= 10
		ss += strconv.FormatInt(int64(d[v]), 32)
	}
	//验证码内容 放到session中
	captchaName := captchaConfig["captchaSessionName"]
	data := map[string]string{
		"captcha": ss,
		"time":    Tool.GetTimeString(),
	}
	durTimeTmp := captchaConfig["captchaDurtime"]
	durTime, _ := strconv.Atoi(durTimeTmp)
	service.SetSessionByName(captchaName, data, durTime, c)
	c.Response().Header().Set("Content-Type", "image/png")
	Tool.NewImage(d, captchaWidth, captchaHeight).WriteTo(c.Response())
	var err error
	return err
}
