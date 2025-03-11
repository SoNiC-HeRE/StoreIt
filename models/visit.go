// models/visit.go

package models

// Visit represents a store visit in a job submission.
// It includes the store identifier, one or more image URLs captured during the visit,
// and the timestamp of the visit.
type Visit struct {
	StoreID   string   `json:"store_id"`   // Unique identifier for the store
	ImageURLs []string `json:"image_url"`  // List of image URLs captured during the visit
	VisitTime string   `json:"visit_time"` // Timestamp of the visit (RFC3339 format recommended)
}
