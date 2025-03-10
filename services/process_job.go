package services

import (
	"context"
	"log"

	"storeit/data"
	"storeit/models"
	"storeit/utils"
)

// ProcessJob processes a job by downloading images, computing perimeters,
// and updating the job status accordingly.
func ProcessJob(jobID string, visits []models.Visit) {
	ctx := context.Background()
	var jobFailed bool
	var jobErrors []models.JobError

	// Process every visit and its images.
	for _, visit := range visits {
		for _, url := range visit.ImageURLs {
			width, height, err := utils.DownloadImage(url)
			if err != nil {
				jobErrors = append(jobErrors, models.JobError{
					StoreID: visit.StoreID,
					Error:   "Image download failed",
				})
				jobFailed = true
				continue
			}

			perimeter := utils.CalculatePerimeter(width, height)
			utils.SimulateProcessingDelay()

			img := models.Image{
				JobID:     jobID,
				StoreID:   visit.StoreID,
				ImageURL:  url,
				Perimeter: int(perimeter),
				Status:    "completed",
			}
			if err := data.SaveImage(ctx, img); err != nil {
				jobErrors = append(jobErrors, models.JobError{
					StoreID: visit.StoreID,
					Error:   "Failed to save image record",
				})
				jobFailed = true
			}
		}
	}

	if jobFailed {
		if err := data.UpdateJobStatusWithErrors(jobID, "failed", jobErrors); err != nil {
			log.Printf("Failed to update job status with errors for job %s: %v", jobID, err)
		}
		log.Printf("Job %s marked as failed due to errors.", jobID)
	} else {
		if err := data.UpdateJobStatus(jobID, "completed"); err != nil {
			log.Printf("Failed to update job status for job %s: %v", jobID, err)
		}
		log.Printf("Job %s completed successfully.", jobID)
	}
}
