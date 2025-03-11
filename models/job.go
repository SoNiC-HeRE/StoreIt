package models

import "time"

// Job represents a job record stored in MongoDB.
// It captures the job status, creation timestamp, and any errors encountered during processing.
type Job struct {
	ID        string     `bson:"_id,omitempty"`         // Unique job identifier
	Status    string     `bson:"status"`                // Job status (e.g., "ongoing", "completed", "failed")
	CreatedAt time.Time  `bson:"created_at"`            // Timestamp when the job was created
	Errors    []JobError `bson:"errors,omitempty"`      // List of errors (if any) encountered during job processing
}

// JobError represents an error encountered during a job processing.
// It associates a store identifier with the error message.
type JobError struct {
	StoreID string `bson:"store_id" json:"store_id"` // Identifier of the store where the error occurred
	Error   string `bson:"error" json:"error"`       // Error message describing the issue
}

// JobStatusResponse defines the JSON structure returned by the job status API endpoint.
type JobStatusResponse struct {
	Status string     `json:"status"`           // Current status of the job
	JobID  string     `json:"job_id"`           // Unique job identifier
	Error  []JobError `json:"error,omitempty"`  // Error details (if any), omitted if there are none
}
