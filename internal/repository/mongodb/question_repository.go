package mongodb

import (
	"context"
	"teaching_assistant/internal/domain/question"
	"teaching_assistant/pkg/pagination"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type questionRepository struct {
	collection *mongo.Collection
}

func NewQuestionRepository(db *mongo.Database) question.QuestionRepository {
	return &questionRepository{
		collection: db.Collection("questions"),
	}
}

func (r *questionRepository) Create(ctx context.Context, question *question.Question) error {
	_, err := r.collection.InsertOne(ctx, question)
	return err
}

func (r *questionRepository) GetQuestions(ctx context.Context, userId string, params pagination.Params, questionType, questionName, subject, grade, difficulty string) ([]*question.Question, int64, error) {
	filter := bson.M{"created_by": userId}
	if questionType != "" {
		filter["type"] = questionType
	}
	if questionName != "" {
		filter["question"] = bson.M{"$regex": questionName, "$options": "i"}
	}
	if subject != "" {
		filter["subject"] = subject
	}
	if grade != "" {
		filter["grade"] = grade
	}
	if difficulty != "" {
		filter["difficulty"] = difficulty
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
	questions := make([]*question.Question, 0)
	for cursor.Next(ctx) {
		var q question.Question
		if err := cursor.Decode(&q); err != nil {
			return nil, 0, err
		}
		questions = append(questions, &q)
	}
	return questions, total, nil
}

func (r *questionRepository) GetQuestionById(ctx context.Context, id primitive.ObjectID) (*question.Question, error) {
	filter := bson.M{"_id": id}
	var question question.Question
	err := r.collection.FindOne(ctx, filter).Decode(&question)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &question, nil
}

func (r *questionRepository) GetQuestionByIds(ctx context.Context, ids []primitive.ObjectID) ([]*question.Question, error) {
	filter := bson.M{"_id": bson.M{"$in": ids}}
	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	questions := make([]*question.Question, 0)
	for cursor.Next(ctx) {
		var q question.Question
		if err := cursor.Decode(&q); err != nil {
			return nil, err
		}
		questions = append(questions, &q)
	}
	return questions, nil
}

func (r *questionRepository) Update(ctx context.Context, q *question.Question) error {
	filter := bson.M{"_id": q.ID}
	_, err := r.collection.ReplaceOne(ctx, filter, q)
	return err
}

func (r *questionRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	filter := bson.M{"_id": id}
	_, err := r.collection.DeleteOne(ctx, filter)
	return err
}
