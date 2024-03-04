package db

import (
	"context"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	utils "newnewmedia.com/commons/utils"
)

// ConnectDB connects to the database
var Client, _ = ConnectDB()

func ConnectDB() (*mongo.Client, error) {

	if err := utils.LoadEnv(); err != nil {
		log.Fatalf("Error loading environment variables: %v", err)
	}
	connectionString := os.Getenv("MONGO_CLUSTER_CONNSTRING")
	// Use the SetServerAPIOptions() method to set the version of the Stable API on the client
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(connectionString).SetServerAPIOptions(serverAPI)

	// Create a new client and connect to the server
	client, err := mongo.Connect(context.TODO(), opts)
	if err != nil {
		panic(err)
	}

	// Send a ping to confirm a successful connection
	if err := client.Database("newnewmedia").RunCommand(context.TODO(), bson.D{{Key: "ping", Value: 1}}).Err(); err != nil {
		panic(err)
	}

	log.Println("MongoDB: OK")
	return client, nil
}
