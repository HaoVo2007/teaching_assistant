package route

import (
	"teaching_assistant/internal/delivery/http/handler"
	"teaching_assistant/internal/delivery/http/middleware"
	"teaching_assistant/internal/domain/user"
	"teaching_assistant/pkg/jwt"

	"github.com/gofiber/fiber/v2"
)

func QuestionSet(api fiber.Router, questionSetH *handler.QuestionSetHandler, jwtManager *jwt.Manager) {
	questionSets := api.Group("/question-sets")
	{
		questionSets.Use(middleware.RequireRole(jwtManager, user.RoleTeacher))
		questionSets.Post("", questionSetH.CreateQuestionSet)
		questionSets.Get("", questionSetH.GetQuestionSets)
		questionSets.Get("/:id", questionSetH.GetQuestionSetById)
		questionSets.Put("/:id", questionSetH.UpdateQuestionSetById)
		questionSets.Delete("/:id", questionSetH.DeleteQuestionSetById)
	}
}
