package questionset

import (
	"context"
	"teaching_assistant/internal/delivery/http/request"
	"teaching_assistant/internal/delivery/http/response"
	"teaching_assistant/pkg/pagination"
)

type QuestionSetService interface {
	CreateQuestionSet(ctx context.Context, userId string, req request.CreateQuestionSetRequest) error
	GetQuestionSets(ctx context.Context, userId string, params pagination.Params) (*response.QuestionSetResponseWithMeta, error)
}
