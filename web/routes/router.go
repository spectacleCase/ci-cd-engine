package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/spectacleCase/ci-cd-engine/core"
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
		v1.GET("/test/dockerTest", api.DockerTest())
		v1.GET("/test/ResponseTest", api.ResponseTest())

		v1.POST("/user/sign", api.Sign())
		v1.POST("/user/login", api.Login())
		v1.POST("/common/captcha", api.Captcha())
	}

	authed := v1.Group("/")
	authed.Use(core.Auth())

	// 需要鉴权
	{
		authed.POST("/project/project", api.Project())
		authed.GET("/project/project", api.GetProject())
		authed.PUT("/project/project", api.Project())
		authed.DELETE("/project/project", api.DeleteProject())
		authed.POST("/project/project/addCrew", api.AddCrew())

		authed.GET("/test/tokenTest", api.TokenTest())
		authed.GET("/user/user", api.GetUser())
		authed.PUT("/user/user", api.Update())
	}
	return r
}
