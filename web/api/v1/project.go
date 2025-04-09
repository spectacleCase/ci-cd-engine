package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/spectacleCase/ci-cd-engine/global"
	"github.com/spectacleCase/ci-cd-engine/models/common/response"
	"github.com/spectacleCase/ci-cd-engine/models/system/request"
	"github.com/spectacleCase/ci-cd-engine/service/system"
	"go.uber.org/zap"
)

// Project 创建/修改项目
func Project() gin.HandlerFunc {
	return func(c *gin.Context) {
		var newProject request.CreateProject
		if err := c.ShouldBind(&newProject); err != nil {
			global.CLog.Error("参数有误", zap.Any("err", err))
			response.FailWithMessage("参数有误", c)
			return
		}
		psv := system.GetProjectSrv()
		err := psv.Project(newProject, c)
		if err != nil {
			response.FailWithMessage(err.Error(), c)
			return
		}
		response.Ok(c)
	}
}

// GetProject 获取项目
func GetProject() gin.HandlerFunc {
	return func(c *gin.Context) {
		psv := system.GetProjectSrv()
		projectList, err := psv.GetProject(c)
		if err != nil {
			response.FailWithMessage(err.Error(), c)
			return
		}
		response.OkWithData(projectList, c)
	}
}

// DeleteProject 删除项目
func DeleteProject() gin.HandlerFunc {
	return func(c *gin.Context) {
		var delProject request.DeleteProject
		if err := c.ShouldBind(&delProject); err != nil {
			global.CLog.Error("参数有误", zap.Any("err", err))
			response.FailWithMessage("参数有误", c)
			return
		}
		psv := system.GetProjectSrv()
		err := psv.DeleteProject(delProject, c)
		if err != nil {
			response.FailWithMessage(err.Error(), c)
			return
		}
		response.Ok(c)

	}
}
