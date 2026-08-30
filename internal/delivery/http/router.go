package http

import (
	"teaching_assistant/internal/delivery/http/handler"
	"teaching_assistant/internal/delivery/http/middleware"
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
	jwtManager *jwt.Manager,
) {
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:5173,https://teachingassistantfe.netlify.app",
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		AllowCredentials: true,
	}))

	api := app.Group("/api/v1")

	// auth routes
	auth := api.Group("/auth")
	{
		auth.Post("/register", userH.Register)
		auth.Post("/login", userH.Login)
		auth.Post("/logout", middleware.AuthMiddleware(jwtManager), userH.Logout)
	}
	// user routes

	//question routes
	question := api.Group("/questions")
	{
		question.Post("", middleware.AuthMiddleware(jwtManager), questionH.CreateQuestion)
		question.Get("", middleware.AuthMiddleware(jwtManager), questionH.GetQuestions)
		question.Get("/:id", middleware.AuthMiddleware(jwtManager), questionH.GetQuestionById)
		question.Put("/:id", middleware.AuthMiddleware(jwtManager), questionH.UpdateQuestionById)
		question.Delete("/:id", middleware.AuthMiddleware(jwtManager), questionH.DeleteQuestionById)
	}
	// question routes

	// question set routes
	questionSet := api.Group("/question-sets")
	{
		questionSet.Post("", middleware.AuthMiddleware(jwtManager), questionSetH.CreateQuestionSet)
		questionSet.Get("", middleware.AuthMiddleware(jwtManager), questionSetH.GetQuestionSets)
		questionSet.Get("/:id", middleware.AuthMiddleware(jwtManager), questionSetH.GetQuestionSetById)
		questionSet.Put("/:id", middleware.AuthMiddleware(jwtManager), questionSetH.UpdateQuestionSetById)
		questionSet.Delete("/:id", middleware.AuthMiddleware(jwtManager), questionSetH.DeleteQuestionSetById)
	}
	// question set routes

	// class routes
	class := api.Group("/classes")
	{
		class.Post("", middleware.AuthMiddleware(jwtManager), classH.CreateClass)
		class.Get("", middleware.AuthMiddleware(jwtManager), classH.GetClasses)
		class.Get("/:id", middleware.AuthMiddleware(jwtManager), classH.GetClassById)
		class.Put("/:id", middleware.AuthMiddleware(jwtManager), classH.UpdateClassById)
		class.Delete("/:id", middleware.AuthMiddleware(jwtManager), classH.DeleteClassById)
	}
	// class routes

	// homework routes
	homework := api.Group("/homeworks")
	{
		homework.Post("", middleware.AuthMiddleware(jwtManager), homeworkH.CreateHomework)
		homework.Get("", middleware.AuthMiddleware(jwtManager), homeworkH.GetHomeworks)
		homework.Get("/:id", middleware.AuthMiddleware(jwtManager), homeworkH.GetHomeworkById)
		homework.Put("/:id", middleware.AuthMiddleware(jwtManager), homeworkH.UpdateHomeworkById)
		homework.Delete("/:id", middleware.AuthMiddleware(jwtManager), homeworkH.DeleteHomeworkById)
		homework.Get("/class/:class_id", middleware.AuthMiddleware(jwtManager), homeworkH.GetHomeworksByClassId)
	}
	// homework routes

	// homework submission routes
	homeworkSubmission := api.Group("/homework-submissions")
	{
		homeworkSubmission.Post("", homeworkSubmissionH.CreateHomeworkSubmission)
		homeworkSubmission.Get("", middleware.AuthMiddleware(jwtManager), homeworkSubmissionH.GetHomeworkSubmissions)
		homeworkSubmission.Get("/:id", middleware.AuthMiddleware(jwtManager), homeworkSubmissionH.GetHomeworkSubmissionById)
		// homeworkSubmission.Put("/:id", middleware.AuthMiddleware(jwtManager), homeworkSubmissionH.UpdateHomeworkSubmissionById)
		// homeworkSubmission.Delete("/:id", middleware.AuthMiddleware(jwtManager), homeworkSubmissionH.DeleteHomeworkSubmissionById)
		homeworkSubmission.Get("/homework/:homework_id", middleware.AuthMiddleware(jwtManager), homeworkSubmissionH.GetHomeworkSubmissionsByHomeworkId)
	}
	// homework submission routes
}
