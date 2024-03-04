package utils

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
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
