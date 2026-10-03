package httperr

import (
	"errors"

	"teaching_assistant/internal/domain/class"
	"teaching_assistant/internal/domain/homework"
	homeworksubmission "teaching_assistant/internal/domain/homework_submission"
	"teaching_assistant/internal/domain/question"
	questionset "teaching_assistant/internal/domain/question_set"
	"teaching_assistant/internal/domain/student"
	"teaching_assistant/internal/domain/user"
	"teaching_assistant/pkg/common"
	"teaching_assistant/pkg/response"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

type mapped struct {
	status  int
	message string
	code    string
}

func Map(err error) (int, string, string) {
	if err == nil {
		return fiber.StatusOK, "", ""
	}

	if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
		return fiber.StatusUnauthorized, "Invalid email or password", "INVALID_CREDENTIALS"
	}

	if errors.Is(err, mongo.ErrNoDocuments) {
		return fiber.StatusNotFound, "Resource not found", "NOT_FOUND"
	}

	if mongo.IsDuplicateKeyError(err) {
		return fiber.StatusConflict, "Resource already exists", "ALREADY_EXISTS"
	}

	if m, ok := lookup(err); ok {
		return m.status, m.message, m.code
	}

	return fiber.StatusInternalServerError, "Internal server error", "INTERNAL_SERVER_ERROR"
}

func lookup(err error) (mapped, bool) {
	var ue user.Error
	if errors.As(err, &ue) {
		return userMapped(ue)
	}
	var ce class.Error
	if errors.As(err, &ce) {
		return classMapped(ce)
	}
	var se student.Error
	if errors.As(err, &se) {
		return studentMapped(se)
	}
	var qe question.Error
	if errors.As(err, &qe) {
		return questionMapped(qe)
	}
	var qse questionset.Error
	if errors.As(err, &qse) {
		return questionSetMapped(qse)
	}
	var he homework.Error
	if errors.As(err, &he) {
		return homeworkMapped(he)
	}
	var hse homeworksubmission.Error
	if errors.As(err, &hse) {
		return submissionMapped(hse)
	}
	var cme common.Error
	if errors.As(err, &cme) {
		return commonMapped(cme)
	}
	return mapped{}, false
}

func userMapped(err user.Error) (mapped, bool) {
	switch err {
	case user.ErrInvalidUsername:
		return mapped{fiber.StatusBadRequest, "Username is required", "INVALID_USERNAME"}, true
	case user.ErrInvalidEmail:
		return mapped{fiber.StatusBadRequest, "Email is required or invalid", "INVALID_EMAIL"}, true
	case user.ErrInvalidPassword:
		return mapped{fiber.StatusBadRequest, "Password is required", "INVALID_PASSWORD"}, true
	case user.ErrInvalidRole:
		return mapped{fiber.StatusBadRequest, "Role must be teacher or parent", "INVALID_ROLE"}, true
	case user.ErrInvalidCredentials:
		return mapped{fiber.StatusUnauthorized, "Invalid email or password", "INVALID_CREDENTIALS"}, true
	case user.ErrUserNotFound:
		return mapped{fiber.StatusNotFound, "User not found", "USER_NOT_FOUND"}, true
	case user.ErrEmailAlreadyExists:
		return mapped{fiber.StatusConflict, "Email is already in use", "EMAIL_ALREADY_EXISTS"}, true
	case user.ErrUnauthorized:
		return mapped{fiber.StatusUnauthorized, "Unauthorized", "UNAUTHORIZED"}, true
	default:
		return mapped{}, false
	}
}

func classMapped(err class.Error) (mapped, bool) {
	switch err {
	case class.ErrInvalidClass:
		return mapped{fiber.StatusBadRequest, "Class name is required", "INVALID_CLASS"}, true
	case class.ErrImageTooLarge:
		return mapped{fiber.StatusBadRequest, "Class image must be 5MB or smaller", "INVALID_IMAGE"}, true
	case class.ErrClassNotFound:
		return mapped{fiber.StatusNotFound, "Class not found", "CLASS_NOT_FOUND"}, true
	case class.ErrUnauthorized:
		return mapped{fiber.StatusForbidden, "You cannot manage this class", "CLASS_FORBIDDEN"}, true
	case class.ErrClassInUse:
		return mapped{fiber.StatusConflict, "Cannot delete class because it still has homework", "CLASS_IN_USE"}, true
	case class.ErrStudentCodeExists:
		return mapped{fiber.StatusConflict, "Student code already exists", "STUDENT_CODE_EXISTS"}, true
	default:
		return mapped{}, false
	}
}

