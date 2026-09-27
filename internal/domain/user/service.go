package user

import (
	"context"
	"teaching_assistant/internal/delivery/http/request"
	"teaching_assistant/internal/delivery/http/response"
	"teaching_assistant/pkg/pagination"
)

type UserService interface {
	Register(ctx context.Context, req request.CreateUserRequest) (*response.AuthResponse, error)
	Login(ctx context.Context, req request.LoginUserRequest) (*response.AuthResponse, error)
	Logout(ctx context.Context, userId string) error
	CreateUser(ctx context.Context, req request.CreateUserRequest) (*response.UserResponse, error)
	GetParents(ctx context.Context, params pagination.Params, q string) (*response.UserResponseWithMeta, error)
}
