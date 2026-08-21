package usecase

import (
	"context"
	"teaching_assistant/internal/delivery/http/request"
	"teaching_assistant/internal/delivery/http/response"
	"teaching_assistant/internal/domain/user"
	"teaching_assistant/pkg/jwt"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

type userService struct {
	userRepo   user.UserRepository
	jwtManager *jwt.Manager
}

func NewUserService(
	userRepo user.UserRepository,
	jwtManager *jwt.Manager,
) user.UserService {
	return &userService{
		userRepo:   userRepo,
		jwtManager: jwtManager,
	}
}

func (s *userService) Register(ctx context.Context, req request.CreateUserRequest) (*response.AuthResponse, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &user.User{
		ID:        primitive.NewObjectID(),
		Username:  req.Username,
		Email:     req.Email,
		Password:  string(hash),
		Role:      user.RoleUser,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	token, err := s.jwtManager.GenerateToken(user.ID.Hex(), user.Username, user.Email, user.Role)
	if err != nil {
		return nil, err
	}

	return &response.AuthResponse{
		Token: token,
		User: response.UserResponse{
			ID:       user.ID.Hex(),
			Username: user.Username,
			Email:    user.Email,
			Role:     string(user.Role),
		},
	}, nil
}
