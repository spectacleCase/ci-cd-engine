package dao

import (
	"context"
	"github.com/spectacleCase/ci-cd-engine/global"
	"github.com/spectacleCase/ci-cd-engine/models/system"
	"gorm.io/gorm"
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

func (dao *ProjectDao) DeleteById(id string) error {
	return dao.DB.Where("id = ?", id).Delete(&system.ProjectPermissions{}).Error
}

func (dao *ProjectDao) GetByIdProjectList(userId string) ([]system.ProjectPermissions, error) {
	var list []system.ProjectPermissions
	err := dao.DB.Where("user_id = ?", userId).Find(&list).Error
	return list, err
}
