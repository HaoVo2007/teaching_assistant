package handler

import (
	"teaching_assistant/internal/delivery/http/middleware"
	"teaching_assistant/internal/delivery/http/request"
	"teaching_assistant/internal/domain/homework"
	"teaching_assistant/pkg/common"
	"teaching_assistant/pkg/pagination"
	"teaching_assistant/pkg/response"

	"github.com/gofiber/fiber/v2"
)

type HomeworkHandler struct {
	homeworkService homework.HomeworkService
}

func NewHomeworkHandler(homeworkService homework.HomeworkService) *HomeworkHandler {
	return &HomeworkHandler{
		homeworkService: homeworkService,
	}
}

func (h *HomeworkHandler) CreateHomework(c *fiber.Ctx) error {
	userId, err := middleware.UserIDFromCtx(c)
	if err != nil {
		return response.Fail(c, fiber.StatusUnauthorized, string(common.ErrUnauthorized), "UNAUTHORIZED")
	}

	var req request.CreateHomeworkRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Fail(c, fiber.StatusBadRequest, err.Error(), "BAD_REQUEST")
	}

	err = h.homeworkService.CreateHomework(c.UserContext(), userId, req)
	if err != nil {
		return response.Fail(c, fiber.StatusInternalServerError, err.Error(), "INTERNAL_SERVER_ERROR")
	}

	return response.OK(c, "Homework created successfully", nil)
}

func (h *HomeworkHandler) GetHomeworks(c *fiber.Ctx) error {
	userId, err := middleware.UserIDFromCtx(c)
	if err != nil {
		return response.Fail(c, fiber.StatusUnauthorized, string(common.ErrUnauthorized), "UNAUTHORIZED")
	}

	pageSize := c.QueryInt("page_size", 10)
	pageIndex := c.QueryInt("page_index", 1)

	params := pagination.New(pageIndex, pageSize)

	homeworks, err := h.homeworkService.GetHomeworks(c.UserContext(), userId, params)
	if err != nil {
		return response.Fail(c, fiber.StatusInternalServerError, err.Error(), "INTERNAL_SERVER_ERROR")
	}

	return response.OK(c, "Homeworks fetched successfully", homeworks)
}

func (h *HomeworkHandler) GetHomeworkById(c *fiber.Ctx) error {
	userId, err := middleware.UserIDFromCtx(c)
	if err != nil {
		return response.Fail(c, fiber.StatusUnauthorized, string(common.ErrUnauthorized), "UNAUTHORIZED")
	}

	id := c.Params("id")

	homework, err := h.homeworkService.GetHomeworkById(c.UserContext(), userId, id)
	if err != nil {
		return response.Fail(c, fiber.StatusInternalServerError, err.Error(), "INTERNAL_SERVER_ERROR")
	}

	return response.OK(c, "Homework fetched successfully", homework)
}

func (h *HomeworkHandler) UpdateHomeworkById(c *fiber.Ctx) error {
	userId, err := middleware.UserIDFromCtx(c)
	if err != nil {
		return response.Fail(c, fiber.StatusUnauthorized, string(common.ErrUnauthorized), "UNAUTHORIZED")
	}

	id := c.Params("id")

	var req request.UpdateHomeworkRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Fail(c, fiber.StatusBadRequest, err.Error(), "BAD_REQUEST")
	}

	err = h.homeworkService.UpdateHomeworkById(c.UserContext(), userId, id, req)
	if err != nil {
		return response.Fail(c, fiber.StatusInternalServerError, err.Error(), "INTERNAL_SERVER_ERROR")
	}

	return response.OK(c, "Homework updated successfully", nil)
}

func (h *HomeworkHandler) DeleteHomeworkById(c *fiber.Ctx) error {
	userId, err := middleware.UserIDFromCtx(c)
	if err != nil {
		return response.Fail(c, fiber.StatusUnauthorized, string(common.ErrUnauthorized), "UNAUTHORIZED")
	}

	id := c.Params("id")

	err = h.homeworkService.DeleteHomeworkById(c.UserContext(), userId, id)
	if err != nil {
		return response.Fail(c, fiber.StatusInternalServerError, err.Error(), "INTERNAL_SERVER_ERROR")
	}

	return response.OK(c, "Homework deleted successfully", nil)
}

func (h *HomeworkHandler) GetHomeworksByClassId(c *fiber.Ctx) error {
	userId, err := middleware.UserIDFromCtx(c)
	if err != nil {
		return response.Fail(c, fiber.StatusUnauthorized, string(common.ErrUnauthorized), "UNAUTHORIZED")
	}

	classId := c.Params("class_id")
	pageSize := c.QueryInt("page_size", 10)
	pageIndex := c.QueryInt("page_index", 1)

	params := pagination.New(pageIndex, pageSize)

	homeworks, err := h.homeworkService.GetHomeworksByClassId(c.UserContext(), userId, classId, params)
	if err != nil {
		return response.Fail(c, fiber.StatusInternalServerError, err.Error(), "INTERNAL_SERVER_ERROR")
	}

	return response.OK(c, "Homeworks fetched successfully", homeworks)
}
