package response

import (
	"github.com/gin-gonic/gin"
	"github.com/spectacleCase/ci-cd-engine/pkg"
	"net/http"
)

type Response struct {
	Code    int         `json:"code"`
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
}

func Result(code int, data interface{}, msg string, c *gin.Context) {
	c.JSON(http.StatusOK, Response{
		code,
		data,
		msg,
	})
}

func Ok(c *gin.Context) {
	Result(pkg.SUCCESS, map[string]interface{}{}, "操作成功", c)
}

func OkWithMessage(message string, c *gin.Context) {
	Result(pkg.SUCCESS, map[string]interface{}{}, message, c)
}

func OkWithData(data interface{}, c *gin.Context) {
	Result(pkg.SUCCESS, data, "操作成功", c)
}

func OkWithDetailed(data interface{}, msg string, c *gin.Context) {
	Result(pkg.SUCCESS, data, msg, c)
}

func Fail(c *gin.Context) {
	Result(pkg.ERROR, map[string]interface{}{}, "操作失败", c)
}

func FailWithMessage(message string, c *gin.Context) {
	Result(pkg.ERROR, map[string]interface{}{}, message, c)
}
func NoAuth(message string, c *gin.Context) {
	c.JSON(http.StatusUnauthorized, Response{
		pkg.ERROR,
		nil,
		message,
	})
}

func FailWithDetailed(data interface{}, message string, c *gin.Context) {
	Result(pkg.ERROR, data, message, c)
}
