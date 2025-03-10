package models

// Store represents a store record from the Store Master.
type Store struct {
	StoreID   string `bson:"store_id" json:"store_id"`
	StoreName string `bson:"store_name" json:"store_name"`
	AreaCode  string `bson:"area_code" json:"area_code"`
}
