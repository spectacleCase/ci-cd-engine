package v1

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/spectacleCase/ci-cd-engine/models/system"
)

func DockerTest() gin.HandlerFunc {
	return func(c *gin.Context) {
		stageMap, err := system.Analyze("file/ci-yaml/.cicd.yaml")
		if err != nil {
			fmt.Println(err)
		} else {
			system.AssemblyLineProject(stageMap["Build"], stageMap["Deploy"])
		}

		//system.AssemblyLinePythonProject()
		c.JSON(200, gin.H{"message": "测试成功"})
	}
}
