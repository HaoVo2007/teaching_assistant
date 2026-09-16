package database

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func EnsureIndexes(ctx context.Context, db *mongo.Database) error {
	indexes := map[string][]mongo.IndexModel{
		"users": {
			{
				Keys:    bson.D{{Key: "email", Value: 1}},
				Options: options.Index().SetUnique(true).SetName("uniq_user_email"),
			},
		},
		"questions": {
			{
				Keys: bson.D{
					{Key: "created_by", Value: 1},
					{Key: "created_at", Value: -1},
				},
				Options: options.Index().SetName("idx_question_created_by_created_at"),
			},
			{
				Keys: bson.D{
					{Key: "created_by", Value: 1},
					{Key: "type", Value: 1},
					{Key: "created_at", Value: -1},
				},
				Options: options.Index().SetName("idx_question_created_by_type_created_at"),
			},
			{
				Keys: bson.D{
					{Key: "created_by", Value: 1},
					{Key: "subject", Value: 1},
					{Key: "created_at", Value: -1},
				},
				Options: options.Index().SetName("idx_question_created_by_subject_created_at"),
			},
			{
				Keys: bson.D{
					{Key: "created_by", Value: 1},
					{Key: "grade", Value: 1},
					{Key: "created_at", Value: -1},
				},
				Options: options.Index().SetName("idx_question_created_by_grade_created_at"),
			},
		},
		"question_sets": {
			{
				Keys: bson.D{
					{Key: "created_by", Value: 1},
					{Key: "created_at", Value: -1},
				},
				Options: options.Index().SetName("idx_question_set_created_by_created_at"),
			},
			{
				Keys: bson.D{
					{Key: "created_by", Value: 1},
					{Key: "question_type", Value: 1},
					{Key: "created_at", Value: -1},
				},
				Options: options.Index().SetName("idx_question_set_created_by_type_created_at"),
			},
			{
				Keys:    bson.D{{Key: "question_ids", Value: 1}},
				Options: options.Index().SetName("idx_question_set_question_ids"),
			},
		},
		"classes": {
			{
				Keys: bson.D{
					{Key: "created_by", Value: 1},
					{Key: "created_at", Value: -1},
				},
				Options: options.Index().SetName("idx_class_created_by_created_at"),
			},
		},
		"students": {
			{
				Keys: bson.D{
					{Key: "created_by", Value: 1},
					{Key: "code", Value: 1},
				},
				Options: options.Index().SetUnique(true).SetName("uniq_student_code_per_teacher"),
			},
			{
				Keys:    bson.D{{Key: "code", Value: 1}},
				Options: options.Index().SetName("idx_student_code"),
			},
			{
				Keys:    bson.D{{Key: "class_id", Value: 1}},
				Options: options.Index().SetName("idx_student_class_id"),
			},
		},
		"guardians": {
			{
				Keys: bson.D{
					{Key: "parent_id", Value: 1},
					{Key: "student_id", Value: 1},
				},
				Options: options.Index().SetUnique(true).SetName("uniq_guardian_parent_id_student_id"),
			},
			{
				Keys:    bson.D{{Key: "parent_id", Value: 1}},
				Options: options.Index().SetName("idx_guardian_parent_id"),
			},
			{
				Keys:    bson.D{{Key: "student_id", Value: 1}},
				Options: options.Index().SetName("idx_guardian_student_id"),
			},
		},
		"homeworks": {
			{
				Keys: bson.D{
					{Key: "created_by", Value: 1},
					{Key: "created_at", Value: -1},
				},
				Options: options.Index().SetName("idx_homework_created_by_created_at"),
			},
			{
				Keys: bson.D{
					{Key: "created_by", Value: 1},
					{Key: "class_id", Value: 1},
					{Key: "created_at", Value: -1},
				},
				Options: options.Index().SetName("idx_homework_created_by_class_created_at"),
			},
			{
				Keys:    bson.D{{Key: "class_id", Value: 1}},
				Options: options.Index().SetName("idx_homework_class_id"),
			},
			{
				Keys:    bson.D{{Key: "questions", Value: 1}},
				Options: options.Index().SetName("idx_homework_questions"),
			},
		},
		"homework_submissions": {
			{
				Keys: bson.D{
					{Key: "teacher_id", Value: 1},
					{Key: "created_at", Value: -1},
				},
				Options: options.Index().SetName("idx_submission_teacher_created_at"),
			},
			{
				Keys: bson.D{
					{Key: "homework_id", Value: 1},
					{Key: "teacher_id", Value: 1},
					{Key: "created_at", Value: -1},
				},
				Options: options.Index().SetName("idx_submission_homework_teacher_created_at"),
			},
			{
				Keys:    bson.D{{Key: "homework_id", Value: 1}},
				Options: options.Index().SetName("idx_submission_homework_id"),
			},
			{
				Keys:    bson.D{{Key: "student_answers.question_id", Value: 1}},
				Options: options.Index().SetName("idx_submission_answer_question_id"),
			},
		},
	}

	for collection, models := range indexes {
		if _, err := db.Collection(collection).Indexes().CreateMany(ctx, models); err != nil {
			return err
		}
	}
	return nil
}
