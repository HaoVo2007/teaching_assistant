package user

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserRepository interface {
	Create(ctx context.Context, user *User) error
	FindById(ctx context.Context, id primitive.ObjectID) (*User, error)
	FindByEmail(ctx context.Context, email string) (*User, error)
}
