package mapper

import (
	"teaching_assistant/internal/delivery/http/response"
	"teaching_assistant/internal/domain/student"
	"teaching_assistant/internal/domain/user"
)

func MapStudentToResponse(item *student.Student, parent *user.User) *response.StudentResponse {
	res := &response.StudentResponse{
		ID:        item.ID.Hex(),
		Name:      item.Name,
		Code:      item.Code,
		Image:     item.Image,
		PublicID:  item.PublicID,
		ClassID:   item.ClassID,
		Status:    string(item.Status),
		CreatedAt: item.CreatedAt,
		UpdatedAt: item.UpdatedAt,
	}
	if parent != nil {
		mapped := MapUserToUserResponse(parent)
		res.Guardian = &mapped
	}
	return res
}
