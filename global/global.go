package global

import (
	"github.com/docker/docker/client"
	"go.uber.org/zap"
)

var (
	DockerCli *client.Client
	C_LOG     *zap.Logger
)
