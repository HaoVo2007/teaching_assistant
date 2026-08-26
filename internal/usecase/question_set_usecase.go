package usecase

import (
	"context"
	"errors"
	"slices"
	"teaching_assistant/internal/delivery/http/mapper"
	"teaching_assistant/internal/delivery/http/request"
	"teaching_assistant/internal/delivery/http/response"
	"teaching_assistant/internal/domain/question"
	questionset "teaching_assistant/internal/domain/question_set"
	"teaching_assistant/pkg/pagination"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type questionSetUsecase struct {
	questionSetRepo questionset.QuestionSetRepository
	questionRepo    question.QuestionRepository
}

func NewQuestionSetService(
	questionSetRepo questionset.QuestionSetRepository,
	questionRepo question.QuestionRepository,
) questionset.QuestionSetService {
	return &questionSetUsecase{
		questionSetRepo: questionSetRepo,
		questionRepo:    questionRepo,
	}
}

func (s *questionSetUsecase) CreateQuestionSet(ctx context.Context, userId string, req request.CreateQuestionSetRequest) error {
	if req.Title == "" {
		return errors.New(string(questionset.ErrInvalidTitle))
	}

	if len(req.Questions) == 0 {
		return errors.New(string(questionset.ErrInvalidQuestions))
	}

	questionSet := &questionset.QuestionSet{
		ID:          primitive.NewObjectID(),
		Title:       req.Title,
		Description: req.Description,
		QuestionIds: req.Questions,
		CreatedBy:   userId,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	err := s.questionSetRepo.Create(ctx, questionSet)
	if err != nil {
		return err
	}

	return nil
}

func (s *questionSetUsecase) GetQuestionSets(ctx context.Context, userId string, params pagination.Params) (*response.QuestionSetResponseWithMeta, error) {
	questionSets, total, err := s.questionSetRepo.GetQuestionSets(ctx, userId, params)
	if err != nil {
		return nil, err
	}

	var questionIds []string
	for _, questionSet := range questionSets {
		questionIds = append(questionIds, questionSet.QuestionIds...)
	}

	var questionIdUnique []string
	for _, questionId := range questionIds {
		if !slices.Contains(questionIdUnique, questionId) {
			questionIdUnique = append(questionIdUnique, questionId)
		}
	}

	var questionIdPrimitive []primitive.ObjectID
	for _, questionId := range questionIdUnique {
		oid, err := primitive.ObjectIDFromHex(questionId)
		if err != nil {
			continue
		}
		questionIdPrimitive = append(questionIdPrimitive, oid)
	}

	questions, err := s.questionRepo.GetQuestionByIds(ctx, questionIdPrimitive)
	if err != nil {
		return nil, err
	}

	var questionMap = make(map[string]*question.Question)
	for _, question := range questions {
		questionMap[question.ID.Hex()] = question
	}

	questionSetResponses := make([]*response.QuestionSetResponse, 0)
	for _, questionSet := range questionSets {
		questions := make([]*response.QuestionResponse, 0)
		for _, questionId := range questionSet.QuestionIds {
			questions = append(questions, mapper.MapQuestionToResponse(questionMap[questionId]))
		}
		questionSetResponses = append(questionSetResponses, &response.QuestionSetResponse{
			ID:          questionSet.ID.Hex(),
			Title:       questionSet.Title,
			Description: questionSet.Description,
			Questions:   questions,
		})
	}

	return &response.QuestionSetResponseWithMeta{
		QuestionSets: questionSetResponses,
		Meta:         pagination.NewMeta(params, total),
	}, nil
}
