package pkg

const (
	StatusPending   = "PENDING"   // 等待中（刚创建）
	StatusQueued    = "QUEUED"    // 已排队
	StatusRunning   = "RUNNING"   // 运行中
	StatusPaused    = "PAUSED"    // 已暂停
	StatusCompleted = "COMPLETED" // 已完成
	StatusFailed    = "Failed"    // 已失败
	StatusCancelled = "Cancelled" // 已取消
	StatusSkipped   = "Skipped"   // 已跳过
)
