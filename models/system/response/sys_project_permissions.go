package response

import (
	"github.com/spectacleCase/ci-cd-engine/models/common/request"
	"github.com/spectacleCase/ci-cd-engine/models/system"
)

type ProjectList struct {
	List     []system.ProjectPermissions `json:"list"`
	PageInfo request.PageInfo            `json:"pageInfo"`
}
