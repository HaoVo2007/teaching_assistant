package usecase

import (
	"context"
	"teaching_assistant/internal/delivery/http/mapper"
	"teaching_assistant/internal/delivery/http/request"
	"teaching_assistant/internal/delivery/http/response"
	"teaching_assistant/internal/domain/student"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type studentUsecase struct {
	studentRepo student.StudentRepository
}

func NewStudentUsecase(studentRepo student.StudentRepository) student.StudentService {
	return &studentUsecase{
		studentRepo: studentRepo,
	}
}

func (u *studentUsecase) ClaimStudent(ctx context.Context, req request.ClaimStudentRequest) error {
	if req.Code == "" {
		return student.ErrInvalidStudentCode
	}

	if req.ParentId == "" {
		return student.ErrInvalidParentID
	}

	studentData, err := u.studentRepo.GetStudentByCode(ctx, req.Code)
	if err != nil {
		return err
	}

	if studentData == nil {
		return student.ErrStudentNotFound
	}

	guardianData, err := u.studentRepo.GetGuardianByParentId(ctx, req.ParentId)
	if err != nil {
		return err
	}

	if guardianData != nil {
		return student.ErrGuardianAlreadyExists
	}

	guardian := &student.Guardian{
		ID:        primitive.NewObjectID(),
		ParentID:  req.ParentId,
		StudentID: studentData.ID.Hex(),
		Role:      "guardian",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := u.studentRepo.CreateGuardian(ctx, guardian); err != nil {
		return err
	}

	return nil
}

func (u *studentUsecase) GetStudentsByGuardian(ctx context.Context, userId string) (*response.StudentResponse, error) {
	guardian, err := u.studentRepo.GetGuardianByParentId(ctx, userId)
	if err != nil {
		return nil, err
	}

	if guardian == nil {
		return nil, student.ErrGuardianNotFound
	}

	objectId, err := primitive.ObjectIDFromHex(guardian.StudentID)
	if err != nil {
		return nil, student.ErrStudentNotFound
	}

	studentData, err := u.studentRepo.GetStudentById(ctx, objectId)
	if err != nil {
		return nil, err
	}

	if studentData == nil {
		return nil, student.ErrStudentNotFound
	}

	return mapper.MapStudentToResponse(studentData, nil), nil
}
