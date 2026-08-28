package class

import (
	"context"
	"teaching_assistant/internal/delivery/http/request"
	"teaching_assistant/internal/delivery/http/response"
	"teaching_assistant/pkg/pagination"
)

type ClassService interface {
	CreateClass(ctx context.Context, userId string, req request.CreateClassRequest) error
	GetClasses(ctx context.Context, userId string, params pagination.Params, name string) (*response.ClassResponseWithMeta, error)
	GetClassById(ctx context.Context, userId string, id string) (*response.ClassResponse, error)
	UpdateClassById(ctx context.Context, userId string, id string, req request.UpdateClassRequest) error
	DeleteClassById(ctx context.Context, userId string, id string) error
}
