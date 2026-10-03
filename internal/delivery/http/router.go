package http

import (
	"teaching_assistant/internal/delivery/http/handler"
	"teaching_assistant/internal/delivery/http/route"
	"teaching_assistant/pkg/jwt"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func NewRouter(
	app *fiber.App,
	userH *handler.UserHandler,
	questionH *handler.QuestionHandler,
	questionSetH *handler.QuestionSetHandler,
	classH *handler.ClassHandler,
	homeworkH *handler.HomeworkHandler,
	homeworkSubmissionH *handler.HomeworkSubmissionHandler,
	studentH *handler.StudentHandler,
	jwtManager *jwt.Manager,
) {
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:5173,https://teachingassistantfe.netlify.app",
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		AllowCredentials: true,
	}))

	api := app.Group("/api/v1")

	route.Auth(api, userH, jwtManager)
	route.User(api, userH, jwtManager)
	route.Student(api, studentH, jwtManager)
	route.Question(api, questionH, jwtManager)
	route.QuestionSet(api, questionSetH, jwtManager)
	route.Class(api, classH, jwtManager)
	route.Homework(api, homeworkH, jwtManager)
	route.HomeworkSubmission(api, homeworkSubmissionH, jwtManager)
}
