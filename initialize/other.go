package initialize

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/songzhibin97/gkit/cache/local_cache"
	"github.com/spectacleCase/ci-cd-engine/config"
	"github.com/spectacleCase/ci-cd-engine/global"
	"github.com/spectacleCase/ci-cd-engine/utils"
	"net/http"
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

	// 启动 prometheus metrics 服务
	go func() {
		http.Handle("/metrics", promhttp.Handler())
		if err := http.ListenAndServe(":2112", nil); err != nil {
			fmt.Println("Prometheus metrics 启动失败:", err)
		}
	}()
}
