package core

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/spectacleCase/ci-cd-engine/models/common/response"
	"github.com/spectacleCase/ci-cd-engine/service/system"
)

// Auth 身份校验
func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		newJwt := system.NewJWT()
		token := newJwt.GetToken(c)
		if token == "" {
			response.NoAuth("未登录或者非法访问", c)
			c.Abort()
			return
		}
		j := system.NewJWT()
		claims, err := j.ParseToken(token)
		if err != nil {
			if errors.Is(err, system.TokenExpired) {
				response.NoAuth("授权已过期", c)
				system.ClearToken(c)
				c.Abort()
				return
			}
			response.NoAuth(err.Error(), c)
			system.ClearToken(c)
			c.Abort()
			return
		}

		c.Set("claims", claims)
		// todo 后续做刷新token
		c.Next()
	}
}
