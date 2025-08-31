// Package database provides functions to connect and disconnect from MongoDB.
package database

import (
	"context"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var MongoDBClient *mongo.Client

func ConnectToMongoDB() {
	mongoURI := os.Getenv("MONGODB_URI")
	if mongoURI == "" {
		log.Fatal("MONGODB_URI environment variable is not set")
	}

	clientOpt := options.Client().ApplyURI(mongoURI)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var err error
	MongoDBClient, err = mongo.Connect(ctx, clientOpt)
	if err != nil {
		log.Fatal("Error connecting to MongoDB:", err)
	}

	err = MongoDBClient.Ping(ctx, nil)
	if err != nil {
		log.Fatal("Error pinging MongoDB:", err)
	}

	log.Println("Connected to MongoDB")
}

func DisconnectFromMongoDB() {
	if MongoDBClient != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := MongoDBClient.Disconnect(ctx); err != nil {
			log.Printf("Error disconnecting from MongoDB: %v", err)
		} else {
			log.Println("Disconnected from MongoDB")
		}
	}
}
