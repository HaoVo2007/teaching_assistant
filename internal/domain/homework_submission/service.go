package homeworksubmission

import (
	"context"
	"teaching_assistant/internal/delivery/http/request"
	"teaching_assistant/internal/delivery/http/response"
	"teaching_assistant/pkg/pagination"
)

type HomeworkSubmissionService interface {
	CreateHomeworkSubmission(ctx context.Context, req request.CreateHomeworkSubmissionRequest) error
	GetHomeworkSubmissions(ctx context.Context, params pagination.Params, userId string) (*response.HomeworkSubmissionResponseWithMeta, error)
	GetHomeworkSubmissionById(ctx context.Context, id string, userId string) (*response.HomeworkSubmissionResponse, error)
	GetHomeworkSubmissionsByHomeworkId(ctx context.Context, homeworkId string, userId string, params pagination.Params) (*response.HomeworkSubmissionResponseWithMeta, error)
}
