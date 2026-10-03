package route

import (
	"teaching_assistant/internal/delivery/http/handler"
	"teaching_assistant/internal/delivery/http/middleware"
	"teaching_assistant/internal/domain/user"
	"teaching_assistant/pkg/jwt"

	"github.com/gofiber/fiber/v2"
)

func Student(api fiber.Router, studentH *handler.StudentHandler, jwtManager *jwt.Manager) {
	students := api.Group("/students")
	{
		students.Post("/claim", middleware.RequireRole(jwtManager, user.RoleTeacher), studentH.ClaimStudent)
		students.Get("/by-guardian", middleware.RequireRole(jwtManager, user.RoleParent), studentH.GetStudentByGuardian)
	}
}
