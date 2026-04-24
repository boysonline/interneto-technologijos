package repository

import (
	"back-end/internal/model"
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type EventRepository struct {
	Events *mongo.Collection
}

func NewEventRepository(db *mongo.Database) *EventRepository {
	return &EventRepository{
		Events: db.Collection("events"),
	}
}

// LECTURES

func (repository *EventRepository) GetLecture(ctx context.Context, id primitive.ObjectID) (*model.LectureEvent, error) {
	var lecture model.LectureEvent
	err := repository.Events.FindOne(ctx, bson.M{"_id": id, "type": bson.M{"$in": model.GetLectureTypes()}}).Decode(&lecture)
	if err != nil {
		return nil, err
	}
	return &lecture, nil
}

func (repository *EventRepository) GetLectures(ctx context.Context, start, end time.Time) ([]model.LectureEvent, error) {
	var lectures []model.LectureEvent
	filter := bson.M{
		"type":       bson.M{"$in": model.GetLectureTypes()},
		"start_time": bson.M{"$gte": start, "$lte": end},
	}
	opts := options.Find().SetSort(bson.D{{Key: "start_time", Value: 1}})

	cursor, err := repository.Events.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	err = cursor.All(ctx, &lectures)
	return lectures, err
}

func (repository *EventRepository) CreateLectures(ctx context.Context, lectures []model.LectureEvent) ([]model.LectureEvent, error) {
	var docs []interface{}
	for i := range lectures {
		if lectures[i].ID.IsZero() {
			lectures[i].ID = primitive.NewObjectID()
		}
		docs = append(docs, lectures[i])
	}
	_, err := repository.Events.InsertMany(ctx, docs)
	return lectures, err
}

func (repository *EventRepository) UpdateLectures(ctx context.Context, groupID primitive.ObjectID, updates model.LectureEvent) ([]model.LectureEvent, error) {
	updateMap, err := ToSetMap(updates)
	if err != nil {
		return nil, err
	}

	delete(updateMap, "start_time")
	delete(updateMap, "end_time")
	delete(updateMap, "group_id")

	_, err = repository.Events.UpdateMany(ctx, bson.M{"group_id": groupID}, bson.M{"$set": updateMap})
	if err != nil {
		return nil, err
	}

	return repository.getLecturesByGroup(ctx, groupID)
}

func (repository *EventRepository) DetachLecture(ctx context.Context, id primitive.ObjectID, lecture model.LectureEvent) (*model.LectureEvent, error) {
	updateMap, err := ToSetMap(lecture)
	if err != nil {
		return nil, err
	}

	updateMap["group_id"] = primitive.NewObjectID()

	var updated model.LectureEvent
	err = repository.Events.FindOneAndUpdate(
		ctx,
		bson.M{"_id": id},
		bson.M{"$set": updateMap},
		options.FindOneAndUpdate().SetReturnDocument(options.After),
	).Decode(&updated)

	return &updated, err
}

// ASSIGNMENTS

func (repository *EventRepository) GetAssignment(ctx context.Context, id primitive.ObjectID) (*model.AssignmentEvent, error) {
	var assignment model.AssignmentEvent
	err := repository.Events.FindOne(ctx, bson.M{"_id": id, "type": "assignment"}).Decode(&assignment)
	if err != nil {
		return nil, err
	}
	return &assignment, nil
}

func (repository *EventRepository) GetAssignments(ctx context.Context, start, end time.Time) ([]model.AssignmentEvent, error) {
	var assignments []model.AssignmentEvent
	filter := bson.M{
		"type":       "assignment",
		"start_time": bson.M{"$gte": start, "$lte": end},
	}
	opts := options.Find().SetSort(bson.D{{Key: "start_time", Value: 1}})

	cursor, err := repository.Events.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	err = cursor.All(ctx, &assignments)
	return assignments, err
}

func (repository *EventRepository) CreateAssignment(ctx context.Context, assignment model.AssignmentEvent) (*model.AssignmentEvent, error) {
	_, err := repository.Events.InsertOne(ctx, assignment)
	if err != nil {
		return nil, err
	}
	return &assignment, nil
}

func (repository *EventRepository) UpdateAssignment(ctx context.Context, assignment model.AssignmentEvent) (*model.AssignmentEvent, error) {
	update, err := ToSetMap(assignment)
	if err != nil {
		return nil, err
	}

	result, err := repository.Events.UpdateOne(
		ctx,
		bson.M{"_id": assignment.ID},
		bson.M{"$set": update},
	)

	if err != nil {
		return nil, err
	}

	if result.MatchedCount == 0 {
		return nil, errors.New("užduotis nerasta")
	}

	return &assignment, nil
}

func (repository *EventRepository) DeleteEvent(ctx context.Context, id primitive.ObjectID) error {
	_, err := repository.Events.DeleteOne(ctx, bson.M{"_id": id})
	return err
}

func (repository *EventRepository) DeleteLectures(ctx context.Context, groupID primitive.ObjectID) error {
	_, err := repository.Events.DeleteMany(ctx, bson.M{"group_id": groupID})
	return err
}

func (repository *EventRepository) getLecturesByGroup(ctx context.Context, groupID primitive.ObjectID) ([]model.LectureEvent, error) {
	var lectures []model.LectureEvent
	cursor, err := repository.Events.Find(ctx, bson.M{"group_id": groupID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	err = cursor.All(ctx, &lectures)
	return lectures, err
}
