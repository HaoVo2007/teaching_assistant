package usecase

import (
	"context"
	"teaching_assistant/internal/domain/user"
)

type userService struct {
	userRepo user.UserRepository
}

func NewUserService(userRepo user.UserRepository) user.UserService {
	return &userService{
		userRepo: userRepo,
	}
}

func (s *userService) CreateUser(ctx context.Context, user *user.User) error {
	return nil
}
