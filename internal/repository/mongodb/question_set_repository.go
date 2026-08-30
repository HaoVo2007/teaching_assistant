package mongodb

import (
	"context"
	questionset "teaching_assistant/internal/domain/question_set"
	"teaching_assistant/pkg/pagination"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type questionSetRepository struct {
	collection *mongo.Collection
}

func NewQuestionSetRepository(db *mongo.Database) questionset.QuestionSetRepository {
	return &questionSetRepository{
		collection: db.Collection("question_sets"),
	}
}

func (r *questionSetRepository) Create(ctx context.Context, questionSet *questionset.QuestionSet) error {
	_, err := r.collection.InsertOne(ctx, questionSet)
	if err != nil {
		return err
	}
	return nil
}

func (r *questionSetRepository) GetQuestionSets(ctx context.Context, userId string, params pagination.Params, title string, questionType string) ([]*questionset.QuestionSet, int64, error) {
	filter := bson.M{"created_by": userId}
	if title != "" {
		filter["title"] = bson.M{"$regex": title, "$options": "i"}
	}
	if questionType != "" {
		filter["question_type"] = questionType
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
	questionSets := make([]*questionset.QuestionSet, 0)
	for cursor.Next(ctx) {
		var questionSet questionset.QuestionSet
		if err := cursor.Decode(&questionSet); err != nil {
			return nil, 0, err
		}
		questionSets = append(questionSets, &questionSet)
	}
	return questionSets, total, nil
}

func (r *questionSetRepository) GetQuestionSetById(ctx context.Context, id primitive.ObjectID) (*questionset.QuestionSet, error) {
	filter := bson.M{"_id": id}
	var questionSet questionset.QuestionSet
	err := r.collection.FindOne(ctx, filter).Decode(&questionSet)
	if err != nil {
		return nil, err
	}
	return &questionSet, nil
}

func (r *questionSetRepository) UpdateQuestionSetById(ctx context.Context, id primitive.ObjectID, questionSet *questionset.QuestionSet) error {
	filter := bson.M{"_id": id}
	_, err := r.collection.UpdateOne(ctx, filter, bson.M{"$set": questionSet})
	if err != nil {
		return err
	}
	return nil
}

func (r *questionSetRepository) DeleteQuestionSetById(ctx context.Context, id primitive.ObjectID) error {
	filter := bson.M{"_id": id}
	_, err := r.collection.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}
	return nil
}

func (r *questionSetRepository) CountByQuestionID(ctx context.Context, questionID string) (int64, error) {
	return r.collection.CountDocuments(ctx, bson.M{"question_ids": questionID})
}
