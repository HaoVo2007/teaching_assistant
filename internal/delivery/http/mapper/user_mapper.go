package mapper

import (
	"teaching_assistant/internal/delivery/http/response"
	"teaching_assistant/internal/domain/user"
)

func MapUserToUserResponse(user *user.User) response.UserResponse {
	return response.UserResponse{
		ID:       user.ID.Hex(),
		Username: user.Username,
		Email:    user.Email,
		Role:     string(user.Role),
	}
}
