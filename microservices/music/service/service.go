package service

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"

	"cloud.google.com/go/storage"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"newnewmedia.com/microservices/music/dao"
	"newnewmedia.com/microservices/music/dto"
	repository "newnewmedia.com/microservices/music/repository"
)

var storageClient *storage.Client // Global variable to hold the GCS client instance

func init() {
	// Initialize the GCS client during application startup
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		panic(fmt.Errorf("failed to initialize GCS client: %v", err))
	}
	storageClient = client
}

// GetAudioFilePath retrieves the audio file path based on song ID
func GetAudioFilePath(songID string) (string, error) {
	// Convert the songID string to a primitive.ObjectID
	objectID, err := primitive.ObjectIDFromHex(songID)
	if err != nil {
		return "", fmt.Errorf("invalid song ID: %v", err)
	}

	// Fetch the music details from the database based on the converted song ID
	music, err := repository.GetMusicById(&objectID)
	if err != nil {
		return "", err
	}

	// Return the audio file path
	return music.Path, nil
}

func GetMusicByPlace(id string) ([]dao.Music, error) {
	placeObjID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	music, err := repository.GetMusicByPlace(&placeObjID)
	if err != nil {
		return nil, err
	}
	return music, nil
}

func CreateMusic(c *fiber.Ctx, musicPayload dto.Music, audioFile *multipart.FileHeader) error {
	// Validate if the music payload contains the necessary fields
	if musicPayload.Name == "" || musicPayload.Artist == "" {
		return fmt.Errorf("Name and Artist are required")
	}

	// Open the uploaded audio file
	audioFileContent, err := audioFile.Open()
	if err != nil {
		return err
	}
	defer audioFileContent.Close()

	// Store the music file in Google Cloud Storage
	storageBucket := "n2media-music"
	objectName := fmt.Sprintf("%s.mp3", primitive.NewObjectID().Hex()) // Generate a unique object name
	if err := storeMusicFileInGCS(c.Context(), storageBucket, objectName, audioFileContent); err != nil {
		return err
	}

	// Store the music details in MongoDB
	music := dao.Music{
		Name:   musicPayload.Name,
		Artist: musicPayload.Artist,
		Path:   fmt.Sprintf("gs://%s/%s", storageBucket, objectName), // GCS object path
		// Add additional fields if needed
	}

	if err := repository.CreateMusic(music); err != nil {
		return err
	}

	return nil
}

// storeMusicFileInGCS stores the music file in Google Cloud Storage
func storeMusicFileInGCS(ctx context.Context, bucketName, objectName string, audioFile io.Reader) error {
	// Create a new GCS client
	client, err := storage.NewClient(ctx)
	if err != nil {
		return err
	}
	defer client.Close()

	// Create a GCS object writer
	wc := client.Bucket(bucketName).Object(objectName).NewWriter(ctx)
	defer wc.Close()

	// Copy the file data to the GCS object writer
	if _, err := io.Copy(wc, audioFile); err != nil {
		return err
	}

	return nil
}

// func storeMusicFileInGCS(bucketName, objectName string, audioFile *fiber.File) error {
// 	ctx := context.Background()

// 	// Create a new GCS client
// 	client, err := storage.NewClient(ctx)
// 	if err != nil {
// 		return err
// 	}
// 	defer client.Close()

// 	// Open the file for reading
// 	file, err := audioFile.Open()
// 	if err != nil {
// 		return err
// 	}
// 	defer file.Close()

// 	// Create a GCS object writer
// 	wc := client.Bucket(bucketName).Object(objectName).NewWriter(ctx)
// 	defer wc.Close()

// 	// Copy the file data to the GCS object writer
// 	if _, err := io.Copy(wc, file); err != nil {
// 		return err
// 	}

// 	return nil
// }

// StreamMusic streams the audio file from Google Cloud Storage to the client
func StreamMusic(c *fiber.Ctx, bucketName, objectName string) error {
	ctx := context.Background()

	// Create a new reader to read the streamed data from Google Cloud Storage
	reader, err := storageClient.Bucket(bucketName).Object(objectName).NewReader(ctx)
	if err != nil {
		return err
	}
	defer reader.Close()

	// Set the Content-Type header
	c.Set("Content-Type", "audio/mpeg")

	// Stream the audio file to the client
	_, err = io.Copy(c, reader)
	if err != nil {
		return err
	}

	return nil
}
