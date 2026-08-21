package response

import (
	"github.com/gofiber/fiber/v2"
)

type Body struct {
	Success bool       `json:"success"`
	Message string     `json:"message"`
	Data    any        `json:"data,omitempty"`
	Error   *ErrorInfo `json:"error,omitempty"`
	Meta    any        `json:"meta,omitempty"`
}

type ErrorInfo struct {
	Code    string `json:"code"`
	Details any    `json:"details,omitempty"`
}

func JSON(c *fiber.Ctx, status int, body Body) error {
	return c.Status(status).JSON(body)
}

func OK(c *fiber.Ctx, message string, data any) error {
	return JSON(c, fiber.StatusOK, Body{
		Success: true,
		Message: message,
		Data:    data,
	})
}

func Created(c *fiber.Ctx, message string, data any) error {
	return JSON(c, fiber.StatusCreated, Body{
		Success: true,
		Message: message,
		Data:    data,
	})
}

func OKWithMeta(c *fiber.Ctx, message string, data any, meta any) error {
	return JSON(c, fiber.StatusOK, Body{
		Success: true,
		Message: message,
		Data:    data,
		Meta:    meta,
	})
}

func Fail(c *fiber.Ctx, status int, message, code string) error {
	return JSON(c, status, Body{
		Success: false,
		Message: message,
		Error: &ErrorInfo{
			Code: code,
		},
	})
}

func FailWithDetails(c *fiber.Ctx, status int, message, code string, details any) error {
	return JSON(c, status, Body{
		Success: false,
		Message: message,
		Error: &ErrorInfo{
			Code:    code,
			Details: details,
		},
	})
}