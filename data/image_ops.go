package data

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"

	"storeit/database"
	"storeit/models"

)

// Image represents an image stored in MongoDB.
type Image struct {
	ID        string    `bson:"_id,omitempty"`
	JobID     string    `bson:"job_id"`
	StoreID   string    `bson:"store_id"`
	ImageURL  string    `bson:"image_url"`
	Perimeter float64   `bson:"perimeter"`
	Status    string    `bson:"status"`    // "completed" or "failed"
	CreatedAt time.Time `bson:"created_at"`
}

// SaveImage stores image processing details in MongoDB.
func SaveImage(ctx context.Context, img models.Image) error {
	// Set the creation time.
	// (Optionally you could also update status here if needed.)
	img.CreatedAt = time.Now()
	_, err := db.ImageCollection.InsertOne(ctx, img)
	if err != nil {
		log.Printf("Failed to store image: %v", err)
	}
	return err
}

// GetImagesByJobID retrieves images related to a specific job.
func GetImagesByJobID(ctx context.Context, jobID string) ([]Image, error) {
	var images []Image
	cursor, err := db.ImageCollection.Find(ctx, bson.M{"job_id": jobID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var img Image
		if err := cursor.Decode(&img); err != nil {
			return nil, err
		}
		images = append(images, img)
	}
	return images, nil
}

// UpdateImageStatus updates the processing status of an image.
func UpdateImageStatus(ctx context.Context, imageID string, status string) error {
	_, err := db.ImageCollection.UpdateOne(ctx, bson.M{"_id": imageID}, bson.M{"$set": bson.M{"status": status}})
	return err
}
