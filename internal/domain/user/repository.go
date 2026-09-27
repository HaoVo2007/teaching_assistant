package user

import (
	"context"

	"teaching_assistant/pkg/pagination"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserRepository interface {
	Create(ctx context.Context, user *User) error
	FindById(ctx context.Context, id primitive.ObjectID) (*User, error)
	FindByIds(ctx context.Context, ids []primitive.ObjectID) ([]*User, error)
	FindByEmail(ctx context.Context, email string) (*User, error)
	FindParents(ctx context.Context, params pagination.Params, q string) ([]*User, int64, error)
}