func studentMapped(err student.Error) (mapped, bool) {
	switch err {
	case student.ErrInvalidStudentCode:
		return mapped{fiber.StatusBadRequest, "Student code is required", "INVALID_STUDENT_CODE"}, true
	case student.ErrInvalidParentID:
		return mapped{fiber.StatusBadRequest, "Parent is required", "INVALID_PARENT_ID"}, true
	case student.ErrNotAParent:
		return mapped{fiber.StatusBadRequest, "The selected user is not a parent", "NOT_A_PARENT"}, true
	case student.ErrStudentNotFound:
		return mapped{fiber.StatusNotFound, "Student not found", "STUDENT_NOT_FOUND"}, true
	case student.ErrGuardianNotFound:
		return mapped{fiber.StatusNotFound, "No student is linked to this parent", "GUARDIAN_NOT_FOUND"}, true
	case student.ErrParentNotFound:
		return mapped{fiber.StatusNotFound, "Parent not found", "PARENT_NOT_FOUND"}, true
	case student.ErrStudentNotInClass:
		return mapped{fiber.StatusForbidden, "Student does not belong to this class", "STUDENT_NOT_IN_CLASS"}, true
	case student.ErrStudentInactive:
		return mapped{fiber.StatusForbidden, "Student is not active", "STUDENT_INACTIVE"}, true
	case student.ErrGuardianAlreadyExists:
		return mapped{fiber.StatusConflict, "This parent already has a student", "PARENT_ALREADY_HAS_STUDENT"}, true
	default:
		return mapped{}, false
	}
}

func questionMapped(err question.Error) (mapped, bool) {
	switch err {
	case question.ErrRequiredQuestionType:
		return mapped{fiber.StatusBadRequest, "Question type is required", "REQUIRED_QUESTION_TYPE"}, true
	case question.ErrRequiredGrade:
		return mapped{fiber.StatusBadRequest, "Grade is required", "REQUIRED_GRADE"}, true
	case question.ErrRequiredSubject:
		return mapped{fiber.StatusBadRequest, "Subject is required", "REQUIRED_SUBJECT"}, true
	case question.ErrRequiredDifficulty:
		return mapped{fiber.StatusBadRequest, "Difficulty is required", "REQUIRED_DIFFICULTY"}, true
	case question.ErrRequiredQuestion:
		return mapped{fiber.StatusBadRequest, "Question content is required", "REQUIRED_QUESTION"}, true
	case question.ErrRequiredCorrectAnswer:
		return mapped{fiber.StatusBadRequest, "Correct answer is required", "REQUIRED_CORRECT_ANSWER"}, true
	case question.ErrInvalidType:
		return mapped{fiber.StatusBadRequest, "Question type must be multiple_choice or true_false", "INVALID_QUESTION_TYPE"}, true
	case question.ErrInvalidSubject:
		return mapped{fiber.StatusBadRequest, "Subject is invalid", "INVALID_SUBJECT"}, true
	case question.ErrInvalidGrade:
		return mapped{fiber.StatusBadRequest, "Grade must be 1 to 5", "INVALID_GRADE"}, true
	case question.ErrInvalidDifficulty:
		return mapped{fiber.StatusBadRequest, "Difficulty must be easy, medium, or hard", "INVALID_DIFFICULTY"}, true
	case question.ErrInvalidCorrectIndex:
		return mapped{fiber.StatusBadRequest, "Correct index must match an option", "INVALID_CORRECT_INDEX"}, true
	case question.ErrInvalidQuantity:
		return mapped{fiber.StatusBadRequest, "Quantity must be between 1 and 20", "INVALID_QUANTITY"}, true
	case question.ErrQuestionNotFound:
		return mapped{fiber.StatusNotFound, "Question not found", "QUESTION_NOT_FOUND"}, true
	case question.ErrUnauthorized:
		return mapped{fiber.StatusForbidden, "You cannot manage this question", "QUESTION_FORBIDDEN"}, true
	case question.ErrQuestionInUse:
		return mapped{fiber.StatusConflict, "Question is used in a set, homework, or submission", "QUESTION_IN_USE"}, true
	case question.ErrLLMFailed:
		return mapped{fiber.StatusBadGateway, "Failed to generate questions", "LLM_FAILED"}, true
	case question.ErrLLMInvalidResponse:
		return mapped{fiber.StatusBadGateway, "Generated questions are invalid, please retry", "LLM_INVALID_RESPONSE"}, true
	default:
		return mapped{}, false
	}
}

