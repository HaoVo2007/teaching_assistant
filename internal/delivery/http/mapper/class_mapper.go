package mapper

import (
	"teaching_assistant/internal/delivery/http/response"
	"teaching_assistant/internal/domain/class"
	"teaching_assistant/internal/domain/student"
	"teaching_assistant/internal/domain/user"
)

func MapClassToResponse(c *class.Class, students []*student.Student, parentsByStudentID map[string]*user.User) *response.ClassResponse {
	mappedStudents := make([]response.StudentResponse, 0, len(students))
	for _, st := range students {
		if st == nil {
			continue
		}
		var parent *user.User
		if parentsByStudentID != nil {
			parent = parentsByStudentID[st.ID.Hex()]
		}
		mappedStudents = append(mappedStudents, *MapStudentToResponse(st, parent))
	}
	return &response.ClassResponse{
		ID:          c.ID.Hex(),
		Name:        c.Name,
		Description: c.Description,
		Image:       c.Image,
		PublicID:    c.PublicID,
		Students:    mappedStudents,
		CreatedBy:   c.CreatedBy,
		CreatedAt:   c.CreatedAt,
		UpdatedAt:   c.UpdatedAt,
	}
}

func MapClassesToResponses(classes []*class.Class, classStudentsMap map[string][]*student.Student, parentsByStudentID map[string]*user.User) []*response.ClassResponse {
	responses := make([]*response.ClassResponse, 0, len(classes))
	for _, c := range classes {
		responses = append(responses, MapClassToResponse(c, classStudentsMap[c.ID.Hex()], parentsByStudentID))
	}
	return responses
}
