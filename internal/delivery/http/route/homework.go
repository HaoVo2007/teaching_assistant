package route

import (
	"teaching_assistant/internal/delivery/http/handler"
	"teaching_assistant/internal/delivery/http/middleware"
	"teaching_assistant/internal/domain/user"
	"teaching_assistant/pkg/jwt"

	"github.com/gofiber/fiber/v2"
)

func Homework(api fiber.Router, homeworkH *handler.HomeworkHandler, jwtManager *jwt.Manager) {
	homeworks := api.Group("/homeworks")
	{
		homeworks.Post("", middleware.RequireRole(jwtManager, user.RoleTeacher), homeworkH.CreateHomework)
		homeworks.Get("", middleware.RequireRole(jwtManager, user.RoleTeacher), homeworkH.GetHomeworks)
		homeworks.Get("/class/:class_id", middleware.RequireRole(jwtManager, user.RoleTeacher), homeworkH.GetHomeworksByClassId)
		homeworks.Get("/student/by-guardian", middleware.RequireRole(jwtManager, user.RoleParent), homeworkH.GetHomeworksByStudentId)
		homeworks.Get("/:id", middleware.RequireRole(jwtManager, user.RoleTeacher), homeworkH.GetHomeworkById)
		homeworks.Put("/:id", middleware.RequireRole(jwtManager, user.RoleTeacher), homeworkH.UpdateHomeworkById)
		homeworks.Delete("/:id", middleware.RequireRole(jwtManager, user.RoleTeacher), homeworkH.DeleteHomeworkById)
	}
}
