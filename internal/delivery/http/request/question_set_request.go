package request

type CreateQuestionSetRequest struct {
	Title       string   `json:"title"`
	Description *string  `json:"description"`
	Questions   []string `json:"questions"`
}
