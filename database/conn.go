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
	mongoURL := os.Getenv("MONGODB_URL")
	if mongoURL == "" {
		log.Fatal("MONGODB_URL environment variable is not set")
	}

	clientOpt := options.Client().ApplyURI(mongoURL)

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
