package utils

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/nats-io/nats.go"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// PrintError prints an error message
func PrintError(err error) {
	fmt.Println(err)
}

func LoadEnv() error {
	environment := os.Getenv("ENVIRONMENT")
	var envFile string
	if environment == "dev" {
		envFile = ".env.dev"
	} else {
		envFile = ".env"
	}

	// Load .env file based on the environment
	err := godotenv.Load(envFile)
	if err != nil {
		return err
	}
	return nil
}

func SendNATSmessage(natsClient *nats.Conn, subject string, message []byte) error {
	if err := natsClient.Publish(subject, message); err != nil {
		return err
	}
	return nil
}

// func SendMsgToPlaceIndexer(natsClient *nats.Conn, subject string, message []byte) error {
// 	if err := natsClient.Publish(subject, message); err != nil {
// 		return err
// 	}
// 	return nil
// }

func ReceiveNATSmessage(natsClient *nats.Conn, subject string) error {
	// Subscribe to the subject to receive messages
	sub, err := natsClient.Subscribe(subject, func(msg *nats.Msg) {
		log.Printf("Received message: %s\n", string(msg.Data))
	})
	if err != nil {
		return err
	}
	defer sub.Unsubscribe()

	return nil
	// // Wait indefinitely to keep the program running
	// select {}
}

func ConvertStringsToObjectIDs(strings []string) ([]primitive.ObjectID, error) {
	var objectIDs []primitive.ObjectID

	for _, str := range strings {
		objectID, err := primitive.ObjectIDFromHex(str)
		if err != nil {
			return nil, err
		}
		objectIDs = append(objectIDs, objectID)
	}

	return objectIDs, nil
}

func ConvertObjectIDsToString(objectIDs []primitive.ObjectID) []string {
	var strings []string

	for _, objectID := range objectIDs {
		strings = append(strings, objectID.Hex())
	}

	return strings
}
