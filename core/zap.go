package core

import (
	"fmt"
	"github.com/spectacleCase/ci-cd-engine/config"
	"github.com/spectacleCase/ci-cd-engine/core/internal"
	"github.com/spectacleCase/ci-cd-engine/utils"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

// Zap 获取 zap.Logger
// Author [SliverHorn](https://github.com/SliverHorn)
func Zap() (logger *zap.Logger) {
	if ok, _ := utils.PathExists(config.Config.Zap.Director); !ok { // 判断是否有Director文件夹
		fmt.Printf("create %v directory\n", config.Config.Zap.Director)
		_ = os.Mkdir(config.Config.Zap.Director, os.ModePerm)
	}
	levels := config.Config.Zap.Levels()
	length := len(levels)
	cores := make([]zapcore.Core, 0, length)
	for i := 0; i < length; i++ {
		core := internal.NewZapCore(levels[i])
		cores = append(cores, core)
	}
	logger = zap.New(zapcore.NewTee(cores...))
	if config.Config.Zap.ShowLine {
		logger = logger.WithOptions(zap.AddCaller())
	}
	return logger
}
