package request

type CreateHomeworkRequest struct {
	ClassID     string   `json:"class_id"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Questions   []string `json:"questions"`
	DueDate     string   `json:"due_date"`
}

type UpdateHomeworkRequest struct {
	ClassID     string   `json:"class_id"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Questions   []string `json:"questions"`
	DueDate     string   `json:"due_date"`
}
