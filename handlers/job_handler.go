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
		Count  int            `json:"count"`
		Visits []models.Visit `json:"visits"`
	}

	// Bind JSON
	if err := c.BindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	if len(request.Visits) != request.Count {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Count does not match the number of visits"})
		return
	}

	// Collect any invalid store errors
	var invalidVisits []models.JobError

	// Validate each visit
	for i, visit := range request.Visits {
		if visit.StoreID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "store_id is mandatory for each visit"})
			return
		}
		if len(visit.ImageURLs) == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "image_url is mandatory for each visit"})
			return
		}
		// Validate store_id exists in the stores collection
		valid, err := data.IsValidStore(visit.StoreID)
		if err != nil || !valid {
			invalidVisits = append(invalidVisits, models.JobError{
				StoreID: visit.StoreID,
				Error:   "Invalid store id",
			})
		}
		// If visit_time is empty, fill with current timestamp
		if visit.VisitTime == "" {
			request.Visits[i].VisitTime = time.Now().Format(time.RFC3339)
		}
	}

	// If there are any invalid store IDs, create a failed job and return its id.
	if len(invalidVisits) > 0 {
		job := models.Job{
			ID:        uuid.NewString(),
			Status:    "failed",
			CreatedAt: time.Now(),
			Errors:    invalidVisits,
		}
		if err := data.CreateJob(job); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create job"})
			return
		}
		c.JSON(http.StatusCreated, gin.H{"job_id": job.ID})
		return
	}

	// Otherwise, create the job record with status "ongoing".
	job := models.Job{
		ID:        uuid.NewString(),
		Status:    "ongoing",
		CreatedAt: time.Now(),
	}
	if err := data.CreateJob(job); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create job"})
		return
	}

	// Start background processing.
	go services.ProcessJob(job.ID, request.Visits)

	c.JSON(http.StatusCreated, gin.H{"job_id": job.ID})
}

// GetJobStatus returns the status (and errors, if any) of a given job.
func GetJobStatus(c *gin.Context) {
	jobID := c.Query("jobid")
	var job models.Job

	err := db.JobCollection.FindOne(context.TODO(), bson.M{"_id": jobID}).Decode(&job)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": "Job not found"})
		return
	}

	resp := models.JobStatusResponse{
		Status: job.Status,
		JobID:  job.ID,
		Error:  job.Errors,
	}
	c.JSON(http.StatusOK, resp)
}
