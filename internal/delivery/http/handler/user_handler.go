package handler

import (
	"teaching_assistant/internal/delivery/http/httperr"
	"teaching_assistant/internal/delivery/http/middleware"
	"teaching_assistant/internal/delivery/http/request"
	"teaching_assistant/internal/domain/user"
	"teaching_assistant/pkg/common"
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
		return response.Fail(c, fiber.StatusBadRequest, string(common.ErrBadRequest), "INVALID_REQUEST_BODY")
	}

	res, err := h.userService.Register(c.UserContext(), req)
	if err != nil {
		return httperr.Fail(c, err)
	}

	return response.OK(c, "User created successfully", res)
}

func (h *UserHandler) Login(c *fiber.Ctx) error {
	var req request.LoginUserRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Fail(c, fiber.StatusBadRequest, string(common.ErrBadRequest), "INVALID_REQUEST_BODY")
	}

	res, err := h.userService.Login(c.UserContext(), req)
	if err != nil {
		return httperr.Fail(c, err)
	}
	return response.OK(c, "Login successful", res)
}

func (h *UserHandler) Logout(c *fiber.Ctx) error {
	userId, err := middleware.UserIDFromCtx(c)
	if err != nil {
		return response.Fail(c, fiber.StatusUnauthorized, string(common.ErrUnauthorized), "UNAUTHORIZED")
	}

	err = h.userService.Logout(c.UserContext(), userId)
	if err != nil {
		return httperr.Fail(c, err)
	}
	return response.OK(c, "Logout successful", nil)
}

func (h *UserHandler) CreateUser(c *fiber.Ctx) error {
	var req request.CreateUserRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Fail(c, fiber.StatusBadRequest, string(common.ErrBadRequest), "INVALID_REQUEST_BODY")
	}

	res, err := h.userService.CreateUser(c.UserContext(), req)
	if err != nil {
		return httperr.Fail(c, err)
	}
	return response.OK(c, "User created successfully", res)
}
