package response

import (
	"teaching_assistant/pkg/pagination"
	"time"
)

type HomeworkResponse struct {
	ID          string             `json:"id"`
	ClassID     string             `json:"class_id"`
	Title       string             `json:"title"`
	Description *string             `json:"description"`
	DueDate     string             `json:"due_date"`
	Questions   []*QuestionResponse `json:"questions"`
	CreatedAt   time.Time          `json:"created_at"`
	UpdatedAt   time.Time          `json:"updated_at"`
}

type HomeworkResponseWithMeta struct {
	Homeworks []HomeworkResponse `json:"homeworks"`
	Meta      pagination.Meta    `json:"meta"`
}
