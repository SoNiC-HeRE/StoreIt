package utils

import (
	"fmt"
	"image"
	_ "image/jpeg" // Enable JPEG image decoding.
	_ "image/png"  // Enable PNG image decoding.
	"math/rand"
	"net/http"
	"time"
)

// DownloadImage retrieves an image from the specified URL and returns its width and height.
// If an error occurs during download or decoding, it returns an error.
func DownloadImage(url string) (int, int, error) {
	// Send HTTP GET request to download the image.
	resp, err := http.Get(url)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to download image: %w", err)
	}
	defer resp.Body.Close()

	// Decode image configuration to get its dimensions.
	imgConfig, _, err := image.DecodeConfig(resp.Body)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to decode image: %w", err)
	}

	return imgConfig.Width, imgConfig.Height, nil
}

// CalculatePerimeter returns the perimeter of an image using its width and height.
// Formula: 2 * (width + height)
func CalculatePerimeter(width, height int) float64 {
	return 2 * float64(width+height)
}

// SimulateProcessingDelay pauses execution for a random duration between 100ms and 400ms.
// This is used to mimic GPU processing delay.
func SimulateProcessingDelay() {
	delay := time.Duration(rand.Intn(300)+100) * time.Millisecond
	time.Sleep(delay)
}
