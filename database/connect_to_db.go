package db

import (
	"context"
	"log"
	"storeit/config"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Global MongoDB collection variables.
var (
	JobCollection   *mongo.Collection
	ImageCollection *mongo.Collection
	StoreCollection *mongo.Collection
)

// ConnectMongoDB establishes a connection to MongoDB using the URI provided
// by the configuration and initializes the global collection variables.
func ConnectMongoDB() {
	// Retrieve MongoDB URI from environment variables via config package.
	uri := config.GetMongoURI()
	clientOptions := options.Client().ApplyURI(uri)

	// Connect to MongoDB.
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal("MongoDB connection failed:", err)
	}

	// Access the specific database.
	database := client.Database("storeit")

	// Initialize collections.
	JobCollection = database.Collection("jobs")
	ImageCollection = database.Collection("images")
	StoreCollection = database.Collection("store_master")

	log.Println("Connected to MongoDB")
}
