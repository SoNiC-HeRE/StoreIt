package data

import (
	"context"

	"storeit/database" // Exposes JobCollection
	"storeit/models"

	"go.mongodb.org/mongo-driver/bson"
)

// CreateJob inserts a new job document into MongoDB.
func CreateJob(job models.Job) error {
	// Use context.TODO() as a placeholder. Consider using a proper context in production.
	_, err := db.JobCollection.InsertOne(context.TODO(), job)
	return err
}

// UpdateJobStatus updates the status field of a job document identified by jobID.
func UpdateJobStatus(jobID, status string) error {
	// Update the "status" field for the job with the given jobID.
	_, err := db.JobCollection.UpdateOne(
		context.TODO(),
		bson.M{"_id": jobID},
		bson.M{"$set": bson.M{"status": status}},
	)
	return err
}

// UpdateJobStatusWithErrors updates both the status and errors fields of a job document.
// The errors parameter is a slice of JobError that details any issues encountered.
func UpdateJobStatusWithErrors(jobID, status string, errors []models.JobError) error {
	// Update the "status" and "errors" fields for the job with the given jobID.
	_, err := db.JobCollection.UpdateOne(
		context.TODO(),
		bson.M{"_id": jobID},
		bson.M{"$set": bson.M{"status": status, "errors": errors}},
	)
	return err
}
