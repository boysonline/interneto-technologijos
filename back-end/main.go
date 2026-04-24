package main

import (
	"back-end/internal/controller"
	"back-end/internal/database"
	"back-end/internal/repository"
	"back-end/internal/routes"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	database.Connect()

	eventRepo := &repository.EventRepository{
		Events: database.Collection("events"),
	}
	subjectRepo := &repository.SubjectRepository{
		Subjects: database.Collection("subjects"),
	}

	r := gin.Default()
	routes.SetupRoutes(r, &controller.SubjectController{Repository: subjectRepo}, &controller.EventController{Repository: eventRepo})
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Fallback for local testing
	}
	if os.Getenv("ENV") == "production" {
		r.Run(":" + port)
	} else {
		log.Fatal(r.RunTLS(":"+port, "./cert.pem", "./key.pem"))
	}

	log.Fatal(r.RunTLS(":8080", "./cert.pem", "./key.pem"))
}
