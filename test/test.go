package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
	"time"
)

// Task 表示一个任务
type Task struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Payload   string    `json:"payload"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// TaskManager 任务管理器
type TaskManager struct {
	tasks        map[string]*Task
	queue        chan *Task
	mu           sync.RWMutex
	repoPath     string
	pollInterval time.Duration
}

// NewTaskManager 创建新的任务管理器
func NewTaskManager(repoPath string, pollInterval time.Duration, queueSize int) *TaskManager {
	return &TaskManager{
		tasks:        make(map[string]*Task),
		queue:        make(chan *Task, queueSize),
		repoPath:     repoPath,
		pollInterval: pollInterval,
	}
}

// Start 启动任务管理器
func (tm *TaskManager) Start(ctx context.Context) {
	// 启动仓库轮询
	go tm.pollRepositoryChanges(ctx)

	// 启动任务消费者
	go tm.consumeTasks(ctx)
}

// AddTask 添加新任务
func (tm *TaskManager) AddTask(task *Task) error {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	if _, exists := tm.tasks[task.ID]; exists {
		return fmt.Errorf("task with ID %s already exists", task.ID)
	}

	task.Status = "queued"
	task.CreatedAt = time.Now()
	task.UpdatedAt = time.Now()

	tm.tasks[task.ID] = task
	tm.queue <- task

	return nil
}

// GetTask 获取任务状态
func (tm *TaskManager) GetTask(id string) (*Task, error) {
	tm.mu.RLock()
	defer tm.mu.RUnlock()

	task, exists := tm.tasks[id]
	if !exists {
		return nil, fmt.Errorf("task not found")
	}

	return task, nil
}

// ListTasks 列出所有任务
func (tm *TaskManager) ListTasks() []*Task {
	tm.mu.RLock()
	defer tm.mu.RUnlock()

	tasks := make([]*Task, 0, len(tm.tasks))
	for _, task := range tm.tasks {
		tasks = append(tasks, task)
	}

	return tasks
}

// pollRepositoryChanges 轮询仓库变化
func (tm *TaskManager) pollRepositoryChanges(ctx context.Context) {
	ticker := time.NewTicker(tm.pollInterval)
	defer ticker.Stop()

	lastHash := ""

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			// 检查仓库变化
			hash, err := getGitRepoHash(tm.repoPath)
			if err != nil {
				log.Printf("Error getting repo hash: %v", err)
				continue
			}

			if lastHash == "" {
				lastHash = hash
				continue
			}

			if hash != lastHash {
				log.Printf("Repository changed (old: %s, new: %s)", lastHash[:8], hash[:8])
				lastHash = hash

				// 创建仓库变化任务
				task := &Task{
					ID:      fmt.Sprintf("repo-%d", time.Now().UnixNano()),
					Name:    "repository_change",
					Payload: fmt.Sprintf(`{"commit_hash": "%s"}`, hash),
				}

				if err := tm.AddTask(task); err != nil {
					log.Printf("Error adding repo change task: %v", err)
				}
			}
		}
	}
}

// consumeTasks 消费任务
func (tm *TaskManager) consumeTasks(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case task := <-tm.queue:
			// 处理任务
			tm.processTask(task)
		}
	}
}

// processTask 处理单个任务
func (tm *TaskManager) processTask(task *Task) {
	tm.mu.Lock()
	task.Status = "processing"
	task.UpdatedAt = time.Now()
	tm.mu.Unlock()

	log.Printf("Processing task %s: %s", task.ID, task.Name)

	// 模拟任务处理
	time.Sleep(2 * time.Second)

	tm.mu.Lock()
	task.Status = "completed"
	task.UpdatedAt = time.Now()
	tm.mu.Unlock()

	log.Printf("Completed task %s: %s", task.ID, task.Name)
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

func main() {
	// 获取仓库路径，默认为当前目录
	repoPath := os.Getenv("REPO_PATH")
	if repoPath == "" {
		repoPath, _ = os.Getwd()
	}

	// 解析为绝对路径
	absRepoPath, err := filepath.Abs(repoPath)
	if err != nil {
		log.Fatalf("Error getting absolute path: %v", err)
	}

	// 创建任务管理器
	tm := NewTaskManager(absRepoPath, 30*time.Second, 100)

	// 创建上下文
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 启动任务管理器
	tm.Start(ctx)

}
