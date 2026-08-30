package homeworksubmission

import (
	"context"
	"teaching_assistant/pkg/pagination"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type HomeworkSubmissionRepository interface {
	CreateHomeworkSubmission(ctx context.Context, submission *HomeworkSubmission) error
	GetHomeworkSubmissionsByUserId(ctx context.Context, userId string, params pagination.Params) ([]*HomeworkSubmission, int64, error)
	GetHomeworkSubmissionById(ctx context.Context, id primitive.ObjectID) (*HomeworkSubmission, error)
	GetHomeworkSubmissionsByHomeworkId(ctx context.Context, homeworkID string, userId string, params pagination.Params) ([]*HomeworkSubmission, int64, error)
	CountByHomeworkID(ctx context.Context, homeworkID string) (int64, error)
	CountByQuestionID(ctx context.Context, questionID string) (int64, error)
}
