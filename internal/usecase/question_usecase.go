package usecase

import (
	"context"
	"encoding/json"
	"strings"
	"time"

	"teaching_assistant/internal/delivery/http/mapper"
	"teaching_assistant/internal/delivery/http/request"
	"teaching_assistant/internal/delivery/http/response"
	"teaching_assistant/internal/domain/homework"
	homeworksubmission "teaching_assistant/internal/domain/homework_submission"
	"teaching_assistant/internal/domain/question"
	questionset "teaching_assistant/internal/domain/question_set"
	"teaching_assistant/pkg/llm"
	"teaching_assistant/pkg/pagination"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type questionUsecase struct {
	questionRepo    question.QuestionRepository
	questionSetRepo questionset.QuestionSetRepository
	homeworkRepo    homework.HomeworkRepository
	submissionRepo  homeworksubmission.HomeworkSubmissionRepository
	llm             llm.LLM
}

func NewQuestionUsecase(
	questionRepo question.QuestionRepository,
	questionSetRepo questionset.QuestionSetRepository,
	homeworkRepo homework.HomeworkRepository,
	submissionRepo homeworksubmission.HomeworkSubmissionRepository,
	llmClient llm.LLM,
) question.QuestionService {
	return &questionUsecase{
		questionRepo:    questionRepo,
		questionSetRepo: questionSetRepo,
		homeworkRepo:    homeworkRepo,
		submissionRepo:  submissionRepo,
		llm:             llmClient,
	}
}

func (s *questionUsecase) CreateQuestion(ctx context.Context, req request.CreateQuestionRequest, userId string) (*question.Question, error) {
	if !question.IsSupportedType(req.Type) {
		return nil, question.ErrInvalidType
	}

	if req.Subject == "" {
		return nil, question.ErrInvalidSubject
	}

	if req.Grade == "" {
		return nil, question.ErrInvalidGrade
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

func (s *questionUsecase) CreateQuestionBatch(ctx context.Context, req []*request.CreateQuestionBatchRequest, userId string) ([]*question.Question, error) {
	items := make([]*question.Question, 0, len(req))
	for _, item := range req {
		if !question.IsSupportedType(item.Type) {
			return nil, question.ErrInvalidType
		}

		if item.Subject == "" {
			return nil, question.ErrInvalidSubject
		}
		
		if item.Grade == "" {
			return nil, question.ErrInvalidGrade
		}

		if item.Difficulty == "" {
			return nil, question.ErrInvalidDifficulty
		}
		
		if item.Question == "" {
			return nil, question.ErrRequiredQuestion
		}

		if question.QuestionType(item.Type) == question.QuestionTypeMultipleChoice {
			if item.CorrectIndex == nil {
				return nil, question.ErrRequiredCorrectAnswer
			}
			if *item.CorrectIndex < 0 || *item.CorrectIndex >= len(item.Options) {
				return nil, question.ErrInvalidCorrectIndex
			}
		}
		
		if question.QuestionType(item.Type) == question.QuestionTypeTrueFalse {
			if item.CorrectBool == nil {
				return nil, question.ErrRequiredCorrectAnswer
			}
		}

		q := &question.Question{
			ID:           primitive.NewObjectID(),
			Type:         item.Type,
			Subject:      item.Subject,
			Grade:        item.Grade,
			Difficulty:   item.Difficulty,
			Question:     item.Question,
			Options:      item.Options,
			CorrectIndex: item.CorrectIndex,
			CorrectBool:  item.CorrectBool,
			Explanation:  item.Explanation,
			CreatedBy:    userId,
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		}

		items = append(items, q)
	}

	if err := s.questionRepo.CreateMany(ctx, items); err != nil {
		return nil, err
	}
	return items, nil
}

func (s *questionUsecase) GetQuestions(ctx context.Context, userId string, params pagination.Params, questionType, questionName, subject, grade, difficulty string) (*response.QuestionResponseWithMeta, error) {
	if questionType != "" && !question.IsSupportedType(questionType) {
		return nil, question.ErrInvalidType
	}

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

func (s *questionUsecase) GetQuestionById(ctx context.Context, id string) (*response.QuestionResponse, error) {
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, question.ErrQuestionNotFound
	}

	item, err := s.questionRepo.GetQuestionById(ctx, objectId)
	if err != nil {
		return nil, err
	}
	if item == nil {
		return nil, question.ErrQuestionNotFound
	}

	return mapper.MapQuestionToResponse(item), nil
}

func (s *questionUsecase) UpdateQuestionById(ctx context.Context, id string, req request.UpdateQuestionRequest, userId string) (*question.Question, error) {
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, question.ErrQuestionNotFound
	}

	questionRes, err := s.questionRepo.GetQuestionById(ctx, objectId)
	if err != nil {
		return nil, err
	}

	if questionRes == nil {
		return nil, question.ErrQuestionNotFound
	}

	if questionRes.CreatedBy != userId {
		return nil, question.ErrUnauthorized
	}

	if req.Type != "" && !question.IsSupportedType(req.Type) {
		return nil, question.ErrInvalidType
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

	questionRes.UpdatedAt = time.Now()

	if err := s.questionRepo.Update(ctx, questionRes); err != nil {
		return nil, err
	}

	return questionRes, nil
}

func (s *questionUsecase) DeleteQuestionById(ctx context.Context, id string, userId string) error {
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return question.ErrQuestionNotFound
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

	return s.questionRepo.Delete(ctx, objectId)
}

func (s *questionUsecase) GenerateQuestionByLLM(ctx context.Context, userId string, questionType, grade, subject, difficulty string, quantity int) ([]*response.QuestionResponse, error) {
	if questionType == "" {
		return nil, question.ErrRequiredQuestionType
	}

	if !question.IsSupportedType(questionType) {
		return nil, question.ErrInvalidType
	}

	if grade == "" {
		return nil, question.ErrRequiredGrade
	}

	if !question.IsSupportedGrade(grade) {
		return nil, question.ErrInvalidGrade
	}

	if subject == "" {
		return nil, question.ErrRequiredSubject
	}

	if !question.IsSupportedSubject(subject) {
		return nil, question.ErrInvalidSubject
	}

	if difficulty == "" {
		return nil, question.ErrRequiredDifficulty
	}

	if !question.IsSupportedDifficulty(difficulty) {
		return nil, question.ErrInvalidDifficulty
	}

	if quantity < 1 || quantity > 20 {
		return nil, question.ErrInvalidQuantity
	}
	
	if s.llm == nil {
		return nil, question.ErrLLMFailed
	}

	raw, err := s.llm.GenerateJSON(ctx, generateQuestionSystemPrompt, generateQuestionUserPrompt(questionType, grade, subject, difficulty, quantity))
	if err != nil {
		return nil, question.ErrLLMFailed
	}

	items, err := parseGeneratedQuestions(raw, questionType, grade, subject, difficulty)
	if err != nil || len(items) == 0 {
		return nil, question.ErrLLMInvalidResponse
	}

	now := time.Now()
	out := make([]*response.QuestionResponse, 0, len(items))
	for _, item := range items {
		item.CreatedBy = userId
		item.CreatedAt = now
		item.UpdatedAt = now
		mapped := mapper.MapQuestionToResponse(item)
		mapped.ID = ""
		out = append(out, mapped)
	}
	return out, nil
}

type llmQuestionPayload struct {
	Questions []llmGeneratedQuestion `json:"questions"`
}

type llmGeneratedQuestion struct {
	Question     string   `json:"question"`
	Options      []string `json:"options"`
	CorrectIndex *int     `json:"correct_index"`
	CorrectBool  *bool    `json:"correct_bool"`
	Explanation  string   `json:"explanation"`
}

func parseGeneratedQuestions(raw, questionType, grade, subject, difficulty string) ([]*question.Question, error) {
	payload, err := unmarshalLLMQuestions(raw)
	if err != nil {
		return nil, err
	}

	items := make([]*question.Question, 0, len(payload.Questions))
	for _, generated := range payload.Questions {
		item, ok := generated.toQuestion(questionType, grade, subject, difficulty)
		if !ok {
			continue
		}
		items = append(items, item)
	}
	if len(items) == 0 {
		return nil, question.ErrLLMInvalidResponse
	}
	return items, nil
}

func unmarshalLLMQuestions(raw string) (llmQuestionPayload, error) {
	var payload llmQuestionPayload
	if err := json.Unmarshal([]byte(extractJSON(raw)), &payload); err == nil && len(payload.Questions) > 0 {
		return payload, nil
	}

	var list []llmGeneratedQuestion
	if err := json.Unmarshal([]byte(extractJSON(raw)), &list); err != nil {
		return llmQuestionPayload{}, err
	}
	return llmQuestionPayload{Questions: list}, nil
}

func extractJSON(raw string) string {
	raw = strings.TrimSpace(raw)
	raw = strings.TrimPrefix(raw, "```json")
	raw = strings.TrimPrefix(raw, "```JSON")
	raw = strings.TrimPrefix(raw, "```")
	raw = strings.TrimSuffix(raw, "```")
	raw = strings.TrimSpace(raw)

	if i := strings.Index(raw, "{"); i >= 0 {
		if j := strings.LastIndex(raw, "}"); j > i {
			return raw[i : j+1]
		}
	}
	if i := strings.Index(raw, "["); i >= 0 {
		if j := strings.LastIndex(raw, "]"); j > i {
			return raw[i : j+1]
		}
	}
	return raw
}

func (g llmGeneratedQuestion) toQuestion(questionType, grade, subject, difficulty string) (*question.Question, bool) {
	text := strings.TrimSpace(g.Question)
	if text == "" {
		return nil, false
	}

	item := &question.Question{
		Type:        questionType,
		Subject:     subject,
		Grade:       grade,
		Difficulty:  difficulty,
		Question:    text,
		Explanation: strings.TrimSpace(g.Explanation),
	}

	switch question.QuestionType(questionType) {
	case question.QuestionTypeMultipleChoice:
		if len(g.Options) < 2 || g.CorrectIndex == nil {
			return nil, false
		}
		if *g.CorrectIndex < 0 || *g.CorrectIndex >= len(g.Options) {
			return nil, false
		}
		item.Options = g.Options
		item.CorrectIndex = g.CorrectIndex
	case question.QuestionTypeTrueFalse:
		if g.CorrectBool == nil {
			return nil, false
		}
		item.CorrectBool = g.CorrectBool
	default:
		return nil, false
	}

	return item, true
}

func (s *questionUsecase) questionIsReferenced(ctx context.Context, questionID string) (bool, error) {
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

func (s *questionUsecase) questionUsedInSubmissions(ctx context.Context, questionID string) (bool, error) {
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
