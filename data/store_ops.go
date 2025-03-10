package data

import (
	"context"

	"storeit/database"
	"storeit/models"

	"go.mongodb.org/mongo-driver/bson"
)

// IsValidStore checks if a given store_id exists in the stores collection.
func IsValidStore(storeID string) (bool, error) {
	var store models.Store
	err := db.StoreCollection.FindOne(context.TODO(), bson.M{"StoreID": storeID}).Decode(&store)
	if err != nil {
		// If err is ErrNoDocuments, return false without error.
		return false, err
	}
	return true, nil
}
