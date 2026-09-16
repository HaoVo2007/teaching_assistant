package response

import "time"

type StudentResponse struct {
	ID        string       `json:"id"`
	Name      string       `json:"name"`
	Code      string       `json:"code"`
	Image     string       `json:"image"`
	PublicID  string       `json:"public_id"`
	ClassID   string       `json:"class_id"`
	Status    string       `json:"status"`
	Guardian  *UserResponse `json:"guardian"`
	CreatedAt time.Time    `json:"created_at"`
	UpdatedAt time.Time    `json:"updated_at"`
}
