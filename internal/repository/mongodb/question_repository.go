package mongodb

import (
	"context"
	"teaching_assistant/internal/domain/question"

	"go.mongodb.org/mongo-driver/mongo"
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
