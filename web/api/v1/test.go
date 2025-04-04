package v1

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/spectacleCase/ci-cd-engine/common"
	"github.com/spectacleCase/ci-cd-engine/initialize"
	moSystem "github.com/spectacleCase/ci-cd-engine/models/system"
	system "github.com/spectacleCase/ci-cd-engine/service/system"
)

func DockerTest() gin.HandlerFunc {
	return func(c *gin.Context) {
		ciCdConfig, err := system.Analyze("file/ci-yaml/.cicd.yaml")
		if err != nil {
			fmt.Println(err)
		} else {
			jsonString, _ := json.Marshal(ciCdConfig)
			task := &moSystem.Task{
				Name:    "后面添加git分支信息",
				Payload: jsonString,
				Status:  common.StatusPending,
			}
			db := initialize.NewDBClient(c)
			err = db.Create(task).Error
			if err != nil {
				return
			}
			_ = system.AddTask(task)

		}

		c.JSON(200, gin.H{"message": "测试成功"})
	}
}
