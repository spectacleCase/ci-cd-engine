package initialize

import (
	"github.com/songzhibin97/gkit/cache/local_cache"
	"github.com/spectacleCase/ci-cd-engine/config"
	"github.com/spectacleCase/ci-cd-engine/global"
	"github.com/spectacleCase/ci-cd-engine/utils"
)

// OtherInit 其他初始化
func OtherInit() {
	dr, err := utils.ParseDuration(config.Config.JWT.ExpiresTime)
	if err != nil {
		panic(err)
	}

	global.BlackCache = local_cache.NewCache(
		local_cache.SetDefaultExpire(dr),
	)
}
