package class

import (
	"context"
	"teaching_assistant/pkg/pagination"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ClassRepository interface {
	Create(ctx context.Context, class *Class) error
	GetClasses(ctx context.Context, userId string, params pagination.Params, name string) ([]*Class, int64, error)
	GetClassById(ctx context.Context, id primitive.ObjectID) (*Class, error)
	UpdateClassById(ctx context.Context, class *Class) error
	DeleteClassById(ctx context.Context, id primitive.ObjectID) error
}
