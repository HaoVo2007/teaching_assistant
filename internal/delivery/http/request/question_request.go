package request

import "mime/multipart"

type PairRequest struct {
	Left      string                `form:"left"`
	LeftKind  string                `form:"left_kind"`
	LeftFile  *multipart.FileHeader `json:"-" form:"-"`
	Right     string                `form:"right"`
	RightKind string                `form:"right_kind"`
	RightFile *multipart.FileHeader `json:"-" form:"-"`
}

type CreateQuestionRequest struct {
	Type         string        `form:"type"`
	Subject      string        `form:"subject"`
	Grade        string        `form:"grade"`
	Difficulty   string        `form:"difficulty"`
	Question     string        `form:"question"`
	Options      []string      `form:"options"`
	Pairs        []PairRequest `form:"pairs"`
	CorrectIndex *int          `form:"correct_index"`
	CorrectBool  *bool         `form:"correct_bool"`
	Explanation  string        `form:"explanation"`
}
