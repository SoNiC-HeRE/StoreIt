package utils

import (
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"math/rand"
	"net/http"
	"time"
)

// DownloadImage downloads an image from a given URL and returns its dimensions.
func DownloadImage(url string) (int, int, error) {
	resp, err := http.Get(url)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to download image: %w", err)
	}
	defer resp.Body.Close()

	img, _, err := image.DecodeConfig(resp.Body)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to decode image: %w", err)
	}

	return img.Width, img.Height, nil
}

// CalculatePerimeter computes the perimeter of an image.
func CalculatePerimeter(width, height int) float64 {
	return 2 * float64(width+height)
}

// SimulateProcessingDelay adds a random delay between 0.1 to 0.4 seconds.
func SimulateProcessingDelay() {
	delay := time.Duration(rand.Intn(300)+100) * time.Millisecond
	time.Sleep(delay)
}
