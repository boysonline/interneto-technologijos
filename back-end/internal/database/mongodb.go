package database

import (
	"context"
	"log"
	"os" // Built-in package to read environment variables
	"time"

	"github.com/joho/godotenv" // The library you just installed
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Client *mongo.Client

func Connect() {
	// 1. Load the .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// 2. Get the value using os.Getenv
	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		log.Fatal("MONGO_URI not found in .env")
	}

	clientOptions := options.Client().ApplyURI(mongoURI)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	Client, err = mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal("Connection failed:", err)
	}

	// ... rest of your ping logic
}
