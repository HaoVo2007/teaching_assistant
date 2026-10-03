package route

import (
	"teaching_assistant/internal/delivery/http/handler"
	"teaching_assistant/internal/delivery/http/middleware"
	"teaching_assistant/pkg/jwt"

	"github.com/gofiber/fiber/v2"
)

func Auth(api fiber.Router, userH *handler.UserHandler, jwtManager *jwt.Manager) {
	auth := api.Group("/auth")
	{
		auth.Post("/register", userH.Register)
		auth.Post("/login", userH.Login)
		auth.Post("/logout", middleware.AuthMiddleware(jwtManager), userH.Logout)
	}
}
