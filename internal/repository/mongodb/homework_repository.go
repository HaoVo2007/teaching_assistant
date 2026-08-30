package mongodb

import (
	"context"
	"teaching_assistant/internal/domain/homework"
	"teaching_assistant/pkg/pagination"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type homeworkRepository struct {
	collection *mongo.Collection
}

func NewHomeworkRepository(db *mongo.Database) homework.HomeworkRepository {
	return &homeworkRepository{
		collection: db.Collection("homeworks"),
	}
}

func (r *homeworkRepository) CreateHomework(ctx context.Context, homework *homework.Homework) error {
	_, err := r.collection.InsertOne(ctx, homework)
	if err != nil {
		return err
	}
	return nil
}

func (r *homeworkRepository) GetHomeworks(ctx context.Context, userId string, classId string, params pagination.Params) ([]*homework.Homework, int64, error) {
	filter := bson.M{
		"created_by": userId,
	}

	if classId != "" {
		filter["class_id"] = classId
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

	homeworks := make([]*homework.Homework, 0)
	for cursor.Next(ctx) {
		var homework homework.Homework
		if err := cursor.Decode(&homework); err != nil {
			return nil, 0, err
		}
		homeworks = append(homeworks, &homework)
	}

	return homeworks, total, nil
}

func (r *homeworkRepository) GetHomeworkById(ctx context.Context, id primitive.ObjectID) (*homework.Homework, error) {
	filter := bson.M{
		"_id": id,
	}

	var homework homework.Homework
	if err := r.collection.FindOne(ctx, filter).Decode(&homework); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &homework, nil
}

func (r *homeworkRepository) GetHomeworkByIds(ctx context.Context, ids []primitive.ObjectID) ([]*homework.Homework, error) {
	if len(ids) == 0 {
		return []*homework.Homework{}, nil
	}

	cursor, err := r.collection.Find(ctx, bson.M{"_id": bson.M{"$in": ids}})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	homeworks := make([]*homework.Homework, 0)
	for cursor.Next(ctx) {
		var hw homework.Homework
		if err := cursor.Decode(&hw); err != nil {
			return nil, err
		}
		homeworks = append(homeworks, &hw)
	}
	return homeworks, nil
}

func (r *homeworkRepository) UpdateHomeworkById(ctx context.Context, id primitive.ObjectID, homework *homework.Homework) error {
	_, err := r.collection.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": homework})
	if err != nil {
		return err
	}
	return nil
}

func (r *homeworkRepository) DeleteHomeworkById(ctx context.Context, id primitive.ObjectID) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return err
	}
	return nil
}

func (r *homeworkRepository) CountByQuestionID(ctx context.Context, questionID string) (int64, error) {
	return r.collection.CountDocuments(ctx, bson.M{"questions": questionID})
}

func (r *homeworkRepository) CountByClassID(ctx context.Context, classID string) (int64, error) {
	return r.collection.CountDocuments(ctx, bson.M{"class_id": classID})
}
