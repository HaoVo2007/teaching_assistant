package http

import (
	"teaching_assistant/internal/delivery/http/handler"
	"github.com/gofiber/fiber/v2"
)

func NewRouter(app *fiber.App, userH *handler.UserHandler) {
	api := app.Group("/api/v1")
	users := api.Group("/users")
	{
		users.Post("/register", userH.Register)
	}

}
