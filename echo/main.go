package main

import (
	"html/template"
	"io"
	"os"
	"strconv"
	"time"

	"./Controller"
	"./Tool"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/labstack/gommon/log"
)

var CurrentPath = ""

//监听端口
var ListenPort = "1323"

//host
var ListenHosts = ""

//日志存放路径
var logFileName = ""

//请求日志
var requestLogFileName = ""

type Template struct {
	templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}
func init() {

	configRuntimeLog := Tool.GetConfig("logFileName")
	logFileName = configRuntimeLog["name"]
	requestLogFileName = configRuntimeLog["requestLogFileName"]
	CurrentPath = Tool.GetCurrentDirectory()

	manConfig := Tool.GetConfig("manConfig")
	ListenPort = manConfig["port"]
	ListenHosts = manConfig["host"]

}
func main() {
	t := &Template{
		templates: template.Must(template.ParseGlob(CurrentPath + "/views/*.html")),
	}
	e := echo.New()

	systemConfig := Tool.GetConfig("systemConfig")

	//开启压缩
	if systemConfig["enablegzip"] != "0" {
		gzipLevel := systemConfig["gzipLevel"]
		gzipLevelInt, _ := strconv.Atoi(gzipLevel)
		e.Use(middleware.GzipWithConfig(middleware.GzipConfig{
			Level: gzipLevelInt,
		}))
	}
	//请求日志开启
	if systemConfig["enablelog"] != "0" {
		//把请求日志放到文件中
		RequestlogFile, err := os.OpenFile(CurrentPath+requestLogFileName, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0755)
		defer RequestlogFile.Close()
		if err != nil {
			e.Logger.Fatal("打开请求日志输出文件失败！ 请检查文件权限。" + requestLogFileName + "  error !!")
		}
		e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
			Format: `{"time_unix":"${time_unix}","time_unix_nano":"${time_unix_nano}","time_rfc3339":"${time_rfc3339}","time_rfc3339_nano":"${time_rfc3339_nano}","id":"${id}","remote_ip":"${remote_ip}","uri":"${uri}","host":"${host}","method":"${method}","path":"${path}","referer":"${referer}","user_agent":"${user_agent}","status":"${status}","latency":"${latency}","latency_human":"${latency_human}","bytes_in":"${bytes_in}","bytes_out":"${bytes_out}","header":"${header}","query":"${query}","form":"${form}"}` + "\n",
			Output: RequestlogFile,
		}))
	}

	//隐藏echo框架log的logo
	e.HideBanner = true
	e.Debug = false

	//把错误日志放到文件中
	logFile, err := os.OpenFile(CurrentPath+logFileName, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0755)
	defer logFile.Close()
	if err != nil {
		e.Logger.Fatal("打开日志输出文件失败！ 请检查文件权限。" + logFileName + "  error !!")
	}

	//日志输出到文件
	e.Logger.SetLevel(log.INFO)
	e.Logger.SetOutput(logFile)
	e.Renderer = t
	e.HTTPErrorHandler = Controller.EchoErrorHandle
	e.DisableHTTP2 = true
	e.Server.ReadTimeout = 10 * time.Second
	e.Server.WriteTimeout = 10 * time.Second

	//静态文件
	e.File("/favicon.ico", CurrentPath+"/static/images/favicon.ico")
	e.Use(middleware.StaticWithConfig(middleware.StaticConfig{
		Root: CurrentPath + "/static/",
	}))

	//route
	e.GET("/api/common/getCaptcha", Controller.Imgcaptcha)
	e.GET("/", Controller.WechatAuth)
	e.POST("/", Controller.WechatAuth)
	e.GET("/userGroup", Controller.WechatUserGroup)
	e.GET("/userGroupDelete", Controller.WechatUserGroup)
	e.POST("/userGroup", Controller.WechatUserGroup)
	e.GET("/login", Controller.AdminLogin)
	e.POST("/loginPost", Controller.AdminLogin)
	e.GET("/main", Controller.AdminMain)
	e.GET("/getUser", Controller.GetWechatUser)
	e.GET("/sendMessage", Controller.SendMessage)
	e.GET("/userGroupEdit",Controller.UserGroupEdit)
	Addr := ListenHosts + ":" + ListenPort
	e.Logger.Fatal(e.Start(Addr))
}
