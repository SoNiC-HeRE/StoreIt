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
	var jobFailed bool

	// Process every visit and its images.
	for _, visit := range visits {
		// Here you can validate the store_id against Store Master.
		// For now, we assume the store exists.
		for _, url := range visit.ImageURLs {
			// Download image dimensions.
			width, height, err := utils.DownloadImage(url)
			if err != nil {
				log.Printf("Image download failed for store %s, URL %s: %v", visit.StoreID, url, err)
				jobFailed = true
				// Continue processing other images.
				continue
			}

			// Calculate perimeter and simulate GPU processing delay.
			perimeter := utils.CalculatePerimeter(width, height)
			utils.SimulateProcessingDelay()

			// Create image record.
			img := models.Image{
				JobID:     jobID,
				StoreID:   visit.StoreID,
				ImageURL:  url,
				Perimeter: int(perimeter),
				Status:    "completed",
			}
			if err := data.SaveImage(ctx, img); err != nil {
				log.Printf("Failed to save image for store %s: %v", visit.StoreID, err)
				jobFailed = true
			}
		}
	}

	// Update job status based on processing outcome.
	if jobFailed {
		data.UpdateJobStatus(jobID, "failed")
		log.Printf("Job %s marked as failed due to one or more errors.", jobID)
	} else {
		data.UpdateJobStatus(jobID, "completed")
		log.Printf("Job %s completed successfully.", jobID)
	}
}
