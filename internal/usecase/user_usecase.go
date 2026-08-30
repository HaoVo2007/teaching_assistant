package usecase

import (
	"context"
	"errors"
	"teaching_assistant/internal/delivery/http/mapper"
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
	if req.Username == "" {
		return nil, errors.New(string(user.ErrInvalidName))
	}

	if req.Email == "" {
		return nil, errors.New(string(user.ErrInvalidEmail))
	}

	if req.Password == "" {
		return nil, errors.New(string(user.ErrInvalidPassword))
	}

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

	token, err := s.jwtManager.GenerateToken(user.ID.Hex(), user.Username, user.Email, string(user.Role))
	if err != nil {
		return nil, err
	}

	return &response.AuthResponse{
		Token: token,
		User:  mapper.MapUserToUserResponse(user),
	}, nil
}

func (s *userService) Login(ctx context.Context, req request.LoginUserRequest) (*response.AuthResponse, error) {
	if req.Email == "" {
		return nil, errors.New(string(user.ErrInvalidEmail))
	}

	if req.Password == "" {
		return nil, errors.New(string(user.ErrInvalidPassword))
	}

	userRes, err := s.userRepo.FindByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}

	if userRes == nil {
		return nil, errors.New(string(user.ErrUserNotFound))
	}

	if err := bcrypt.CompareHashAndPassword([]byte(userRes.Password), []byte(req.Password)); err != nil {
		return nil, err
	}

	token, err := s.jwtManager.GenerateToken(userRes.ID.Hex(), userRes.Username, userRes.Email, string(userRes.Role))
	if err != nil {
		return nil, err
	}

	return &response.AuthResponse{
		Token: token,
		User:  mapper.MapUserToUserResponse(userRes),
	}, nil
}

func (s *userService) Logout(ctx context.Context, userId string) error {
	if userId == "" {
		return errors.New(string(user.ErrUnauthorized))
	}

	objectId, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		return err
	}

	userRes, err := s.userRepo.FindById(ctx, objectId)
	if err != nil {
		return err
	}

	if userRes == nil {
		return errors.New(string(user.ErrUserNotFound))
	}

	return nil
}
