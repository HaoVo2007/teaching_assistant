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
	ID           string    `bson:"_id" json:"id"`
	Type         string    `bson:"type" json:"type"` // multiple_choice, true_false, matching
	Subject      string    `bson:"subject" json:"subject"`
	Grade        string    `bson:"grade" json:"grade"`
	Difficulty   string    `bson:"difficulty" json:"difficulty"`
	Question     string    `bson:"question" json:"question"`
	Options      []string  `bson:"options,omitempty" json:"options"`
	CorrectIndex *int      `bson:"correct_index,omitempty" json:"correct_index"`
	CorrectBool  *bool     `bson:"correct_bool,omitempty" json:"correct_bool"`
	Pairs        []Pair    `bson:"pairs,omitempty" json:"pairs"`
	Explanation  string    `bson:"explanation,omitempty" json:"explanation"`
	CreatedBy    string    `bson:"created_by" json:"created_by"`
	CreatedAt    time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt    time.Time `bson:"updated_at" json:"updated_at"`
}

type Pair struct {
	Left          string `bson:"left"`
	LeftPublicID  string `bson:"left_public_id"`
	LeftKind      string `bson:"left_kind"`
	Right         string `bson:"right"`
	RightPublicID string `bson:"right_public_id"`
	RightKind     string `bson:"right_kind"`
}
