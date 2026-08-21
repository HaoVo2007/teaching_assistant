package handler

import (
	"teaching_assistant/internal/delivery/http/request"
	"teaching_assistant/internal/domain/user"
	"teaching_assistant/pkg/response"

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

func (h *UserHandler) Register(c *fiber.Ctx) error {
	var req request.CreateUserRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Fail(c, fiber.StatusBadRequest, "Invalid request body", "INVALID_REQUEST_BODY")
	}

	res, err := h.userService.Register(c.UserContext(), req)
	if err != nil {
		return response.Fail(c, fiber.StatusInternalServerError, "Failed to create user", "FAILED_TO_CREATE_USER")
	}

	return response.OK(c, "User created successfully", res)
}
