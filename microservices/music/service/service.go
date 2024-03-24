package service

import (
	"context"
	"fmt"
	"io"
	"log"
	"mime/multipart"

	"cloud.google.com/go/storage"
	"github.com/gofiber/fiber/v2"
	"github.com/zmb3/spotify"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/oauth2"
	"newnew.media/microservices/music/dao"
	"newnew.media/microservices/music/dto"
	repository "newnew.media/microservices/music/repository"
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
		ID:     dao.ID.Hex(),
		Name:   dao.Name,
		Artist: dao.Artist,
		Path:   dao.Path,
	}

	// Return the audio file path
	return song, nil
}

func GetMusicByPlace(id string) ([]dao.Song, error) {
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

func CreateMusic(c *fiber.Ctx, musicPayload dto.CreateMusic, audioFile *multipart.FileHeader, storageClient *storage.Client) (dto.MusicRetrieve, error) {
	// Validate if the music payload contains the necessary fields
	if musicPayload.Name == "" || musicPayload.Artist == "" {
		return dto.MusicRetrieve{}, fmt.Errorf("name and Artist are required")
	}

	// Open the uploaded audio file
	audioFileContent, err := audioFile.Open()
	if err != nil {
		return dto.MusicRetrieve{}, err
	}
	defer audioFileContent.Close()

	// Store the music file in Google Cloud Storage
	storageBucket := "n2media-music"
	objectName := fmt.Sprintf("%s.mp3", primitive.NewObjectID().Hex()) // Generate a unique object name
	if err := storeMusicFileInGCS(c.Context(), storageBucket, objectName, audioFileContent, storageClient); err != nil {
		return dto.MusicRetrieve{}, err
	}
	// Store the music details in MongoDB
	music := dao.Song{
		Name:   musicPayload.Name,
		Artist: musicPayload.Artist,
		Path:   fmt.Sprintf("gs://%s/%s", storageBucket, objectName), // GCS object path
		Genres: musicPayload.Genres,
		// Add additional fields if needed
	}
	log.Print(music)
	var song dao.Song
	if song, err = repository.CreateMusic(music); err != nil {
		return dto.MusicRetrieve{}, err
	}

	createSongDTO := dto.MusicRetrieve{
		ID:     song.ID.Hex(),
		Name:   song.Name,
		Artist: song.Artist,
		Path:   song.Path,
		Genres: song.Genres,
	}

	return createSongDTO, nil
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

// FetchRecentlyPlayedSongs fetches the user's recently played songs.
func FetchRecentlyPlayedSongs(token *oauth2.Token) ([]spotify.RecentlyPlayedItem, error) {
	// Create Spotify client with the obtained access token
	client := spotify.NewAuthenticator("").NewClient(token)

	// Retrieve user's recently played tracks
	tracks, err := client.PlayerRecentlyPlayed()
	if err != nil {
		return nil, err
	}

	return tracks, nil
}

// FetchUserPlaylists fetches the user's playlists.
func FetchUserPlaylists(token *oauth2.Token) ([]spotify.SimplePlaylist, error) {
	// Create Spotify client with the obtained access token
	client := spotify.NewAuthenticator("").NewClient(token)

	// Retrieve user's playlists
	playlists, err := client.CurrentUsersPlaylists()
	if err != nil {
		return nil, err
	}

	return playlists.Playlists, nil
}

// GenerateGenreAnalysis generates genre analysis based on the user's listening history on spotify.
//
//	genreStats := map[string]int{
//	    "rock":    10,
//	    "pop":     5,
//	    "jazz":    3,
//	    // Other genres...
//	}

func GenerateGenreAnalysis(tracks []spotify.RecentlyPlayedItem, token *oauth2.Token) map[string]int {
	genreStats := make(map[string]int)

	// Iterate through each recently played track
	for _, track := range tracks {
		// Fetch detailed information about the track's artists
		for _, artist := range track.Track.Artists {
			// Fetch genres for the artist
			artistInfo, err := fetchArtistInfo(token, string(artist.ID))
			if err != nil {
				// Log or handle error
				continue
			}

			// Increment genre count for each genre of the artist
			for _, genre := range artistInfo.Genres {
				genreStats[genre]++
			}
		}
	}

	return genreStats
}

// fetchArtistInfo retrieves detailed information about an artist from the Spotify API,
// including genres.
func fetchArtistInfo(token *oauth2.Token, artistID string) (*spotify.FullArtist, error) {
	client := spotify.NewAuthenticator("").NewClient(token)

	// Make a request to the Spotify API to fetch detailed information about the artist
	// using the artistID.
	artist, err := client.GetArtist(spotify.ID(artistID))
	if err != nil {
		return nil, err
	}
	return artist, nil
}
