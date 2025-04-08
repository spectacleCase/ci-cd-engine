package routes

import (
	"github.com/gin-gonic/gin"
	api "github.com/spectacleCase/ci-cd-engine/web/api/v1"
)

func NewRouter() *gin.Engine {
	r := gin.Default()

	v1 := r.Group("api/v1")

	{
		v1.GET("ping", func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "测试成功"})
		})
	}

	{
		// docker 测试
		v1.GET("/dockerTest", api.DockerTest())
		v1.GET("/ResponseTest", api.ResponseTest())

		v1.POST("/user/sign", api.Sign())
		v1.POST("/user/login", api.Login())
		v1.POST("/common/captcha", api.Captcha())
	}

	//authed := v1.Group("/")

	// 需要鉴权
	//v1.POST("/user/", api.DockerTest())
	return r
}
