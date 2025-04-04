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

	// docker 测试
	v1.GET("/dockerTest", api.DockerTest())

	return r
}
