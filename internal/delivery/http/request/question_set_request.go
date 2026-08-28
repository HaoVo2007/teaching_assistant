package request

type CreateQuestionSetRequest struct {
	Title        string   `json:"title"`
	QuestionType string   `json:"question_type"`
	Description  *string  `json:"description"`
	Questions    []string `json:"questions"`
}

type UpdateQuestionSetRequest struct {
	Title       *string  `json:"title"`
	Description *string  `json:"description"`
	Questions   []string `json:"questions"`
}
