package system

import (
	"github.com/spectacleCase/ci-cd-engine/config"
	"github.com/spectacleCase/ci-cd-engine/models"
	"github.com/spectacleCase/ci-cd-engine/utils"
	"time"
)

// Users 用户表
type Users struct {
	models.BaseMODEL
	Username     string    `json:"username"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"password_hash"`
	Salt         string    `json:"salt"`
	IsAdmin      bool      `json:"is_admin" gorm:"type:tinyint(1)"`
	IsActive     bool      `json:"is_active" gorm:"type:tinyint(1)"`
	LastLoginAt  time.Time `json:"last_login_at"`
	LastLoginIP  string    `json:"last_login_ip"`
}

func (u *Users) SetPasswordHash(password string) error {
	salt, err := utils.GenerateSalt(config.Config.System.SaltLength)
	if err != nil {
		return err
	}
	u.Salt = salt
	u.PasswordHash = utils.BcryptHash(password + salt)
	return nil
}
