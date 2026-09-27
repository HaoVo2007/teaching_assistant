package response

import "teaching_assistant/pkg/pagination"

type AuthResponse struct {
	Token string       `json:"token"`
	User  UserResponse `json:"user"`
}

type UserResponse struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Role     string `json:"role"`
}

type UserResponseWithMeta struct {
	Parents []*UserResponse `json:"parents"`
	Meta    pagination.Meta `json:"meta"`
}
