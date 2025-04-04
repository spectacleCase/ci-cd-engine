package global

import (
	"github.com/docker/docker/client"
	system "github.com/spectacleCase/ci-cd-engine/models/system"
	"go.uber.org/zap"
)

var (
	DockerCli    *client.Client
	CLog         *zap.Logger
	CTaskManager *system.TaskManager
)
