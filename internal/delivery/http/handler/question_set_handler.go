package handler

import (
	"teaching_assistant/internal/delivery/http/middleware"
	"teaching_assistant/internal/delivery/http/request"
	questionset "teaching_assistant/internal/domain/question_set"
	"teaching_assistant/pkg/common"
	"teaching_assistant/pkg/pagination"
	"teaching_assistant/pkg/response"

	"github.com/gofiber/fiber/v2"
)

type QuestionSetHandler struct {
	questionSetService questionset.QuestionSetService
}

func NewQuestionSetHandler(questionSetService questionset.QuestionSetService) *QuestionSetHandler {
	return &QuestionSetHandler{
		questionSetService: questionSetService,
	}
}

func (h *QuestionSetHandler) CreateQuestionSet(c *fiber.Ctx) error {
	userId, err := middleware.UserIDFromCtx(c)
	if err != nil {
		return response.Fail(c, fiber.StatusUnauthorized, string(common.ErrUnauthorized), "UNAUTHORIZED")
	}

	var req request.CreateQuestionSetRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Fail(c, fiber.StatusBadRequest, string(common.ErrBadRequest), "INVALID_REQUEST_BODY")
	}

	err = h.questionSetService.CreateQuestionSet(c.Context(), userId, req)
	if err != nil {
		return response.Fail(c, fiber.StatusInternalServerError, err.Error(), "INTERNAL_SERVER_ERROR")
	}

	return response.OK(c, "Question set created successfully", nil)
}

func (h *QuestionSetHandler) GetQuestionSets(c *fiber.Ctx) error {
	userId, err := middleware.UserIDFromCtx(c)
	if err != nil {
		return response.Fail(c, fiber.StatusUnauthorized, string(common.ErrUnauthorized), "UNAUTHORIZED")
	}

	pageSize := c.QueryInt("page_size", 10)
	pageIndex := c.QueryInt("page_index", 1)
	params := pagination.New(pageIndex, pageSize)

	questionSets, err := h.questionSetService.GetQuestionSets(c.UserContext(), userId, params)
	if err != nil {
		return response.Fail(c, fiber.StatusInternalServerError, err.Error(), "INTERNAL_SERVER_ERROR")
	}

	return response.OK(c, "Question sets fetched successfully", questionSets)
}
