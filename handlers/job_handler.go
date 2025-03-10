package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"storeit/database"       // Exposes JobCollection
	"storeit/models"
	"storeit/data"     // Contains CreateJob and UpdateJobStatus functions
	"storeit/services" // Contains ProcessJob
)

// SubmitJob handles the job submission endpoint.
func SubmitJob(c *gin.Context) {
	var request struct {
		Count  int `json:"count"`
		Visits []struct {
			StoreID   string   `json:"store_id"`
			ImageURLs []string `json:"image_url"`
			VisitTime string   `json:"visit_time"`
		} `json:"visits"`
	}

	if err := c.BindJSON(&request); err != nil || len(request.Visits) != request.Count {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	job := models.Job{
		ID:        uuid.NewString(),
		Status:    "ongoing",
		CreatedAt: time.Now(),
	}

	// Create job record in the database.
	if err := data.CreateJob(job); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create job"})
		return
	}

	// Start background processing of images.
	go services.ProcessJob(job.ID, request.Visits)

	c.JSON(http.StatusCreated, gin.H{"job_id": job.ID})
}

// GetJobStatus returns the status (and errors, if any) of a given job.
func GetJobStatus(c *gin.Context) {
	jobID := c.Query("jobid")
	var job models.Job

	err := db.JobCollection.FindOne(context.TODO(), bson.M{"_id": jobID}).Decode(&job)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Job not found"})
		return
	}

	resp := models.JobStatusResponse{
		Status: job.Status,
		JobID:  job.ID,
		Error:  job.Errors, // or omit if nil
	}

	c.JSON(http.StatusOK, resp)
}
