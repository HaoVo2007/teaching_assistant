package homework

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Homework struct {
	ID          primitive.ObjectID `bson:"_id"`
	ClassID     string             `bson:"class_id"`
	Title       string             `bson:"title"`
	Description *string            `bson:"description"`
	Questions   []string           `bson:"questions"`
	DueDate     time.Time          `bson:"due_date"`
	CreatedBy   string             `bson:"created_by"`
	CreatedAt   time.Time          `bson:"created_at"`
	UpdatedAt   time.Time          `bson:"updated_at"`
}
