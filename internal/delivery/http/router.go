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
}
