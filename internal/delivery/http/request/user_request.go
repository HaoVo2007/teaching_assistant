package request

type CreateUserRequest struct {
	Username string `json:"username" validate:"required,min=3,max=32"`
	Email    string `json:"email" validate:"email"`
	Password string `json:"password" validate:"required,min=8"`
}

type LoginUserRequest struct {
	Email    string `json:"email" validate:"email"`
	Password string `json:"password" validate:"required,min=8"`
}
