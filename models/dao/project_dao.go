package dao

import (
	"context"
	"errors"
	"github.com/spectacleCase/ci-cd-engine/global"
	commonReq "github.com/spectacleCase/ci-cd-engine/models/common/request"
	"github.com/spectacleCase/ci-cd-engine/models/system"
	"github.com/spectacleCase/ci-cd-engine/models/system/request"
	"gorm.io/gorm"
	"strconv"
)

type ProjectDao struct {
	*gorm.DB
}

func NewProjectDao(ctx context.Context) *ProjectDao {
	return &ProjectDao{global.NewDBClient(ctx)}
}

// CreateProject 创建项目
func (dao *ProjectDao) CreateProject(pro *system.ProjectPermissions) error {
	return dao.DB.Model(&system.ProjectPermissions{}).Create(&pro).Error
}

func (dao *ProjectDao) UpdateProject(pro *system.ProjectPermissions) error {
	return dao.DB.Model(pro).Updates(pro).Error
}

// ExistOrNotByProjectId 根据项目标识查询是否存在
func (dao *ProjectDao) ExistOrNotByProjectId(email string) (user *system.ProjectPermissions, exist bool, err error) {
	var count int64
	err = dao.DB.Model(&system.ProjectPermissions{}).Where("project_id = ?", email).Count(&count).Error
	if count == 0 {
		return user, false, nil
	}
	err = dao.DB.Model(&system.ProjectPermissions{}).Where("project_id = ?", email).First(&user).Error
	if err != nil {
		return user, false, err
	}
	return user, true, nil
}

// ExistOrNotByProjectName 根据项目名称查询是否存在
func (dao *ProjectDao) ExistOrNotByProjectName(name string) (user *system.ProjectPermissions, exist bool, err error) {
	var count int64
	err = dao.DB.Model(&system.ProjectPermissions{}).Where("project_name = ?", name).Count(&count).Error
	if count == 0 {
		return user, false, nil
	}
	err = dao.DB.Model(&system.ProjectPermissions{}).Where("project_name = ?", name).First(&user).Error
	if err != nil {
		return user, false, err
	}
	return user, true, nil
}

// ExistOrNotById 根据项目路径查询是否存在
func (dao *ProjectDao) ExistOrNotById(id string) (user *system.ProjectPermissions, exist bool, err error) {
	var count int64
	err = dao.DB.Model(&system.ProjectPermissions{}).Where("id = ?", id).Count(&count).Error
	if count == 0 {
		return user, false, nil
	}
	err = dao.DB.Model(&system.ProjectPermissions{}).Where("id = ?", id).First(&user).Error
	if err != nil {
		return user, false, err
	}
	return user, true, nil
}

// DeleteById 根据id删除
func (dao *ProjectDao) DeleteById(id string) error {
	return dao.DB.Where("id = ?", id).Delete(&system.ProjectPermissions{}).Error
}

// GetByIdProjectList  获取项目列表
func (dao *ProjectDao) GetByIdProjectList(userId string, pageInfo commonReq.PageInfo) ([]system.ProjectPermissions, commonReq.PageInfo, error) {
	var list []system.ProjectPermissions
	var total int64

	// 查询总数
	if err := dao.DB.Model(&system.ProjectPermissions{}).Where("user_id = ?", userId).Count(&total).Error; err != nil {
		return nil, commonReq.PageInfo{}, err
	}

	// 查询分页数据
	err := dao.DB.Where("user_id = ?", userId).
		Offset(pageInfo.Offset()).
		Limit(pageInfo.PageSize).
		Find(&list).Error

	if err != nil {
		return nil, commonReq.PageInfo{}, err
	}

	return list, pageInfo, nil
}

// AddCrew 添加组员
func (dao *ProjectDao) AddCrew(crew request.AddCrew) (err error) {
	var count int64
	var pro system.ProjectPermissions
	tx := dao.DB.Model(&system.ProjectPermissions{}).Where("project_id = ?", crew.ProjectId)
	err = tx.Count(&count).Error
	if err != nil {
		err = errors.New("系统错误")
		return err
	}
	if count > 0 {
		err = tx.First(&pro).Error
		if err != nil {
			return err
		}
	} else {
		err = errors.New("项目不存在")
		return err
	}

	count = 0
	err = dao.DB.Model(&system.ProjectPermissions{}).Where("project_id = ? and user_id = ?", crew.ProjectId, crew.UserId).Count(&count).Error
	if count != 0 {
		err = errors.New("用户已经存在")
		return err
	}
	uid, err := strconv.ParseUint(crew.UserId, 10, 64)
	if err != nil {
		err = errors.New("系统错误")
		return err
	}

	return dao.DB.Model(&system.ProjectPermissions{}).Create(&system.ProjectPermissions{
		UserId:          uint(uid),
		ProjectId:       crew.ProjectId,
		PermissionLevel: crew.PermissionLevel,
		ProjectName:     pro.ProjectName,
	}).Error
}
