package student

import (
	"context"
	"teaching_assistant/internal/delivery/http/request"
	"teaching_assistant/internal/delivery/http/response"
)

type StudentService interface {
	ClaimStudent(ctx context.Context, req request.ClaimStudentRequest) error
	GetStudentsByGuardian(ctx context.Context, userId string) (*response.StudentResponse, error)
}
