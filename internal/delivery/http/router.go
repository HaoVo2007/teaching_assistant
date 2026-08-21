package http

import (
	"teaching_assistant/internal/delivery/http/handler"
	"teaching_assistant/internal/delivery/http/middleware"
	"teaching_assistant/pkg/jwt"

	"github.com/gofiber/fiber/v2"
)

func NewRouter(app *fiber.App, userH *handler.UserHandler, jwtManager *jwt.Manager) {
	api := app.Group("/api/v1")
	auth := api.Group("/auth")
	{
		auth.Post("/register", userH.Register)
		auth.Post("/login", userH.Login)
		auth.Post("/logout", userH.Logout, middleware.AuthMiddleware(jwtManager))
	}

}
