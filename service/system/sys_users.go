package system

import (
	"context"
	"errors"
	"github.com/spectacleCase/ci-cd-engine/global"
	"github.com/spectacleCase/ci-cd-engine/models/dao"
	"github.com/spectacleCase/ci-cd-engine/models/system"
	"github.com/spectacleCase/ci-cd-engine/models/system/request"
	"go.uber.org/zap"
	"sync"
	"time"
)

var UserSrvIns *UserSrv
var UserSrvOnce sync.Once

type UserSrv struct {
}

func GetUserSrv() *UserSrv {
	UserSrvOnce.Do(func() {
		UserSrvIns = &UserSrv{}
	})
	return UserSrvIns
}

// Sign 注册
func (usr *UserSrv) Sign(c context.Context, user request.Users) (err error) {
	userDao := dao.NewUserDao(c)

	// 检查邮箱是否已注册
	_, exist, err := userDao.ExistOrNotByEmail(user.Email)
	if err != nil {
		err = errors.New("系统错误")
		return
	}
	if exist {
		err = errors.New("该邮箱已被注册")
		return
	}

	_, exist, err = userDao.ExistOrNotByUserName(user.Username)
	if err != nil {
		err = errors.New("系统错误")
		return
	}
	if exist {
		err = errors.New("用户名已经存在了")
		return
	}

	// 注册新用户
	newUser := &system.Users{
		Username:    user.Username,
		Email:       user.Email,
		IsActive:    false,
		IsAdmin:     false,
		LastLoginAt: time.Now(),
	}
	err = newUser.SetPasswordHash(user.Password)
	if err != nil {
		return err
	}

	err = userDao.CreateUser(newUser)
	if err != nil {
		global.CLog.Error("注册失败", zap.Any("err", err))
		err = errors.New("注册失败")
		return
	}
	return
}
