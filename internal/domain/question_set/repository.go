package questionset

import (
	"context"
	"teaching_assistant/pkg/pagination"
)

type QuestionSetRepository interface {
	Create(ctx context.Context, questionSet *QuestionSet) error
	GetQuestionSets(ctx context.Context, userId string, params pagination.Params) ([]*QuestionSet, int64, error)
}
