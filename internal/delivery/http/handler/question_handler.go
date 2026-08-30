package handler

import (
	"errors"
	"fmt"
	"mime/multipart"

	"teaching_assistant/internal/delivery/http/middleware"
	"teaching_assistant/internal/delivery/http/request"
	"teaching_assistant/internal/domain/question"
	"teaching_assistant/pkg/common"
	"teaching_assistant/pkg/pagination"
	"teaching_assistant/pkg/response"

	"github.com/gofiber/fiber/v2"
)

const maxPairImageSize = 5 << 20

type QuestionHandler struct {
	questionService question.QuestionService
}

func NewQuestionHandler(questionService question.QuestionService) *QuestionHandler {
	return &QuestionHandler{
		questionService: questionService,
	}
}

func (h *QuestionHandler) CreateQuestion(c *fiber.Ctx) error {
	userId, err := middleware.UserIDFromCtx(c)
	if err != nil {
		return response.Fail(c, fiber.StatusUnauthorized, string(common.ErrUnauthorized), "UNAUTHORIZED")
	}
	var req request.CreateQuestionRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Fail(c, fiber.StatusBadRequest, string(common.ErrBadRequest), "INVALID_REQUEST_BODY")
	}

	if question.QuestionType(req.Type) == question.QuestionTypeMatching {
		for i := range req.Pairs {
			if req.Pairs[i].LeftKind == string(question.Image) {
				header, err := pairImage(c, i, "left")
				if err != nil {
					return response.Fail(c, fiber.StatusBadRequest, err.Error(), "INVALID_PAIRS")
				}
				req.Pairs[i].LeftFile = header
				req.Pairs[i].Left = ""
			}
			if req.Pairs[i].RightKind == string(question.Image) {
				header, err := pairImage(c, i, "right")
				if err != nil {
					return response.Fail(c, fiber.StatusBadRequest, err.Error(), "INVALID_PAIRS")
				}
				req.Pairs[i].RightFile = header
				req.Pairs[i].Right = ""
			}
		}
	}

	_, err = h.questionService.CreateQuestion(c.UserContext(), req, userId)
	if err != nil {
		return response.Fail(c, fiber.StatusInternalServerError, err.Error(), "INTERNAL_SERVER_ERROR")
	}

	return response.Created(c, "Question created successfully", nil)
}

func (h *QuestionHandler) GetQuestions(c *fiber.Ctx) error {
	pageSize := c.QueryInt("page_size", 10)
	pageIndex := c.QueryInt("page_index", 1)
	userId, err := middleware.UserIDFromCtx(c)
	if err != nil {
		return response.Fail(c, fiber.StatusUnauthorized, string(common.ErrUnauthorized), "UNAUTHORIZED")
	}
	questionType := c.Query("question_type")
	questionName := c.Query("question_name")
	subject := c.Query("subject")
	grade := c.Query("grade")
	difficulty := c.Query("difficulty")
	params := pagination.New(pageIndex, pageSize)
	questions, err := h.questionService.GetQuestions(c.UserContext(), userId, params, questionType, questionName, subject, grade, difficulty)
	if err != nil {
		return response.Fail(c, fiber.StatusInternalServerError, err.Error(), "INTERNAL_SERVER_ERROR")
	}

	return response.OK(c, "Questions fetched successfully", questions)
}

func (h *QuestionHandler) GetQuestionById(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return response.Fail(c, fiber.StatusBadRequest, string(common.ErrBadRequest), "INVALID_QUESTION_ID")
	}

	question, err := h.questionService.GetQuestionById(c.UserContext(), id)
	if err != nil {
		return response.Fail(c, fiber.StatusInternalServerError, err.Error(), "INTERNAL_SERVER_ERROR")
	}

	return response.OK(c, "Question fetched successfully", question)
}

func (h *QuestionHandler) UpdateQuestionById(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return response.Fail(c, fiber.StatusBadRequest, string(common.ErrBadRequest), "INVALID_QUESTION_ID")
	}

	userId, err := middleware.UserIDFromCtx(c)
	if err != nil {
		return response.Fail(c, fiber.StatusUnauthorized, string(common.ErrUnauthorized), "UNAUTHORIZED")
	}

	var req request.UpdateQuestionRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Fail(c, fiber.StatusBadRequest, string(common.ErrBadRequest), "INVALID_REQUEST_BODY")
	}

	if question.QuestionType(req.Type) == question.QuestionTypeMatching {
		for i := range req.Pairs {
			if req.Pairs[i].LeftKind == string(question.Image) {
				header, err := pairImage(c, i, "left")
				if err != nil {
					return response.Fail(c, fiber.StatusBadRequest, err.Error(), "INVALID_PAIRS")
				}
				req.Pairs[i].LeftFile = header
				req.Pairs[i].Left = ""
			}
			if req.Pairs[i].RightKind == string(question.Image) {
				header, err := pairImage(c, i, "right")
				if err != nil {
					return response.Fail(c, fiber.StatusBadRequest, err.Error(), "INVALID_PAIRS")
				}
				req.Pairs[i].RightFile = header
				req.Pairs[i].Right = ""
			}
		}
	}

	_, err = h.questionService.UpdateQuestionById(c.UserContext(), id, req, userId)
	if err != nil {
		return response.Fail(c, fiber.StatusInternalServerError, err.Error(), "INTERNAL_SERVER_ERROR")
	}

	return response.OK(c, "Question updated successfully", nil)
}

func (h *QuestionHandler) DeleteQuestionById(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return response.Fail(c, fiber.StatusBadRequest, string(common.ErrBadRequest), "INVALID_QUESTION_ID")
	}

	userId, err := middleware.UserIDFromCtx(c)
	if err != nil {
		return response.Fail(c, fiber.StatusUnauthorized, string(common.ErrUnauthorized), "UNAUTHORIZED")
	}

	err = h.questionService.DeleteQuestionById(c.UserContext(), id, userId)
	if err != nil {
		return response.Fail(c, fiber.StatusInternalServerError, err.Error(), "INTERNAL_SERVER_ERROR")
	}

	return response.OK(c, "Question deleted successfully", nil)
}

func pairImage(c *fiber.Ctx, index int, side string) (*multipart.FileHeader, error) {
	header, err := c.FormFile(fmt.Sprintf("pairs[%d].%s_image", index, side))
	if err != nil {
		return nil, nil
	}
	if header.Size > maxPairImageSize {
		return nil, errors.New(string(question.ErrImageTooLarge))
	}
	return header, nil
}
