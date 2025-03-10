package models

import "time"

type Job struct {
    ID        string    `bson:"_id,omitempty"`
    Status    string    `bson:"status"`
    CreatedAt time.Time `bson:"created_at"`
}
