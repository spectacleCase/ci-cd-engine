package system

import (
	"github.com/spectacleCase/ci-cd-engine/global"
	commonReq "github.com/spectacleCase/ci-cd-engine/models/common/request"
	"net"
	"time"

	"github.com/gin-gonic/gin"
)

func ClearToken(c *gin.Context) {
	// 增加cookie x-token 向来源的web添加
	host, _, err := net.SplitHostPort(c.Request.Host)
	if err != nil {
		host = c.Request.Host
	}

	if net.ParseIP(host) != nil {
		c.SetCookie("x-token", "", -1, "/", "", false, false)
	} else {
		c.SetCookie("x-token", "", -1, "/", host, false, false)
	}
}

func SetToken(c *gin.Context, token string, maxAge int) {
	// 增加cookie x-token 向来源的web添加
	host, _, err := net.SplitHostPort(c.Request.Host)
	if err != nil {
		host = c.Request.Host
	}

	if net.ParseIP(host) != nil {
		c.SetCookie("x-token", token, maxAge, "/", "", false, false)
	} else {
		c.SetCookie("x-token", token, maxAge, "/", host, false, false)
	}
}

func GetToken(c *gin.Context) string {
	token := c.Request.Header.Get("x-token")
	if token == "" {
		j := NewJWT()
		token, _ = c.Cookie("x-token")
		claims, err := j.ParseToken(token)
		if err != nil {
			global.CLog.Error("重新写入cookie token失败,未能成功解析token,请检查请求头是否存在x-token且claims是否为规定结构")
			return token
		}
		SetToken(c, token, int((claims.ExpiresAt.Unix()-time.Now().Unix())/60))
	}
	return token
}

func GetClaims(c *gin.Context) (*commonReq.CustomClaims, error) {
	token := GetToken(c)
	j := NewJWT()
	claims, err := j.ParseToken(token)
	if err != nil {
		global.CLog.Error("从Gin的Context中获取从jwt解析信息失败, 请检查请求头是否存在x-token且claims是否为规定结构")
	}
	return claims, err
}

// GetUserID 从Gin的Context中获取从jwt解析出来的用户ID
func GetUserID(c *gin.Context) uint {
	if claims, exists := c.Get("claims"); !exists {
		if cl, err := GetClaims(c); err != nil {
			return 0
		} else {
			return cl.BaseClaims.ID
		}
	} else {
		waitUse := claims.(*commonReq.CustomClaims)
		return waitUse.BaseClaims.ID
	}
}

// GetUserUuid 从Gin的Context中获取从jwt解析出来的用户UUID
//func GetUserUuid(c *gin.Context) uuid.UUID {
//	if claims, exists := c.Get("claims"); !exists {
//		if cl, err := GetClaims(c); err != nil {
//			return int{}
//		} else {
//			return cl.ID
//		}
//	} else {
//		waitUse := claims.(*commonReq.CustomClaims)
//		return waitUse.ID
//	}
//}

// GetUserAuthorityId 从Gin的Context中获取从jwt解析出来的用户角色id
//func GetUserAuthorityId(c *gin.Context) uint {
//	if claims, exists := c.Get("claims"); !exists {
//		if cl, err := GetClaims(c); err != nil {
//			return 0
//		} else {
//			return cl.AuthorityId
//		}
//	} else {
//		waitUse := claims.(*commonReq.CustomClaims)
//		return waitUse.AuthorityId
//	}
//}

// GetUserInfo 从Gin的Context中获取从jwt解析出来的用户角色id
func GetUserInfo(c *gin.Context) *commonReq.CustomClaims {
	if claims, exists := c.Get("claims"); !exists {
		if cl, err := GetClaims(c); err != nil {
			return nil
		} else {
			return cl
		}
	} else {
		waitUse := claims.(*commonReq.CustomClaims)
		return waitUse
	}
}

//// GetUserName
