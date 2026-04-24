package database

import (
	"context"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Client *mongo.Client

func Connect() {
	var err error

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
	err = Client.Ping(ctx, nil)
	if err != nil {
		log.Fatal("Could not ping MongoDB:", err)
	}
}

func Collection(collection string) *mongo.Collection {
	database := os.Getenv("DATABASE_NAME")
	if database == "" {
		log.Fatal("DATABASE_NAME variable is missing from .env")
	}
	return Client.Database(database).Collection(collection)
}
