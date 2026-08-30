package homeworksubmission

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type HomeworkSubmission struct {
	ID             primitive.ObjectID `bson:"_id" json:"id"`
	HomeworkID     string             `bson:"homework_id" json:"homework_id"`
	StudentName    string             `bson:"student_name" json:"student_name"`
	IsSubmitted    bool               `bson:"is_submitted" json:"is_submitted"`
	StudentAnswers []StudentAnswer    `bson:"student_answers" json:"student_answers"`
	SubmittedAt    time.Time          `bson:"submitted_at" json:"submitted_at"`
	TeacherID      string             `bson:"teacher_id" json:"teacher_id"`
	CreatedAt      time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt      time.Time          `bson:"updated_at" json:"updated_at"`
}

type StudentAnswer struct {
	QuestionID    string `bson:"question_id" json:"question_id"`
	SelectedIndex *int   `bson:"selected_index" json:"selected_index"`
	SelectedBool  *bool  `bson:"selected_bool" json:"selected_bool"`
}
