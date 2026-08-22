package handler

import (
	"fmt"
	"mime/multipart"

	"teaching_assistant/internal/delivery/http/request"
	"teaching_assistant/internal/domain/question"
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
	var req request.CreateQuestionRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Fail(c, fiber.StatusBadRequest, "Invalid request body", "INVALID_REQUEST_BODY")
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

	_, err := h.questionService.CreateQuestion(c.UserContext(), req)
	if err != nil {
		return response.Fail(c, fiber.StatusInternalServerError, "Internal server error", "INTERNAL_SERVER_ERROR")
	}

	return response.Created(c, "Question created successfully", nil)
}

func pairImage(c *fiber.Ctx, index int, side string) (*multipart.FileHeader, error) {
	header, err := c.FormFile(fmt.Sprintf("pairs[%d].%s_image", index, side))
	if err != nil {
		return nil, question.ErrInvalidPairs
	}
	if header.Size > maxPairImageSize {
		return nil, question.ErrImageTooLarge
	}
	return header, nil
}
