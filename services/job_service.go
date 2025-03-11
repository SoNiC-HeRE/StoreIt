package services

import (
	"context"
	"log"

	"storeit/data"
	"storeit/models"
	"storeit/utils"
)

// ProcessJob processes a job by downloading images from each visit,
// computing the perimeter of each image, and updating the job status accordingly.
//
// For each visit, the function iterates over all provided image URLs, and for each image:
// 1. Downloads the image to get its dimensions.
// 2. Calculates the image perimeter using the formula 2 * (height + width).
// 3. Simulates a processing delay to mimic GPU processing.
// 4. Saves the processed image record in the database.
// 
// If any error occurs (during image download or saving), the job is marked as "failed"
// and the corresponding errors (store IDs and error messages) are recorded.
// If no errors occur, the job status is updated to "completed".
func ProcessJob(jobID string, visits []models.Visit) {
	// Create a context for DB operations.
	ctx := context.Background()

	// jobFailed indicates if any errors occurred during processing.
	var jobFailed bool
	// jobErrors collects all errors encountered during processing.
	var jobErrors []models.JobError

	// Process each visit.
	for _, visit := range visits {
		// Process each image URL in the visit.
		for _, imageURL := range visit.ImageURLs {
			// Download the image and get its dimensions.
			width, height, err := utils.DownloadImage(imageURL)
			if err != nil {
				jobErrors = append(jobErrors, models.JobError{
					StoreID: visit.StoreID,
					Error:   "Image download failed",
				})
				jobFailed = true
				// Continue with next image if download fails.
				continue
			}

			// Calculate the perimeter: 2 * (height + width).
			perimeter := utils.CalculatePerimeter(width, height)

			// Simulate a random delay (between 100ms and 400ms) to mimic GPU processing.
			utils.SimulateProcessingDelay()

			// Create an image record.
			img := models.Image{
				JobID:     jobID,
				StoreID:   visit.StoreID,
				ImageURL:  imageURL,
				Perimeter: int(perimeter),
				Status:    "completed",
			}

			// Save the image record in the database.
			if err := data.SaveImage(ctx, img); err != nil {
				jobErrors = append(jobErrors, models.JobError{
					StoreID: visit.StoreID,
					Error:   "Failed to save image record",
				})
				jobFailed = true
			}
		}
	}

	// Update the job status in the database based on the processing outcome.
	if jobFailed {
		// Update the job status as "failed" along with error details.
		if err := data.UpdateJobStatusWithErrors(jobID, "failed", jobErrors); err != nil {
			log.Printf("Failed to update job status with errors for job %s: %v", jobID, err)
		}
		log.Printf("Job %s marked as failed due to errors.", jobID)
	} else {
		// Update the job status as "completed".
		if err := data.UpdateJobStatus(jobID, "completed"); err != nil {
			log.Printf("Failed to update job status for job %s: %v", jobID, err)
		}
		log.Printf("Job %s completed successfully.", jobID)
	}
}
