package request

type CreateProject struct {
	Id          string `json:"id"`
	ProjectName string `json:"project_name"`
	ProjectId   string `json:"project_id"`
}

type DeleteProject struct {
	Id string `json:"id"`
}

type AddCrew struct {
	UserId          string `json:"user_id"`
	ProjectId       string `json:"project_id"`
	PermissionLevel string `json:"permission_level"`
}
