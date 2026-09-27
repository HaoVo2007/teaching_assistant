package mapper

import (
	"teaching_assistant/internal/delivery/http/response"
	"teaching_assistant/internal/domain/question"
)

func MapQuestionsToResponses(questions []*question.Question) []*response.QuestionResponse {
	responses := make([]*response.QuestionResponse, 0)
	for _, q := range questions {
		responses = append(responses, MapQuestionToResponse(q))
	}
	return responses
}

func MapQuestionToResponse(item *question.Question) *response.QuestionResponse {
	return &response.QuestionResponse{
		ID:           item.ID.Hex(),
		Type:         item.Type,
		Subject:      item.Subject,
		Grade:        item.Grade,
		Difficulty:   item.Difficulty,
		Question:     item.Question,
		Options:      item.Options,
		CorrectIndex: item.CorrectIndex,
		CorrectBool:  item.CorrectBool,
		Explanation:  item.Explanation,
		CreatedBy:    item.CreatedBy,
		CreatedAt:    item.CreatedAt,
		UpdatedAt:    item.UpdatedAt,
	}
}

func HideQuestionAnswers(q *response.QuestionResponse) {
	if q == nil {
		return
	}
	q.CorrectIndex = nil
	q.CorrectBool = nil
	q.Explanation = ""
}

func HideHomeworkAnswers(hw *response.HomeworkResponse) {
	if hw == nil {
		return
	}
	for _, q := range hw.Questions {
		HideQuestionAnswers(q)
	}
}
