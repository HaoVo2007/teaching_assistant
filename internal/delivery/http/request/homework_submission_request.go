package request

type CreateHomeworkSubmissionRequest struct {
	HomeworkID     string          `json:"homework_id"`
	StudentName    string          `json:"student_name"`
	StudentAnswers []StudentAnswer `json:"student_answers"`
}

type StudentAnswer struct {
	QuestionID    string `json:"question_id"`
	SelectedIndex *int   `json:"selected_index"`
	SelectedBool  *bool  `json:"selected_bool"`
	IsCorrect     bool   `json:"is_correct"`
	Score         int    `json:"score"`
}
