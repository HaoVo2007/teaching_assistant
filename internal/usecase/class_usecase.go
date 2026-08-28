package usecase

import (
	"context"
	"errors"
	"mime/multipart"
	"teaching_assistant/internal/delivery/http/mapper"
	"teaching_assistant/internal/delivery/http/request"
	"teaching_assistant/internal/delivery/http/response"
	"teaching_assistant/internal/domain/class"
	infrastructureCloudinary "teaching_assistant/internal/infrastructure/cloudinary"
	"teaching_assistant/pkg/pagination"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type classUsecase struct {
	classRepo  class.ClassRepository
	cloudinary *infrastructureCloudinary.CloudinaryUploader
}

func NewClassUsecase(
	classRepo class.ClassRepository,
	cloudinary *infrastructureCloudinary.CloudinaryUploader,
) class.ClassService {
	return &classUsecase{
		classRepo:  classRepo,
		cloudinary: cloudinary,
	}
}

func (u *classUsecase) CreateClass(ctx context.Context, userId string, req request.CreateClassRequest) error {
	if req.Name == "" {
		return errors.New(string(class.ErrInvalidClass))
	}

	image, publicID, err := u.uploadClassImage(ctx, req.Image)
	if err != nil {
		return err
	}

	item := &class.Class{
		ID:          primitive.NewObjectID(),
		Name:        req.Name,
		Description: req.Description,
		Image:       image,
		PublicID:    publicID,
		Students:    req.Students,
		CreatedBy:   userId,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	return u.classRepo.Create(ctx, item)
}

func (u *classUsecase) GetClasses(ctx context.Context, userId string, params pagination.Params, name string) (*response.ClassResponseWithMeta, error) {
	classes, total, err := u.classRepo.GetClasses(ctx, userId, params, name)
	if err != nil {
		return nil, err
	}

	return &response.ClassResponseWithMeta{
		Classes: mapper.MapClassesToResponses(classes),
		Meta:    pagination.NewMeta(params, total),
	}, nil
}

func (u *classUsecase) GetClassById(ctx context.Context, userId string, id string) (*response.ClassResponse, error) {
	item, err := u.getOwnedClass(ctx, userId, id)
	if err != nil {
		return nil, err
	}
	return mapper.MapClassToResponse(item), nil
}

func (u *classUsecase) UpdateClassById(ctx context.Context, userId string, id string, req request.UpdateClassRequest) error {
	item, err := u.getOwnedClass(ctx, userId, id)
	if err != nil {
		return err
	}

	if req.Name != nil {
		if *req.Name == "" {
			return errors.New(string(class.ErrInvalidClass))
		}
		item.Name = *req.Name
	}

	if req.Description != nil {
		item.Description = *req.Description
	}

	if req.Students != nil {
		item.Students = req.Students
	}

	if req.Image != nil {
		image, publicID, err := u.uploadClassImage(ctx, req.Image)
		if err != nil {
			return err
		}
		oldPublicID := item.PublicID
		item.Image = image
		item.PublicID = publicID
		if oldPublicID != "" && oldPublicID != publicID {
			_ = u.cloudinary.DeleteImage(ctx, oldPublicID)
		}
	}

	item.UpdatedAt = time.Now()
	return u.classRepo.UpdateClassById(ctx, item)
}

func (u *classUsecase) DeleteClassById(ctx context.Context, userId string, id string) error {
	item, err := u.getOwnedClass(ctx, userId, id)
	if err != nil {
		return err
	}

	if err := u.classRepo.DeleteClassById(ctx, item.ID); err != nil {
		return err
	}

	if item.PublicID != "" {
		_ = u.cloudinary.DeleteImage(ctx, item.PublicID)
	}

	return nil
}

func (u *classUsecase) getOwnedClass(ctx context.Context, userId, id string) (*class.Class, error) {
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	item, err := u.classRepo.GetClassById(ctx, objectId)
	if err != nil {
		return nil, errors.New(string(class.ErrClassNotFound))
	}

	if item.CreatedBy != userId {
		return nil, errors.New(string(class.ErrUnauthorized))
	}

	return item, nil
}

func (u *classUsecase) uploadClassImage(ctx context.Context, header *multipart.FileHeader) (string, string, error) {
	if header == nil {
		return "", "", nil
	}

	src, err := header.Open()
	if err != nil {
		return "", "", err
	}
	defer src.Close()

	return u.cloudinary.UploadImage(ctx, src, "classes")
}
