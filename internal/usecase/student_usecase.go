package usecase

import (
	"context"
	"errors"
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

func (u *studentUsecase) ClaimStudent(ctx context.Context, userId string, req request.ClaimStudentRequest) error {
	if req.Code == "" {
		return errors.New(string(student.ErrInvalidStudentCode))
	}

	studentData, err := u.studentRepo.GetStudentByCode(ctx, req.Code)
	if err != nil {
		return err
	}

	if studentData == nil {
		return errors.New(string(student.ErrStudentNotFound))
	}

	guardian := &student.Guardian{
		ID:        primitive.NewObjectID(),
		ParentID:  userId,
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
		return nil, errors.New(string(student.ErrGuardianNotFound))
	}

	objectId, err := primitive.ObjectIDFromHex(guardian.StudentID)
	if err != nil {
		return nil, err
	}

	studentData, err := u.studentRepo.GetStudentById(ctx, objectId)
	if err != nil {
		return nil, err
	}

	if studentData == nil {
		return nil, errors.New(string(student.ErrStudentNotFound))
	}

	studentResponse := mapper.MapStudentToResponse(studentData, nil)

	return studentResponse, nil
}
