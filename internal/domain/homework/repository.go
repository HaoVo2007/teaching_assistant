package homework

import (
	"context"
	"teaching_assistant/pkg/pagination"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type HomeworkRepository interface {
	CreateHomework(ctx context.Context, homework *Homework) error
	GetHomeworks(ctx context.Context, userId string, classId string, params pagination.Params) ([]*Homework, int64, error)
	GetHomeworkById(ctx context.Context, id primitive.ObjectID) (*Homework, error)
	GetHomeworkByIds(ctx context.Context, ids []primitive.ObjectID) ([]*Homework, error)
	UpdateHomeworkById(ctx context.Context, id primitive.ObjectID, homework *Homework) error
	DeleteHomeworkById(ctx context.Context, id primitive.ObjectID) error
	CountByQuestionID(ctx context.Context, questionID string) (int64, error)
	CountByClassID(ctx context.Context, classID string) (int64, error)
}
