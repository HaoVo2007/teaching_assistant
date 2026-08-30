package mapper

import (
	"teaching_assistant/internal/delivery/http/response"
	"teaching_assistant/internal/domain/homework"
	"teaching_assistant/internal/domain/question"
)

func QuestionResponseMap(questions []*question.Question) map[string]*response.QuestionResponse {
	m := make(map[string]*response.QuestionResponse, len(questions))
	for _, q := range questions {
		if q == nil {
			continue
		}
		m[q.ID.Hex()] = MapQuestionToResponse(q)
	}
	return m
}

func MapHomeworkToResponse(hw *homework.Homework, questionMap map[string]*response.QuestionResponse) *response.HomeworkResponse {
	questions := make([]*response.QuestionResponse, 0, len(hw.Questions))
	for _, questionID := range hw.Questions {
		q, ok := questionMap[questionID]
		if !ok {
			continue
		}
		questions = append(questions, q)
	}

	return &response.HomeworkResponse{
		ID:          hw.ID.Hex(),
		ClassID:     hw.ClassID,
		Title:       hw.Title,
		Description: hw.Description,
		DueDate:     hw.DueDate.Format("2006-01-02"),
		Questions:   questions,
		CreatedAt:   hw.CreatedAt,
		UpdatedAt:   hw.UpdatedAt,
	}
}

func MapHomeworksToResponses(homeworks []*homework.Homework, questionMap map[string]*response.QuestionResponse) []response.HomeworkResponse {
	responses := make([]response.HomeworkResponse, 0, len(homeworks))
	for _, hw := range homeworks {
		if hw == nil {
			continue
		}
		responses = append(responses, *MapHomeworkToResponse(hw, questionMap))
	}
	return responses
}
