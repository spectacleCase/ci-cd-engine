package system

import (
	"sync"
	"time"
)

// Task 表示一个任务
type Task struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Payload   []byte    `json:"payload"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// TaskManager 任务管理器
type TaskManager struct {
	Tasks        map[string]*Task
	Queue        chan *Task
	Mu           sync.RWMutex
	RepoPath     string
	PollInterval time.Duration
}

// NewTaskManager 创建新的任务管理器
func NewTaskManager(repoPath string, pollInterval time.Duration, queueSize int) *TaskManager {
	return &TaskManager{
		Tasks:        make(map[string]*Task),
		Queue:        make(chan *Task, queueSize),
		RepoPath:     repoPath,
		PollInterval: pollInterval,
	}
}
