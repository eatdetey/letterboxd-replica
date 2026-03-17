package mapper

import userpb "github.com/eatdetey/letterboxd-replica/source/api-gateway/gen/go/user/v1"

type UserResponse struct {
	ID        int64  `json:"id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	Bio       string `json:"bio"`
	AvatarURL string `json:"avatar_url"`
	Status    string `json:"status"`
	Role      string `json:"role"`
}

func UserFromPB(u *userpb.User) UserResponse {
	if u == nil {
		return UserResponse{}
	}
	return UserResponse{
		ID:        u.Id,
		Username:  u.Username,
		Email:     u.Email,
		Bio:       u.Bio,
		AvatarURL: u.AvatarUrl,
		Status:    u.Status,
		Role:      u.Role.String(),
	}
}
