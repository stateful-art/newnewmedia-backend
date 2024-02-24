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

// GetAudioFilePath retrieves the audio file path based on song ID
func GetAudioFilePath(songID string) (string, error) {
	// Convert the songID string to a primitive.ObjectID
	objectID, err := primitive.ObjectIDFromHex(songID)
	if err != nil {
		return "", fmt.Errorf("invalid song ID: %v", err)
	}

	// Fetch the music details from the database based on the converted song ID
	music, err := repository.GetMusicById(objectID)
	if err != nil {
		return "", err
	}

	// Return the audio file path
	return music.Path, nil
}

func GetSong(songID string) (dto.MusicRetrieve, error) {
	objectID, err := primitive.ObjectIDFromHex(songID)
	if err != nil {
		return dto.MusicRetrieve{}, fmt.Errorf("invalid song ID: %v", err)
	}

	// Fetch the music details from the database based on the converted song ID
	dao, err := repository.GetMusicById(objectID)
	if err != nil {
		return dto.MusicRetrieve{}, err
	}

	// Create a DTO instance
	song := dto.MusicRetrieve{
		ID:     dao.ID,
		Name:   dao.Name,
		Artist: dao.Artist,
		Path:   dao.Path,
	}

	// Return the audio file path
	return song, nil
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

func CreateMusic(c *fiber.Ctx, musicPayload dto.MusicCreate, audioFile *multipart.FileHeader, storageClient *storage.Client) error {
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
	if err := storeMusicFileInGCS(c.Context(), storageBucket, objectName, audioFileContent, storageClient); err != nil {
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
func storeMusicFileInGCS(ctx context.Context, bucketName, objectName string, audioFile io.Reader, storageClient *storage.Client) error {
	// Create a new GCS client

	// Create a GCS object writer
	writer := storageClient.Bucket(bucketName).Object(objectName).NewWriter(ctx)
	defer writer.Close()

	// Copy the file data to the GCS object writer
	if _, err := io.Copy(writer, audioFile); err != nil {
		return err
	}

	return nil
}

// StreamMusic streams the audio file from Google Cloud Storage to the client
func StreamMusic(c *fiber.Ctx, bucketName, objectName string, storageClient *storage.Client) error {
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
