package usecase

import (
	"context"
	"errors"
	"time"

	"teaching_assistant/internal/delivery/http/mapper"
	"teaching_assistant/internal/delivery/http/request"
	"teaching_assistant/internal/delivery/http/response"
	"teaching_assistant/internal/domain/homework"
	homeworksubmission "teaching_assistant/internal/domain/homework_submission"
	"teaching_assistant/internal/domain/question"
	"teaching_assistant/pkg/pagination"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type homeworkSubmissionService struct {
	homeworkSubmissionRepository homeworksubmission.HomeworkSubmissionRepository
	homeworkRepository           homework.HomeworkRepository
	questionRepository           question.QuestionRepository
}

func NewHomeworkSubmissionService(
	homeworkSubmissionRepository homeworksubmission.HomeworkSubmissionRepository,
	homeworkRepository homework.HomeworkRepository,
	questionRepository question.QuestionRepository,
) homeworksubmission.HomeworkSubmissionService {
	return &homeworkSubmissionService{
		homeworkSubmissionRepository: homeworkSubmissionRepository,
		homeworkRepository:           homeworkRepository,
		questionRepository:           questionRepository,
	}
}

func (s *homeworkSubmissionService) CreateHomeworkSubmission(ctx context.Context, req request.CreateHomeworkSubmissionRequest) error {
	if req.HomeworkID == "" {
		return errors.New(string(homeworksubmission.ErrInvalidHomeworkSubmission))
	}

	objectId, err := primitive.ObjectIDFromHex(req.HomeworkID)
	if err != nil {
		return errors.New(string(homeworksubmission.ErrInvalidHomeworkSubmission))
	}

	if req.StudentName == "" {
		return errors.New(string(homeworksubmission.ErrInvalidHomeworkSubmission))
	}

	if len(req.StudentAnswers) == 0 {
		return errors.New(string(homeworksubmission.ErrInvalidHomeworkSubmission))
	}

	hw, err := s.homeworkRepository.GetHomeworkById(ctx, objectId)
	if err != nil {
		return errors.New(string(homeworksubmission.ErrHomeworkSubmissionNotFound))
	}

	if hw == nil {
		return errors.New(string(homeworksubmission.ErrHomeworkSubmissionNotFound))
	}

	questions, err := s.loadHomeworkQuestions(ctx, hw.Questions)
	if err != nil {
		return err
	}

	if err := validateStudentAnswers(hw.Questions, req.StudentAnswers, questions); err != nil {
		return err
	}

	studentAnswers := make([]homeworksubmission.StudentAnswer, 0, len(req.StudentAnswers))
	for _, answer := range req.StudentAnswers {
		studentAnswers = append(studentAnswers, homeworksubmission.StudentAnswer{
			QuestionID:    answer.QuestionID,
			SelectedIndex: answer.SelectedIndex,
			SelectedBool:  answer.SelectedBool,
		})
	}

	submission := &homeworksubmission.HomeworkSubmission{
		ID:             primitive.NewObjectID(),
		HomeworkID:     req.HomeworkID,
		StudentName:    req.StudentName,
		IsSubmitted:    true,
		StudentAnswers: studentAnswers,
		TeacherID:      hw.CreatedBy,
		SubmittedAt:    time.Now(),
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	err = s.homeworkSubmissionRepository.CreateHomeworkSubmission(ctx, submission)
	if err != nil {
		return errors.New(string(homeworksubmission.ErrHomeworkSubmissionInternal))
	}

	return nil
}

func (s *homeworkSubmissionService) GetHomeworkSubmissions(ctx context.Context, params pagination.Params, userId string) (*response.HomeworkSubmissionResponseWithMeta, error) {
	homeworkSubmissions, total, err := s.homeworkSubmissionRepository.GetHomeworkSubmissionsByUserId(ctx, userId, params)
	if err != nil {
		return nil, err
	}

	homeworksByID, questionsByID, err := s.loadSubmissionRelations(ctx, homeworkSubmissions)
	if err != nil {
		return nil, err
	}

	return &response.HomeworkSubmissionResponseWithMeta{
		HomeworkSubmissions: mapper.MapHomeworkSubmissionsToResponses(homeworkSubmissions, homeworksByID, questionsByID),
		Meta:                pagination.NewMeta(params, total),
	}, nil
}

func (s *homeworkSubmissionService) GetHomeworkSubmissionById(ctx context.Context, id string, userId string) (*response.HomeworkSubmissionResponse, error) {
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.New(string(homeworksubmission.ErrHomeworkSubmissionNotFound))
	}

	submission, err := s.homeworkSubmissionRepository.GetHomeworkSubmissionById(ctx, objectId)
	if err != nil {
		return nil, err
	}

	if submission == nil {
		return nil, errors.New(string(homeworksubmission.ErrHomeworkSubmissionNotFound))
	}

	if submission.TeacherID != userId {
		return nil, errors.New(string(homeworksubmission.ErrHomeworkSubmissionNotFound))
	}

	homeworkId, err := primitive.ObjectIDFromHex(submission.HomeworkID)
	if err != nil {
		return nil, errors.New(string(homeworksubmission.ErrHomeworkSubmissionNotFound))
	}

	homework, err := s.homeworkRepository.GetHomeworkById(ctx, homeworkId)
	if err != nil {
		return nil, err
	}

	if homework == nil {
		return nil, errors.New(string(homeworksubmission.ErrHomeworkSubmissionNotFound))
	}

	questionsByID, err := s.loadHomeworkQuestions(ctx, homework.Questions)
	if err != nil {
		return nil, err
	}

	response := mapper.MapHomeworkSubmissionToResponse(submission, homework, questionsByID)

	return response, nil
}

func (s *homeworkSubmissionService) GetHomeworkSubmissionsByHomeworkId(ctx context.Context, homeworkId string, userId string, params pagination.Params) (*response.HomeworkSubmissionResponseWithMeta, error) {
	objectId, err := primitive.ObjectIDFromHex(homeworkId)
	if err != nil {
		return nil, errors.New(string(homeworksubmission.ErrHomeworkSubmissionNotFound))
	}

	homework, err := s.homeworkRepository.GetHomeworkById(ctx, objectId)
	if err != nil {
		return nil, err
	}

	if homework == nil {
		return nil, errors.New(string(homeworksubmission.ErrHomeworkSubmissionNotFound))
	}

	if homework.CreatedBy != userId {
		return nil, errors.New(string(homeworksubmission.ErrHomeworkSubmissionNotFound))
	}

	submissions, total, err := s.homeworkSubmissionRepository.GetHomeworkSubmissionsByHomeworkId(ctx, homeworkId, userId, params)
	if err != nil {
		return nil, err
	}

	homeworksByID, questionsByID, err := s.loadSubmissionRelations(ctx, submissions)
	if err != nil {
		return nil, err
	}

	return &response.HomeworkSubmissionResponseWithMeta{
		HomeworkSubmissions: mapper.MapHomeworkSubmissionsToResponses(submissions, homeworksByID, questionsByID),
		Meta:                pagination.NewMeta(params, total),
	}, nil

}

func (s *homeworkSubmissionService) loadSubmissionRelations(
	ctx context.Context,
	submissions []*homeworksubmission.HomeworkSubmission,
) (map[string]*homework.Homework, map[string]*question.Question, error) {
	homeworkIDSet := make(map[string]struct{})
	questionIDSet := make(map[string]struct{})
	for _, submission := range submissions {
		if submission == nil {
			continue
		}
		if submission.HomeworkID != "" {
			homeworkIDSet[submission.HomeworkID] = struct{}{}
		}
		for _, answer := range submission.StudentAnswers {
			if answer.QuestionID != "" {
				questionIDSet[answer.QuestionID] = struct{}{}
			}
		}
	}

	homeworkOIDs := make([]primitive.ObjectID, 0, len(homeworkIDSet))
	for id := range homeworkIDSet {
		oid, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			continue
		}
		homeworkOIDs = append(homeworkOIDs, oid)
	}

	homeworks, err := s.homeworkRepository.GetHomeworkByIds(ctx, homeworkOIDs)
	if err != nil {
		return nil, nil, err
	}

	homeworksByID := make(map[string]*homework.Homework, len(homeworks))
	for _, hw := range homeworks {
		if hw == nil {
			continue
		}
		homeworksByID[hw.ID.Hex()] = hw
		for _, questionID := range hw.Questions {
			questionIDSet[questionID] = struct{}{}
		}
	}

	questionsByID, err := s.loadHomeworkQuestions(ctx, setKeys(questionIDSet))
	if err != nil {
		return nil, nil, err
	}

	return homeworksByID, questionsByID, nil
}

