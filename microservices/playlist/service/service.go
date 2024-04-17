package service

import (
	"time"

	"github.com/gofiber/fiber/v2/log"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"newnew.media/microservices/playlist/dao"
	"newnew.media/microservices/playlist/dto"

	"newnew.media/microservices/playlist/repository"
)

type PlaylistService struct {
	playlistRepo *repository.PlaylistRepository
}

func NewPlaylistService(playlistRepo *repository.PlaylistRepository) *PlaylistService {
	return &PlaylistService{playlistRepo: playlistRepo}
}

func (s *PlaylistService) CreatePlaylist(playlist dto.CreatePlaylist) error {
	newPlaylist := ConvertPlaylistDTOToDAO(playlist)
	return s.playlistRepo.CreatePlaylist(*newPlaylist)
}

func (s *PlaylistService) GetPlaylistByID(id primitive.ObjectID) (dao.Playlist, error) {
	return s.playlistRepo.GetPlaylistByID(id)
}

func (s *PlaylistService) GetPlaylists() ([]dao.Playlist, error) {
	return s.playlistRepo.GetPlaylists()
}

func (s *PlaylistService) GetPlaylistsByOwner(ownerID primitive.ObjectID) ([]dao.Playlist, error) {
	return s.playlistRepo.GetPlaylistsByOwner(ownerID)
}

func (s *PlaylistService) UpdatePlaylist(id primitive.ObjectID, playlist dao.Playlist) error {
	return s.playlistRepo.UpdatePlaylist(id, playlist)
}

func (s *PlaylistService) DeletePlaylist(id primitive.ObjectID) error {
	return s.playlistRepo.DeletePlaylist(id)
}

func (s *PlaylistService) AddSongsToPlaylist(playlistID primitive.ObjectID, songIDs []primitive.ObjectID) error {
	// Check if the playlist exists
	playlist, err := s.playlistRepo.GetPlaylistByID(playlistID)
	if err != nil {
		return err
	}

	// Create Song objects for each songID and append them to the playlist
	for _, songID := range songIDs {
		newSong := dao.Song{ID: songID}
		playlist.Songs = append(playlist.Songs, newSong)
	}

	// Update the playlist with the new songs
	err = s.playlistRepo.UpdatePlaylist(playlistID, playlist)
	if err != nil {
		return err
	}

	return nil
}

func (s *PlaylistService) RemoveSongsFromPlaylist(playlistID primitive.ObjectID, songIDs []primitive.ObjectID) error {
	// Check if the playlist exists
	playlist, err := s.playlistRepo.GetPlaylistByID(playlistID)
	if err != nil {
		return err
	}

	// Create a map to store the songIDs to be removed for efficient lookup
	songIDMap := make(map[primitive.ObjectID]bool)
	for _, songID := range songIDs {
		songIDMap[songID] = true
	}

	// Filter out the songs to be removed from the playlist's list of songs
	var filteredSongs []dao.Song
	for _, song := range playlist.Songs {
		if !songIDMap[song.ID] {
			filteredSongs = append(filteredSongs, song)
		}
	}

	// Update the playlist with the filtered list of songs
	playlist.Songs = filteredSongs
	err = s.playlistRepo.UpdatePlaylist(playlistID, playlist)
	if err != nil {
		return err
	}

	return nil
}

func ConvertPlaylistDTOToDAO(playlist dto.CreatePlaylist) *dao.Playlist {

	ownerID, err := primitive.ObjectIDFromHex(playlist.Owner)
	if err != nil {
		log.Error("owner id not valid.")
	}
	songs := ConvertSongDTOToDAO(playlist.Songs)
	return &dao.Playlist{
		ID:                   primitive.NilObjectID, // Will be generated automatically by MongoDB
		Name:                 playlist.Name,
		Description:          playlist.Description,
		Owner:                ownerID, // Update with actual owner ID from authentication
		Type:                 dao.PlaylistType(playlist.Type),
		RevenueSharingModel:  dao.RevenueSharingModel(playlist.RevenueSharingModel),
		RevenueCutPercentage: playlist.RevenueCutPercentage,
		Songs:                songs,
		CreatedAt:            time.Now(),
		UpdatedAt:            time.Now(),
	}
}

func ConvertPlaylistDAOToDTO(playlist dao.Playlist) *dto.GetPlaylist {
	songs := ConvertSongDAOToDTO(playlist.Songs)
	return &dto.GetPlaylist{
		ID:                   playlist.ID.Hex(),
		Name:                 playlist.Name,
		Description:          playlist.Description,
		Owner:                playlist.Owner.Hex(),
		Type:                 dto.PlaylistType(playlist.Type),
		RevenueSharingModel:  dto.RevenueSharingModel(playlist.RevenueSharingModel),
		RevenueCutPercentage: playlist.RevenueCutPercentage,
		Songs:                songs,
	}
}

func ConvertSongDTOToDAO(songs []dto.Song) []dao.Song {
	var daoSongs []dao.Song
	for _, song := range songs {
		daoSong := dao.Song{
			ID:        primitive.NilObjectID, // Will be generated automatically by MongoDB
			Name:      song.Name,
			Artist:    primitive.NilObjectID, // Update with actual artist ID from your database
			PlayCount: song.PlayCount,
		}
		daoSongs = append(daoSongs, daoSong)
	}
	return daoSongs
}

func ConvertSongDAOToDTO(songs []dao.Song) []dto.Song {
	var dtoSongs []dto.Song
	for _, song := range songs {
		dtoSong := dto.Song{
			ID:        song.ID.Hex(), // Will be generated automatically by MongoDB
			Name:      song.Name,
			Artist:    song.Artist.Hex(), // Update with actual artist ID from your database
			PlayCount: song.PlayCount,
		}
		dtoSongs = append(dtoSongs, dtoSong)
	}
	return dtoSongs
}
