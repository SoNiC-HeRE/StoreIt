package data

import (
	"context"

	"storeit/database"
	"storeit/models"

	"go.mongodb.org/mongo-driver/bson"
)

// CreateJob inserts a new job document into MongoDB.
func CreateJob(job models.Job) error {
	_, err := db.JobCollection.InsertOne(context.TODO(), job)
	return err
}

// UpdateJobStatus updates the job status in MongoDB.
func UpdateJobStatus(jobID, status string) error {
	_, err := db.JobCollection.UpdateOne(context.TODO(), bson.M{"_id": jobID}, bson.M{"$set": bson.M{"status": status}})
	return err
}

// UpdateJobStatusWithErrors updates the job status and saves error details.
func UpdateJobStatusWithErrors(jobID, status string, errors []models.JobError) error {
	_, err := db.JobCollection.UpdateOne(
		context.TODO(),
		bson.M{"_id": jobID},
		bson.M{"$set": bson.M{"status": status, "errors": errors}},
	)
	return err
}