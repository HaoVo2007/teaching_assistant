package questionset

import (
	"context"
	"teaching_assistant/internal/delivery/http/request"
	"teaching_assistant/internal/delivery/http/response"
	"teaching_assistant/pkg/pagination"
)

type QuestionSetService interface {
	CreateQuestionSet(ctx context.Context, userId string, req request.CreateQuestionSetRequest) error
	GetQuestionSets(ctx context.Context, userId string, params pagination.Params, title string, questionType string) (*response.QuestionSetResponseWithMeta, error)
	GetQuestionSetById(ctx context.Context, userId string, id string) (*response.QuestionSetResponse, error)
	UpdateQuestionSetById(ctx context.Context, userId string, id string, req request.UpdateQuestionSetRequest) error
	DeleteQuestionSetById(ctx context.Context, userId string, id string) error
}
