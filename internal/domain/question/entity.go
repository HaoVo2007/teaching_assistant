package question

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type QuestionType string

const (
	QuestionTypeMultipleChoice QuestionType = "multiple_choice"
	QuestionTypeTrueFalse      QuestionType = "true_false"
	QuestionTypeMatching       QuestionType = "matching"
)

type KindPair string

const (
	Text  KindPair = "text"
	Image KindPair = "image"
)

type Subject string

const (
	SubjectVietnamese             Subject = "vietnamese"
	SubjectMathematics            Subject = "mathematics"
	SubjectEthics                 Subject = "ethics"
	SubjectEnglish                Subject = "english"
	SubjectNatureAndSociety       Subject = "nature_and_society"
	SubjectHistoryAndGeography    Subject = "history_and_geography"
	SubjectScience                Subject = "science"
	SubjectInformatics            Subject = "informatics"
	SubjectTechnology             Subject = "technology"
	SubjectPhysicalEducation      Subject = "physical_education"
	SubjectMusic                  Subject = "music"
	SubjectArt                    Subject = "art"
	SubjectExperientialActivities Subject = "experiential_activities"
)

type Grade string

const (
	Grade1 Grade = "1"
	Grade2 Grade = "2"
	Grade3 Grade = "3"
	Grade4 Grade = "4"
	Grade5 Grade = "5"
)

type Question struct {
	ID           primitive.ObjectID `bson:"_id"`
	Type         string             `bson:"type"` // multiple_choice, true_false, matching
	Subject      string             `bson:"subject"`
	Grade        string             `bson:"grade"`
	Difficulty   string             `bson:"difficulty"`
	Question     string             `bson:"question"`
	Options      []string           `bson:"options,omitempty"`
	CorrectIndex *int               `bson:"correct_index,omitempty"`
	CorrectBool  *bool              `bson:"correct_bool,omitempty"`
	Pairs        []Pair             `bson:"pairs,omitempty"`
	Explanation  string             `bson:"explanation,omitempty"`
	CreatedBy    string             `bson:"created_by"`
	CreatedAt    time.Time          `bson:"created_at"`
	UpdatedAt    time.Time          `bson:"updated_at"`
}

type Pair struct {
	Left          string `bson:"left"`
	LeftPublicID  string `bson:"left_public_id"`
	LeftKind      string `bson:"left_kind"`
	Right         string `bson:"right"`
	RightPublicID string `bson:"right_public_id"`
	RightKind     string `bson:"right_kind"`
}
