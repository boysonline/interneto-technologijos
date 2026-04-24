package routes

import (
	"back-end/internal/controller"
	"back-end/internal/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine, subjectController *controller.SubjectController, eventController *controller.EventController) {
	api := r.Group("/api")
	api.Use(middleware.APIKeyMiddleware())
	{
		subjects := api.Group("/subjects")
		{
			subjects.GET("/", subjectController.GetSubjects)
			subjects.GET("/:id", subjectController.GetSubject)
			subjects.POST("/", subjectController.CreateSubject)
			subjects.PATCH("/:id", subjectController.UpdateSubject)
			subjects.DELETE("/:id", subjectController.DeleteSubject)
		}

		events := api.Group("/events")
		{
			events.GET("/", eventController.GetEvents)
			events.DELETE("/:id", eventController.DeleteEvent)
		}

		lectures := api.Group("/lectures")
		{
			lectures.GET("/:id", eventController.GetLecture)
			lectures.POST("/", eventController.CreateLectures)
			lectures.PATCH("/:id", eventController.UpdateLecture)
			lectures.DELETE("/group/:group_id", eventController.DeleteLectures)
		}

		assignments := api.Group("/assignments")
		{
			assignments.GET("/:id", eventController.GetAssignment)
			assignments.POST("/", eventController.CreateAssignment)
			assignments.PATCH("/:id", eventController.UpdateAssignment)
		}
	}
}
