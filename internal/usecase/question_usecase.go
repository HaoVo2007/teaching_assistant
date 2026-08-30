package usecase

import (
	"context"
	"errors"
	"mime/multipart"
	"time"

	"teaching_assistant/internal/delivery/http/mapper"
	"teaching_assistant/internal/delivery/http/request"
	"teaching_assistant/internal/delivery/http/response"
	"teaching_assistant/internal/domain/homework"
	homeworksubmission "teaching_assistant/internal/domain/homework_submission"
	"teaching_assistant/internal/domain/question"
	questionset "teaching_assistant/internal/domain/question_set"
	"teaching_assistant/pkg/pagination"

	infrastructureCloudinary "teaching_assistant/internal/infrastructure/cloudinary"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type questionService struct {
	questionRepo    question.QuestionRepository
	questionSetRepo questionset.QuestionSetRepository
	homeworkRepo    homework.HomeworkRepository
	submissionRepo  homeworksubmission.HomeworkSubmissionRepository
	cloudinary      *infrastructureCloudinary.CloudinaryUploader
}

func NewQuestionService(
	questionRepo question.QuestionRepository,
	questionSetRepo questionset.QuestionSetRepository,
	homeworkRepo homework.HomeworkRepository,
	submissionRepo homeworksubmission.HomeworkSubmissionRepository,
	cloudinary *infrastructureCloudinary.CloudinaryUploader,
) question.QuestionService {
	return &questionService{
		questionRepo:    questionRepo,
		questionSetRepo: questionSetRepo,
		homeworkRepo:    homeworkRepo,
		submissionRepo:  submissionRepo,
		cloudinary:      cloudinary,
	}
}

func (s *questionService) CreateQuestion(ctx context.Context, req request.CreateQuestionRequest, userId string) (*question.Question, error) {
	if req.Type == "" {
		return nil, errors.New(string(question.ErrInvalidType))
	}

	if req.Subject == "" {
		return nil, errors.New(string(question.ErrInvalidSubject))
	}

	if req.Grade == "" {
		return nil, errors.New(string(question.ErrInvalidGrade))
	}

	pairs := make([]question.Pair, 0, len(req.Pairs))
	for _, p := range req.Pairs {
		pair, _, err := s.buildPair(ctx, p, question.Pair{})
		if err != nil {
			return nil, err
		}
		pairs = append(pairs, pair)
	}

	q := &question.Question{
		ID:           primitive.NewObjectID(),
		Type:         req.Type,
		Subject:      req.Subject,
		Grade:        req.Grade,
		Difficulty:   req.Difficulty,
		Question:     req.Question,
		Options:      req.Options,
		CorrectIndex: req.CorrectIndex,
		CorrectBool:  req.CorrectBool,
		Pairs:        pairs,
		Explanation:  req.Explanation,
		CreatedBy:    userId,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := s.questionRepo.Create(ctx, q); err != nil {
		return nil, err
	}
	return q, nil
}

func (s *questionService) GetQuestions(ctx context.Context, userId string, params pagination.Params, questionType, questionName, subject, grade, difficulty string) (*response.QuestionResponseWithMeta, error) {
	questions, total, err := s.questionRepo.GetQuestions(ctx, userId, params, questionType, questionName, subject, grade, difficulty)
	if err != nil {
		return nil, err
	}

	responses := mapper.MapQuestionsToResponses(questions)
	meta := pagination.NewMeta(params, total)
	return &response.QuestionResponseWithMeta{
		Questions: responses,
		Meta:      meta,
	}, nil
}

func (s *questionService) GetQuestionById(ctx context.Context, id string) (*response.QuestionResponse, error) {
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	question, err := s.questionRepo.GetQuestionById(ctx, objectId)
	if err != nil {
		return nil, err
	}

	return mapper.MapQuestionToResponse(question), nil
}

func (s *questionService) UpdateQuestionById(ctx context.Context, id string, req request.UpdateQuestionRequest, userId string) (*question.Question, error) {
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	questionRes, err := s.questionRepo.GetQuestionById(ctx, objectId)
	if err != nil {
		return nil, err
	}

	if questionRes == nil {
		return nil, errors.New(string(question.ErrQuestionNotFound))
	}

	if questionRes.CreatedBy != userId {
		return nil, question.ErrUnauthorized
	}

	if scoringFieldsChanging(questionRes, req) {
		inUse, err := s.questionUsedInSubmissions(ctx, questionRes.ID.Hex())
		if err != nil {
			return nil, err
		}
		if inUse {
			return nil, question.ErrQuestionInUse
		}
	}

	if req.Type != "" {
		questionRes.Type = req.Type
	}

	if req.Subject != "" {
		questionRes.Subject = req.Subject
	}

	if req.Grade != "" {
		questionRes.Grade = req.Grade
	}

	if req.Difficulty != "" {
		questionRes.Difficulty = req.Difficulty
	}

	if req.Question != "" {
		questionRes.Question = req.Question
	}

	if len(req.Options) > 0 {
		questionRes.Options = req.Options
	}

	if req.CorrectIndex != nil {
		questionRes.CorrectIndex = req.CorrectIndex
	}

	if req.CorrectBool != nil {
		questionRes.CorrectBool = req.CorrectBool
	}

	if req.Explanation != "" {
		questionRes.Explanation = req.Explanation
	}

	var toDelete []string
	if len(req.Pairs) > 0 {
		oldPairs := questionRes.Pairs
		newPairs := make([]question.Pair, 0, len(req.Pairs))
		for i, p := range req.Pairs {
			var old question.Pair
			if i < len(oldPairs) {
				old = oldPairs[i]
			}
			pair, staleIDs, err := s.buildPair(ctx, p, old)
			if err != nil {
				return nil, err
			}
			newPairs = append(newPairs, pair)
			toDelete = append(toDelete, staleIDs...)
		}
		questionRes.Pairs = newPairs
	}

	questionRes.UpdatedAt = time.Now()

	if err := s.questionRepo.Update(ctx, questionRes); err != nil {
		return nil, err
	}

	for _, publicID := range toDelete {
		if publicID == "" {
			continue
		}
		_ = s.cloudinary.DeleteImage(ctx, publicID)
	}

	return questionRes, nil
}

func (s *questionService) DeleteQuestionById(ctx context.Context, id string, userId string) error {
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	questionRes, err := s.questionRepo.GetQuestionById(ctx, objectId)
	if err != nil {
		return err
	}

	if questionRes == nil {
		return question.ErrQuestionNotFound
	}

	if questionRes.CreatedBy != userId {
		return question.ErrUnauthorized
	}

	inUse, err := s.questionIsReferenced(ctx, questionRes.ID.Hex())
	if err != nil {
		return err
	}
	if inUse {
		return question.ErrQuestionInUse
	}

	if err := s.questionRepo.Delete(ctx, objectId); err != nil {
		return err
	}

	if questionRes.Type == string(question.QuestionTypeMatching) {
		for _, p := range questionRes.Pairs {
			if p.LeftPublicID != "" {
				_ = s.cloudinary.DeleteImage(ctx, p.LeftPublicID)
			}
			if p.RightPublicID != "" {
				_ = s.cloudinary.DeleteImage(ctx, p.RightPublicID)
			}
		}
	}

	return nil
}

func (s *questionService) buildPair(ctx context.Context, p request.PairRequest, old question.Pair) (question.Pair, []string, error) {
	pair := question.Pair{}
	var stale []string

	left, leftStale, err := s.resolveSide(ctx, p.LeftKind, p.Left, p.LeftFile, old.Left, old.LeftPublicID, old.LeftKind)
	if err != nil {
		return pair, nil, err
	}
	pair.Left = left.value
	pair.LeftPublicID = left.publicID
	pair.LeftKind = left.kind
	stale = append(stale, leftStale...)

	right, rightStale, err := s.resolveSide(ctx, p.RightKind, p.Right, p.RightFile, old.Right, old.RightPublicID, old.RightKind)
	if err != nil {
		return pair, nil, err
	}
	pair.Right = right.value
	pair.RightPublicID = right.publicID
	pair.RightKind = right.kind
	stale = append(stale, rightStale...)

	return pair, stale, nil
}

type pairSide struct {
	value    string
	publicID string
	kind     string
}

func (s *questionService) resolveSide(
	ctx context.Context,
	kind, text string,
	file *multipart.FileHeader,
	oldValue, oldPublicID, oldKind string,
) (pairSide, []string, error) {
	if kind == string(question.Image) {
		if file != nil {
			url, publicID, err := s.uploadFile(ctx, file)
			if err != nil {
				return pairSide{}, nil, err
			}
			var stale []string
			if oldPublicID != "" && oldPublicID != publicID {
				stale = append(stale, oldPublicID)
			}
			return pairSide{value: url, publicID: publicID, kind: string(question.Image)}, stale, nil
		}
		if oldKind == string(question.Image) && oldValue != "" {
			return pairSide{value: oldValue, publicID: oldPublicID, kind: string(question.Image)}, nil, nil
		}
		return pairSide{}, nil, errors.New(string(question.ErrInvalidPairs))
	}

	var stale []string
	if oldPublicID != "" {
		stale = append(stale, oldPublicID)
	}
	return pairSide{value: text, publicID: "", kind: string(question.Text)}, stale, nil
}

func (s *questionService) uploadFile(ctx context.Context, header *multipart.FileHeader) (string, string, error) {
	src, err := header.Open()
	if err != nil {
		return "", "", err
	}
	defer src.Close()
	return s.cloudinary.UploadImage(ctx, src, "questions")
}

func (s *questionService) questionIsReferenced(ctx context.Context, questionID string) (bool, error) {
	setCount, err := s.questionSetRepo.CountByQuestionID(ctx, questionID)
	if err != nil {
		return false, err
	}
	if setCount > 0 {
		return true, nil
	}

	homeworkCount, err := s.homeworkRepo.CountByQuestionID(ctx, questionID)
	if err != nil {
		return false, err
	}
	if homeworkCount > 0 {
		return true, nil
	}

	return s.questionUsedInSubmissions(ctx, questionID)
}

func (s *questionService) questionUsedInSubmissions(ctx context.Context, questionID string) (bool, error) {
	count, err := s.submissionRepo.CountByQuestionID(ctx, questionID)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func scoringFieldsChanging(q *question.Question, req request.UpdateQuestionRequest) bool {
	if req.Type != "" && req.Type != q.Type {
		return true
	}
	if req.Difficulty != "" && req.Difficulty != q.Difficulty {
		return true
	}
	if len(req.Options) > 0 && !sameStrings(req.Options, q.Options) {
		return true
	}
	if req.CorrectIndex != nil && !sameIntPtr(req.CorrectIndex, q.CorrectIndex) {
		return true
	}
	if req.CorrectBool != nil && !sameBoolPtr(req.CorrectBool, q.CorrectBool) {
		return true
	}
	if len(req.Pairs) > 0 {
		return true
	}
	return false
}

func sameStrings(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func sameIntPtr(a, b *int) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	return *a == *b
}

func sameBoolPtr(a, b *bool) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	return *a == *b
}