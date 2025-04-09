package system

import "github.com/spectacleCase/ci-cd-engine/models"

// ProjectPermissions 项目表
type ProjectPermissions struct {
	models.BaseMODEL
	UserId          uint   `json:"user_id"`
	ProjectId       string `json:"project_id"`
	PermissionLevel string `json:"permission_level"` //read,write,admin
	ProjectName     string `json:"project_name"`
}
