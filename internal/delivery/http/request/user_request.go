package request

type CreateUserRequest struct {
	Username string `json:"username" validate:"required,min=3,max=32"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
	Role     string `json:"role" validate:"required,oneof=admin teacher student parent"`
}

type LoginUserRequest struct {
	Email    string `json:"email" validate:"email"`
	Password string `json:"password" validate:"required,min=8"`
}
