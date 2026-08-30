package handler

import (
	"teaching_assistant/internal/delivery/http/middleware"
	"teaching_assistant/internal/delivery/http/request"
	homeworksubmission "teaching_assistant/internal/domain/homework_submission"
	"teaching_assistant/pkg/common"
	"teaching_assistant/pkg/pagination"
	"teaching_assistant/pkg/response"

	"github.com/gofiber/fiber/v2"
)

type HomeworkSubmissionHandler struct {
	homeworkSubmissionService homeworksubmission.HomeworkSubmissionService
}

func NewHomeworkSubmissionHandler(homeworkSubmissionService homeworksubmission.HomeworkSubmissionService) *HomeworkSubmissionHandler {
	return &HomeworkSubmissionHandler{
		homeworkSubmissionService: homeworkSubmissionService,
	}
}

func (h *HomeworkSubmissionHandler) CreateHomeworkSubmission(c *fiber.Ctx) error {
	var req request.CreateHomeworkSubmissionRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Fail(c, fiber.StatusBadRequest, err.Error(), "BAD_REQUEST")
	}

	err := h.homeworkSubmissionService.CreateHomeworkSubmission(c.UserContext(), req)
	if err != nil {
		return response.Fail(c, fiber.StatusInternalServerError, err.Error(), "INTERNAL_SERVER_ERROR")
	}

	return response.OK(c, "Homework submission created successfully", nil)
}

func (h *HomeworkSubmissionHandler) GetHomeworkSubmissions(c *fiber.Ctx) error {
	userId, err := middleware.UserIDFromCtx(c)
	if err != nil {
		return response.Fail(c, fiber.StatusUnauthorized, string(common.ErrUnauthorized), "UNAUTHORIZED")
	}

	pageSize := c.QueryInt("page_size", 10)
	pageIndex := c.QueryInt("page_index", 1)

	params := pagination.New(pageIndex, pageSize)

	homeworkSubmissions, err := h.homeworkSubmissionService.GetHomeworkSubmissions(c.UserContext(), params, userId)
	if err != nil {
		return response.Fail(c, fiber.StatusInternalServerError, err.Error(), "INTERNAL_SERVER_ERROR")
	}

	return response.OK(c, "Homework submissions fetched successfully", homeworkSubmissions)
}

func (h *HomeworkSubmissionHandler) GetHomeworkSubmissionById(c *fiber.Ctx) error {
	userId, err := middleware.UserIDFromCtx(c)
	if err != nil {
		return response.Fail(c, fiber.StatusUnauthorized, string(common.ErrUnauthorized), "UNAUTHORIZED")
	}

	homeworkSubmissionId := c.Params("id")
	if homeworkSubmissionId == "" {
		return response.Fail(c, fiber.StatusBadRequest, "Homework submission ID is required", "BAD_REQUEST")
	}

	homeworkSubmission, err := h.homeworkSubmissionService.GetHomeworkSubmissionById(c.UserContext(), homeworkSubmissionId, userId)
	if err != nil {
		return response.Fail(c, fiber.StatusInternalServerError, err.Error(), "INTERNAL_SERVER_ERROR")
	}

	return response.OK(c, "Homework submission fetched successfully", homeworkSubmission)
}

func (h *HomeworkSubmissionHandler) GetHomeworkSubmissionsByHomeworkId(c *fiber.Ctx) error {
	userId, err := middleware.UserIDFromCtx(c)
	if err != nil {
		return response.Fail(c, fiber.StatusUnauthorized, string(common.ErrUnauthorized), "UNAUTHORIZED")
	}

	homeworkId := c.Params("homework_id")
	if homeworkId == "" {
		return response.Fail(c, fiber.StatusBadRequest, "Homework ID is required", "BAD_REQUEST")
	}

	pageSize := c.QueryInt("page_size", 10)
	pageIndex := c.QueryInt("page_index", 1)

	params := pagination.New(pageIndex, pageSize)

	homeworkSubmissions, err := h.homeworkSubmissionService.GetHomeworkSubmissionsByHomeworkId(c.UserContext(), homeworkId, userId, params)
	if err != nil {
		return response.Fail(c, fiber.StatusInternalServerError, err.Error(), "INTERNAL_SERVER_ERROR")
	}

	return response.OK(c, "Homework submissions fetched successfully", homeworkSubmissions)
}
