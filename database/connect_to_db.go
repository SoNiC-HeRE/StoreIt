package db

import (
	"context"
	"log"
	"storeit/config"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var JobCollection *mongo.Collection
var ImageCollection *mongo.Collection

func ConnectMongoDB() {
	clientOptions := options.Client().ApplyURI(config.GetMongoURI())
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal("MongoDB connection failed:", err)
	}

	database := client.Database("storeit")
	JobCollection = database.Collection("jobs")
	ImageCollection = database.Collection("images")

	log.Println("Connected to MongoDB")
}
