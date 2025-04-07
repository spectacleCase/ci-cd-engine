package system

import (
	"github.com/spectacleCase/ci-cd-engine/models"
	"sync"
	"time"
)

// Task 表示一个任务
type Task struct {
	models.BaseMODEL
	Name    string `json:"name"`
	Payload []byte `json:"payload"`
	Status  string `json:"status"`
}

// TaskManager 任务管理器
type TaskManager struct {
	Tasks        map[string]*Task
	Queue        chan *Task
	Mu           sync.RWMutex
	RepoPath     []string
	PollInterval time.Duration
}
