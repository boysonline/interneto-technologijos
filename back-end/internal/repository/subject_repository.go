package repository

import (
	"back-end/internal/model"
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type SubjectRepository struct {
	Subjects *mongo.Collection
}

func NewSubjectRepository(db *mongo.Database) *SubjectRepository {
	return &SubjectRepository{
		Subjects: db.Collection("subjects"),
	}
}
func (repository *SubjectRepository) GetAll(ctx context.Context) ([]model.Subject, error) {
	var subjects []model.Subject
	cursor, err := repository.Subjects.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	err = cursor.All(ctx, &subjects)
	return subjects, err
}

func (repository *SubjectRepository) GetSubject(ctx context.Context, id primitive.ObjectID) (model.Subject, error) {
	var subject model.Subject
	err := repository.Subjects.FindOne(ctx, bson.M{"_id": id}).Decode(&subject)
	return subject, err
}

func (repository *SubjectRepository) CreateSubject(ctx context.Context, subject model.Subject) (*model.Subject, error) {
	result, err := repository.Subjects.InsertOne(ctx, subject)
	if err != nil {
		return nil, err
	}
	if oid, ok := result.InsertedID.(primitive.ObjectID); ok {
		subject.ID = oid
	}

	return &subject, nil
}

func (repository *SubjectRepository) UpdateSubject(ctx context.Context, subject model.Subject) (*model.Subject, error) {
	update, err := ToSetMap(subject)
	if err != nil {
		return nil, err
	}

	result, err := repository.Subjects.UpdateOne(
		ctx,
		bson.M{"_id": subject.ID},
		bson.M{"$set": update},
	)

	if err != nil {
		return nil, err
	}

	if result.MatchedCount == 0 {
		return nil, errors.New("dokumentas nerastas")
	}

	return &subject, nil
}

func (repository *SubjectRepository) DeleteSubject(ctx context.Context, id primitive.ObjectID) error {
	_, err := repository.Subjects.DeleteOne(ctx, bson.M{"_id": id})
	return err
}

func ToSetMap(doc interface{}) (bson.M, error) {
	data, err := bson.Marshal(doc)
	if err != nil {
		return nil, err
	}
	var m bson.M
	if err := bson.Unmarshal(data, &m); err != nil {
		return nil, err
	}
	delete(m, "_id")
	return m, nil
}
