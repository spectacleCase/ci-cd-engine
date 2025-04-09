package global

import (
	"context"
	"github.com/docker/docker/client"
	"github.com/songzhibin97/gkit/cache/local_cache"
	system "github.com/spectacleCase/ci-cd-engine/models/system"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

var (
	CDB          *gorm.DB
	DockerCli    *client.Client
	CLog         *zap.Logger
	CTaskManager *system.TaskManager
	BlackCache   local_cache.Cache
)

func NewDBClient(ctx context.Context) *gorm.DB {
	db := CDB
	return db.WithContext(ctx)
}
