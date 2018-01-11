package Controller

import (
	"../Tool"
	"github.com/labstack/echo"
)

/**
 *错误处理
 */
func EchoErrorHandle(err error, c echo.Context) {
	errorPages := Tool.GetCurrentDirectory() + "/views/error.html"
	errs := c.File(errorPages)
	if errs != nil {
		c.Logger().Error(err)
	}
}
