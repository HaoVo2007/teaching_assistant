package question

import (
	"context"

	"teaching_assistant/internal/delivery/http/request"
)

type QuestionService interface {
	CreateQuestion(ctx context.Context, req request.CreateQuestionRequest) (*Question, error)
}
