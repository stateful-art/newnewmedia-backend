package utils

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"

	"github.com/gofiber/fiber/v2"
	uuid "github.com/google/uuid"
)

func SaveMusicFile(c *fiber.Ctx) (string, error) {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	serverPathEnv := os.Getenv("BACKEND_ORIGIN")

	musicFile, err := c.FormFile("music")
	if err != nil {
		return "", err
	}

	musicUUID := uuid.New().String()
	musicNameFormat := fmt.Sprintf("./public/music/%s%s", musicUUID, ".mp3")
	saveErr := c.SaveFile(musicFile, musicNameFormat)
	if saveErr != nil {
		return "", saveErr
	}

	savePath := fmt.Sprintf("%s/music/%s%s%s", serverPathEnv, "file/", musicUUID, ".mp3")

	return savePath, nil

}
