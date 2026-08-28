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
	if req.QuestionType == "" {
		return errors.New(string(questionset.ErrInvalidQuestionType))
	}

	if req.Title == "" {
		return errors.New(string(questionset.ErrInvalidTitle))
	}

	if len(req.Questions) == 0 {
		return errors.New(string(questionset.ErrInvalidQuestions))
	}

	var questionIdsPrimitive []primitive.ObjectID
	for _, questionId := range req.Questions {
		oid, err := primitive.ObjectIDFromHex(questionId)
		if err != nil {
			return errors.New(string(questionset.ErrInvalidQuestions))
		}
		questionIdsPrimitive = append(questionIdsPrimitive, oid)
	}

	questions, err := s.questionRepo.GetQuestionByIds(ctx, questionIdsPrimitive)
	if err != nil {
		return err
	}

	for _, question := range questions {
		if questionset.QuestionSetType(question.Type) != questionset.QuestionSetType(req.QuestionType) {
			return errors.New(string(questionset.ErrInvalidQuestionTypeForQuestion))
		}
	}

	questionSet := &questionset.QuestionSet{
		ID:           primitive.NewObjectID(),
		Title:        req.Title,
		Description:  req.Description,
		QuestionType: questionset.QuestionSetType(req.QuestionType),
		QuestionIds:  req.Questions,
		CreatedBy:    userId,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	err = s.questionSetRepo.Create(ctx, questionSet)
	if err != nil {
		return err
	}

	return nil
}

func (s *questionSetUsecase) GetQuestionSets(ctx context.Context, userId string, params pagination.Params, title string, questionType string) (*response.QuestionSetResponseWithMeta, error) {
	questionSets, total, err := s.questionSetRepo.GetQuestionSets(ctx, userId, params, title, questionType)
	if err != nil {
		return nil, err
	}

	var questionIds []string
	for _, questionSet := range questionSets {
		questionIds = append(questionIds, questionSet.QuestionIds...)
	}

	if len(questionIds) == 0 {
		return &response.QuestionSetResponseWithMeta{
			QuestionSets: make([]*response.QuestionSetResponse, 0),
			Meta:         pagination.NewMeta(params, 0),
		}, nil
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
			mapped := questionMap[questionId]
			if mapped == nil {
				continue
			}
			questions = append(questions, mapper.MapQuestionToResponse(mapped))
		}
		questionSetResponses = append(questionSetResponses, &response.QuestionSetResponse{
			ID:          questionSet.ID.Hex(),
			Title:       questionSet.Title,
			Description: questionSet.Description,
			Questions:   questions,
			CreatedAt:   questionSet.CreatedAt,
			UpdatedAt:   questionSet.UpdatedAt,
		})
	}

	return &response.QuestionSetResponseWithMeta{
		QuestionSets: questionSetResponses,
		Meta:         pagination.NewMeta(params, total),
	}, nil
}

func (s *questionSetUsecase) GetQuestionSetById(ctx context.Context, userId string, id string) (*response.QuestionSetResponse, error) {
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	questionSet, err := s.questionSetRepo.GetQuestionSetById(ctx, objectId)
	if err != nil {
		return nil, errors.New(string(questionset.ErrQuestionSetNotFound))
	}

	if questionSet.CreatedBy != userId {
		return nil, errors.New(string(questionset.ErrUnauthorized))
	}

	var questionIds []string
	for _, questionId := range questionSet.QuestionIds {
		questionIds = append(questionIds, questionId)
	}

	var questionIdPrimitive []primitive.ObjectID
	for _, questionId := range questionIds {
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

	questionResponses := make([]*response.QuestionResponse, 0)
	for _, questionId := range questionSet.QuestionIds {
		mapped := questionMap[questionId]
		if mapped == nil {
			continue
		}
		questionResponses = append(questionResponses, mapper.MapQuestionToResponse(mapped))
	}

	return &response.QuestionSetResponse{
		ID:           questionSet.ID.Hex(),
		Title:        questionSet.Title,
		Description:  questionSet.Description,
		QuestionType: string(questionSet.QuestionType),
		Questions:    questionResponses,
		CreatedAt:    questionSet.CreatedAt,
		UpdatedAt:    questionSet.UpdatedAt,
	}, nil
}

func (s *questionSetUsecase) UpdateQuestionSetById(ctx context.Context, userId string, id string, req request.UpdateQuestionSetRequest) error {
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	questionSet, err := s.questionSetRepo.GetQuestionSetById(ctx, objectId)
	if err != nil {
		return errors.New(string(questionset.ErrQuestionSetNotFound))
	}

	if questionSet.CreatedBy != userId {
		return errors.New(string(questionset.ErrUnauthorized))
	}

	if req.Title != nil {
		questionSet.Title = *req.Title
	}

	if req.Description != nil {
		questionSet.Description = req.Description
	}

	if req.Questions != nil {
		questionSet.QuestionIds = req.Questions
	}

	questionSet.UpdatedAt = time.Now()

	err = s.questionSetRepo.UpdateQuestionSetById(ctx, objectId, questionSet)
	if err != nil {
		return err
	}

	return nil
}

func (s *questionSetUsecase) DeleteQuestionSetById(ctx context.Context, userId string, id string) error {
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	questionSet, err := s.questionSetRepo.GetQuestionSetById(ctx, objectId)
	if err != nil {
		return errors.New(string(questionset.ErrQuestionSetNotFound))
	}

	if questionSet.CreatedBy != userId {
		return errors.New(string(questionset.ErrUnauthorized))
	}

	err = s.questionSetRepo.DeleteQuestionSetById(ctx, objectId)
	if err != nil {
		return err
	}

	return nil
}
