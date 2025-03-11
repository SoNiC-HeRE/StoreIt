package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"storeit/database"   // Exposes JobCollection
	"storeit/data"       // Contains CreateJob, UpdateJobStatus functions, and store validation
	"storeit/models"     // Contains Job, Visit, and JobError definitions
	"storeit/services"   // Contains ProcessJob for background processing
)

// SubmitJob handles the job submission endpoint.
// It validates the incoming request, creates a job record (with failed status if any store is invalid),
// and triggers background processing for valid jobs.
func SubmitJob(c *gin.Context) {
	// Define the request structure.
	var request struct {
		Count  int            `json:"count"`
		Visits []models.Visit `json:"visits"`
	}

	// Bind incoming JSON to the request structure.
	if err := c.BindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// Validate that the count matches the number of visits.
	if len(request.Visits) != request.Count {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Count does not match the number of visits"})
		return
	}

	// Prepare a slice to collect invalid store errors.
	var invalidVisits []models.JobError

	// Iterate over each visit to validate required fields.
	for i, visit := range request.Visits {
		// Check that store_id is provided.
		if visit.StoreID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "store_id is mandatory for each visit"})
			return
		}

		// Check that at least one image URL is provided.
		if len(visit.ImageURLs) == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "image_url is mandatory for each visit"})
			return
		}

		// Validate that the store exists.
		valid, err := data.IsValidStore(visit.StoreID)
		if err != nil || !valid {
			invalidVisits = append(invalidVisits, models.JobError{
				StoreID: visit.StoreID,
				Error:   "Invalid store id",
			})
		}

		// If visit_time is not provided, assign the current timestamp.
		if visit.VisitTime == "" {
			request.Visits[i].VisitTime = time.Now().Format(time.RFC3339)
		}
	}

	// If any store IDs are invalid, create a job with "failed" status and return its ID.
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

	// Create a job record with status "ongoing" for valid requests.
	job := models.Job{
		ID:        uuid.NewString(),
		Status:    "ongoing",
		CreatedAt: time.Now(),
	}
	if err := data.CreateJob(job); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create job"})
		return
	}

	// Trigger background processing of the job.
	go services.ProcessJob(job.ID, request.Visits)

	// Return the job ID in the response.
	c.JSON(http.StatusCreated, gin.H{"job_id": job.ID})
}

// GetJobStatus returns the status and any error details of a given job.
// It retrieves the job from MongoDB using the provided job ID.
func GetJobStatus(c *gin.Context) {
	// Extract the jobid query parameter.
	jobID := c.Query("jobid")
	var job models.Job

	// Find the job in the database.
	err := db.JobCollection.FindOne(context.TODO(), bson.M{"_id": jobID}).Decode(&job)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Job not found"})
		return
	}

	// Build the response structure.
	resp := models.JobStatusResponse{
		Status: job.Status,
		JobID:  job.ID,
		Error:  job.Errors,
	}
	c.JSON(http.StatusOK, resp)
}
