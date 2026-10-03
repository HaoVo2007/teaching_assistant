package question

type Error string

const (
	ErrInvalidType           Error = "invalid question type"
	ErrInvalidSubject        Error = "invalid subject"
	ErrInvalidGrade          Error = "invalid grade"
	ErrInvalidDifficulty     Error = "invalid difficulty"
	ErrInvalidCorrectIndex   Error = "correct index is invalid"
	ErrInvalidQuantity       Error = "quantity must be between 1 and 20"
	ErrRequiredQuestionType  Error = "question type is required"
	ErrRequiredGrade         Error = "grade is required"
	ErrRequiredSubject       Error = "subject is required"
	ErrRequiredDifficulty    Error = "difficulty is required"
	ErrRequiredQuestion      Error = "question is required"
	ErrRequiredCorrectAnswer Error = "correct answer is required"
	ErrQuestionNotFound      Error = "question not found"
	ErrUnauthorized          Error = "unauthorized"
	ErrQuestionInUse         Error = "question is used in a set, homework, or submission"
	ErrLLMFailed             Error = "failed to generate questions"
	ErrLLMInvalidResponse    Error = "llm returned invalid questions"
)

func (e Error) Error() string {
	return string(e)
}
