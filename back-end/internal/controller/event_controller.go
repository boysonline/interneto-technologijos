package controller

import (
	"back-end/internal/model"
	"back-end/internal/repository"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type EventController struct {
	Repository *repository.EventRepository
}

func NewEventController(repo *repository.EventRepository) *EventController {
	return &EventController{Repository: repo}
}

// GET /events?start=2024-05-01T00:00:00Z&end=2024-06-01T00:00:00Z
func (controller *EventController) GetEvents(ginContext *gin.Context) {
	start := time.Now().AddDate(0, 0, -1)
	end := time.Now().AddDate(0, 0, 30)

	if startQuery := ginContext.Query("start"); startQuery != "" {
		parsedStart, err := time.Parse(time.RFC3339, startQuery)
		if err != nil {
			ginContext.JSON(http.StatusBadRequest, gin.H{"error": "Invalid start date format. Use 2024-05-01T00:00:00Z"})
			return
		}
		start = parsedStart
	}

	if endQuery := ginContext.Query("end"); endQuery != "" {
		parsedEnd, err := time.Parse(time.RFC3339, endQuery)
		if err != nil {
			ginContext.JSON(http.StatusBadRequest, gin.H{"error": "Invalid end date format. Use 2024-05-01T00:00:00Z"})
			return
		}
		end = parsedEnd
	}

	lectures, err := controller.Repository.GetLectures(ginContext.Request.Context(), start, end)
	if err != nil {
		ginContext.JSON(http.StatusInternalServerError, gin.H{"error": "Nepavyko užkrauti paskaitų"})
		return
	}

	assignments, err := controller.Repository.GetAssignments(ginContext.Request.Context(), start, end)
	if err != nil {
		ginContext.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ginContext.JSON(http.StatusOK, gin.H{
		"lectures":    lectures,
		"assignments": assignments,
	})
}

// GET /lectures/:id
func (controller *EventController) GetLecture(ginContext *gin.Context) {
	id, ok := ParseID(ginContext)
	if !ok {
		return
	}

	lecture, err := controller.Repository.GetLecture(ginContext.Request.Context(), id)
	if err != nil {
		ginContext.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	ginContext.JSON(http.StatusOK, lecture)
}

// POST /lectures
type LectureRequest struct {
	Title         string             `json:"title" binding:"required"`
	Type          string             `json:"type" binding:"required"` // e.g., "lecture", "practice"
	SubjectID     primitive.ObjectID `json:"subject_id" binding:"required"`
	Start         time.Time          `json:"start" binding:"required"`
	DurationMins  int                `json:"duration_mins" binding:"required"`
	Auditorium    *string            `json:"auditorium"`
	CourseURL     *string            `json:"course_url"`
	LectureURL    *string            `json:"lecture_url"`
	AttendanceURL *string            `json:"attendance_url"`
	Intervals     []int              `json:"intervals"` // e.g., [7] for weekly
	Until         time.Time          `json:"until" binding:"required"`
}

func (controller *EventController) CreateLectures(ginContext *gin.Context) {
	var request LectureRequest
	if err := ginContext.ShouldBindJSON(&request); err != nil {
		ginContext.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if !model.IsValidLectureType(request.Type) {
		ginContext.JSON(http.StatusBadRequest, gin.H{
			"error":       "Netinkamas paskaitos tipas",
			"valid_types": model.GetLectureTypes(),
		})
		return
	}

	duration := time.Duration(request.DurationMins) * time.Minute

	lectures, err := model.NewLectures(
		request.Title,
		request.Type,
		request.SubjectID,
		request.Start,
		duration,
		request.Auditorium,
		request.CourseURL,
		request.LectureURL,
		request.AttendanceURL,
		request.Intervals,
		request.Until,
	)

	if err != nil {
		ginContext.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	created, err := controller.Repository.CreateLectures(ginContext.Request.Context(), lectures)
	if err != nil {
		ginContext.JSON(http.StatusInternalServerError, gin.H{"error": "Nepavyko sukurti paskaitų"})
		return
	}

	ginContext.JSON(http.StatusCreated, created)
}

// PATCH /lectures/:id
func (controller *EventController) UpdateLecture(ginContext *gin.Context) {
	id, ok := ParseID(ginContext)
	if !ok {
		return
	}

	var input struct {
		Scope string             `json:"scope"` // "single" or "series"
		Data  model.LectureEvent `json:"data"`
	}

	if err := ginContext.ShouldBindJSON(&input); err != nil {
		ginContext.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if input.Scope == "single" {
		updated, err := controller.Repository.DetachLecture(ginContext.Request.Context(), id, input.Data)
		if err != nil {
			ginContext.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update single lecture"})
			return
		}
		ginContext.JSON(http.StatusOK, updated)
		return
	}

	if input.Scope == "series" {
		if input.Data.GroupID.IsZero() {
			ginContext.JSON(http.StatusBadRequest, gin.H{"error": "group_id is required for series update"})
			return
		}

		updatedList, err := controller.Repository.UpdateLectures(ginContext.Request.Context(), input.Data.GroupID, input.Data)
		if err != nil {
			ginContext.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		ginContext.JSON(http.StatusOK, updatedList)
		return
	}

	ginContext.JSON(http.StatusBadRequest, gin.H{"error": "invalid scope (must be 'single' or 'series')"})
}

// DELETE /lectures/group/:group_id
func (controller *EventController) DeleteLectures(ginContext *gin.Context) {
	groupID, err := primitive.ObjectIDFromHex(ginContext.Param("group_id"))
	if err != nil {
		ginContext.JSON(http.StatusBadRequest, gin.H{"error": "invalid group ID"})
		return
	}

	if err := controller.Repository.DeleteLectures(ginContext.Request.Context(), groupID); err != nil {
		ginContext.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ginContext.JSON(http.StatusOK, gin.H{"message": "lecture series deleted"})
}

// GET /assignments/:id
func (controller *EventController) GetAssignment(ginContext *gin.Context) {
	id, ok := ParseID(ginContext)
	if !ok {
		return
	}

	assignment, err := controller.Repository.GetAssignment(ginContext.Request.Context(), id)
	if err != nil {
		ginContext.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	ginContext.JSON(http.StatusOK, assignment)
}

// POST /assignments
type AssignmentRequest struct {
	Title     string              `json:"title" binding:"required,min=3,max=50"`
	SubjectID primitive.ObjectID  `json:"subject_id" binding:"required"`
	Due       time.Time           `json:"due" binding:"required"`
	DefenceID *primitive.ObjectID `json:"defence_id"`
}

func (controller *EventController) CreateAssignment(ginContext *gin.Context) {
	var request AssignmentRequest

	if err := ginContext.ShouldBindJSON(&request); err != nil {
		ginContext.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var defenceID primitive.ObjectID
	if request.DefenceID != nil {
		defenceID = *request.DefenceID
	}
	assignment := model.NewAssignment(
		request.Title,
		request.SubjectID,
		request.Due,
		defenceID,
	)

	created, err := controller.Repository.CreateAssignment(ginContext.Request.Context(), assignment)
	if err != nil {
		ginContext.JSON(http.StatusInternalServerError, gin.H{"error": "Nepavyko sukurti užduoties"})
		return
	}

	ginContext.JSON(http.StatusCreated, created)
}

// PATCH /assignments/:id
func (controller *EventController) UpdateAssignment(ginContext *gin.Context) {
	id, ok := ParseID(ginContext)
	if !ok {
		return
	}

	var assignment model.AssignmentEvent
	if err := ginContext.ShouldBindJSON(&assignment); err != nil {
		ginContext.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	assignment.ID = id

	updated, err := controller.Repository.UpdateAssignment(ginContext.Request.Context(), assignment)
	if err != nil {
		ginContext.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ginContext.JSON(http.StatusOK, updated)
}

// DELETE /events/:id
func (controller *EventController) DeleteEvent(ginContext *gin.Context) {
	id, ok := ParseID(ginContext)
	if !ok {
		return
	}

	if err := controller.Repository.DeleteEvent(ginContext.Request.Context(), id); err != nil {
		ginContext.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ginContext.JSON(http.StatusOK, gin.H{"message": "deleted successfully"})
}

func ParseID(ginContext *gin.Context) (primitive.ObjectID, bool) {
	id, err := primitive.ObjectIDFromHex(ginContext.Param("id"))
	if err != nil {
		ginContext.JSON(http.StatusBadRequest, gin.H{"error": "Netinkamas ID"})
		return primitive.NilObjectID, false
	}
	return id, true
}
