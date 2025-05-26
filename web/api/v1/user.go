package v1

import (
	"github.com/gin-gonic/gin"
	commonReq "github.com/spectacleCase/ci-cd-engine/models/common/request"
	systemRes "github.com/spectacleCase/ci-cd-engine/models/system/response"
	"github.com/spectacleCase/ci-cd-engine/utils"
	"time"

	"github.com/spectacleCase/ci-cd-engine/global"
	"github.com/spectacleCase/ci-cd-engine/models/common/response"
	"github.com/spectacleCase/ci-cd-engine/models/system/request"
	"github.com/spectacleCase/ci-cd-engine/service/system"
	"go.uber.org/zap"
)

// Sign 用户注册
func Sign() gin.HandlerFunc {
	return func(c *gin.Context) {
		var user request.Sign
		if err := c.ShouldBindJSON(&user); err != nil || !utils.EmailVerify(user.Email) || len(user.Username) < 6 {
			global.CLog.Error("参数有误", zap.Any("err", err))
			response.FailWithMessage("参数有误", c)
			return
		}
		userSer := system.GetUserSrv()
		err := userSer.Sign(c.Request.Context(), user)
		if err != nil {
			response.FailWithMessage(err.Error(), c)
			return
		}
		response.OkWithMessage("注册成功", c)
	}
}

// Login 用户登录
func Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		var user request.Login
		if err := c.ShouldBindJSON(&user); err != nil || !utils.EmailVerify(user.Email) {
			global.CLog.Error("参数有误", zap.Any("err", err))
			response.FailWithMessage("参数有误", c)
			return
		}
		// 验证码校验
		if user.CaptchaId != "" && user.Captcha != "" && store.Verify(user.CaptchaId, user.Captcha, false) {
			userSer := system.GetUserSrv()
			claims, loginResponse, err := userSer.Login(c.Request.Context(), user)
			if err != nil {
				response.FailWithMessage(err.Error(), c)
				return
			}
			system.SetToken(c, loginResponse.Token, int(claims.RegisteredClaims.ExpiresAt.Unix()-time.Now().Unix()))
			response.OkWithDetailed(systemRes.LoginResponse{
				User:      loginResponse.User,
				Token:     loginResponse.Token,
				ExpiresAt: claims.RegisteredClaims.ExpiresAt.Unix() * 1000,
			}, "登录成功", c)
			return

		}
		response.FailWithMessage("验证码错误", c)
		return

	}
}

// GetUser 获取用户
func GetUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		var pageInfo commonReq.PageInfo
		pageInfo = commonReq.NewPageInfo(c)

		userSer := system.GetUserSrv()
		username := c.DefaultQuery("username", "")
		email := c.DefaultQuery("email", "")
		if !utils.EmailVerify(email) {
			response.FailWithMessage("邮箱格式有问题", c)
			return
		}
		user, err := userSer.GetUser(c.Request.Context(), pageInfo, username, email)
		if err != nil {
			response.FailWithMessage(err.Error(), c)
			return
		}
		response.OkWithData(user, c)
		return

	}
}

// Update 修改用户
func Update() gin.HandlerFunc {
	return func(c *gin.Context) {
		var update request.Update
		if err := c.ShouldBindJSON(&update); err != nil {
			global.CLog.Error("参数有误", zap.Any("err", err))
			response.FailWithMessage("参数有误", c)
			return
		}
		userSer := system.GetUserSrv()
		user, err := userSer.Update(c.Request.Context(), update)
		if err != nil {
			response.FailWithMessage(err.Error(), c)
			return
		}
		response.OkWithData(user, c)
	}
}
