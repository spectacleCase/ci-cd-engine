package system

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/spectacleCase/ci-cd-engine/common"
	"github.com/spectacleCase/ci-cd-engine/global"
	system "github.com/spectacleCase/ci-cd-engine/models/system"
	"go.uber.org/zap"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

// NewTaskManager 创建新的任务管理器
func NewTaskManager(repoPath string, pollInterval time.Duration, queueSize int) *system.TaskManager {
	if global.CTaskManager == nil {
		return &system.TaskManager{
			Tasks:        make(map[string]*system.Task),
			Queue:        make(chan *system.Task, queueSize),
			RepoPath:     repoPath,
			PollInterval: pollInterval,
		}
	}
	return global.CTaskManager
}

// Start 启动任务管理器
func Start(ctx context.Context) {
	// 启动仓库轮询
	go func() {
		err := pollRepositoryChanges(ctx)
		if err != nil {

		}
	}()

	// 启动任务消费者
	go consumeTasks(ctx)
}

// AddTask 添加新任务
func AddTask(task *system.Task) error {
	global.CTaskManager.Mu.Lock()

	defer global.CTaskManager.Mu.Unlock()

	if _, exists := global.CTaskManager.Tasks[task.ID]; exists {
		return errors.New("task already exists")

	}

	task.Status = common.StatusQueued
	task.CreatedAt = time.Now()
	task.UpdatedAt = time.Now()

	global.CTaskManager.Tasks[task.ID] = task
	global.CTaskManager.Queue <- task
	return nil
}

// GetTask 获取任务状态
func GetTask(id string) (*system.Task, error) {
	global.CTaskManager.Mu.RLock()
	defer global.CTaskManager.Mu.RUnlock()

	task, exists := global.CTaskManager.Tasks[id]
	if !exists {
		return nil, errors.New("task not found")
	}

	return task, nil
}

// ListTasks 列出所有任务
func ListTasks() []*system.Task {
	global.CTaskManager.Mu.RLock()
	defer global.CTaskManager.Mu.RUnlock()

	tasks := make([]*system.Task, 0, len(global.CTaskManager.Tasks))
	for _, task := range global.CTaskManager.Tasks {
		tasks = append(tasks, task)
	}

	return tasks
}

// pollRepositoryChanges 轮询仓库变化
func pollRepositoryChanges(ctx context.Context) error {
	ticker := time.NewTicker(global.CTaskManager.PollInterval)
	defer ticker.Stop()

	lastHash := ""

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			// 检查仓库变化
			hash, err := getGitRepoHash(global.CTaskManager.RepoPath)
			if err != nil {
				continue
			}

			if lastHash == "" {
				lastHash = hash
				continue
			}

			if hash != lastHash {
				lastHash = hash

				// 创建仓库变化任务
				task := &system.Task{
					ID:   fmt.Sprintf("repo-%d", time.Now().UnixNano()),
					Name: "repository_change",
					// todo 添加任务
					//Payload: fmt.Sprintf(`{"commit_hash": "%s"}`, hash),
				}

				if err := AddTask(task); err != nil {
					return err
				}
			}
		}
	}
}

// consumeTasks 消费任务
func consumeTasks(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case task := <-global.CTaskManager.Queue:
			// 处理任务
			processTask(task)
		}
	}
}

// processTask 处理单个任务
func processTask(task *system.Task) {
	global.CLog.Info("process task", zap.String("id", task.ID))
	global.CTaskManager.Mu.Lock()
	task.Status = common.StatusRunning
	task.UpdatedAt = time.Now()
	global.CTaskManager.Mu.Unlock()

	// 模拟任务处理
	time.Sleep(2 * time.Second)
	var newConfig system.CiCdConfig
	if err := json.Unmarshal(task.Payload, &newConfig); err != nil {
		global.CLog.Error("JSON反序列化失败", zap.String("payload", string(task.Payload)))
	}
	stageMap, _ := AnalyzeToMap(newConfig)
	AssemblyLineProject(stageMap["Build"], stageMap["Deploy"])
	global.CTaskManager.Mu.Lock()
	task.Status = common.StatusCompleted
	task.UpdatedAt = time.Now()
	global.CTaskManager.Mu.Unlock()
	global.CLog.Info("执行成功")
}

// getGitRepoHash 获取Git仓库当前hash
func getGitRepoHash(repoPath string) (string, error) {
	cmd := exec.Command("git", "rev-parse", "HEAD")
	cmd.Dir = repoPath

	output, err := cmd.Output()
	if err != nil {
		return "", err
	}

	return string(output[:40]), nil
}

// InitTaskManager 初始化任务管理器
func InitTaskManager() (context.Context, context.CancelFunc) {
	// 获取仓库路径，默认为当前目录
	repoPath := os.Getenv("REPO_PATH")
	if repoPath == "" {
		repoPath, _ = os.Getwd()
	}

	// 解析为绝对路径
	absRepoPath, err := filepath.Abs(repoPath)
	if err != nil {
		global.CLog.Error("Error getting absolute path", zap.Error(err))
	}

	// 创建任务管理器
	tm := NewTaskManager(absRepoPath, 30*time.Second, 100)
	global.CTaskManager = tm

	// 创建上下文
	ctx, cancel := context.WithCancel(context.Background())
	return ctx, cancel

}
