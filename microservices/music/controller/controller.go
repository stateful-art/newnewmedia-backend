package controller

import (
	"errors"
	"log"
	"strings"

	"cloud.google.com/go/storage"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/oauth2"
	dto "newnewmedia.com/microservices/music/dto"
	service "newnewmedia.com/microservices/music/service"
)

func GetSong(c *fiber.Ctx, storageClient *storage.Client) (dto.MusicRetrieve, error) {
	songID := c.Params("id")
	log.Print(songID)

	// Fetch the audio file path for the given song ID
	song, err := service.GetSong(songID)
	if err != nil {
		return dto.MusicRetrieve{}, err
	}

	// Return the song details as JSON
	c.JSON(song)

	// Since we've already sent the response, return an empty DTO and nil error
	return dto.MusicRetrieve{}, nil
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
	var musicPayload dto.MusicCreate

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

// SPOTIFY RELATED ENDPOINTS
// RecentlyPlayedSongs retrieves the user's recently played songs
func RecentlyPlayedSongs(c *fiber.Ctx) error {

	// accessToken := c.Get("x-spotify-token")
	// if accessToken == "" {
	// 	// Handle the case where the access token is not provided in the header
	// 	return fiber.NewError(fiber.StatusBadRequest, "Access token not provided")
	// }

	// Convert the access token to an oauth2.Token struct
	// token := &oauth2.Token{
	// 	AccessToken: accessToken,
	// }

	// // Set the token in the Fiber context
	// c.Locals("oauth_token", token)

	token, err := getOauthFromHeader(c)
	// token := c.Locals("access_token").(*oauth2.Token)
	tracks, err := service.FetchRecentlyPlayedSongs(token)
	if err != nil {
		// Handle error
		return err
	}
	return c.JSON(tracks)
}

// UserPlaylists retrieves the user's playlists
func UserPlaylists(c *fiber.Ctx) error {
	token, err := getOauthFromHeader(c)

	playlists, err := service.FetchUserPlaylists(token)
	if err != nil {
		// Handle error
		return err
	}
	return c.JSON(playlists)
}

// GenreAnalysis generates genre analysis based on the user's listening history
func GenreAnalysis(c *fiber.Ctx) error {
	// Retrieve user's access token from the context
	token, err := getOauthFromHeader(c)

	// Fetch user's recently played tracks from Spotify
	tracks, err := service.FetchRecentlyPlayedSongs(token)
	if err != nil {
		// Handle error
		return err
	}

	// Generate genre analysis based on the user's listening history
	genreAnalysis := service.GenerateGenreAnalysis(tracks, token)
	return c.JSON(fiber.Map{"genreAnalysis": genreAnalysis})
}

// getSpotifyAccessToken retrieves the Spotify access token from the request headers
func getOauthFromHeader(c *fiber.Ctx) (*oauth2.Token, error) {
	accessToken := c.Get("x-spotify-token")
	if accessToken == "" {
		return nil, fiber.NewError(fiber.StatusBadRequest, "Access token not provided")
	}

	token := &oauth2.Token{
		AccessToken: accessToken,
	}

	return token, nil
}
