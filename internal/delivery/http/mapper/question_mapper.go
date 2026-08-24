package mapper

import (
	"teaching_assistant/internal/delivery/http/response"
	"teaching_assistant/internal/domain/question"
)

func MapQuestionsToResponses(questions []*question.Question) []*response.QuestionResponse {
	responses := make([]*response.QuestionResponse, 0)
	for _, q := range questions {
		pairs := make([]response.Pair, 0)
		for _, p := range q.Pairs {
			pairs = append(pairs, response.Pair{
				Left:          p.Left,
				LeftPublicID:  p.LeftPublicID,
				LeftKind:      p.LeftKind,
				Right:         p.Right,
				RightPublicID: p.RightPublicID,
				RightKind:     p.RightKind,
			})
		}
		responses = append(responses, &response.QuestionResponse{
			ID:           q.ID.Hex(),
			Type:         q.Type,
			Subject:      q.Subject,
			Grade:        q.Grade,
			Difficulty:   q.Difficulty,
			Question:     q.Question,
			Options:      q.Options,
			CorrectIndex: q.CorrectIndex,
			CorrectBool:  q.CorrectBool,
			Pairs:        pairs,
			Explanation:  q.Explanation,
			CreatedBy:    q.CreatedBy,
			CreatedAt:    q.CreatedAt,
			UpdatedAt:    q.UpdatedAt,
		})
	}
	return responses
}

func MapQuestionToResponse(question *question.Question) *response.QuestionResponse {
	pairs := make([]response.Pair, 0)
	for _, p := range question.Pairs {
		pairs = append(pairs, response.Pair{
			Left:          p.Left,
			LeftPublicID:  p.LeftPublicID,
			LeftKind:      p.LeftKind,
			Right:         p.Right,
			RightPublicID: p.RightPublicID,
			RightKind:     p.RightKind,
		})
	}
	return &response.QuestionResponse{
		ID:           question.ID.Hex(),
		Type:         question.Type,
		Subject:      question.Subject,
		Grade:        question.Grade,
		Difficulty:   question.Difficulty,
		Question:     question.Question,
		Options:      question.Options,
		CorrectIndex: question.CorrectIndex,
		CorrectBool:  question.CorrectBool,
		Pairs:        pairs,
		Explanation:  question.Explanation,
		CreatedBy:    question.CreatedBy,
		CreatedAt:    question.CreatedAt,
		UpdatedAt:    question.UpdatedAt,
	}
}
