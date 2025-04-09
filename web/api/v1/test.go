package v1

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/spectacleCase/ci-cd-engine/global"
	"github.com/spectacleCase/ci-cd-engine/models/common/response"
	moSystem "github.com/spectacleCase/ci-cd-engine/models/system"
	"github.com/spectacleCase/ci-cd-engine/pkg"
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
				Status:  pkg.StatusPending,
			}
			err = global.CDB.Create(task).Error
			if err != nil {
				return
			}
			_ = system.AddTask(c, task)

		}

		response.Ok(c)
	}
}

func ResponseTest() gin.HandlerFunc {
	return func(c *gin.Context) {
		response.Ok(c)
	}
}

func TokenTest() gin.HandlerFunc {
	return func(c *gin.Context) {
		response.OkWithMessage("token校验成功", c)
	}
}
