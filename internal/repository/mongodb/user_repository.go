package mongodb

import (
	"context"

	"teaching_assistant/internal/domain/user"

	"go.mongodb.org/mongo-driver/mongo"
)

type userRepository struct {
	collection *mongo.Collection
}

func NewUserRepository(db *mongo.Database) user.UserRepository {
	return &userRepository{
		collection: db.Collection("users"),
	}
}

func (r *userRepository) CreateUser(ctx context.Context, user *user.User) error {
	_, err := r.collection.InsertOne(ctx, user)
	return err
}
