package model

import (
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Subject struct {
	ID    primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Title string             `bson:"title" json:"title" binding:"required"`
}

type CalendarEvent struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	GroupID   primitive.ObjectID `bson:"group_id,omitempty" json:"group_id"`
	SubjectID primitive.ObjectID `bson:"subject_id" json:"subject_id" binding:"required"`
	Title     string             `bson:"title" json:"title" binding:"required"`
	StartTime time.Time          `bson:"start_time" json:"start_time" binding:"required"`
	EndTime   time.Time          `bson:"end_time" json:"end_time" binding:"required"`
	Type      string             `bson:"type" json:"type" binding:"required"`
}

type LectureEvent struct {
	CalendarEvent `bson:",inline"`
	Auditorium    *string `bson:"auditorium,omitempty" json:"auditorium,omitempty"`
	CourseURL     *string `bson:"course_url,omitempty" json:"course_url,omitempty"`
	LectureURL    *string `bson:"lecture_url,omitempty" json:"lecture_url,omitempty"`
	AttendanceURL *string `bson:"attendance_url,omitempty" json:"attendance_url,omitempty"`
}

type AssignmentEvent struct {
	CalendarEvent `bson:",inline"`
	DefenceID     primitive.ObjectID `bson:"defence_id,omitempty" json:"defence_id,omitempty"`
}

// SUBJECT
func NewSubject(title string) Subject {
	return Subject{
		ID:    primitive.NewObjectID(),
		Title: title,
	}
}

// LECTURE
const (
	TypeLecture      = "lecture"
	TypeConsultation = "consultation"
	TypeTest         = "test"
	TypeExam         = "exam"
	TypeDefence      = "defence"
)

func GetLectureTypes() []string {
	return []string{
		TypeLecture,
		TypeConsultation,
		TypeTest,
		TypeExam,
		TypeDefence,
	}
}

func IsValidLectureType(t string) bool {
	for _, valid := range GetLectureTypes() {
		if t == valid {
			return true
		}
	}
	return false
}

func NewLectures(title string, lectureType string, subjectID primitive.ObjectID, start time.Time, duration time.Duration,
	auditorium *string, courseURL *string, lectureURL *string, attendanceURL *string, intervals []int, until time.Time) ([]LectureEvent, error) {
	groupID := primitive.NewObjectID()

	lectures := make([]LectureEvent, 0)
	if until.Before(start) {
		return nil, errors.New("the 'until' date must be after the start date")
	}
	if duration <= 0 {
		return nil, errors.New("lecture duration must be positive")
	}
	if len(intervals) == 0 {
		return append(lectures, NewLecture(groupID, title, subjectID, start, duration, auditorium, courseURL, lectureURL, attendanceURL, lectureType)), nil
	}

	date := start
	for !date.After(until) {
		for _, interval := range intervals {
			lecture := NewLecture(groupID, title, subjectID, date, duration, auditorium, courseURL, lectureURL, attendanceURL, lectureType)
			lectures = append(lectures, lecture)
			date = date.AddDate(0, 0, interval)
			if date.After(until) {
				break
			}
		}
	}

	return lectures, nil
}

func NewLecture(groupID primitive.ObjectID, title string, subjectID primitive.ObjectID, start time.Time, duration time.Duration,
	auditorium *string, courseURL *string, lectureURL *string, attendanceURL *string, lectureType string) LectureEvent {

	lecture := LectureEvent{
		CalendarEvent: CalendarEvent{
			ID:        primitive.NewObjectID(),
			GroupID:   groupID,
			SubjectID: subjectID,
			Title:     title,
			StartTime: start,
			EndTime:   start.Add(duration),
			Type:      lectureType,
		},
		Auditorium:    auditorium,
		CourseURL:     courseURL,
		LectureURL:    lectureURL,
		AttendanceURL: attendanceURL,
	}
	return lecture
}

// ASSIGNMENT
func NewAssignment(title string, subjectID primitive.ObjectID, due time.Time, defenceID primitive.ObjectID) AssignmentEvent {
	assignment := AssignmentEvent{
		CalendarEvent: CalendarEvent{
			ID:        primitive.NewObjectID(),
			GroupID:   primitive.NewObjectID(),
			SubjectID: subjectID,
			Title:     title,
			StartTime: due,
			EndTime:   due,
			Type:      "assignment",
		},

		DefenceID: defenceID,
	}
	return assignment
}
