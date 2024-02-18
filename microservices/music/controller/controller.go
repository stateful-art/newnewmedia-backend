package controller

import (
	"errors"
	"log"
	"strings"

	"cloud.google.com/go/storage"
	"github.com/gofiber/fiber/v2"
	dto "newnewmedia.com/microservices/music/dto"
	service "newnewmedia.com/microservices/music/service"
)

// GetMusic gets all music
func GetMusic(c *fiber.Ctx) error {
	return c.SendFile("./public/music/test.mp3")
}

func GetMusicFile(c *fiber.Ctx) error {
	fileName := c.Params("id")
	return c.SendFile("./public/music/" + fileName)
}

// PlayMusic streams the audio file based on song ID
func PlayMusic(c *fiber.Ctx, storageClient *storage.Client) error {
	// Get the song ID from the request parameters
	songID := c.Params("id")

	// Fetch the audio file path for the given song ID
	audioFilePath, err := service.GetAudioFilePath(songID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	log.Println(audioFilePath)
	// Extract bucketName and objectName from audioFilePath
	bucketName, objectName, err := extractBucketAndObjectName(audioFilePath)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	// Stream the audio file
	err = service.StreamMusic(c, bucketName, objectName, storageClient)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return nil
}

func CreateMusic(c *fiber.Ctx, storageClient *storage.Client) error {
	var musicPayload dto.Music

	// Parse the request body into musicPayload
	if err := c.BodyParser(&musicPayload); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	// Get the uploaded audio file from the form data
	audioFile, err := c.FormFile("music")
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "Audio file is required",
		})
	}

	// Create music entry and store the music file
	err = service.CreateMusic(c, musicPayload, audioFile, storageClient)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Music created successfully",
	})
}

func GetMusicByPlace(c *fiber.Ctx) error {
	id := c.Params("id")
	music, err := service.GetMusicByPlace(id)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	return c.JSON(fiber.Map{
		"message": "Music fetched successfully",
		"data":    music,
	})
}

// extractBucketAndObjectName extracts bucketName and objectName from audioFilePath
func extractBucketAndObjectName(audioFilePath string) (string, string, error) {
	// Split audioFilePath into bucketName and objectName
	// Example: audioFilePath = "gs://your-bucket-name/your-object-name.mp3"
	log.Println(audioFilePath)
	parts := strings.SplitN(audioFilePath, "/", 4)
	if len(parts) < 4 || parts[0] != "gs:" {
		return "", "", errors.New("invalid audioFilePath")
	}
	return parts[2], parts[3], nil
}
