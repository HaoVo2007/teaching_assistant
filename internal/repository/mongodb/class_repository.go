package mongodb

import (
	"context"
	"teaching_assistant/internal/domain/class"
	"teaching_assistant/pkg/pagination"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type classRepository struct {
	collection *mongo.Collection
}

func NewClassRepository(db *mongo.Database) class.ClassRepository {
	return &classRepository{
		collection: db.Collection("classes"),
	}
}

func (r *classRepository) Create(ctx context.Context, class *class.Class) error {
	_, err := r.collection.InsertOne(ctx, class)
	if err != nil {
		return err
	}

	return nil
}

func (r *classRepository) GetClasses(ctx context.Context, userId string, params pagination.Params, name string) ([]*class.Class, int64, error) {
	filter := bson.M{"created_by": userId}
	if name != "" {
		filter["name"] = bson.M{"$regex": name, "$options": "i"}
	}

	opts := options.Find().SetSkip(params.Skip()).SetLimit(params.Limit64())
	opts.SetSort(bson.D{{Key: "created_at", Value: -1}})

	total, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	classes := make([]*class.Class, 0)
	for cursor.Next(ctx) {
		var item class.Class
		if err := cursor.Decode(&item); err != nil {
			return nil, 0, err
		}
		classes = append(classes, &item)
	}
	return classes, total, nil
}

func (r *classRepository) GetClassById(ctx context.Context, id primitive.ObjectID) (*class.Class, error) {
	filter := bson.M{"_id": id}
	var item class.Class
	err := r.collection.FindOne(ctx, filter).Decode(&item)
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *classRepository) UpdateClassById(ctx context.Context, c *class.Class) error {
	filter := bson.M{"_id": c.ID}
	_, err := r.collection.ReplaceOne(ctx, filter, c)
	return err
}

func (r *classRepository) DeleteClassById(ctx context.Context, id primitive.ObjectID) error {
	filter := bson.M{"_id": id}
	_, err := r.collection.DeleteOne(ctx, filter)
	return err
}
