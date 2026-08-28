package http

import (
	"teaching_assistant/internal/delivery/http/handler"
	"teaching_assistant/internal/delivery/http/middleware"
	"teaching_assistant/pkg/jwt"

	"github.com/gofiber/fiber/v2"
)

func NewRouter(
	app *fiber.App,
	userH *handler.UserHandler,
	questionH *handler.QuestionHandler,
	questionSetH *handler.QuestionSetHandler,
	classH *handler.ClassHandler,
	jwtManager *jwt.Manager,
) {
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
}
