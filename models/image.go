package models

import "time"

// Image represents an image record stored in MongoDB.
// It includes metadata about the image processing job, the store associated with the image,
// the calculated perimeter of the image, and its processing status.
type Image struct {
	ID        string    `bson:"_id,omitempty"`         // Unique identifier (auto-generated by MongoDB)
	JobID     string    `bson:"job_id"`                // ID of the associated job
	StoreID   string    `bson:"store_id"`              // Identifier of the store where the image was captured
	ImageURL  string    `bson:"image_url"`             // URL of the image file
	Perimeter int       `bson:"perimeter,omitempty"`   // Calculated perimeter of the image (2*(height+width))
	Status    string    `bson:"status"`                // Processing status: "completed", "failed", etc.
	CreatedAt time.Time `bson:"created_at"`            // Timestamp when the record was created
}
