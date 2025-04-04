package main

import (
	"fmt"
	"github.com/spectacleCase/ci-cd-engine/config"
	"github.com/spectacleCase/ci-cd-engine/core"
	"github.com/spectacleCase/ci-cd-engine/global"
	"github.com/spectacleCase/ci-cd-engine/service/system"
	"github.com/spectacleCase/ci-cd-engine/web/routes"
	"go.uber.org/zap"
	"os"
)

func main() {
	loading()
	router := routes.NewRouter()
	banner, err := os.ReadFile("banner.txt")
	if err != nil {
		panic(err)
	} else {
		fmt.Printf(string(banner))
	}
	global.CLog.Info("ci-cd-engine启动")
	_ = router.Run(config.Config.System.HttpPort)

}

// 初始化配置
func loading() {
	config.InitConfig()
	system.InitDockerCli()
	global.CLog = core.Zap() // 初始化zap日志库
	zap.ReplaceGlobals(global.CLog)
	ctx, _ := system.InitTaskManager()
	system.Start(ctx)
	//defer cancel()

}
