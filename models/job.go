package models

import "time"

// Job represents the job document stored in MongoDB.
type Job struct {
	ID        string    `bson:"_id,omitempty"`
	Status    string    `bson:"status"`
	CreatedAt time.Time `bson:"created_at"`
}
