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
	case user.ErrInvalidName, user.ErrInvalidUsername:
		return mapped{fiber.StatusBadRequest, "Username is required", "INVALID_USERNAME"}, true
	case user.ErrInvalidEmail:
		return mapped{fiber.StatusBadRequest, "Email is required or invalid", "INVALID_EMAIL"}, true
	case user.ErrInvalidPassword, user.ErrPasswordTooWeak:
		return mapped{fiber.StatusBadRequest, "Password is required", "INVALID_PASSWORD"}, true
	case user.ErrInvalidRole:
		return mapped{fiber.StatusBadRequest, "Role must be teacher or parent", "INVALID_ROLE"}, true
	case user.ErrPasswordMismatch:
		return mapped{fiber.StatusBadRequest, "Password confirmation does not match", "PASSWORD_MISMATCH"}, true
	case user.ErrInvalidCredentials, user.ErrWrongPassword:
		return mapped{fiber.StatusUnauthorized, "Invalid email or password", "INVALID_CREDENTIALS"}, true
	case user.ErrUserNotFound:
		return mapped{fiber.StatusNotFound, "User not found", "USER_NOT_FOUND"}, true
	case user.ErrUserAlreadyExists:
		return mapped{fiber.StatusConflict, "User already exists", "USER_ALREADY_EXISTS"}, true
	case user.ErrEmailAlreadyExists:
		return mapped{fiber.StatusConflict, "Email is already in use", "EMAIL_ALREADY_EXISTS"}, true
	case user.ErrUnauthorized:
		return mapped{fiber.StatusUnauthorized, "Unauthorized", "UNAUTHORIZED"}, true
	case user.ErrForbidden, user.ErrInsufficientPermission:
		return mapped{fiber.StatusForbidden, "You do not have permission", "FORBIDDEN"}, true
	default:
		return mapped{}, false
	}
}

func classMapped(err class.Error) (mapped, bool) {
	switch err {
	case class.ErrInvalidClass:
		return mapped{fiber.StatusBadRequest, "Class name is required", "INVALID_CLASS"}, true
	case class.ErrImageTooLarge:
		return mapped{fiber.StatusBadRequest, "Class image is too large", "INVALID_IMAGE"}, true
	case class.ErrClassNotFound:
		return mapped{fiber.StatusNotFound, "Class not found", "CLASS_NOT_FOUND"}, true
	case class.ErrClassAlreadyExists:
		return mapped{fiber.StatusConflict, "Class already exists", "CLASS_ALREADY_EXISTS"}, true
	case class.ErrClassNotAuthorized, class.ErrUnauthorized:
		return mapped{fiber.StatusForbidden, "You cannot manage this class", "CLASS_FORBIDDEN"}, true
	case class.ErrClassInUse:
		return mapped{fiber.StatusConflict, "Cannot delete class: a claimed student already has a submission", "CLASS_IN_USE"}, true
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
	case student.ErrStudentNotFound:
		return mapped{fiber.StatusNotFound, "Student not found", "STUDENT_NOT_FOUND"}, true
	case student.ErrGuardianNotFound:
		return mapped{fiber.StatusNotFound, "No student is linked to this parent", "GUARDIAN_NOT_FOUND"}, true
	case student.ErrGuardianAlreadyExists:
		return mapped{fiber.StatusConflict, "This parent already has a student", "PARENT_ALREADY_HAS_STUDENT"}, true
	case student.ErrStudentNotInClass:
		return mapped{fiber.StatusForbidden, "Student does not belong to this class", "STUDENT_NOT_IN_CLASS"}, true
	default:
		return mapped{}, false
	}
}

