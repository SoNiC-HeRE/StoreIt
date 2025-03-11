package data

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"

	"storeit/database"
	"storeit/models"
)

// Image represents the structure for an image record stored in MongoDB.
type Image struct {
	ID        string    `bson:"_id,omitempty"` // Unique identifier (auto-generated)
	JobID     string    `bson:"job_id"`        // ID of the associated job
	StoreID   string    `bson:"store_id"`      // Store identifier
	ImageURL  string    `bson:"image_url"`     // URL of the image
	Perimeter float64   `bson:"perimeter"`     // Calculated perimeter of the image
	Status    string    `bson:"status"`        // Processing status: "completed" or "failed"
	CreatedAt time.Time `bson:"created_at"`    // Timestamp when the record was created
}

// SaveImage stores the provided image record into MongoDB.
// It sets the current time as the CreatedAt timestamp.
func SaveImage(ctx context.Context, img models.Image) error {
	// Set the creation time for the image record.
	img.CreatedAt = time.Now()

	// Insert the image record into the "images" collection.
	_, err := db.ImageCollection.InsertOne(ctx, img)
	if err != nil {
		log.Printf("Failed to store image: %v", err)
	}
	return err
}

// GetImagesByJobID retrieves all image records associated with the specified jobID.
// It returns a slice of Image and an error (if any).
func GetImagesByJobID(ctx context.Context, jobID string) ([]Image, error) {
	var images []Image

	// Find all images with the matching job_id.
	cursor, err := db.ImageCollection.Find(ctx, bson.M{"job_id": jobID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	// Decode each document into an Image struct and append to the slice.
	for cursor.Next(ctx) {
		var img Image
		if err := cursor.Decode(&img); err != nil {
			return nil, err
		}
		images = append(images, img)
	}

	return images, nil
}

// UpdateImageStatus updates the processing status of a specific image record.
// imageID: the unique identifier of the image record.
// status: the new status value ("completed" or "failed").
func UpdateImageStatus(ctx context.Context, imageID string, status string) error {
	// Update the status field of the document with the specified imageID.
	_, err := db.ImageCollection.UpdateOne(ctx, bson.M{"_id": imageID}, bson.M{"$set": bson.M{"status": status}})
	return err
}
