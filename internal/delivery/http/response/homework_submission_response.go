package response

import (
	"teaching_assistant/pkg/pagination"
	"time"
)

type HomeworkSubmissionResponseWithMeta struct {
	Meta                pagination.Meta               `json:"meta"`
	HomeworkSubmissions []*HomeworkSubmissionResponse `json:"homework_submissions"`
}

type HomeworkSubmissionResponse struct {
	ID             string          `json:"id"`
	HomeworkID     string          `json:"homework_id"`
	StudentID      string          `json:"student_id"`
	IsSubmitted    bool            `json:"is_submitted"`
	StudentAnswers []StudentAnswer `json:"student_answers"`
	TotalScore     float64         `json:"total_score"`
	MaxScore       float64         `json:"max_score"`
	SubmittedBy    string          `json:"submitted_by"`
	SubmittedAt    time.Time       `json:"submitted_at"`
	CreatedAt      time.Time       `json:"created_at"`
	UpdatedAt      time.Time       `json:"updated_at"`
}

type StudentAnswer struct {
	Question      QuestionResponse `json:"question"`
	SelectedIndex *int             `json:"selected_index"`
	SelectedBool  *bool            `json:"selected_bool"`
	IsCorrect     bool             `json:"is_correct"`
	Score         float64          `json:"score"`
	MaxScore      float64          `json:"max_score"`
}
