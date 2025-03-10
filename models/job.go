package models

import "time"

type Job struct {
	ID        string        `bson:"_id,omitempty"`
	Status    string        `bson:"status"`
	CreatedAt time.Time     `bson:"created_at"`
	Errors    []JobError    `bson:"errors,omitempty"` // New field to capture errors
}

type JobError struct {
	StoreID string `bson:"store_id"`
	Error   string `bson:"error"`
}
