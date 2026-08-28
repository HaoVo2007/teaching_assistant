package mapper

import (
	"teaching_assistant/internal/delivery/http/response"
	"teaching_assistant/internal/domain/class"
)

func MapClassToResponse(c *class.Class) *response.ClassResponse {
	students := c.Students
	if students == nil {
		students = make([]string, 0)
	}
	return &response.ClassResponse{
		ID:          c.ID.Hex(),
		Name:        c.Name,
		Description: c.Description,
		Image:       c.Image,
		PublicID:    c.PublicID,
		Students:    students,
		CreatedBy:   c.CreatedBy,
		CreatedAt:   c.CreatedAt,
		UpdatedAt:   c.UpdatedAt,
	}
}

func MapClassesToResponses(classes []*class.Class) []*response.ClassResponse {
	responses := make([]*response.ClassResponse, 0, len(classes))
	for _, c := range classes {
		responses = append(responses, MapClassToResponse(c))
	}
	return responses
}
