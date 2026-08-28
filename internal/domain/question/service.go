package question

import (
	"context"

	"teaching_assistant/internal/delivery/http/request"
	"teaching_assistant/internal/delivery/http/response"
	"teaching_assistant/pkg/pagination"
)

type QuestionService interface {
	CreateQuestion(ctx context.Context, req request.CreateQuestionRequest, userId string) (*Question, error)
	GetQuestions(ctx context.Context, userId string, params pagination.Params, questionType, questionName, subject, grade, difficulty string) (*response.QuestionResponseWithMeta, error)
	GetQuestionById(ctx context.Context, id string) (*response.QuestionResponse, error)
	UpdateQuestionById(ctx context.Context, id string, req request.UpdateQuestionRequest, userId string) (*Question, error)
	DeleteQuestionById(ctx context.Context, id string, userId string) error
}
