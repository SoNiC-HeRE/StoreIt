package services

import (
	"context"
	"log"

	"storeit/data"
	"storeit/models"
	"storeit/utils"
)

// VisitInput holds the information for each store visit.
type VisitInput struct {
	StoreID   string   `json:"store_id"`
	ImageURLs []string `json:"image_url"`
	VisitTime string   `json:"visit_time"`
}

// ProcessJob processes a job by downloading images, computing perimeters,
// and updating the job status accordingly. If any image fails to process,
// the job is marked as "failed" and errors (failed store_ids) are logged.
func ProcessJob(jobID string, visits []struct {
	StoreID   string   `json:"store_id"`
	ImageURLs []string `json:"image_url"`
	VisitTime string   `json:"visit_time"`
}) {
	ctx := context.Background()
	var jobErrors []models.JobError

	// Process every visit and its images.
	for _, visit := range visits {
		for _, url := range visit.ImageURLs {
			width, height, err := utils.DownloadImage(url)
			if err != nil {
				// Append an error record for this store.
				jobErrors = append(jobErrors, models.JobError{
					StoreID: visit.StoreID,
					Error:   "Image download failed", // You can include err.Error() if needed
				})
				continue
			}

			// Process image: calculate perimeter and simulate delay.
			perimeter := utils.CalculatePerimeter(width, height)
			utils.SimulateProcessingDelay()

			// Create and save image record.
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
			}
		}
	}

	// Update job status based on errors.
	if len(jobErrors) > 0 {
		// Update job as failed with error details.
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
