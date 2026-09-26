package mongodb

import (
	"context"

	"teaching_assistant/internal/domain/user"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

func (r *userRepository) Create(ctx context.Context, item *user.User) error {
	_, err := r.collection.InsertOne(ctx, item)
	if mongo.IsDuplicateKeyError(err) {
		return user.ErrEmailAlreadyExists
	}
	return err
}

func (r *userRepository) FindByEmail(ctx context.Context, email string) (*user.User, error) {
	var item user.User
	err := r.collection.FindOne(ctx, bson.M{"email": email}).Decode(&item)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, user.ErrUserNotFound
		}
		return nil, err
	}
	return &item, err
}

func (r *userRepository) FindById(ctx context.Context, id primitive.ObjectID) (*user.User, error) {
	filter := bson.M{"_id": id}
	var item user.User
	err := r.collection.FindOne(ctx, filter).Decode(&item)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, user.ErrUserNotFound
		}
		return nil, err
	}
	return &item, nil
}

func (r *userRepository) FindByIds(ctx context.Context, ids []primitive.ObjectID) ([]*user.User, error) {
	if len(ids) == 0 {
		return []*user.User{}, nil
	}

	cursor, err := r.collection.Find(ctx, bson.M{"_id": bson.M{"$in": ids}})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	users := make([]*user.User, 0)
	for cursor.Next(ctx) {
		var item user.User
		if err := cursor.Decode(&item); err != nil {
			return nil, err
		}
		users = append(users, &item)
	}
	return users, nil
}
