package response

import (
	"teaching_assistant/pkg/pagination"
	"time"
)

type QuestionResponseWithMeta struct {
	Questions []*QuestionResponse `json:"questions"`
	Meta      pagination.Meta     `json:"meta"`
}

type QuestionResponse struct {
	ID           string    `json:"id"`
	Type         string    `json:"type"`
	Subject      string    `json:"subject"`
	Grade        string    `json:"grade"`
	Difficulty   string    `json:"difficulty"`
	Question     string    `json:"question"`
	Options      []string  `json:"options,omitempty"`
	CorrectIndex *int      `json:"correct_index,omitempty"`
	CorrectBool  *bool     `json:"correct_bool,omitempty"`
	Explanation  string    `json:"explanation,omitempty"`
	CreatedBy    string    `json:"created_by"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
