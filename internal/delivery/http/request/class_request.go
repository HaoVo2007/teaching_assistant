package request

import "mime/multipart"

type CreateClassRequest struct {
	Name        string                `form:"name" validate:"required,min=3,max=100"`
	Description string                `form:"description" validate:"required,min=3,max=1000"`
	Image       *multipart.FileHeader `form:"image"`
	Students    []string              `form:"students"`
}

type UpdateClassRequest struct {
	Name        *string               `form:"name"`
	Description *string               `form:"description"`
	Image       *multipart.FileHeader `form:"image"`
	Students    []string              `form:"students"`
}
