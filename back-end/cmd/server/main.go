package main

import (
	"back-end/internal/database"
	"back-end/internal/model"
	"time"
)

func main() {

	database.Connect()

	// Example: Create a test lecture
	newLecture := model.Lecture{
		CalendarEvent: model.CalendarEvent{
			Title:     "TEST LECTURE",
			StartTime: time.Now(),
			EndTime:   time.Now().Add(time.Hour * 2),
		},
		Type:       "Lecture",
		Auditorium: "TEST 0",
	}

	// Save it!
	database.CreateLecture(newLecture)
}

// myliu tave
