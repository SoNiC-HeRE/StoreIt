package data

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// Image represents an image stored in MongoDB
type Image struct {
	ID        string    `bson:"_id,omitempty"`
	JobID     string    `bson:"job_id"`
	StoreID   string    `bson:"store_id"`
	ImageURL  string    `bson:"image_url"`
	Perimeter float64   `bson:"perimeter"`
	Processed bool      `bson:"processed"`
	CreatedAt time.Time `bson:"created_at"`
}

// ImageStorage handles image-related DB operations
type ImageStorage struct {
	Collection *mongo.Collection
}

// NewImageStorage initializes a new ImageStorage instance
func NewImageStorage(db *mongo.Database) *ImageStorage {
	return &ImageStorage{
		Collection: db.Collection("images"),
	}
}

// SaveImage stores image processing details in MongoDB
func (s *ImageStorage) SaveImage(ctx context.Context, img Image) error {
	img.CreatedAt = time.Now()
	_, err := s.Collection.InsertOne(ctx, img)
	if err != nil {
		log.Printf("Failed to store image: %v", err)
	}
	return err
}

// GetImagesByJobID retrieves images related to a specific job
func (s *ImageStorage) GetImagesByJobID(ctx context.Context, jobID string) ([]Image, error) {
	var images []Image
	cursor, err := s.Collection.Find(ctx, bson.M{"job_id": jobID})
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

// UpdateImageStatus updates the processing status of an image
func (s *ImageStorage) UpdateImageStatus(ctx context.Context, imageID string, processed bool) error {
	_, err := s.Collection.UpdateOne(ctx, bson.M{"_id": imageID}, bson.M{"$set": bson.M{"processed": processed}})
	return err
}
