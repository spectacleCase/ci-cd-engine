package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/spectacleCase/ci-cd-engine/global"
	commonReq "github.com/spectacleCase/ci-cd-engine/models/common/request"
	commonRes "github.com/spectacleCase/ci-cd-engine/models/common/response"
	"github.com/spectacleCase/ci-cd-engine/models/system/request"
	"github.com/spectacleCase/ci-cd-engine/models/system/response"
	"github.com/spectacleCase/ci-cd-engine/service/system"
	"go.uber.org/zap"
)

// Project 创建/修改项目
func Project() gin.HandlerFunc {
	return func(c *gin.Context) {
		var newProject request.CreateProject
		if err := c.ShouldBind(&newProject); err != nil {
			global.CLog.Error("参数有误", zap.Any("err", err))
			commonRes.FailWithMessage("参数有误", c)
			return
		}
		psv := system.GetProjectSrv()
		err := psv.Project(newProject, c)
		if err != nil {
			commonRes.FailWithMessage(err.Error(), c)
			return
		}
		commonRes.Ok(c)
	}
}

// GetProject 获取项目
func GetProject() gin.HandlerFunc {
	return func(c *gin.Context) {
		var pageInfo commonReq.PageInfo
		pageInfo = commonReq.NewPageInfo(c)

		psv := system.GetProjectSrv()
		project, projectList, err := psv.GetProject(pageInfo, c)
		if err != nil {
			commonRes.FailWithMessage(err.Error(), c)
			return
		}
		commonRes.OkWithData(response.ProjectList{
			List:     projectList,
			PageInfo: project,
		}, c)
	}
}

// DeleteProject 删除项目
func DeleteProject() gin.HandlerFunc {
	return func(c *gin.Context) {
		var delProject request.DeleteProject
		if err := c.ShouldBind(&delProject); err != nil {
			global.CLog.Error("参数有误", zap.Any("err", err))
			commonRes.FailWithMessage("参数有误", c)
			return
		}
		psv := system.GetProjectSrv()
		err := psv.DeleteProject(delProject, c)
		if err != nil {
			commonRes.FailWithMessage(err.Error(), c)
			return
		}
		commonRes.Ok(c)

	}
}

// AddCrew 添加组员
func AddCrew() gin.HandlerFunc {
	return func(c *gin.Context) {
		var crew request.AddCrew
		if err := c.ShouldBind(&crew); err != nil {
			global.CLog.Error("参数有误", zap.Any("err", err))
			commonRes.FailWithMessage("参数有误", c)
			return
		}

		psv := system.GetProjectSrv()
		err := psv.AddGrew(c, crew)
		if err != nil {
			commonRes.FailWithMessage(err.Error(), c)
			return
		}
		commonRes.Ok(c)
	}
}
