package route

import (
	"teaching_assistant/internal/delivery/http/handler"
	"teaching_assistant/internal/delivery/http/middleware"
	"teaching_assistant/internal/domain/user"
	"teaching_assistant/pkg/jwt"

	"github.com/gofiber/fiber/v2"
)

func HomeworkSubmission(api fiber.Router, homeworkSubmissionH *handler.HomeworkSubmissionHandler, jwtManager *jwt.Manager) {
	submissions := api.Group("/homework-submissions")
	{
		submissions.Post("", middleware.RequireRole(jwtManager, user.RoleParent), homeworkSubmissionH.CreateHomeworkSubmission)
		submissions.Get("", middleware.RequireRole(jwtManager, user.RoleTeacher), homeworkSubmissionH.GetHomeworkSubmissions)
		submissions.Get("/homework/:homework_id", middleware.RequireRole(jwtManager, user.RoleTeacher), homeworkSubmissionH.GetHomeworkSubmissionsByHomeworkId)
		submissions.Get("/student/by-guardian/:homework_id", middleware.RequireRole(jwtManager, user.RoleParent), homeworkSubmissionH.GetHomeworkSubmissionsByHomeworkIdByGuardian)
		submissions.Get("/:id", middleware.RequireRole(jwtManager, user.RoleTeacher), homeworkSubmissionH.GetHomeworkSubmissionById)
	}
}
