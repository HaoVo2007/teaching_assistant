package usecase

import (
	"context"
	"crypto/rand"
	"encoding/base32"
	"errors"
	"mime/multipart"
	"strings"
	"teaching_assistant/internal/delivery/http/mapper"
	"teaching_assistant/internal/delivery/http/request"
	"teaching_assistant/internal/delivery/http/response"
	"teaching_assistant/internal/domain/class"
	"teaching_assistant/internal/domain/homework"
	"teaching_assistant/internal/domain/student"
	"teaching_assistant/internal/domain/user"
	infrastructureCloudinary "teaching_assistant/internal/infrastructure/cloudinary"
	"teaching_assistant/pkg/pagination"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type classUsecase struct {
	classRepo    class.ClassRepository
	studentRepo  student.StudentRepository
	userRepo     user.UserRepository
	homeworkRepo homework.HomeworkRepository
	cloudinary   *infrastructureCloudinary.CloudinaryUploader
}

func NewClassUsecase(
	classRepo class.ClassRepository,
	studentRepo student.StudentRepository,
	userRepo user.UserRepository,
	homeworkRepo homework.HomeworkRepository,
	cloudinary *infrastructureCloudinary.CloudinaryUploader,
) class.ClassService {
	return &classUsecase{
		classRepo:    classRepo,
		studentRepo:  studentRepo,
		userRepo:     userRepo,
		homeworkRepo: homeworkRepo,
		cloudinary:   cloudinary,
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

	names := studentNames(req.Students)

	now := time.Now()
	classID := primitive.NewObjectID()
	students := make([]*student.Student, 0, len(names))
	studentIDs := make([]string, 0, len(names))
	for _, name := range names {
		id := primitive.NewObjectID()
		studentIDs = append(studentIDs, id.Hex())
		students = append(students, &student.Student{
			ID:        id,
			Code:      newStudentCode(),
			Name:      name,
			ClassID:   classID.Hex(),
			Status:    student.StudentStatusActive,
			CreatedBy: userId,
			CreatedAt: now,
			UpdatedAt: now,
		})
	}

	item := &class.Class{
		ID:          classID,
		Name:        req.Name,
		Description: req.Description,
		Image:       image,
		PublicID:    publicID,
		Students:    studentIDs,
		CreatedBy:   userId,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	if err := u.classRepo.Create(ctx, item); err != nil {
		return err
	}

	if err := u.studentRepo.CreateMany(ctx, students); err != nil {
		return err
	}

	return nil
}

func (u *classUsecase) GetClasses(ctx context.Context, userId string, params pagination.Params, name string) (*response.ClassResponseWithMeta, error) {
	classes, total, err := u.classRepo.GetClasses(ctx, userId, params, name)
	if err != nil {
		return nil, err
	}

	studentIdUniqueMap := make(map[string]bool)
	for _, class := range classes {
		for _, studentId := range class.Students {
			if _, ok := studentIdUniqueMap[studentId]; ok {
				continue
			}
			studentIdUniqueMap[studentId] = true
		}
	}

	studentIds := make([]string, 0, len(studentIdUniqueMap))
	for studentId := range studentIdUniqueMap {
		studentIds = append(studentIds, studentId)
	}

	studentIdsPrimitive := make([]primitive.ObjectID, 0, len(studentIds))
	for _, studentId := range studentIds {
		studentIdPrimitive, err := primitive.ObjectIDFromHex(studentId)
		if err != nil {
			return nil, err
		}
		studentIdsPrimitive = append(studentIdsPrimitive, studentIdPrimitive)
	}

	students, err := u.studentRepo.GetStudentsByIds(ctx, studentIdsPrimitive)
	if err != nil {
		return nil, err
	}

	studentIdMap := make(map[string]*student.Student)
	for _, student := range students {
		studentIdMap[student.ID.Hex()] = student
	}

	classStudentsMap := make(map[string][]*student.Student)
	for _, classRes := range classes {
		classStudentsMap[classRes.ID.Hex()] = make([]*student.Student, 0)
		for _, studentId := range classRes.Students {
			classStudentsMap[classRes.ID.Hex()] = append(classStudentsMap[classRes.ID.Hex()], studentIdMap[studentId])
		}
	}

	parentsByStudentID, err := u.parentsByStudentIDs(ctx, studentIds)
	if err != nil {
		return nil, err
	}

	return &response.ClassResponseWithMeta{
		Classes: mapper.MapClassesToResponses(classes, classStudentsMap, parentsByStudentID),
		Meta:    pagination.NewMeta(params, total),
	}, nil
}

func (u *classUsecase) GetClassById(ctx context.Context, userId string, id string) (*response.ClassResponse, error) {
	item, err := u.getOwnedClass(ctx, userId, id)
	if err != nil {
		return nil, err
	}

	studentIdsPrimitive := make([]primitive.ObjectID, 0, len(item.Students))
	for _, studentId := range item.Students {
		studentIdPrimitive, err := primitive.ObjectIDFromHex(studentId)
		if err != nil {
			return nil, err
		}
		studentIdsPrimitive = append(studentIdsPrimitive, studentIdPrimitive)
	}

	students, err := u.studentRepo.GetStudentsByIds(ctx, studentIdsPrimitive)
	if err != nil {
		return nil, err
	}

	studentIdMap := make(map[string]*student.Student, len(students))
	for _, st := range students {
		studentIdMap[st.ID.Hex()] = st
	}

	ordered := make([]*student.Student, 0, len(item.Students))
	for _, studentId := range item.Students {
		if st := studentIdMap[studentId]; st != nil {
			ordered = append(ordered, st)
		}
	}

	parentsByStudentID, err := u.parentsByStudentIDs(ctx, item.Students)
	if err != nil {
		return nil, err
	}

	return mapper.MapClassToResponse(item, ordered, parentsByStudentID), nil
}

func (u *classUsecase) parentsByStudentIDs(ctx context.Context, studentIDs []string) (map[string]*user.User, error) {
	parentsByStudentID := make(map[string]*user.User)
	if len(studentIDs) == 0 {
		return parentsByStudentID, nil
	}

	guardians, err := u.studentRepo.GetGuardiansByStudentIds(ctx, studentIDs)
	if err != nil {
		return nil, err
	}

	parentOIDSet := make(map[string]primitive.ObjectID)
	for _, guardian := range guardians {
		if guardian == nil || guardian.ParentID == "" {
			continue
		}
		if _, ok := parentOIDSet[guardian.ParentID]; ok {
			continue
		}
		oid, err := primitive.ObjectIDFromHex(guardian.ParentID)
		if err != nil {
			continue
		}
		parentOIDSet[guardian.ParentID] = oid
	}

	parentOIDs := make([]primitive.ObjectID, 0, len(parentOIDSet))
	for _, oid := range parentOIDSet {
		parentOIDs = append(parentOIDs, oid)
	}

	parents, err := u.userRepo.FindByIds(ctx, parentOIDs)
	if err != nil {
		return nil, err
	}

	parentByID := make(map[string]*user.User, len(parents))
	for _, parent := range parents {
		if parent == nil {
			continue
		}
		parentByID[parent.ID.Hex()] = parent
	}

	for _, guardian := range guardians {
		if guardian == nil {
			continue
		}
		if _, exists := parentsByStudentID[guardian.StudentID]; exists {
			continue
		}
		if parent := parentByID[guardian.ParentID]; parent != nil {
			parentsByStudentID[guardian.StudentID] = parent
		}
	}

	return parentsByStudentID, nil
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
		studentIDs, err := u.syncClassStudents(ctx, userId, item, req.Students)
		if err != nil {
			return err
		}
		item.Students = studentIDs
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

func (u *classUsecase) syncClassStudents(ctx context.Context, userId string, item *class.Class, tokens []string) ([]string, error) {
	oldIDs := make(map[string]struct{}, len(item.Students))
	for _, id := range item.Students {
		oldIDs[id] = struct{}{}
	}

	lookupOIDs := make([]primitive.ObjectID, 0)
	lookupSeen := make(map[string]struct{})
	for _, token := range tokens {
		token = strings.TrimSpace(token)
		oid, err := primitive.ObjectIDFromHex(token)
		if err != nil {
			continue
		}
		if _, ok := lookupSeen[token]; ok {
			continue
		}
		lookupSeen[token] = struct{}{}
		lookupOIDs = append(lookupOIDs, oid)
	}

	existing := make(map[string]*student.Student, len(lookupOIDs))
	if len(lookupOIDs) > 0 {
		students, err := u.studentRepo.GetStudentsByIds(ctx, lookupOIDs)
		if err != nil {
			return nil, err
		}
		for _, st := range students {
			existing[st.ID.Hex()] = st
		}
	}

	seen := make(map[string]struct{})
	keptIDs := make([]string, 0, len(tokens))
	keepObjectIDs := make([]primitive.ObjectID, 0)
	newStudents := make([]*student.Student, 0)
	now := time.Now()

	for _, token := range tokens {
		token = strings.TrimSpace(token)
		if token == "" {
			continue
		}

		oid, err := primitive.ObjectIDFromHex(token)
		if err == nil {
			if _, dup := seen[token]; dup {
				continue
			}
			st, ok := existing[token]
			if !ok {
				return nil, student.ErrStudentNotFound
			}
			if st.ClassID != item.ID.Hex() || st.CreatedBy != userId {
				return nil, student.ErrStudentNotInClass
			}
			seen[token] = struct{}{}
			keptIDs = append(keptIDs, token)
			keepObjectIDs = append(keepObjectIDs, oid)
			continue
		}

		id := primitive.NewObjectID()
		keptIDs = append(keptIDs, id.Hex())
		newStudents = append(newStudents, &student.Student{
			ID:        id,
			Code:      newStudentCode(),
			Name:      token,
			ClassID:   item.ID.Hex(),
			Status:    student.StudentStatusActive,
			CreatedBy: userId,
			CreatedAt: now,
			UpdatedAt: now,
		})
	}

	removed := make([]primitive.ObjectID, 0)
	for id := range oldIDs {
		if _, keep := seen[id]; keep {
			continue
		}
		oid, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			continue
		}
		removed = append(removed, oid)
	}

	if err := u.studentRepo.CreateMany(ctx, newStudents); err != nil {
		return nil, err
	}
	if err := u.studentRepo.ActivateByIDs(ctx, keepObjectIDs); err != nil {
		return nil, err
	}
	if err := u.studentRepo.DeactivateByIDs(ctx, removed); err != nil {
		return nil, err
	}

	return keptIDs, nil
}

func (u *classUsecase) DeleteClassById(ctx context.Context, userId string, id string) error {
	item, err := u.getOwnedClass(ctx, userId, id)
	if err != nil {
		return err
	}

	homeworkCount, err := u.homeworkRepo.CountByClassID(ctx, item.ID.Hex())
	if err != nil {
		return err
	}
	if homeworkCount > 0 {
		return class.ErrClassInUse
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

func newStudentCode() string {
	var b [3]byte
	if _, err := rand.Read(b[:]); err != nil {
		return "HS" + time.Now().Format("060102150405")
	}

	suffix := strings.ToUpper(base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(b[:]))
	if len(suffix) > 4 {
		suffix = suffix[:4]
	}
	return "HS" + time.Now().Format("060102") + "-" + suffix
}

func studentNames(names []string) []string {
	out := make([]string, 0, len(names))
	for _, name := range names {
		name = strings.TrimSpace(name)
		if name == "" {
			continue
		}
		out = append(out, name)
	}
	return out
}
