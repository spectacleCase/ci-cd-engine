package system

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/spectacleCase/ci-cd-engine/common"
	"github.com/spectacleCase/ci-cd-engine/global"
	system "github.com/spectacleCase/ci-cd-engine/models/system"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"strconv"
	"time"
)

// NewTaskManager 创建新的任务管理器
func NewTaskManager(repoPath []string, pollInterval time.Duration, queueSize int) *system.TaskManager {
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
func Start(ctx *context.Context) {
	// 启动仓库轮询
	go func() {
		err := pollRepositoryChanges(*ctx)
		if err != nil {

		}
	}()

	// 启动任务消费者
	go consumeTasks(*ctx)
}

// AddTask 添加新任务
func AddTask(c context.Context, task *system.Task) error {
	global.CTaskManager.Mu.Lock()
	defer global.CTaskManager.Mu.Unlock()

	// 检查内存中是否已存在
	if _, exists := global.CTaskManager.Tasks[strconv.Itoa(int(task.ID))]; exists {
		return errors.New("task already exists in manager")
	}

	// 设置新状态
	task.Status = common.StatusQueued

	err := global.NewDBClient(c).Model(system.Task{}).Where("id = ?", task.ID).Update("status", task.Status).Error
	if err != nil {
		return err
	}

	// 添加到管理器
	global.CTaskManager.Tasks[strconv.Itoa(int(task.ID))] = task
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
			commit, err := getGitCommitInfo(global.CTaskManager.RepoPath[0])
			if err != nil {
				global.CLog.Error("获取Git提交信息失败", zap.Error(err))
			}
			if err != nil {
				continue
			}

			if lastHash == "" {
				lastHash = commit.Hash.String()
				continue
			}

			if commit.Hash.String() != lastHash {
				lastHash = commit.Hash.String()
				global.CLog.Info("Git提交信息",
					zap.String("hash", commit.Hash.String()),
					zap.String("author", commit.Author.Name),
					zap.Time("date", commit.Author.When),
					zap.String("message", commit.Message),
				)
				ciCdConfig, err := Analyze("file/ci-yaml/.cicd.yaml")
				jsonString, _ := json.Marshal(ciCdConfig)
				task := &system.Task{
					Name:    commit.Message,
					Payload: jsonString,
					Status:  common.StatusPending,
				}
				err = global.CDB.Create(task).Error
				if err != nil {

				}
				if err := AddTask(ctx, task); err != nil {
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
			err := processTask(global.NewDBClient(ctx), task)
			if err != nil {
				return
			}
		}
	}
}

// processTask 处理单个任务
func processTask(db *gorm.DB, task *system.Task) error {
	if task.Status == common.StatusQueued {
		global.CLog.Info("process task", zap.Uint("id", task.ID))
		global.CTaskManager.Mu.Lock()
		task.Status = common.StatusRunning
		err := db.Model(&system.Task{}).Where("id = ?", task.ID).Update("status", task.Status).Error
		if err != nil {
			return err
		}
		var newConfig system.CiCdConfig
		if err := json.Unmarshal(task.Payload, &newConfig); err != nil {
			global.CLog.Error("JSON反序列化失败", zap.String("payload", string(task.Payload)))
		}
		stageMap, _ := AnalyzeToMap(newConfig)
		AssemblyLineProject(stageMap["Build"], stageMap["Deploy"])
		task.Status = common.StatusCompleted
		err = db.Model(&system.Task{}).Where("id = ?", task.ID).Update("status", task.Status).Error
		if err != nil {
			return err
		}
		global.CTaskManager.Mu.Unlock()
		delete(global.CTaskManager.Tasks, strconv.Itoa(int(task.ID)))
		global.CLog.Info("执行成功")

	}
	return nil
}

func getGitCommitInfo(repoPath string) (*object.Commit, error) {
	repo, err := git.PlainOpen(repoPath)
	if err != nil {
		return nil, err
	}
	ref, err := repo.Head()
	if err != nil {
		return nil, err
	}
	return repo.CommitObject(ref.Hash())
}

// InitTaskManager 初始化任务管理器
func InitTaskManager() (context.Context, context.CancelFunc) {
	// 获取仓库路径，默认为当前目录
	//repoPath := os.Getenv("REPO_PATH")
	//if repoPath == "" {
	//	repoPath, _ = os.Getwd()
	//}
	//
	//// 解析为绝对路径
	//absRepoPath, err := filepath.Abs(repoPath)
	//if err != nil {
	//	global.CLog.Error("Error getting absolute path", zap.Error(err))
	//}
	//absRepoPathList := []string{absRepoPath}
	absRepoPathList := []string{
		"D:\\Workspace\\Python\\cicdDemo",
	}
	// 创建任务管理器
	tm := NewTaskManager(absRepoPathList, 5*time.Second, 100)
	global.CTaskManager = tm

	// 创建上下文
	ctx, cancel := context.WithCancel(context.Background())
	return ctx, cancel

}
