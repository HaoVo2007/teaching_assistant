package usecase

import (
	"context"
	"teaching_assistant/internal/delivery/http/mapper"
	"teaching_assistant/internal/delivery/http/request"
	"teaching_assistant/internal/delivery/http/response"
	"teaching_assistant/internal/domain/class"
	"teaching_assistant/internal/domain/student"
	"teaching_assistant/internal/domain/user"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type studentUsecase struct {
	studentRepo student.StudentRepository
	userRepo    user.UserRepository
	classRepo   class.ClassRepository
}

func NewStudentUsecase(
	studentRepo student.StudentRepository,
	userRepo user.UserRepository,
	classRepo class.ClassRepository,
) student.StudentService {
	return &studentUsecase{
		studentRepo: studentRepo,
		userRepo:    userRepo,
		classRepo:   classRepo,
	}
}

func (u *studentUsecase) ClaimStudent(ctx context.Context, teacherId string, req request.ClaimStudentRequest) error {
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
	if studentData.Status != student.StudentStatusActive {
		return student.ErrStudentInactive
	}
	if err := u.ensureStudentOwnedByTeacher(ctx, teacherId, studentData); err != nil {
		return err
	}

	parentOID, err := primitive.ObjectIDFromHex(req.ParentId)
	if err != nil {
		return student.ErrInvalidParentID
	}
	parent, err := u.userRepo.FindById(ctx, parentOID)
	if err != nil {
		if err == user.ErrUserNotFound {
			return student.ErrParentNotFound
		}
		return err
	}
	if parent == nil {
		return student.ErrParentNotFound
	}
	if parent.Role != user.RoleParent {
		return student.ErrNotAParent
	}

	studentID := studentData.ID.Hex()
	existingByStudent, err := u.studentRepo.GetGuardianByStudentId(ctx, studentID)
	if err != nil {
		return err
	}
	existingByParent, err := u.studentRepo.GetGuardianByParentId(ctx, req.ParentId)
	if err != nil {
		return err
	}

	if existingByStudent == nil {
		if existingByParent != nil {
			return student.ErrGuardianAlreadyExists
		}
		return u.studentRepo.CreateGuardian(ctx, &student.Guardian{
			ID:        primitive.NewObjectID(),
			ParentID:  req.ParentId,
			StudentID: studentID,
			Role:      "guardian",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		})
	}

	if existingByStudent.ParentID == req.ParentId {
		return nil
	}
	if existingByParent != nil && existingByParent.StudentID != studentID {
		return student.ErrGuardianAlreadyExists
	}

	return u.studentRepo.UpdateGuardianParent(ctx, studentID, req.ParentId)
}

func (u *studentUsecase) ensureStudentOwnedByTeacher(ctx context.Context, teacherId string, st *student.Student) error {
	if st.CreatedBy != teacherId {
		return student.ErrStudentNotInClass
	}
	classOID, err := primitive.ObjectIDFromHex(st.ClassID)
	if err != nil {
		return student.ErrStudentNotInClass
	}
	cls, err := u.classRepo.GetClassById(ctx, classOID)
	if err != nil || cls == nil {
		return student.ErrStudentNotInClass
	}
	if cls.CreatedBy != teacherId {
		return student.ErrStudentNotInClass
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
