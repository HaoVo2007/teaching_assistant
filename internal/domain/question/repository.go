package question

import (
	"context"
	"teaching_assistant/pkg/pagination"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type QuestionRepository interface {
	Create(ctx context.Context, q *Question) error
	GetQuestions(ctx context.Context, userId string, params pagination.Params) ([]*Question, int64, error)
	GetQuestionById(ctx context.Context, id primitive.ObjectID) (*Question, error)
	GetQuestionByIds(ctx context.Context, ids []primitive.ObjectID) ([]*Question, error)
	Update(ctx context.Context, q *Question) error
	Delete(ctx context.Context, id primitive.ObjectID) error
}
