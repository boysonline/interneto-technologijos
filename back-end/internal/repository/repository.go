package repository

import (
	"back-end/internal/model"
	"context"

	"go.mongodb.org/mongo-driver/mongo"
)

type Repository struct {
	Collection *mongo.Collection
}

func (repository *Repository) CreateLecture(lecture model.Lecture) error {
	_, err := repository.Collection.InsertOne(context.TODO(), lecture)
	return err
}

func (repository *Repository) CreateAssignment(assignment model.Assignment) error {
	_, err := repository.Collection.InsertOne(context.TODO(), assignment)
	return err
}
