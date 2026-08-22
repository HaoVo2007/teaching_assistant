package usecase

import (
	"context"
	"time"

	"teaching_assistant/internal/delivery/http/request"
	"teaching_assistant/internal/domain/question"

	infrastructureCloudinary "teaching_assistant/internal/infrastructure/cloudinary"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type questionService struct {
	questionRepo question.QuestionRepository
	cloudinary   *infrastructureCloudinary.CloudinaryUploader
}

func NewQuestionService(
	questionRepo question.QuestionRepository,
	cloudinary *infrastructureCloudinary.CloudinaryUploader,
) question.QuestionService {
	return &questionService{
		questionRepo: questionRepo,
		cloudinary:   cloudinary,
	}
}

func (s *questionService) CreateQuestion(ctx context.Context, req request.CreateQuestionRequest) (*question.Question, error) {
	if req.Type == "" {
		return nil, question.ErrInvalidType
	}

	if req.Subject == "" {
		return nil, question.ErrInvalidSubject
	}

	if req.Grade == "" {
		return nil, question.ErrInvalidGrade
	}

	pairs := make([]question.Pair, 0, len(req.Pairs))
	for _, p := range req.Pairs {
		pair := question.Pair{}
		if p.LeftKind == string(question.Image) {
			if p.LeftFile == nil {
				return nil, question.ErrInvalidPairs
			}
			src, err := p.LeftFile.Open()
			if err != nil {
				return nil, err
			}
			url, publicID, err := s.cloudinary.UploadImage(ctx, src, "questions")
			src.Close()
			if err != nil {
				return nil, err
			}
			pair.Left = url
			pair.LeftPublicID = publicID
			pair.LeftKind = string(question.Image)
		}
		if p.RightKind == string(question.Image) {
			if p.RightFile == nil {
				return nil, question.ErrInvalidPairs
			}
			src, err := p.RightFile.Open()
			if err != nil {
				return nil, err
			}
			url, publicID, err := s.cloudinary.UploadImage(ctx, src, "questions")
			src.Close()
			if err != nil {
				return nil, err
			}
			pair.Right = url
			pair.RightPublicID = publicID
			pair.RightKind = string(question.Image)
		}
		if p.LeftKind == string(question.Text) {
			pair.Left = p.Left
			pair.LeftKind = string(question.Text)
		}
		if p.RightKind == string(question.Text) {
			pair.Right = p.Right
			pair.RightKind = string(question.Text)
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
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := s.questionRepo.Create(ctx, q); err != nil {
		return nil, err
	}
	return q, nil
}
