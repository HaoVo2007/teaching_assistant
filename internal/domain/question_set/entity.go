package questionset

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type QuestionSet struct {
	ID          primitive.ObjectID `bson:"_id" json:"id"`
	Title       string             `bson:"title" json:"title"`
	Description *string            `bson:"description" json:"description"`
	QuestionIds []string           `bson:"question_ids" json:"question_ids"`
	CreatedBy   string             `bson:"created_by" json:"created_by"`
	CreatedAt   time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt   time.Time          `bson:"updated_at" json:"updated_at"`
}
