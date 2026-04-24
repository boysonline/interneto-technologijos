package controller

import (
	"back-end/internal/model"
	"back-end/internal/repository"
	"net/http"

	"github.com/gin-gonic/gin"
)

type SubjectController struct {
	Repository *repository.SubjectRepository
}

func NewSubjectController(repo *repository.SubjectRepository) *SubjectController {
	return &SubjectController{Repository: repo}
}

// GET /subjects/
func (controller *SubjectController) GetSubjects(ginContext *gin.Context) {
	subjects, err := controller.Repository.GetAll(ginContext.Request.Context())
	if err != nil {
		ginContext.JSON(http.StatusInternalServerError, gin.H{"error": "Nepavyko gauti dalykų sąrašo"})
		return
	}

	if subjects == nil {
		subjects = []model.Subject{}
	}

	ginContext.JSON(http.StatusOK, subjects)
}

// GET /subjects/:id
func (controller *SubjectController) GetSubject(ginContext *gin.Context) {
	id, ok := ParseID(ginContext)
	if !ok {
		return
	}

	subject, err := controller.Repository.GetSubject(ginContext.Request.Context(), id)
	if err != nil {
		ginContext.JSON(http.StatusNotFound, gin.H{"error": "Dalykas nerastas"})
		return
	}

	ginContext.JSON(http.StatusOK, subject)
}

// POST /subjects
type SubjectRequest struct {
	Title string `json:"title" binding:"required,min=1,max=50"`
}

func (controller *SubjectController) CreateSubject(ginContext *gin.Context) {
	var request SubjectRequest
	if err := ginContext.ShouldBindJSON(&request); err != nil {
		ginContext.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	subject := model.NewSubject(request.Title)

	createdSubject, err := controller.Repository.CreateSubject(ginContext.Request.Context(), subject)
	if err != nil {
		ginContext.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ginContext.JSON(http.StatusCreated, createdSubject)
}

// PATCH /subjects/:id
func (controller *SubjectController) UpdateSubject(ginContext *gin.Context) {
	id, ok := ParseID(ginContext)
	if !ok {
		return
	}

	var request SubjectRequest
	if err := ginContext.ShouldBindJSON(&request); err != nil {
		ginContext.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	subject := model.Subject{
		ID:    id,
		Title: request.Title,
	}

	updatedSubject, err := controller.Repository.UpdateSubject(ginContext.Request.Context(), subject)
	if err != nil {
		ginContext.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ginContext.JSON(http.StatusOK, updatedSubject)
}

// DELETE /subjects/:id
func (controller *SubjectController) DeleteSubject(ginContext *gin.Context) {
	id, ok := ParseID(ginContext)
	if !ok {
		return
	}

	if err := controller.Repository.DeleteSubject(ginContext.Request.Context(), id); err != nil {
		ginContext.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ginContext.Status(http.StatusNoContent)
}