func questionMapped(err question.Error) (mapped, bool) {
	switch err {
	case question.ErrInvalidType:
		return mapped{fiber.StatusBadRequest, "Question type is invalid", "INVALID_QUESTION_TYPE"}, true
	case question.ErrInvalidQuestion:
		return mapped{fiber.StatusBadRequest, "Question content is invalid", "INVALID_QUESTION"}, true
	case question.ErrInvalidOptions:
		return mapped{fiber.StatusBadRequest, "Question options are invalid", "INVALID_OPTIONS"}, true
	case question.ErrInvalidCorrect:
		return mapped{fiber.StatusBadRequest, "Correct answer is invalid", "INVALID_CORRECT_ANSWER"}, true
	case question.ErrInvalidPairs:
		return mapped{fiber.StatusBadRequest, "Matching pairs are invalid", "INVALID_PAIRS"}, true
	case question.ErrImageTooLarge:
		return mapped{fiber.StatusBadRequest, "Image is too large", "INVALID_IMAGE"}, true
	case question.ErrInvalidSubject:
		return mapped{fiber.StatusBadRequest, "Subject is invalid", "INVALID_SUBJECT"}, true
	case question.ErrInvalidGrade:
		return mapped{fiber.StatusBadRequest, "Grade is invalid", "INVALID_GRADE"}, true
	case question.ErrQuestionNotFound:
		return mapped{fiber.StatusNotFound, "Question not found", "QUESTION_NOT_FOUND"}, true
	case question.ErrUnauthorized:
		return mapped{fiber.StatusForbidden, "You cannot manage this question", "QUESTION_FORBIDDEN"}, true
	case question.ErrQuestionInUse:
		return mapped{fiber.StatusConflict, "Question is used in a set, homework, or submission", "QUESTION_IN_USE"}, true
	default:
		return mapped{}, false
	}
}

func questionSetMapped(err questionset.Error) (mapped, bool) {
	switch err {
	case questionset.ErrInvalidTitle:
		return mapped{fiber.StatusBadRequest, "Question set title is required", "INVALID_TITLE"}, true
	case questionset.ErrInvalidQuestionType:
		return mapped{fiber.StatusBadRequest, "Question set type is invalid", "INVALID_QUESTION_TYPE"}, true
	case questionset.ErrInvalidQuestions:
		return mapped{fiber.StatusBadRequest, "Question list is invalid", "INVALID_QUESTIONS"}, true
	case questionset.ErrInvalidQuestionTypeForQuestion:
		return mapped{fiber.StatusBadRequest, "All questions must match the set type", "QUESTION_TYPE_MISMATCH"}, true
	case questionset.ErrQuestionSetNotFound:
		return mapped{fiber.StatusNotFound, "Question set not found", "QUESTION_SET_NOT_FOUND"}, true
	case questionset.ErrQuestionSetAlreadyExists:
		return mapped{fiber.StatusConflict, "Question set already exists", "QUESTION_SET_ALREADY_EXISTS"}, true
	case questionset.ErrQuestionSetNotAuthorized, questionset.ErrUnauthorized:
		return mapped{fiber.StatusForbidden, "You cannot manage this question set", "QUESTION_SET_FORBIDDEN"}, true
	default:
		return mapped{}, false
	}
}

func homeworkMapped(err homework.Error) (mapped, bool) {
	switch err {
	case homework.ErrInvalidHomework:
		return mapped{fiber.StatusBadRequest, "Homework is invalid", "INVALID_HOMEWORK"}, true
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
	case homework.ErrHomeworkAlreadyExists:
		return mapped{fiber.StatusConflict, "Homework already exists", "HOMEWORK_ALREADY_EXISTS"}, true
	case homework.ErrHomeworkNotAuthorized, homework.ErrUnauthorized:
		return mapped{fiber.StatusForbidden, "You cannot manage this homework", "HOMEWORK_FORBIDDEN"}, true
	case homework.ErrHomeworkInUse:
		return mapped{fiber.StatusConflict, "Cannot change or delete homework that already has submissions", "HOMEWORK_IN_USE"}, true
	default:
		return mapped{}, false
	}
}

func submissionMapped(err homeworksubmission.Error) (mapped, bool) {
	switch err {
	case homeworksubmission.ErrInvalidHomeworkSubmission, homeworksubmission.ErrHomeworkSubmissionInvalid:
		return mapped{fiber.StatusBadRequest, "Submission is invalid", "INVALID_SUBMISSION"}, true
	case homeworksubmission.ErrQuestionMismatch:
		return mapped{fiber.StatusBadRequest, "Answers must match every homework question", "QUESTION_MISMATCH"}, true
	case homeworksubmission.ErrInvalidStudentAnswer:
		return mapped{fiber.StatusBadRequest, "Answer does not match the question type", "INVALID_STUDENT_ANSWER"}, true
	case homeworksubmission.ErrHomeworkSubmissionDueDateExpired:
		return mapped{fiber.StatusBadRequest, "Homework due date has passed", "DUE_DATE_EXPIRED"}, true
	case homeworksubmission.ErrHomeworkSubmissionNotFound, homeworksubmission.ErrHomeworkSubmissionNotSubmitted:
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
	default:
		return mapped{}, false
	}
}

func Fail(c *fiber.Ctx, err error) error {
	status, message, code := Map(err)
	return response.Fail(c, status, message, code)
}
