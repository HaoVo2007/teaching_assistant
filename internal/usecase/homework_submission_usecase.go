package usecase

import (
	"context"
	"time"

	"teaching_assistant/internal/delivery/http/mapper"
	"teaching_assistant/internal/delivery/http/request"
	"teaching_assistant/internal/delivery/http/response"
	"teaching_assistant/internal/domain/homework"
	homeworksubmission "teaching_assistant/internal/domain/homework_submission"
	"teaching_assistant/internal/domain/question"
	"teaching_assistant/internal/domain/student"
	"teaching_assistant/pkg/pagination"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type homeworkSubmissionUsecase struct {
	homeworkSubmissionRepository homeworksubmission.HomeworkSubmissionRepository
	homeworkRepository           homework.HomeworkRepository
	questionRepository           question.QuestionRepository
	studentRepository            student.StudentRepository
}

func NewHomeworkSubmissionUsecase(
	homeworkSubmissionRepository homeworksubmission.HomeworkSubmissionRepository,
	homeworkRepository homework.HomeworkRepository,
	questionRepository question.QuestionRepository,
	studentRepository student.StudentRepository,
) homeworksubmission.HomeworkSubmissionService {
	return &homeworkSubmissionUsecase{
		homeworkSubmissionRepository: homeworkSubmissionRepository,
		homeworkRepository:           homeworkRepository,
		questionRepository:           questionRepository,
		studentRepository:            studentRepository,
	}
}

func (s *homeworkSubmissionUsecase) CreateHomeworkSubmission(ctx context.Context, req request.CreateHomeworkSubmissionRequest, userId string) error {
	if req.HomeworkID == "" {
		return homeworksubmission.ErrInvalidHomeworkSubmission
	}

	objectId, err := primitive.ObjectIDFromHex(req.HomeworkID)
	if err != nil {
		return homeworksubmission.ErrInvalidHomeworkSubmission
	}

	guardian, err := s.studentRepository.GetGuardianByParentId(ctx, userId)
	if err != nil {
		return student.ErrGuardianNotFound
	}

	if guardian == nil {
		return student.ErrGuardianNotFound
	}

	studentId, err := primitive.ObjectIDFromHex(guardian.StudentID)
	if err != nil {
		return student.ErrStudentNotFound
	}

	studentRes, err := s.studentRepository.GetStudentById(ctx, studentId)
	if err != nil {
		return student.ErrStudentNotFound
	}

	if studentRes == nil {
		return student.ErrStudentNotFound
	}

	if len(req.StudentAnswers) == 0 {
		return homeworksubmission.ErrInvalidHomeworkSubmission
	}

	hw, err := s.homeworkRepository.GetHomeworkById(ctx, objectId)
	if err != nil {
		return homework.ErrHomeworkNotFound
	}

	if hw == nil {
		return homework.ErrHomeworkNotFound
	}

	if studentRes.ClassID != hw.ClassID {
		return homeworksubmission.ErrStudentNotInClass
	}

	dueDateUtc := hw.DueDate.UTC()
	nowUtc := time.Now().UTC()
	if nowUtc.After(dueDateUtc) {
		return homeworksubmission.ErrHomeworkSubmissionDueDateExpired
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
		StudentID:      studentId.Hex(),
		IsSubmitted:    true,
		StudentAnswers: studentAnswers,
		TeacherID:      hw.CreatedBy,
		SubmittedBy:    userId,
		SubmittedAt:    time.Now(),
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	return s.homeworkSubmissionRepository.CreateHomeworkSubmission(ctx, submission)
}

func (s *homeworkSubmissionUsecase) GetHomeworkSubmissions(ctx context.Context, params pagination.Params, userId string) (*response.HomeworkSubmissionResponseWithMeta, error) {
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

func (s *homeworkSubmissionUsecase) GetHomeworkSubmissionById(ctx context.Context, id string, userId string) (*response.HomeworkSubmissionResponse, error) {
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, homeworksubmission.ErrHomeworkSubmissionNotFound
	}

	submission, err := s.homeworkSubmissionRepository.GetHomeworkSubmissionById(ctx, objectId)
	if err != nil {
		return nil, err
	}

	if submission == nil {
		return nil, homeworksubmission.ErrHomeworkSubmissionNotFound
	}

	if submission.TeacherID != userId {
		return nil, homeworksubmission.ErrHomeworkSubmissionNotFound
	}

	homeworkId, err := primitive.ObjectIDFromHex(submission.HomeworkID)
	if err != nil {
		return nil, homeworksubmission.ErrHomeworkSubmissionNotFound
	}

	homework, err := s.homeworkRepository.GetHomeworkById(ctx, homeworkId)
	if err != nil {
		return nil, err
	}

	if homework == nil {
		return nil, homeworksubmission.ErrHomeworkSubmissionNotFound
	}

	questionsByID, err := s.loadHomeworkQuestions(ctx, homework.Questions)
	if err != nil {
		return nil, err
	}

	response := mapper.MapHomeworkSubmissionToResponse(submission, homework, questionsByID)

	return response, nil
}

func (s *homeworkSubmissionUsecase) GetHomeworkSubmissionsByHomeworkId(ctx context.Context, homeworkId string, userId string, params pagination.Params) (*response.HomeworkSubmissionResponseWithMeta, error) {
	objectId, err := primitive.ObjectIDFromHex(homeworkId)
	if err != nil {
		return nil, homeworksubmission.ErrHomeworkSubmissionNotFound
	}

	homework, err := s.homeworkRepository.GetHomeworkById(ctx, objectId)
	if err != nil {
		return nil, err
	}

	if homework == nil {
		return nil, homeworksubmission.ErrHomeworkSubmissionNotFound
	}

	if homework.CreatedBy != userId {
		return nil, homeworksubmission.ErrHomeworkSubmissionNotFound
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

func (s *homeworkSubmissionUsecase) GetHomeworkSubmissionsByHomeworkIdByGuardian(ctx context.Context, homeworkId string, userId string) (*response.HomeworkSubmissionResponse, error) {
	guardian, err := s.studentRepository.GetGuardianByParentId(ctx, userId)
	if err != nil {
		return nil, err
	}

	if guardian == nil {
		return nil, student.ErrGuardianNotFound
	}

	studentId, err := primitive.ObjectIDFromHex(guardian.StudentID)
	if err != nil {
		return nil, student.ErrStudentNotFound
	}

	studentRes, err := s.studentRepository.GetStudentById(ctx, studentId)
	if err != nil {
		return nil, student.ErrStudentNotFound
	}

	if studentRes == nil {
		return nil, student.ErrStudentNotFound
	}

	submission, err := s.homeworkSubmissionRepository.GetHomeworkSubmissionsOfStudentByHomeworkId(ctx, homeworkId, studentRes.ID.Hex())
	if err != nil {
		return nil, err
	}

	if submission != nil {
		if submission.SubmittedBy != userId {
			if submission.SubmittedBy != "" {
				return nil, homeworksubmission.ErrHomeworkSubmissionNotFound
			}
		}
	}

	homeworkIdObjectID, err := primitive.ObjectIDFromHex(homeworkId)
	if err != nil {
		return nil, homeworksubmission.ErrHomeworkSubmissionNotFound
	}

	homework, err := s.homeworkRepository.GetHomeworkById(ctx, homeworkIdObjectID)
	if err != nil {
		return nil, err
	}

	if homework == nil {
		return nil, homeworksubmission.ErrHomeworkSubmissionNotFound
	}

	questionsByID, err := s.loadHomeworkQuestions(ctx, homework.Questions)
	if err != nil {
		return nil, err
	}

	response := mapper.MapHomeworkSubmissionToResponse(submission, homework, questionsByID)

	return response, nil
}

func (s *homeworkSubmissionUsecase) loadSubmissionRelations(
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

func (s *homeworkSubmissionUsecase) loadHomeworkQuestions(ctx context.Context, questionIDs []string) (map[string]*question.Question, error) {
	oids := make([]primitive.ObjectID, 0, len(questionIDs))
	for _, id := range questionIDs {
		oid, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			return nil, homeworksubmission.ErrQuestionMismatch
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
		return homeworksubmission.ErrQuestionMismatch
	}

	homeworkSet := make(map[string]struct{}, len(homeworkQuestionIDs))
	for _, id := range homeworkQuestionIDs {
		homeworkSet[id] = struct{}{}
	}

	seen := make(map[string]struct{}, len(answers))
	for _, answer := range answers {
		if answer.QuestionID == "" {
			return homeworksubmission.ErrQuestionMismatch
		}
		if _, ok := homeworkSet[answer.QuestionID]; !ok {
			return homeworksubmission.ErrQuestionMismatch
		}
		if _, dup := seen[answer.QuestionID]; dup {
			return homeworksubmission.ErrQuestionMismatch
		}
		seen[answer.QuestionID] = struct{}{}

		q, ok := questionsByID[answer.QuestionID]
		if !ok {
			return homeworksubmission.ErrQuestionMismatch
		}

		if err := validateAnswerByType(q.Type, answer); err != nil {
			return err
		}
	}

	if len(seen) != len(homeworkSet) {
		return homeworksubmission.ErrQuestionMismatch
	}

	return nil
}

func validateAnswerByType(questionType string, answer request.StudentAnswer) error {
	switch question.QuestionType(questionType) {
	case question.QuestionTypeMultipleChoice:
		if answer.SelectedIndex == nil {
			return homeworksubmission.ErrInvalidStudentAnswer
		}
	case question.QuestionTypeTrueFalse:
		if answer.SelectedBool == nil {
			return homeworksubmission.ErrInvalidStudentAnswer
		}
	default:
		return homeworksubmission.ErrInvalidStudentAnswer
	}
	return nil
}
