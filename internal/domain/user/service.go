package user

import (
	"context"
	"teaching_assistant/internal/delivery/http/request"
	"teaching_assistant/internal/delivery/http/response"
)

type UserService interface {
	Register(ctx context.Context, req request.CreateUserRequest) (*response.AuthResponse, error)
}
