package storage

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func ConnectToMongoDB() (*mongo.Client, error) {
	// Set the MongoDB connection options
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")

	// Connect to MongoDB
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		return nil, err
	}

	// Ping the MongoDB server to check the connection
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		return nil, err
	}

	fmt.Println("Connected to MongoDB!")
	return client, nil
}

func CreateDatabaseAndCollections(client *mongo.Client, dbName string) (*mongo.Database, error) {
	// Check if the database already exists
	databases, err := client.ListDatabaseNames(context.TODO(), bson.M{})
	if err != nil {
		return nil, err
	}

	for _, db := range databases {
		if db == dbName {
			fmt.Printf("Database '%s' already exists\n", dbName)
			return client.Database(dbName), nil
		}
	}

	// Create the new database
	db := client.Database(dbName)

	// Create a collection for users
	err = db.CreateCollection(context.TODO(), "users", nil)
	if err != nil {
		return nil, err
	}

	// Create a collection for tasks
	err = db.CreateCollection(context.TODO(), "tasks", nil)
	if err != nil {
		return nil, err
	}

	fmt.Printf("Database '%s' created successfully\n", dbName)
	return db, nil
}
