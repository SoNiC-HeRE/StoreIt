package handlers

import (
    "context"
    "net/http"
    "github.com/gin-gonic/gin"
    "retail_pulse/models"
    "retail_pulse/repository"
)

func SubmitJob(c *gin.Context) {
    var request struct {
        Count  int `json:"count"`
        Visits []struct {
            StoreID   string   `json:"store_id"`
            ImageURLs []string `json:"image_url"`
        } `json:"visits"`
    }

    if err := c.BindJSON(&request); err != nil || len(request.Visits) != request.Count {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
        return
    }

    job := models.Job{ID: uuid.NewString(), Status: "ongoing", CreatedAt: time.Now()}
    repository.CreateJob(job)

    c.JSON(http.StatusCreated, gin.H{"job_id": job.ID})
}

func GetJobStatus(c *gin.Context) {
    jobID := c.Query("jobid")
    var job models.Job

    err := db.JobCollection.FindOne(context.TODO(), bson.M{"_id": jobID}).Decode(&job)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Job not found"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"status": job.Status, "job_id": job.ID})
}
