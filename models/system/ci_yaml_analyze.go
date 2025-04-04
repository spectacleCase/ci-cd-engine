package system

import (
	"fmt"
	"github.com/spectacleCase/ci-cd-engine/global"
	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
	"os"
	"os/exec"
)

// Stage 定义阶段配置
type Stage struct {
	Name             string   `yaml:"name"`
	Image            string   `yaml:"image"`
	WorkingDirectory string   `yaml:"working_directory"`
	Commands         []string `yaml:"commands"`
	Volumes          []string `yaml:"volumes"`
	Ports            []string `yaml:"ports"`
}

// CiCdConfig 定义 CI/CD 配置结构体
type CiCdConfig struct {
	Stages []Stage `yaml:"stages"`
}

// Analyze 解析文件
func Analyze(filename string) (map[string]Stage, error) {
	file, err := os.Open(filename)
	if err != nil {
		global.C_LOG.Error("failed to open file:", zap.Error(err))
		return nil, err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			global.C_LOG.Error("failed to close file", zap.Error(err))
		}
	}(file)

	var config CiCdConfig
	if err := yaml.NewDecoder(file).Decode(&config); err != nil {
		global.C_LOG.Error("failed to open file:", zap.Error(err))
		return nil, err
	}
	stageMap := make(map[string]Stage)

	for _, stage := range config.Stages {
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
func executeStages(stages []Stage) {
	for _, stage := range stages {
		fmt.Printf("准备执行阶段: %s\n", stage.Name)
		global.C_LOG.Info("准备执行阶段",
			zap.String("stage", stage.Name))
		for _, cmd := range stage.Commands {
			runCommand(cmd)
		}
	}
}
