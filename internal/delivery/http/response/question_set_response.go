package response

import (
	"teaching_assistant/pkg/pagination"
	"time"
)

type QuestionSetResponse struct {
	ID           string              `json:"id"`
	Title        string              `json:"title"`
	QuestionType string              `json:"question_type"`
	Description  *string             `json:"description"`
	Questions    []*QuestionResponse `json:"questions"`
	CreatedAt    time.Time           `json:"created_at"`
	UpdatedAt    time.Time           `json:"updated_at"`
}

type QuestionSetResponseWithMeta struct {
	QuestionSets []*QuestionSetResponse `json:"question_sets"`
	Meta         pagination.Meta        `json:"meta"`
}
