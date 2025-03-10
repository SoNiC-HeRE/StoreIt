package models

import "time"

type Image struct {
	ID        string    `bson:"_id,omitempty"`
	JobID     string    `bson:"job_id"`
	StoreID   string    `bson:"store_id"`
	ImageURL  string    `bson:"image_url"`
	Perimeter int       `bson:"perimeter,omitempty"`
	Status    string    `bson:"status"`
	CreatedAt time.Time `bson:"created_at"` 
}
