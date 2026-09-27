package request

type CreateQuestionRequest struct {
	Type         string   `form:"type"`
	Subject      string   `form:"subject"`
	Grade        string   `form:"grade"`
	Difficulty   string   `form:"difficulty"`
	Question     string   `form:"question"`
	Options      []string `form:"options"`
	CorrectIndex *int     `form:"correct_index"`
	CorrectBool  *bool    `form:"correct_bool"`
	Explanation  string   `form:"explanation"`
}

type UpdateQuestionRequest struct {
	Type         string   `form:"type"`
	Subject      string   `form:"subject"`
	Grade        string   `form:"grade"`
	Difficulty   string   `form:"difficulty"`
	Question     string   `form:"question"`
	Options      []string `form:"options"`
	CorrectIndex *int     `form:"correct_index"`
	CorrectBool  *bool    `form:"correct_bool"`
	Explanation  string   `form:"explanation"`
}
