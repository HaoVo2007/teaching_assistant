package handler

import (
	"mime/multipart"

	"teaching_assistant/internal/delivery/http/middleware"
	"teaching_assistant/internal/delivery/http/request"
	"teaching_assistant/internal/domain/class"
	"teaching_assistant/pkg/common"
	"teaching_assistant/pkg/pagination"
	"teaching_assistant/pkg/response"

	"github.com/gofiber/fiber/v2"
)

const maxClassImageSize = 5 << 20

type ClassHandler struct {
	classService class.ClassService
}

func NewClassHandler(classService class.ClassService) *ClassHandler {
	return &ClassHandler{
		classService: classService,
	}
}

func (h *ClassHandler) CreateClass(c *fiber.Ctx) error {
	userId, err := middleware.UserIDFromCtx(c)
	if err != nil {
		return response.Fail(c, fiber.StatusUnauthorized, string(common.ErrUnauthorized), "UNAUTHORIZED")
	}

	var req request.CreateClassRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Fail(c, fiber.StatusBadRequest, string(common.ErrBadRequest), "INVALID_REQUEST_BODY")
	}

	if err := bindClassImage(c, &req.Image); err != nil {
		return response.Fail(c, fiber.StatusBadRequest, err.Error(), "INVALID_IMAGE")
	}

	err = h.classService.CreateClass(c.UserContext(), userId, req)
	if err != nil {
		return response.Fail(c, fiber.StatusInternalServerError, err.Error(), "INTERNAL_SERVER_ERROR")
	}

	return response.Created(c, "Class created successfully", nil)
}

func (h *ClassHandler) GetClasses(c *fiber.Ctx) error {
	userId, err := middleware.UserIDFromCtx(c)
	if err != nil {
		return response.Fail(c, fiber.StatusUnauthorized, string(common.ErrUnauthorized), "UNAUTHORIZED")
	}

	name := c.Query("name")
	pageSize := c.QueryInt("page_size", 10)
	pageIndex := c.QueryInt("page_index", 1)
	params := pagination.New(pageIndex, pageSize)

	classes, err := h.classService.GetClasses(c.UserContext(), userId, params, name)
	if err != nil {
		return response.Fail(c, fiber.StatusInternalServerError, err.Error(), "INTERNAL_SERVER_ERROR")
	}

	return response.OK(c, "Classes fetched successfully", classes)
}

func (h *ClassHandler) GetClassById(c *fiber.Ctx) error {
	id := c.Params("id")
	userId, err := middleware.UserIDFromCtx(c)
	if err != nil {
		return response.Fail(c, fiber.StatusUnauthorized, string(common.ErrUnauthorized), "UNAUTHORIZED")
	}

	item, err := h.classService.GetClassById(c.UserContext(), userId, id)
	if err != nil {
		return response.Fail(c, fiber.StatusInternalServerError, err.Error(), "INTERNAL_SERVER_ERROR")
	}

	return response.OK(c, "Class fetched successfully", item)
}

func (h *ClassHandler) UpdateClassById(c *fiber.Ctx) error {
	id := c.Params("id")
	userId, err := middleware.UserIDFromCtx(c)
	if err != nil {
		return response.Fail(c, fiber.StatusUnauthorized, string(common.ErrUnauthorized), "UNAUTHORIZED")
	}

	var req request.UpdateClassRequest
	if err := c.BodyParser(&req); err != nil {
		return response.Fail(c, fiber.StatusBadRequest, string(common.ErrBadRequest), "INVALID_REQUEST_BODY")
	}

	if err := bindClassImage(c, &req.Image); err != nil {
		return response.Fail(c, fiber.StatusBadRequest, err.Error(), "INVALID_IMAGE")
	}

	err = h.classService.UpdateClassById(c.UserContext(), userId, id, req)
	if err != nil {
		return response.Fail(c, fiber.StatusInternalServerError, err.Error(), "INTERNAL_SERVER_ERROR")
	}

	return response.OK(c, "Class updated successfully", nil)
}

func (h *ClassHandler) DeleteClassById(c *fiber.Ctx) error {
	id := c.Params("id")
	userId, err := middleware.UserIDFromCtx(c)
	if err != nil {
		return response.Fail(c, fiber.StatusUnauthorized, string(common.ErrUnauthorized), "UNAUTHORIZED")
	}

	err = h.classService.DeleteClassById(c.UserContext(), userId, id)
	if err != nil {
		return response.Fail(c, fiber.StatusInternalServerError, err.Error(), "INTERNAL_SERVER_ERROR")
	}

	return response.OK(c, "Class deleted successfully", nil)
}

func bindClassImage(c *fiber.Ctx, dest **multipart.FileHeader) error {
	header, err := c.FormFile("image")
	if err != nil {
		return nil
	}
	if header.Size > maxClassImageSize {
		return class.ErrImageTooLarge
	}
	*dest = header
	return nil
}
