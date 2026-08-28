package response

import (
	"teaching_assistant/pkg/pagination"
	"time"
)

type ClassResponse struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Image       string    `json:"image"`
	PublicID    string    `json:"public_id"`
	Students    []string  `json:"students"`
	CreatedBy   string    `json:"created_by"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type ClassResponseWithMeta struct {
	Classes []*ClassResponse `json:"classes"`
	Meta    pagination.Meta  `json:"meta"`
}
