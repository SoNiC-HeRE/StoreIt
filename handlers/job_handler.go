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

// JobStatusResponse defines the structure for the job status API response.
type JobStatusResponse struct {
	Status string            `json:"status"`
	JobID  string            `json:"job_id"`
	Error  []models.JobError `json:"error,omitempty"`
}

// SubmitJob handles the job submission endpoint.
func SubmitJob(c *gin.Context) {
	var request struct {
		Count  int           `json:"count"`
		Visits []models.Visit `json:"visits"`
	}

	// Bind the incoming JSON.
	if err := c.BindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// Check that the count matches the number of visits.
	if len(request.Visits) != request.Count {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Count does not match the number of visits"})
		return
	}

	// Validate each visit: store_id and image_url must be provided.
	for i, visit := range request.Visits {
		if visit.StoreID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "store_id is mandatory"})
			return
		}
		if len(visit.ImageURLs) == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "image_url is mandatory for each visit"})
			return
		}
		// If visit_time is empty, fill with current timestamp in RFC3339 format.
		if visit.VisitTime == "" {
			request.Visits[i].VisitTime = time.Now().Format(time.RFC3339)
		}
	}

	// Create a new job record.
	job := models.Job{
		ID:        uuid.NewString(),
		Status:    "ongoing",
		CreatedAt: time.Now(),
	}

	// Insert the job record into the database.
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
		c.JSON(http.StatusBadRequest, gin.H{})
		return
	}

	resp := JobStatusResponse{
		Status: job.Status,
		JobID:  job.ID,
		Error:  job.Errors,
	}
	c.JSON(http.StatusOK, resp)
}
