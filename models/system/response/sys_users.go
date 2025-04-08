package response

import "github.com/spectacleCase/ci-cd-engine/models/system"

type LoginResponse struct {
	User      system.Users `json:"user"`
	Token     string       `json:"token"`
	ExpiresAt int64        `json:"expiresAt"`
}
