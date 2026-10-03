package route

import (
	"teaching_assistant/internal/delivery/http/handler"
	"teaching_assistant/internal/delivery/http/middleware"
	"teaching_assistant/internal/domain/user"
	"teaching_assistant/pkg/jwt"

	"github.com/gofiber/fiber/v2"
)

func User(api fiber.Router, userH *handler.UserHandler, jwtManager *jwt.Manager) {
	adminUsers := api.Group("/admin/users")
	adminUsers.Use(middleware.RequireRole(jwtManager, user.RoleAdmin))
	{
		adminUsers.Post("/create", userH.CreateUser)
	}

	users := api.Group("/users")
	{
		users.Get("/parents", middleware.RequireRole(jwtManager, user.RoleTeacher), userH.GetParents)
	}
}
