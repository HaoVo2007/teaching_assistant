package route

import (
	"teaching_assistant/internal/delivery/http/handler"
	"teaching_assistant/internal/delivery/http/middleware"
	"teaching_assistant/internal/domain/user"
	"teaching_assistant/pkg/jwt"

	"github.com/gofiber/fiber/v2"
)

func Question(api fiber.Router, questionH *handler.QuestionHandler, jwtManager *jwt.Manager) {
	questions := api.Group("/questions")
	questions.Use(middleware.RequireRole(jwtManager, user.RoleTeacher))
	{
		questions.Post("", questionH.CreateQuestion)
		questions.Get("", questionH.GetQuestions)
		questions.Get("/:id", questionH.GetQuestionById)
		questions.Put("/:id", questionH.UpdateQuestionById)
		questions.Delete("/:id", questionH.DeleteQuestionById)
		questions.Post("/create/batch", questionH.CreateQuestionBatch)
	}

	questionLLM := api.Group("/question-llm")
	questionLLM.Use(middleware.RequireRole(jwtManager, user.RoleTeacher))
	{
		questionLLM.Get("/generate", questionH.GenerateQuestionByLLM)
	}
}
