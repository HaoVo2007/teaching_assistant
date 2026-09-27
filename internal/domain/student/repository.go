package student

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type StudentRepository interface {
	CreateMany(ctx context.Context, students []*Student) error
	GetStudentById(ctx context.Context, id primitive.ObjectID) (*Student, error)
	GetStudentsByIds(ctx context.Context, ids []primitive.ObjectID) ([]*Student, error)
	GetStudentByCode(ctx context.Context, code string) (*Student, error)
	CreateGuardian(ctx context.Context, guardian *Guardian) error
	GetGuardianByParentId(ctx context.Context, parentId string) (*Guardian, error)
	GetGuardianByStudentId(ctx context.Context, studentId string) (*Guardian, error)
	UpdateGuardianParent(ctx context.Context, studentId, parentId string) error
	GetGuardiansByStudentIds(ctx context.Context, studentIds []string) ([]*Guardian, error)
	DeactivateByIDs(ctx context.Context, ids []primitive.ObjectID) error
	ActivateByIDs(ctx context.Context, ids []primitive.ObjectID) error
}