func setKeys(set map[string]struct{}) []string {
	keys := make([]string, 0, len(set))
	for key := range set {
		keys = append(keys, key)
	}
	return keys
}

func (s *homeworkSubmissionService) loadHomeworkQuestions(ctx context.Context, questionIDs []string) (map[string]*question.Question, error) {
	oids := make([]primitive.ObjectID, 0, len(questionIDs))
	for _, id := range questionIDs {
		oid, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			return nil, errors.New(string(homeworksubmission.ErrQuestionMismatch))
		}
		oids = append(oids, oid)
	}

	items, err := s.questionRepository.GetQuestionByIds(ctx, oids)
	if err != nil {
		return nil, err
	}

	byID := make(map[string]*question.Question, len(items))
	for _, q := range items {
		if q == nil {
			continue
		}
		byID[q.ID.Hex()] = q
	}
	return byID, nil
}

func validateStudentAnswers(
	homeworkQuestionIDs []string,
	answers []request.StudentAnswer,
	questionsByID map[string]*question.Question,
) error {
	if len(answers) != len(homeworkQuestionIDs) {
		return errors.New(string(homeworksubmission.ErrQuestionMismatch))
	}

	homeworkSet := make(map[string]struct{}, len(homeworkQuestionIDs))
	for _, id := range homeworkQuestionIDs {
		homeworkSet[id] = struct{}{}
	}

	seen := make(map[string]struct{}, len(answers))
	for _, answer := range answers {
		if answer.QuestionID == "" {
			return errors.New(string(homeworksubmission.ErrQuestionMismatch))
		}
		if _, ok := homeworkSet[answer.QuestionID]; !ok {
			return errors.New(string(homeworksubmission.ErrQuestionMismatch))
		}
		if _, dup := seen[answer.QuestionID]; dup {
			return errors.New(string(homeworksubmission.ErrQuestionMismatch))
		}
		seen[answer.QuestionID] = struct{}{}

		q, ok := questionsByID[answer.QuestionID]
		if !ok {
			return errors.New(string(homeworksubmission.ErrQuestionMismatch))
		}

		if err := validateAnswerByType(q.Type, answer); err != nil {
			return err
		}
	}

	if len(seen) != len(homeworkSet) {
		return errors.New(string(homeworksubmission.ErrQuestionMismatch))
	}

	return nil
}

func validateAnswerByType(questionType string, answer request.StudentAnswer) error {
	switch question.QuestionType(questionType) {
	case question.QuestionTypeMultipleChoice:
		if answer.SelectedIndex == nil {
			return errors.New(string(homeworksubmission.ErrInvalidStudentAnswer))
		}
	case question.QuestionTypeTrueFalse:
		if answer.SelectedBool == nil {
			return errors.New(string(homeworksubmission.ErrInvalidStudentAnswer))
		}
	default:
		return errors.New(string(homeworksubmission.ErrInvalidStudentAnswer))
	}
	return nil
}
