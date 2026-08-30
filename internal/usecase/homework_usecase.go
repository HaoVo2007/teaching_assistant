package usecase

import (
	"context"
	"errors"
	"teaching_assistant/internal/delivery/http/mapper"
	"teaching_assistant/internal/delivery/http/request"
	"teaching_assistant/internal/delivery/http/response"
	"teaching_assistant/internal/domain/class"
	"teaching_assistant/internal/domain/homework"
	homeworksubmission "teaching_assistant/internal/domain/homework_submission"
	"teaching_assistant/internal/domain/question"
	"teaching_assistant/pkg/pagination"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type homeworkService struct {
	homeworkRepo   homework.HomeworkRepository
	questionRepo   question.QuestionRepository
	classRepo      class.ClassRepository
	submissionRepo homeworksubmission.HomeworkSubmissionRepository
}

func NewHomeworkService(
	homeworkRepo homework.HomeworkRepository,
	questionRepo question.QuestionRepository,
	classRepo class.ClassRepository,
	submissionRepo homeworksubmission.HomeworkSubmissionRepository,
) homework.HomeworkService {
	return &homeworkService{
		homeworkRepo:   homeworkRepo,
		questionRepo:   questionRepo,
		classRepo:      classRepo,
		submissionRepo: submissionRepo,
	}
}

func (u *homeworkService) CreateHomework(ctx context.Context, userId string, req request.CreateHomeworkRequest) error {
	if req.Title == "" {
		return homework.ErrInvalidTitle
	}

	if req.ClassID == "" {
		return homework.ErrInvalidClassID
	}

	if len(req.Questions) == 0 {
		return homework.ErrInvalidQuestions
	}

	if req.DueDate == "" {
		return homework.ErrInvalidDueDate
	}

	dueDate, err := time.Parse("2006-01-02", req.DueDate)
	if err != nil {
		return homework.ErrInvalidDueDate
	}

	if err := u.ensureOwnedClass(ctx, userId, req.ClassID); err != nil {
		return err
	}

	if err := u.ensureQuestionsExist(ctx, req.Questions); err != nil {
		return err
	}

	dueDateUtc := dueDate.UTC()

	homework := &homework.Homework{
		ID:          primitive.NewObjectID(),
		ClassID:     req.ClassID,
		Title:       req.Title,
		Description: &req.Description,
		Questions:   req.Questions,
		DueDate:     dueDateUtc,
		CreatedBy:   userId,
		CreatedAt:   time.Now().UTC(),
		UpdatedAt:   time.Now().UTC(),
	}

	if err := u.homeworkRepo.CreateHomework(ctx, homework); err != nil {
		return err
	}

	return nil
}

func (u *homeworkService) GetHomeworks(ctx context.Context, userId string, params pagination.Params) (*response.HomeworkResponseWithMeta, error) {
	homeworks, total, err := u.homeworkRepo.GetHomeworks(ctx, userId, "", params)
	if err != nil {
		return nil, err
	}

	var questionIdsSet = map[string]bool{}
	for _, homework := range homeworks {
		for _, questionId := range homework.Questions {
			questionIdsSet[questionId] = true
		}
	}

	questionIds := make([]string, 0, len(questionIdsSet))
	for questionId := range questionIdsSet {
		questionIds = append(questionIds, questionId)
	}

	questionIdsObjectIDs := make([]primitive.ObjectID, 0, len(questionIds))
	for _, questionId := range questionIds {
		questionIdObjectID, err := primitive.ObjectIDFromHex(questionId)
		if err != nil {
			return nil, err
		}
		questionIdsObjectIDs = append(questionIdsObjectIDs, questionIdObjectID)
	}

	var questions []*question.Question
	if len(questionIdsObjectIDs) > 0 {
		questions, err = u.questionRepo.GetQuestionByIds(ctx, questionIdsObjectIDs)
		if err != nil {
			return nil, err
		}
	}

	return &response.HomeworkResponseWithMeta{
		Homeworks: mapper.MapHomeworksToResponses(homeworks, mapper.QuestionResponseMap(questions)),
		Meta:      pagination.NewMeta(params, total),
	}, nil
}

func (u *homeworkService) GetHomeworkById(ctx context.Context, userId string, id string) (*response.HomeworkResponse, error) {
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	homeworkRes, err := u.homeworkRepo.GetHomeworkById(ctx, objectId)
	if err != nil {
		return nil, err
	}

	if homeworkRes.CreatedBy != userId {
		return nil, errors.New(string(homework.ErrHomeworkNotAuthorized))
	}

	questionIds := homeworkRes.Questions
	questionIdsObjectIDs := make([]primitive.ObjectID, 0, len(questionIds))
	for _, questionId := range questionIds {
		questionIdObjectID, err := primitive.ObjectIDFromHex(questionId)
		if err != nil {
			return nil, err
		}
		questionIdsObjectIDs = append(questionIdsObjectIDs, questionIdObjectID)
	}

	var questions []*question.Question
	if len(questionIdsObjectIDs) > 0 {
		questions, err = u.questionRepo.GetQuestionByIds(ctx, questionIdsObjectIDs)
		if err != nil {
			return nil, err
		}
	}

	return mapper.MapHomeworkToResponse(homeworkRes, mapper.QuestionResponseMap(questions)), nil
}

func (u *homeworkService) UpdateHomeworkById(ctx context.Context, userId string, id string, req request.UpdateHomeworkRequest) error {
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	homeworkRes, err := u.homeworkRepo.GetHomeworkById(ctx, objectId)
	if err != nil {
		return err
	}

	if homeworkRes.CreatedBy != userId {
		return homework.ErrHomeworkNotAuthorized
	}

	if req.Title != "" {
		homeworkRes.Title = req.Title
	}

	if req.Description != "" {
		homeworkRes.Description = &req.Description
	}

	if req.ClassID != "" && req.ClassID != homeworkRes.ClassID {
		if err := u.ensureOwnedClass(ctx, userId, req.ClassID); err != nil {
			return err
		}
		homeworkRes.ClassID = req.ClassID
	}

	if req.DueDate != "" {
		dueDate, err := time.Parse("2006-01-02", req.DueDate)
		if err != nil {
			return homework.ErrInvalidDueDate
		}
		dueDateUtc := dueDate.UTC()
		homeworkRes.DueDate = dueDateUtc
	}

	if len(req.Questions) > 0 && !sameStrings(req.Questions, homeworkRes.Questions) {
		inUse, err := u.homeworkHasSubmissions(ctx, homeworkRes.ID.Hex())
		if err != nil {
			return err
		}
		if inUse {
			return homework.ErrHomeworkInUse
		}
		if err := u.ensureQuestionsExist(ctx, req.Questions); err != nil {
			return err
		}
		homeworkRes.Questions = req.Questions
	}

	homeworkRes.UpdatedAt = time.Now().UTC()

	if err := u.homeworkRepo.UpdateHomeworkById(ctx, objectId, homeworkRes); err != nil {
		return err
	}

	return nil
}

func (u *homeworkService) DeleteHomeworkById(ctx context.Context, userId string, id string) error {
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	homeworkRes, err := u.homeworkRepo.GetHomeworkById(ctx, objectId)
	if err != nil {
		return err
	}

	if homeworkRes.CreatedBy != userId {
		return homework.ErrHomeworkNotAuthorized
	}

	inUse, err := u.homeworkHasSubmissions(ctx, homeworkRes.ID.Hex())
	if err != nil {
		return err
	}
	if inUse {
		return homework.ErrHomeworkInUse
	}

	if err := u.homeworkRepo.DeleteHomeworkById(ctx, objectId); err != nil {
		return err
	}

	return nil
}

func (u *homeworkService) GetHomeworksByClassId(ctx context.Context, userId string, classId string, params pagination.Params) (*response.HomeworkResponseWithMeta, error) {
	homeworks, total, err := u.homeworkRepo.GetHomeworks(ctx, userId, classId, params)
	if err != nil {
		return nil, err
	}

	var questionIdsSet = map[string]bool{}
	for _, homework := range homeworks {
		for _, questionId := range homework.Questions {
			questionIdsSet[questionId] = true
		}
	}

	questionIds := make([]string, 0, len(questionIdsSet))
	for questionId := range questionIdsSet {
		questionIds = append(questionIds, questionId)
	}

	questionIdsObjectIDs := make([]primitive.ObjectID, 0, len(questionIds))
	for _, questionId := range questionIds {
		questionIdObjectID, err := primitive.ObjectIDFromHex(questionId)
		if err != nil {
			return nil, err
		}
		questionIdsObjectIDs = append(questionIdsObjectIDs, questionIdObjectID)
	}

	var questions []*question.Question
	if len(questionIdsObjectIDs) > 0 {
		questions, err = u.questionRepo.GetQuestionByIds(ctx, questionIdsObjectIDs)
		if err != nil {
			return nil, err
		}
	}

	return &response.HomeworkResponseWithMeta{
		Homeworks: mapper.MapHomeworksToResponses(homeworks, mapper.QuestionResponseMap(questions)),
		Meta:      pagination.NewMeta(params, total),
	}, nil
}

func (u *homeworkService) ensureOwnedClass(ctx context.Context, userId, classID string) error {
	objectId, err := primitive.ObjectIDFromHex(classID)
	if err != nil {
		return homework.ErrInvalidClassID
	}

	item, err := u.classRepo.GetClassById(ctx, objectId)
	if err != nil {
		return homework.ErrInvalidClassID
	}
	if item.CreatedBy != userId {
		return homework.ErrHomeworkNotAuthorized
	}
	return nil
}

func (u *homeworkService) ensureQuestionsExist(ctx context.Context, questionIDs []string) error {
	unique := uniqueStrings(questionIDs)
	objectIDs := make([]primitive.ObjectID, 0, len(unique))
	for _, id := range unique {
		oid, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			return homework.ErrInvalidQuestions
		}
		objectIDs = append(objectIDs, oid)
	}

	questions, err := u.questionRepo.GetQuestionByIds(ctx, objectIDs)
	if err != nil {
		return err
	}
	if len(questions) != len(unique) {
		return homework.ErrInvalidQuestions
	}
	return nil
}

func (u *homeworkService) homeworkHasSubmissions(ctx context.Context, homeworkID string) (bool, error) {
	count, err := u.submissionRepo.CountByHomeworkID(ctx, homeworkID)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func uniqueStrings(ids []string) []string {
	seen := make(map[string]struct{}, len(ids))
	out := make([]string, 0, len(ids))
	for _, id := range ids {
		if _, ok := seen[id]; ok {
			continue
		}
		seen[id] = struct{}{}
		out = append(out, id)
	}
	return out
}
