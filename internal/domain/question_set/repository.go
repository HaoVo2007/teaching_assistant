package questionset

import (
	"context"
	"teaching_assistant/pkg/pagination"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type QuestionSetRepository interface {
	Create(ctx context.Context, questionSet *QuestionSet) error
	GetQuestionSets(ctx context.Context, userId string, params pagination.Params, title string, questionType string) ([]*QuestionSet, int64, error)
	GetQuestionSetById(ctx context.Context, id primitive.ObjectID) (*QuestionSet, error)
	UpdateQuestionSetById(ctx context.Context, id primitive.ObjectID, questionSet *QuestionSet) error
	DeleteQuestionSetById(ctx context.Context, id primitive.ObjectID) error
	CountByQuestionID(ctx context.Context, questionID string) (int64, error)
}
