package utils

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/nats-io/nats.go"
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

func SendNATSmessage(natsClient *nats.Conn, subject string, content string) error {
	// Send a message with the subject
	message := []byte(content)
	if err := natsClient.Publish(subject, message); err != nil {
		return err
	}
	log.Printf("Published message: %s\n", string(message))

	// // Wait indefinitely to keep the program running
	// select {}

	return nil
}

func SendMsgToPlaceIndexer(natsClient *nats.Conn, subject string, message []byte) error {
	// Send a message with the subject
	if err := natsClient.Publish(subject, message); err != nil {
		return err
	}
	log.Printf("Published message: %s\n", string(message))

	// // Wait indefinitely to keep the program running
	// select {}

	return nil
}

// func SendNATSmessageToElastic(natsClient *nats.Conn, subject string, place dto.Place) error {
// 	// Send a message with the subject
// 	message := []byte(place)
// 	if err := natsClient.Publish(subject, message); err != nil {
// 		return err
// 	}
// 	log.Printf("Published message: %v\n", dto.Place(message))

// 	// // Wait indefinitely to keep the program running
// 	// select {}

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