func questionSetMapped(err questionset.Error) (mapped, bool) {
	switch err {
	case questionset.ErrInvalidTitle:
		return mapped{fiber.StatusBadRequest, "Question set title is required", "INVALID_TITLE"}, true
	case questionset.ErrInvalidQuestionType:
		return mapped{fiber.StatusBadRequest, "Question set type must be multiple_choice or true_false", "INVALID_QUESTION_TYPE"}, true
	case questionset.ErrInvalidQuestions:
		return mapped{fiber.StatusBadRequest, "Question set must have at least one valid question", "INVALID_QUESTIONS"}, true
	case questionset.ErrInvalidQuestionTypeForQuestion:
		return mapped{fiber.StatusBadRequest, "All questions must match the set type", "QUESTION_TYPE_MISMATCH"}, true
	case questionset.ErrQuestionSetNotFound:
		return mapped{fiber.StatusNotFound, "Question set not found", "QUESTION_SET_NOT_FOUND"}, true
	case questionset.ErrUnauthorized:
		return mapped{fiber.StatusForbidden, "You cannot manage this question set", "QUESTION_SET_FORBIDDEN"}, true
	default:
		return mapped{}, false
	}
}

func homeworkMapped(err homework.Error) (mapped, bool) {
	switch err {
	case homework.ErrInvalidTitle:
		return mapped{fiber.StatusBadRequest, "Homework title is required", "INVALID_TITLE"}, true
	case homework.ErrInvalidClassID:
		return mapped{fiber.StatusBadRequest, "Class is invalid", "INVALID_CLASS_ID"}, true
	case homework.ErrInvalidQuestions:
		return mapped{fiber.StatusBadRequest, "Homework must have at least one question", "INVALID_QUESTIONS"}, true
	case homework.ErrInvalidDueDate:
		return mapped{fiber.StatusBadRequest, "Due date is invalid (use YYYY-MM-DD UTC)", "INVALID_DUE_DATE"}, true
	case homework.ErrHomeworkNotFound:
		return mapped{fiber.StatusNotFound, "Homework not found", "HOMEWORK_NOT_FOUND"}, true
	case homework.ErrHomeworkNotAuthorized:
		return mapped{fiber.StatusForbidden, "You cannot manage this homework", "HOMEWORK_FORBIDDEN"}, true
	case homework.ErrHomeworkInUse:
		return mapped{fiber.StatusConflict, "Cannot change or delete homework that already has submissions", "HOMEWORK_IN_USE"}, true
	default:
		return mapped{}, false
	}
}

func submissionMapped(err homeworksubmission.Error) (mapped, bool) {
	switch err {
	case homeworksubmission.ErrInvalidHomeworkSubmission:
		return mapped{fiber.StatusBadRequest, "Homework and answers are required", "INVALID_SUBMISSION"}, true
	case homeworksubmission.ErrQuestionMismatch:
		return mapped{fiber.StatusBadRequest, "Answers must match every homework question", "QUESTION_MISMATCH"}, true
	case homeworksubmission.ErrInvalidStudentAnswer:
		return mapped{fiber.StatusBadRequest, "Answer does not match the question type", "INVALID_STUDENT_ANSWER"}, true
	case homeworksubmission.ErrHomeworkSubmissionDueDateExpired:
		return mapped{fiber.StatusBadRequest, "Homework due date has passed", "DUE_DATE_EXPIRED"}, true
	case homeworksubmission.ErrHomeworkSubmissionNotFound:
		return mapped{fiber.StatusNotFound, "Submission not found", "SUBMISSION_NOT_FOUND"}, true
	case homeworksubmission.ErrStudentNotInClass:
		return mapped{fiber.StatusForbidden, "Student is not in the homework class", "STUDENT_NOT_IN_CLASS"}, true
	case homeworksubmission.ErrHomeworkSubmissionAlreadyExists:
		return mapped{fiber.StatusConflict, "This student has already submitted this homework", "SUBMISSION_ALREADY_EXISTS"}, true
	default:
		return mapped{}, false
	}
}

func commonMapped(err common.Error) (mapped, bool) {
	switch err {
	case common.ErrUnauthorized:
		return mapped{fiber.StatusUnauthorized, "Unauthorized", "UNAUTHORIZED"}, true
	case common.ErrBadRequest:
		return mapped{fiber.StatusBadRequest, "Invalid request", "BAD_REQUEST"}, true
	case common.ErrNotFound:
		return mapped{fiber.StatusNotFound, "Not found", "NOT_FOUND"}, true
	case common.ErrAlreadyExists:
		return mapped{fiber.StatusConflict, "Already exists", "ALREADY_EXISTS"}, true
	case common.ErrInternalServerError:
		return mapped{fiber.StatusInternalServerError, "Internal server error", "INTERNAL_SERVER_ERROR"}, true
	default:
		return mapped{}, false
	}
}

func Fail(c *fiber.Ctx, err error) error {
	status, message, code := Map(err)
	return response.Fail(c, status, message, code)
}
