package system

import (
	"context"
	"errors"
	"github.com/spectacleCase/ci-cd-engine/global"
	commonReq "github.com/spectacleCase/ci-cd-engine/models/common/request"
	"github.com/spectacleCase/ci-cd-engine/models/dao"
	"github.com/spectacleCase/ci-cd-engine/models/system"
	"github.com/spectacleCase/ci-cd-engine/models/system/request"
	systemRes "github.com/spectacleCase/ci-cd-engine/models/system/response"
	"github.com/spectacleCase/ci-cd-engine/utils"
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
func (usr *UserSrv) Sign(c context.Context, user request.Sign) (err error) {
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

// Login 登录
func (usr *UserSrv) Login(c context.Context, loginUser request.Login) (commonReq.CustomClaims, systemRes.LoginResponse, error) {
	userDao := dao.NewUserDao(c)
	// 是否存在该用户
	user, exist, err := userDao.ExistOrNotByEmail(loginUser.Email)
	if err != nil {
		err = errors.New("系统错误")
		return commonReq.CustomClaims{}, systemRes.LoginResponse{}, err
	}
	if !exist {
		err = errors.New("用户名不存在")
		return commonReq.CustomClaims{}, systemRes.LoginResponse{}, err
	}

	// 密码校验
	if utils.BcryptCheck(user.PasswordHash, loginUser.Password) {
		err = errors.New("密码错误")
		return commonReq.CustomClaims{}, systemRes.LoginResponse{}, err
	}

	// 签发jwt
	return usr.TokenNext(user)
}

// TokenNext 登录成功之后签发jwt
func (usr *UserSrv) TokenNext(user *system.Users) (commonReq.CustomClaims, systemRes.LoginResponse, error) {
	j := NewJWT()
	claims := j.CreateClaims(commonReq.BaseClaims{
		ID: user.ID,
	})
	token, err := j.CreateToken(claims)
	if err != nil {
		global.CLog.Error("获取token失败!", zap.Error(err))
		return commonReq.CustomClaims{}, systemRes.LoginResponse{}, err
	}

	logR := systemRes.LoginResponse{
		User:      *user,
		Token:     token,
		ExpiresAt: claims.RegisteredClaims.ExpiresAt.Unix() * 1000,
	}
	return claims, logR, err
}

// GetUser 获取用户
func (usr *UserSrv) GetUser(c context.Context, pageInfo commonReq.PageInfo, username, email string) (user *systemRes.User, err error) {
	userDao := dao.NewUserDao(c)
	users, exist, err := userDao.ExistOrNotByEmailAndName(username, email)
	if err != nil {
		err = errors.New("系统错误")
		return
	}
	if exist {
		responseUser := &systemRes.User{
			Id:       users.ID,
			Username: users.Username,
			Email:    users.Email,
			IsActive: users.IsActive,
			IsAdmin:  users.IsAdmin,
		}
		return responseUser, nil
	}

	return nil, nil

}

// Update 修改用户
func (usr *UserSrv) Update(c context.Context, update request.Update) (user *systemRes.User, err error) {
	//userDao := dao.NewUserDao(c)
	// 修改权限 -> admin
	if update.IsAdmin || update.IsActive {
		//userDao := dao.NewUserDao(c)
	}

	//users, exist, err := userDao.ExistOrNotByEmailAndName(username, email)

	return &systemRes.User{}, nil
}
