package v1

import "C"
import (
	"github.com/gin-gonic/gin"
	"github.com/mojocn/base64Captcha"
	"github.com/spectacleCase/ci-cd-engine/config"
	"github.com/spectacleCase/ci-cd-engine/global"
	"github.com/spectacleCase/ci-cd-engine/models/common/response"
	systemRes "github.com/spectacleCase/ci-cd-engine/models/system/response"
	"go.uber.org/zap"
)

var store = base64Captcha.DefaultMemStore

// Captcha  生成验证码
func Captcha() gin.HandlerFunc {
	return func(c *gin.Context) {
		//openCaptchaTimeOut := 3600
		//key := c.ClientIP()
		//v, ok := global.BlackCache.Get(key)
		//if !ok {
		//	global.BlackCache.Set(key, 1, time.Second*time.Duration(openCaptchaTimeOut))
		//}

		var oc bool
		// 字符,公式,验证码配置
		// 生成默认数字的driver
		driver := base64Captcha.NewDriverDigit(config.Config.Captcha.ImgHeight, config.Config.Captcha.ImgWidth, config.Config.Captcha.KeyLong, 0.7, 80)
		// cp := base64Captcha.NewCaptcha(driver, store.UseWithCtx(c))   // v8下使用redis
		cp := base64Captcha.NewCaptcha(driver, store)
		id, b64s, _, err := cp.Generate()
		if err != nil {
			global.CLog.Error("验证码获取失败!", zap.Error(err))
			response.FailWithMessage("验证码获取失败", c)
			return
		}
		response.OkWithDetailed(systemRes.SysCaptchaResponse{
			CaptchaId:     id,
			PicPath:       b64s,
			CaptchaLength: config.Config.Captcha.KeyLong,
			OpenCaptcha:   oc,
		}, "验证码获取成功", c)
	}
}
