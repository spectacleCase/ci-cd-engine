package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/spectacleCase/ci-cd-engine/utils"

	"github.com/spectacleCase/ci-cd-engine/global"
	"github.com/spectacleCase/ci-cd-engine/models/common/response"
	"github.com/spectacleCase/ci-cd-engine/models/system/request"
	"github.com/spectacleCase/ci-cd-engine/service/system"
	"go.uber.org/zap"
)

// Sign 用户注册
func Sign() gin.HandlerFunc {
	return func(c *gin.Context) {
		var user request.Users
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
func Login(c *gin.Context) {

}
