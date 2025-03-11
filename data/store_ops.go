package data

import (
	"context"

	"storeit/database"
	"storeit/models"

	"go.mongodb.org/mongo-driver/bson"
)

// IsValidStore checks if a store with the given storeID exists in the stores collection.
// It returns true if the store is found; otherwise, it returns false along with an error.
func IsValidStore(storeID string) (bool, error) {
	// Create a context. In a production system, you might want to use a context with timeout.
	ctx := context.Background()

	// Query the stores collection for a document with the specified StoreID.
	var store models.Store
	err := db.StoreCollection.FindOne(ctx, bson.M{"StoreID": storeID}).Decode(&store)
	if err != nil {
		// If an error occurs (e.g., no document found), return false and the error.
		return false, err
	}

	// Store found, return true.
	return true, nil
}
