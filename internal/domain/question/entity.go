package question

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type QuestionType string

const (
	QuestionTypeMultipleChoice QuestionType = "multiple_choice"
	QuestionTypeTrueFalse      QuestionType = "true_false"
)

func IsSupportedType(t string) bool {
	return QuestionType(t) == QuestionTypeMultipleChoice || QuestionType(t) == QuestionTypeTrueFalse
}

func IsSupportedSubject(s string) bool {
	switch Subject(s) {
	case SubjectVietnamese, SubjectMathematics, SubjectEthics, SubjectEnglish,
		SubjectNatureAndSociety, SubjectHistoryAndGeography, SubjectScience,
		SubjectInformatics, SubjectTechnology, SubjectPhysicalEducation,
		SubjectMusic, SubjectArt, SubjectExperientialActivities:
		return true
	default:
		return false
	}
}

func IsSupportedGrade(g string) bool {
	switch Grade(g) {
	case Grade1, Grade2, Grade3, Grade4, Grade5:
		return true
	default:
		return false
	}
}

func IsSupportedDifficulty(d string) bool {
	switch Difficulty(d) {
	case DifficultyEasy, DifficultyMedium, DifficultyHard:
		return true
	default:
		return false
	}
}

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

type Difficulty string

const (
	DifficultyEasy   Difficulty = "easy"
	DifficultyMedium Difficulty = "medium"
	DifficultyHard   Difficulty = "hard"
)

type Question struct {
	ID           primitive.ObjectID `bson:"_id"`
	Type         string             `bson:"type"`
	Subject      string             `bson:"subject"`
	Grade        string             `bson:"grade"`
	Difficulty   string             `bson:"difficulty"`
	Question     string             `bson:"question"`
	Options      []string           `bson:"options,omitempty"`
	CorrectIndex *int               `bson:"correct_index,omitempty"`
	CorrectBool  *bool              `bson:"correct_bool,omitempty"`
	Explanation  string             `bson:"explanation,omitempty"`
	CreatedBy    string             `bson:"created_by"`
	CreatedAt    time.Time          `bson:"created_at"`
	UpdatedAt    time.Time          `bson:"updated_at"`
}
