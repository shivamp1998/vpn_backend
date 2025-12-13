package database

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func InitializeIndexes(ctx context.Context) error {
	if DB == nil {
		return nil
	}

	if err := CreateUserIndexes(ctx); err != nil {
		return err
	}

	log.Println("Database indexes initialized")
	return nil
}

func CreateUserIndexes(ctx context.Context) error {
	userCollection := DB.Collection("users")

	emailIndexModel := mongo.IndexModel{
		Keys: map[string]interface{}{
			"email": 1,
		},
		Options: options.Index().SetUnique(true).SetName("email_unique"),
	}

	indexes := []mongo.IndexModel{emailIndexModel}

	_, err := userCollection.Indexes().CreateMany(ctx, indexes)
	if err != nil {
		if !isDuplicateIndexError(err) {
			return err
		}

		log.Println("Index already exist, skipping...")
	}
	return nil
}

func isDuplicateIndexError(err error) bool {
	return err != nil && (err.Error() == "index already exists" || err.Error() == "IndexOptionsConflict")
}
