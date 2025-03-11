package models

// Store represents a store record from the Store Master.
// It includes the store's unique identifier, name, and area code.
type Store struct {
	StoreID   string `bson:"store_id" json:"store_id"`     // Unique identifier for the store
	StoreName string `bson:"store_name" json:"store_name"` // Name of the store
	AreaCode  string `bson:"area_code" json:"area_code"`   // Area code associated with the store
}
