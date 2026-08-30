package homework

import (
	"context"
	"teaching_assistant/internal/delivery/http/request"
	"teaching_assistant/internal/delivery/http/response"
	"teaching_assistant/pkg/pagination"
)

type HomeworkService interface {
	CreateHomework(ctx context.Context, userId string, req request.CreateHomeworkRequest) error
	GetHomeworks(ctx context.Context, userId string, params pagination.Params) (*response.HomeworkResponseWithMeta, error)
	GetHomeworkById(ctx context.Context, userId string, id string) (*response.HomeworkResponse, error)
	UpdateHomeworkById(ctx context.Context, userId string, id string, req request.UpdateHomeworkRequest) error
	DeleteHomeworkById(ctx context.Context, userId string, id string) error
	GetHomeworksByClassId(ctx context.Context, userId string, classId string, params pagination.Params) (*response.HomeworkResponseWithMeta, error)
}
