package user

import (
	"context"
	"teaching_assistant/internal/delivery/http/request"
	"teaching_assistant/internal/delivery/http/response"
)

type UserService interface {
	Register(ctx context.Context, req request.CreateUserRequest) (*response.AuthResponse, error)
	Login(ctx context.Context, req request.LoginUserRequest) (*response.AuthResponse, error)
	Logout(ctx context.Context, userId string) error
}
