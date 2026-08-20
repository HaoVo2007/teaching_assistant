package handler

import (
	"teaching_assistant/internal/domain/user"

	"github.com/gofiber/fiber/v2"
)

type UserHandler struct {
	userService user.UserService
}

func NewUserHandler(userService user.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

func (h *UserHandler) CreateUser(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"message": "User created successfully",
	})
}
