package mongodb

import (
	"context"
	"errors"
	"time"

	"teaching_assistant/internal/domain/class"
	"teaching_assistant/internal/domain/student"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type studentRepository struct {
	collection         *mongo.Collection
	guardianCollection *mongo.Collection
}

func NewStudentRepository(db *mongo.Database) student.StudentRepository {
	return &studentRepository{
		collection:         db.Collection("students"),
		guardianCollection: db.Collection("guardians"),
	}
}

func (r *studentRepository) CreateMany(ctx context.Context, students []*student.Student) error {
	if len(students) == 0 {
		return nil
	}

	docs := make([]interface{}, 0, len(students))
	for _, student := range students {
		docs = append(docs, student)
	}

	_, err := r.collection.InsertMany(ctx, docs)
	if mongo.IsDuplicateKeyError(err) {
		return class.ErrStudentCodeExists
	}
	return err
}

func (r *studentRepository) GetStudentById(ctx context.Context, id primitive.ObjectID) (*student.Student, error) {
	filter := bson.M{"_id": id}
	var student student.Student
	err := r.collection.FindOne(ctx, filter).Decode(&student)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &student, nil
}

func (r *studentRepository) GetStudentsByIds(ctx context.Context, ids []primitive.ObjectID) ([]*student.Student, error) {
	filter := bson.M{"_id": bson.M{"$in": ids}}
	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	students := make([]*student.Student, 0)
	for cursor.Next(ctx) {
		var student student.Student
		if err := cursor.Decode(&student); err != nil {
			return nil, err
		}
		students = append(students, &student)
	}
	return students, nil
}

func (r *studentRepository) GetStudentByCode(ctx context.Context, code string) (*student.Student, error) {
	filter := bson.M{"code": code}
	var student student.Student
	err := r.collection.FindOne(ctx, filter).Decode(&student)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &student, nil
}

func (r *studentRepository) CreateGuardian(ctx context.Context, guardian *student.Guardian) error {
	_, err := r.guardianCollection.InsertOne(ctx, guardian)
	if mongo.IsDuplicateKeyError(err) {
		return errors.New(string(student.ErrGuardianAlreadyExists))
	}
	return err
}

func (r *studentRepository) GetGuardianByParentId(ctx context.Context, parentId string) (*student.Guardian, error) {
	filter := bson.M{"parent_id": parentId}
	var guardian student.Guardian
	err := r.guardianCollection.FindOne(ctx, filter).Decode(&guardian)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &guardian, nil
}

func (r *studentRepository) DeactivateByIDs(ctx context.Context, ids []primitive.ObjectID) error {
	return r.updateStatusByIDs(ctx, ids, student.StudentStatusInactive)
}

func (r *studentRepository) ActivateByIDs(ctx context.Context, ids []primitive.ObjectID) error {
	return r.updateStatusByIDs(ctx, ids, student.StudentStatusActive)
}

func (r *studentRepository) updateStatusByIDs(ctx context.Context, ids []primitive.ObjectID, status student.StudentStatus) error {
	if len(ids) == 0 {
		return nil
	}

	_, err := r.collection.UpdateMany(
		ctx,
		bson.M{"_id": bson.M{"$in": ids}},
		bson.M{"$set": bson.M{"status": status, "updated_at": time.Now()}},
	)
	return err
}

func (r *studentRepository) GetGuardiansByStudentIds(ctx context.Context, studentIds []string) ([]*student.Guardian, error) {
	if len(studentIds) == 0 {
		return []*student.Guardian{}, nil
	}

	filter := bson.M{"student_id": bson.M{"$in": studentIds}}
	cursor, err := r.guardianCollection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	guardians := make([]*student.Guardian, 0)
	for cursor.Next(ctx) {
		var guardian student.Guardian
		if err := cursor.Decode(&guardian); err != nil {
			return nil, err
		}
		guardians = append(guardians, &guardian)
	}
	return guardians, nil
}
