package service

import (
	"context"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"strconv"
	"strings"

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

func CreateMusic(c *fiber.Ctx, musicPayload dto.Music, audioFile *multipart.FileHeader, storageClient *storage.Client) error {
	// Validate if the music payload contains the necessary fields
	if musicPayload.Name == "" || musicPayload.Artist == "" {
		return fmt.Errorf("Name and Artist are required")
	}

	fmt.Print(musicPayload)

	// Open the uploaded audio file
	audioFileContent, err := audioFile.Open()
	if err != nil {
		return err
	}
	defer audioFileContent.Close()

	storageBucket := "n2media-music"
	playlistStorageBucket := "n2media-playlists"

	objectName := fmt.Sprintf("%s.mp3", primitive.NewObjectID().Hex())           // Generate a unique object name
	songObjectName := fmt.Sprintf("songs/%s.mp3", primitive.NewObjectID().Hex()) // Generate a unique object name
	// Extract the audio file name without any directory prefix
	// audioObjectName := objectName

	if err := storeMusicFileInGCS(c.Context(), storageBucket, songObjectName, audioFileContent, storageClient); err != nil {
		return err
	}

	// Store the HLS playlists in Google Cloud Storage under "playlists" folder and get their paths
	playlistPaths, err := generateHLSPlaylists(c.Context(), playlistStorageBucket, objectName, storageClient)
	if err != nil {
		return err
	}

	// Store the music details in MongoDB, including the playlist paths
	music := dao.Music{
		Name:          musicPayload.Name,
		Artist:        musicPayload.Artist,
		Path:          fmt.Sprintf("gs://%s/%s", storageBucket, songObjectName),                            // Update with absolute URI without duplicate "songs"
		PlaylistPaths: generatePlaylistAbsolutePaths(playlistStorageBucket, songObjectName, playlistPaths), // Generate absolute paths for playlists
		// Add additional fields if needed
	}

	if err := repository.CreateMusic(music); err != nil {
		return err
	}

	return nil
}

// generatePlaylistAbsolutePaths generates absolute paths for HLS playlists with the "/playlists" prefix
func generatePlaylistAbsolutePaths(bucketName, audioObjectName string, playlistPaths []string) []string {
	var absolutePaths []string
	for _, path := range playlistPaths {
		// Prepend the "/playlists" prefix to each playlist path
		absolutePath := fmt.Sprintf("gs://%s/%s", bucketName, path)
		absolutePaths = append(absolutePaths, absolutePath)
	}
	return absolutePaths
}

// storeMusicFileInGCS stores the music file in Google Cloud Storage
func storeMusicFileInGCS(ctx context.Context, bucketName, objectName string, audioFile io.Reader, storageClient *storage.Client) error {
	// Create a GCS object writer
	writer := storageClient.Bucket(bucketName).Object(objectName).NewWriter(ctx)
	defer writer.Close()

	// Copy the file data to the GCS object writer
	if _, err := io.Copy(writer, audioFile); err != nil {
		return err
	}

	return nil
}

// generateHLSPlaylists generates HLS playlists for different bitrates
func generateHLSPlaylists(ctx context.Context, bucketName, audioObjectName string, storageClient *storage.Client) ([]string, error) {
	// Create a new GCS object writer for each HLS playlist
	var playlistPaths []string
	bitrates := []int{128, 256, 512, 1024, 2048, 4096}
	for _, bitrate := range bitrates {
		playlistName := fmt.Sprintf("%s_%dk.m3u8", audioObjectName, bitrate)
		path, err := generateHLSPlaylist(ctx, bucketName, audioObjectName, playlistName, bitrate, storageClient)
		if err != nil {
			return nil, err
		}
		playlistPaths = append(playlistPaths, path)
	}
	return playlistPaths, nil
}

// generateHLSPlaylist generates a single HLS playlist with the specified bitrate
func generateHLSPlaylist(ctx context.Context, bucketName, audioObjectName, playlistName string, bitrate int, storageClient *storage.Client) (string, error) {
	// Create a new GCS object writer for the playlist under the "playlists" directory
	writer := storageClient.Bucket(bucketName).Object(playlistName).NewWriter(ctx)
	defer writer.Close()

	// Write the HLS playlist content with the specified bitrate
	playlistContent := []byte(fmt.Sprintf(`#EXTM3U
											#EXT-X-VERSION:3
											#EXT-X-STREAM-INF:BANDWIDTH=%dk
											%s_%dk.m3u8`,
		bitrate, audioObjectName, bitrate))
	_, err := writer.Write(playlistContent)
	if err != nil {
		return "", err
	}

	// Return the playlist name (without the "playlists/" prefix) to be stored in the database
	return playlistName, nil
}

// func StreamMusic(c *fiber.Ctx, bucketName string, playlistPath string, storageClient *storage.Client) error {
// 	ctx := context.Background()

// 	reader, err := storageClient.Bucket(bucketName).Object("playlists/" + playlistPath).NewReader(ctx) // Update playlist path
// 	if err != nil {
// 		return err
// 	}
// 	defer reader.Close()

// 	c.Set("Content-Type", "application/vnd.apple.mpegurl")
// 	c.Set("Accept-Ranges", "bytes")

// 	rangeHeader := c.Get("Range")
// 	if rangeHeader != "" {
// 		startByte, endByte, totalFileSize := parseRangeHeader(rangeHeader)

// 		c.Status(fiber.StatusPartialContent)
// 		c.Set("Content-Range", fmt.Sprintf("bytes %d-%d/%d", startByte, endByte, totalFileSize))

// 		bytesToRead := endByte - startByte + 1
// 		_, err = io.CopyN(c, reader, bytesToRead)
// 		if err != nil {
// 			return err
// 		}
// 	} else {
// 		_, err = io.Copy(c, reader)
// 		if err != nil {
// 			return err
// 		}
// 	}

//		return nil
//	}

// // parseRangeHeader parses the Range header and returns the start byte, end byte, and total file size.
// func parseRangeHeader(rangeHeader string) (startByte, endByte, totalFileSize int64) {
// 	parts := strings.Split(rangeHeader, "=")
// 	if len(parts) != 2 || parts[0] != "bytes" {
// 		return 0, 0, 0
// 	}

// 	byteRange := strings.Split(parts[1], "-")
// 	if len(byteRange) != 2 {
// 		return 0, 0, 0
// 	}

// 	var err error
// 	startByte, err = strconv.ParseInt(byteRange[0], 10, 64)
// 	if err != nil {
// 		return 0, 0, 0
// 	}

// 	endByte, err = strconv.ParseInt(byteRange[1], 10, 64)
// 	if err != nil {
// 		return 0, 0, 0
// 	}

// 	return startByte, endByte, endByte - startByte + 1
// }

// func StreamMusic(c *fiber.Ctx, bucketName string, objectName string, storageClient *storage.Client) error {
// 	ctx := context.Background()
// 	// Get the Range header from the request
// 	rangeHeader := c.Get("Range")

// 	// Create a new reader to read the streamed data from Google Cloud Storage for the specified range
// 	attrs, err := storageClient.Bucket(bucketName).Object(objectName).Attrs(ctx)
// 	if err != nil {
// 		log.Println("Error getting object attributes:", err)
// 		return err
// 	}

// 	var start, end int64

// 	if rangeHeader != "" {
// 		start, end, err = parseRangeHeader(rangeHeader)
// 		if err != nil {
// 			log.Println("Error parsing range header:", err)
// 			return err
// 		}
// 	} else {
// 		// If no range is specified, serve the entire file
// 		start = 0
// 		end = attrs.Size - 1
// 	}

// 	// Set the Content-Type header
// 	c.Set("Content-Type", "application/vnd.apple.mpegurl")
// 	c.Set("Accept-Ranges", "bytes")

// 	// Set the status code to 206 (Partial Content)
// 	c.Status(fiber.StatusPartialContent)

// 	// Set the Content-Range header for the partial content
// 	c.Set("Content-Range", fmt.Sprintf("bytes %d-%d/%d", start, end, attrs.Size))

// 	// Create a new reader to read the specified range of bytes from Google Cloud Storage
// 	reader, err := storageClient.Bucket(bucketName).Object(objectName).NewRangeReader(ctx, start, end-start+1)
// 	if err != nil {
// 		log.Println("Error creating range reader:", err)
// 		return err
// 	}
// 	defer reader.Close()

// 	// Copy the specified range of bytes to the response
// 	if _, err := io.Copy(c, reader); err != nil {
// 		log.Println("Error copying bytes:", err)
// 		return err
// 	}

// 	return nil
// }

func StreamMusic(c *fiber.Ctx, bucketName string, objectName string, storageClient *storage.Client) error {
	ctx := context.Background()
	// Get the Range header from the request
	rangeHeader := c.Get("Range")

	// Create a new reader to read the streamed data from Google Cloud Storage for the specified range
	attrs, err := storageClient.Bucket(bucketName).Object(objectName).Attrs(ctx)
	if err != nil {
		log.Println("Error getting object attributes:", err)
		return err
	}

	var end int64

	if rangeHeader != "" {
		_, end, err = parseRangeHeader(rangeHeader)
		if err != nil {
			log.Println("Error parsing range header:", err)
			return err
		}
	} else {
		// If no range is specified, serve the entire file
		end = attrs.Size - 1
	}

	// Set the Content-Type header
	c.Set("Content-Type", "application/vnd.apple.mpegurl")
	c.Set("Accept-Ranges", "bytes")

	// Set the status code to 206 (Partial Content)
	c.Status(fiber.StatusPartialContent)

	// Set the Content-Range header for the partial content
	c.Set("Content-Range", fmt.Sprintf("bytes 0-%d/%d", end, attrs.Size))

	// Create a new reader to read the specified range of bytes from Google Cloud Storage
	reader, err := storageClient.Bucket(bucketName).Object(objectName).NewRangeReader(ctx, 0, end+1)
	if err != nil {
		log.Println("Error creating range reader:", err)
		return err
	}
	defer reader.Close()

	// Copy the specified range of bytes to the response
	if _, err := io.Copy(c, reader); err != nil {
		log.Println("Error copying bytes:", err)
		return err
	}

	// TODO:// currently this bucket cannot be made public since project-level access is overwriting it to private.
	// Include HLS URL in the response
	hlsURL := fmt.Sprintf("https://storage.cloud.google.com/n2media-playlists/%s", objectName)
	//https://storage.cloud.google.com/n2media-playlists/65d3e3e20b6fb299a0854112.mp3_1024k.m3u8
	_, err = c.WriteString(hlsURL)
	if err != nil {
		log.Println("Error writing HLS URL to response:", err)
		return err
	}

	return nil
}

// parseRangeHeader parses the Range header value and returns the start and end bytes of the requested range.
func parseRangeHeader(rangeHeader string) (start, end int64, err error) {
	const prefix = "bytes="
	r := strings.SplitN(rangeHeader[len(prefix):], "-", 2)
	start, err = strconv.ParseInt(r[0], 10, 64)
	if err != nil {
		return 0, 0, err
	}
	if len(r) > 1 && r[1] != "" {
		end, err = strconv.ParseInt(r[1], 10, 64)
		if err != nil {
			return 0, 0, err
		}
	}
	return start, end, nil
}
