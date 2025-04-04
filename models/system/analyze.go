package system

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
