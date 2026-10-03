package route

import (
	"teaching_assistant/internal/delivery/http/handler"
	"teaching_assistant/internal/delivery/http/middleware"
	"teaching_assistant/internal/domain/user"
	"teaching_assistant/pkg/jwt"

	"github.com/gofiber/fiber/v2"
)

func Class(api fiber.Router, classH *handler.ClassHandler, jwtManager *jwt.Manager) {
	classes := api.Group("/classes")
	classes.Use(middleware.RequireRole(jwtManager, user.RoleTeacher))
	{
		classes.Post("", classH.CreateClass)
		classes.Get("", classH.GetClasses)
		classes.Get("/:id", classH.GetClassById)
		classes.Put("/:id", classH.UpdateClassById)
		classes.Delete("/:id", classH.DeleteClassById)
	}
}
