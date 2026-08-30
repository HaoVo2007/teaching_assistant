package mongodb

import (
	"context"
	homeworksubmission "teaching_assistant/internal/domain/homework_submission"
	"teaching_assistant/pkg/pagination"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type homeworkSubmissionRepository struct {
	collection *mongo.Collection
}

func NewHomeworkSubmissionRepository(db *mongo.Database) homeworksubmission.HomeworkSubmissionRepository {
	return &homeworkSubmissionRepository{
		collection: db.Collection("homework_submissions"),
	}
}

func (r *homeworkSubmissionRepository) CreateHomeworkSubmission(ctx context.Context, submission *homeworksubmission.HomeworkSubmission) error {
	_, err := r.collection.InsertOne(ctx, submission)
	if err != nil {
		return err
	}
	return nil
}

func (r *homeworkSubmissionRepository) GetHomeworkSubmissionsByUserId(ctx context.Context, userId string, params pagination.Params) ([]*homeworksubmission.HomeworkSubmission, int64, error) {
	filter := bson.M{
		"teacher_id": userId,
	}

	total, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	options := options.Find().SetSort(bson.M{
		"created_at": -1,
	}).SetSkip(params.Skip()).SetLimit(params.Limit64())

	cursor, err := r.collection.Find(ctx, filter, options)
	if err != nil {
		return nil, 0, err
	}

	defer cursor.Close(ctx)

	homeworkSubmissions := make([]*homeworksubmission.HomeworkSubmission, 0)
	for cursor.Next(ctx) {
		var submission homeworksubmission.HomeworkSubmission
		if err := cursor.Decode(&submission); err != nil {
			return nil, 0, err
		}
		homeworkSubmissions = append(homeworkSubmissions, &submission)
	}

	return homeworkSubmissions, total, nil

}

func (r *homeworkSubmissionRepository) GetHomeworkSubmissionById(ctx context.Context, id primitive.ObjectID) (*homeworksubmission.HomeworkSubmission, error) {
	filter := bson.M{
		"_id": id,
	}

	var submission homeworksubmission.HomeworkSubmission
	if err := r.collection.FindOne(ctx, filter).Decode(&submission); err != nil {
		return nil, err
	}

	return &submission, nil
}

func (r *homeworkSubmissionRepository) GetHomeworkSubmissionsByHomeworkId(ctx context.Context, homeworkID string, userId string, params pagination.Params) ([]*homeworksubmission.HomeworkSubmission, int64, error) {
	filter := bson.M{
		"homework_id": homeworkID,
		"teacher_id":  userId,
	}

	total, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	options := options.Find().SetSort(bson.M{
		"created_at": -1,
	}).SetSkip(params.Skip()).SetLimit(params.Limit64())

	cursor, err := r.collection.Find(ctx, filter, options)
	if err != nil {
		return nil, 0, err
	}

	defer cursor.Close(ctx)

	submissions := make([]*homeworksubmission.HomeworkSubmission, 0)
	for cursor.Next(ctx) {
		var submission homeworksubmission.HomeworkSubmission
		if err := cursor.Decode(&submission); err != nil {
			return nil, 0, err
		}
		submissions = append(submissions, &submission)
	}

	return submissions, total, nil
}

func (r *homeworkSubmissionRepository) CountByHomeworkID(ctx context.Context, homeworkID string) (int64, error) {
	return r.collection.CountDocuments(ctx, bson.M{"homework_id": homeworkID})
}

func (r *homeworkSubmissionRepository) CountByQuestionID(ctx context.Context, questionID string) (int64, error) {
	return r.collection.CountDocuments(ctx, bson.M{"student_answers.question_id": questionID})
}
