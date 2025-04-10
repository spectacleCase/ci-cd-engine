package system

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/spectacleCase/ci-cd-engine/global"
	commonReq "github.com/spectacleCase/ci-cd-engine/models/common/request"
	"github.com/spectacleCase/ci-cd-engine/models/dao"
	"github.com/spectacleCase/ci-cd-engine/models/system"
	"github.com/spectacleCase/ci-cd-engine/models/system/request"
	"go.uber.org/zap"
	"strconv"
	"sync"
)

var ProjectSrvIns *ProjectSrv
var ProjectSrvOnce sync.Once

type ProjectSrv struct{}

func GetProjectSrv() *ProjectSrv {
	ProjectSrvOnce.Do(func() {
		ProjectSrvIns = &ProjectSrv{}
	})
	return ProjectSrvIns
}

// Project 创建/修改项目
func (psv *ProjectSrv) Project(project request.CreateProject, ctx *gin.Context) (err error) {
	proDao := dao.NewProjectDao(ctx)
	if project.Id == "" {
		// 新增
		pro := &system.ProjectPermissions{
			UserId:          GetUserID(ctx),
			ProjectId:       project.ProjectId,
			ProjectName:     project.ProjectName,
			PermissionLevel: "admin",
		}

		err = proDao.CreateProject(pro)
		if err != nil {
			err = errors.New("创建项目失败")
			global.CLog.Error("创建项目失败", zap.Any("err", err))
			return err
		}
	} else {
		// 修改
		pro, exist, err := proDao.ExistOrNotById(project.Id)
		if err != nil {
			err = errors.New("系统错误")
			return err
		}
		if exist {
			pro.ProjectId = project.ProjectId
			pro.ProjectName = project.ProjectName
			err = proDao.UpdateProject(pro)
			if err != nil {
				err = errors.New("修改项目失败")
				global.CLog.Error("修改项目失败", zap.Any("err", err))
				return err
			}
			return nil
		}
		err = errors.New("没有该项目")
		global.CLog.Error("没有该项目", zap.Any("err", err))
		return err
	}
	return err

}

// DeleteProject 删除项目
func (psv *ProjectSrv) DeleteProject(project request.DeleteProject, ctx *gin.Context) (err error) {
	proDao := dao.NewProjectDao(ctx)
	err = proDao.DeleteById(project.Id)
	if err != nil {
		err = errors.New("删除项目失败")
		global.CLog.Error("删除项目失败", zap.Any("err", err))
		return err
	}
	return nil
}

// GetProject 返回项目权限列表
func (psv *ProjectSrv) GetProject(pageInfo commonReq.PageInfo, ctx *gin.Context) (commonReq.PageInfo, []system.ProjectPermissions, error) {
	proDao := dao.NewProjectDao(ctx)
	id := GetUserID(ctx)
	list, info, err := proDao.GetByIdProjectList(strconv.Itoa(int(id)), pageInfo)
	if err != nil {
		global.CLog.Error("获取项目列表失败", zap.Any("err", err))
		return commonReq.PageInfo{}, nil, errors.New("获取项目列表失败")
	}
	return info, list, nil
}

// AddGrew 添加组员
func (psv *ProjectSrv) AddGrew(ctx *gin.Context, crew request.AddCrew) (err error) {
	proDao := dao.NewProjectDao(ctx)
	err = proDao.AddCrew(crew)
	if err != nil {
		global.CLog.Error(err.Error())
		return errors.New(err.Error())
	}
	return nil
}
