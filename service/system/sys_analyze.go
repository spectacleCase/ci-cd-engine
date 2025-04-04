package system

import (
	"fmt"
	"github.com/spectacleCase/ci-cd-engine/global"
	system "github.com/spectacleCase/ci-cd-engine/models/system"
	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
	"os"
	"os/exec"
)

// Analyze 解析文件
func Analyze(filename string) (*system.CiCdConfig, error) {
	file, err := os.Open(filename)
	if err != nil {
		global.CLog.Error("failed to open file:", zap.Error(err))
		return nil, err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			global.CLog.Error("failed to close file", zap.Error(err))
		}
	}(file)

	var config system.CiCdConfig
	if err := yaml.NewDecoder(file).Decode(&config); err != nil {
		global.CLog.Error("failed to open file:", zap.Error(err))
		return nil, err
	}
	return &config, nil
}

func AnalyzeToMap(conf system.CiCdConfig) (map[string]system.Stage, error) {
	stageMap := make(map[string]system.Stage)

	for _, stage := range conf.Stages {
		stageMap[stage.Name] = stage
	}
	return stageMap, nil
}

// 执行命令
func runCommand(command string) {
	cmd := exec.Command("sh", "-c", command)
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("命令执行失败: %v\n输出: %s", err, output)
	}
	fmt.Printf("命令输出: %s\n", output)
}

// 根据阶段执行任务
func executeStages(stages []system.Stage) {
	for _, stage := range stages {
		fmt.Printf("准备执行阶段: %s\n", stage.Name)
		global.CLog.Info("准备执行阶段",
			zap.String("stage", stage.Name))
		for _, cmd := range stage.Commands {
			runCommand(cmd)
		}
	}
}
