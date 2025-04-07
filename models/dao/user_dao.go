package dao

import (
	"context"
	"github.com/spectacleCase/ci-cd-engine/global"
	"github.com/spectacleCase/ci-cd-engine/models/system"
	"gorm.io/gorm"
)

type UserDao struct {
	*gorm.DB
}

func NewUserDao(ctx context.Context) *UserDao {
	return &UserDao{global.NewDBClient(ctx)}
}

// ExistOrNotByEmail 根据邮箱查询是否存在
func (dao *UserDao) ExistOrNotByEmail(email string) (user *system.Users, exist bool, err error) {
	var count int64
	err = dao.DB.Model(&system.Users{}).Where("email = ?", email).Count(&count).Error
	if count == 0 {
		return user, false, nil
	}
	err = dao.DB.Model(&system.Users{}).Where("email = ?", email).First(&user).Error
	if err != nil {
		return user, false, err
	}
	return user, true, nil
}

// ExistOrNotByUserName 根据username判断是否存在该名字
func (dao *UserDao) ExistOrNotByUserName(userName string) (user *system.Users, exist bool, err error) {
	var count int64
	err = dao.DB.Model(&system.Users{}).Where("username = ?", userName).Count(&count).Error
	if count == 0 {
		return user, false, err
	}
	err = dao.DB.Model(&system.Users{}).Where("username = ?", userName).First(&user).Error
	if err != nil {
		return user, false, err
	}
	return user, true, nil
}

// CreateUser 创建用户
func (dao *UserDao) CreateUser(user *system.Users) error {
	return dao.DB.Model(&system.Users{}).Create(&user).Error
}
