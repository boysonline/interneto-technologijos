package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Subject struct {
	ID    primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Title string             `bson:"title" json:"title"`
}

type CalendarEvent struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Title     string             `bson:"title" json:"title"`
	StartTime time.Time          `bson:"start_time" json:"start_time"`
	EndTime   time.Time          `bson:"end_time" json:"end_time"`
}

type Lecture struct {
	CalendarEvent `bson:",inline"`
	Subject       Subject `bson:"subject" json:"subject"`
	Auditorium    string  `bson:"auditorium" json:"auditorium"`
	LectureURL    string  `bson:"lecture_url" json:"lecture_url"`
	AttendanceURL string  `bson:"attendance_url" json:"attendance_url"`
	Type          string  `bson:"type" json:"type"` // Lecture, Consultation, Test, Exam, Defence, ...
}

type Assignment struct {
	CalendarEvent `bson:",inline"`
	Subject       Subject  `bson:"subject" json:"subject"`
	Defence       *Lecture `bson:"defence,omitempty" json:"defence,omitempty"`
}
