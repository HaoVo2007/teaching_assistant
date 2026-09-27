package usecase

import (
	"context"
	"teaching_assistant/internal/delivery/http/mapper"
	"teaching_assistant/internal/delivery/http/request"
	"teaching_assistant/internal/delivery/http/response"
	"teaching_assistant/internal/domain/user"
	"teaching_assistant/pkg/jwt"
	"teaching_assistant/pkg/pagination"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

type userUsecase struct {
	userRepo   user.UserRepository
	jwtManager *jwt.Manager
}

func NewUserUsecase(
	userRepo user.UserRepository,
	jwtManager *jwt.Manager,
) user.UserService {
	return &userUsecase{
		userRepo:   userRepo,
		jwtManager: jwtManager,
	}
}

func (s *userUsecase) Register(ctx context.Context, req request.CreateUserRequest) (*response.AuthResponse, error) {
	if req.Username == "" {
		return nil, user.ErrInvalidName
	}

	if req.Email == "" {
		return nil, user.ErrInvalidEmail
	}

	if req.Password == "" {
		return nil, user.ErrInvalidPassword
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	item := &user.User{
		ID:        primitive.NewObjectID(),
		Username:  req.Username,
		Email:     req.Email,
		Password:  string(hash),
		Role:      user.RoleTeacher,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.userRepo.Create(ctx, item); err != nil {
		return nil, err
	}

	token, err := s.jwtManager.GenerateToken(item.ID.Hex(), item.Username, item.Email, string(item.Role))
	if err != nil {
		return nil, err
	}

	return &response.AuthResponse{
		Token: token,
		User:  *mapper.MapUserToUserResponse(item),
	}, nil
}

func (s *userUsecase) Login(ctx context.Context, req request.LoginUserRequest) (*response.AuthResponse, error) {
	if req.Email == "" {
		return nil, user.ErrInvalidEmail
	}

	if req.Password == "" {
		return nil, user.ErrInvalidPassword
	}

	userRes, err := s.userRepo.FindByEmail(ctx, req.Email)
	if err != nil {
		if err == user.ErrUserNotFound {
			return nil, user.ErrInvalidCredentials
		}
		return nil, err
	}

	if userRes == nil {
		return nil, user.ErrInvalidCredentials
	}

	if err := bcrypt.CompareHashAndPassword([]byte(userRes.Password), []byte(req.Password)); err != nil {
		return nil, user.ErrInvalidCredentials
	}

	token, err := s.jwtManager.GenerateToken(userRes.ID.Hex(), userRes.Username, userRes.Email, string(userRes.Role))
	if err != nil {
		return nil, err
	}

	return &response.AuthResponse{
		Token: token,
		User:  *mapper.MapUserToUserResponse(userRes),
	}, nil
}

func (s *userUsecase) Logout(ctx context.Context, userId string) error {
	if userId == "" {
		return user.ErrUnauthorized
	}

	objectId, err := primitive.ObjectIDFromHex(userId)
	if err != nil {
		return user.ErrUnauthorized
	}

	userRes, err := s.userRepo.FindById(ctx, objectId)
	if err != nil {
		return err
	}

	if userRes == nil {
		return user.ErrUserNotFound
	}

	return nil
}

func (s *userUsecase) CreateUser(ctx context.Context, req request.CreateUserRequest) (*response.UserResponse, error) {
	if req.Username == "" {
		return nil, user.ErrInvalidName
	}

	if req.Email == "" {
		return nil, user.ErrInvalidEmail
	}

	if req.Password == "" {
		return nil, user.ErrInvalidPassword
	}

	if req.Role != string(user.RoleTeacher) && req.Role != string(user.RoleParent) {
		return nil, user.ErrInvalidRole
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	item := &user.User{
		ID:        primitive.NewObjectID(),
		Username:  req.Username,
		Email:     req.Email,
		Password:  string(hash),
		Role:      user.Role(req.Role),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.userRepo.Create(ctx, item); err != nil {
		return nil, err
	}

	return mapper.MapUserToUserResponse(item), nil
}

func (s *userUsecase) GetParents(ctx context.Context, params pagination.Params, q string) (*response.UserResponseWithMeta, error) {
	parents, total, err := s.userRepo.FindParents(ctx, params, q)
	if err != nil {
		return nil, err
	}

	return &response.UserResponseWithMeta{
		Parents: mapper.MapUsersToResponses(parents),
		Meta:    pagination.NewMeta(params, total),
	}, nil
}
