package student

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type StudentStatus string

const (
	StudentStatusActive   StudentStatus = "active"
	StudentStatusInactive StudentStatus = "inactive"
)

type Student struct {
	ID        primitive.ObjectID `bson:"_id" json:"id"`
	Code      string             `bson:"code" json:"code"`
	Image     string             `bson:"image" json:"image"`
	PublicID  string             `bson:"public_id" json:"public_id"`
	Name      string             `bson:"name" json:"name"`
	ClassID   string             `bson:"class_id" json:"class_id"`
	Status    StudentStatus      `bson:"status" json:"status"`
	CreatedBy string             `bson:"created_by" json:"created_by"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at" json:"updated_at"`
}

type Guardian struct {
	ID        primitive.ObjectID `bson:"_id"`
	ParentID  string             `bson:"parent_id"`
	StudentID string             `bson:"student_id"`
	Role      string             `bson:"role"` // father, mother, guardian
	CreatedAt time.Time          `bson:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at"`
}
