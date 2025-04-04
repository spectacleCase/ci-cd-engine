package models

import (
	"time"
)

type BaseMODEL struct {
	ID        uint      `gorm:"primarykey;autoIncrement" json:"id"` // 主键ID
	CreatedAt time.Time `json:"created_at"`                         // 创建时间
	UpdatedAt time.Time `json:"updated_at"`                         // 更新时间
	IsDelete  int       // 逻辑删除
}
