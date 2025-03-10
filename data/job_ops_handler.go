package data

import (
    "context"
    "storeit/database"
    "storeit/models"
)

func CreateJob(job models.Job) error {
    _, err := db.JobCollection.InsertOne(context.TODO(), job)
    return err
}

func UpdateJobStatus(jobID, status string) error {
    _, err := db.JobCollection.UpdateOne(context.TODO(), bson.M{"_id": jobID}, bson.M{"$set": bson.M{"status": status}})
    return err
}
