package service

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"newnew.media/microservices/playlist/dao"

	"newnew.media/microservices/playlist/repository"
)

type PlaylistService struct {
	playlistRepo *repository.PlaylistRepository
}

func NewPlaylistService(playlistRepo *repository.PlaylistRepository) *PlaylistService {
	return &PlaylistService{playlistRepo: playlistRepo}
}

func (s *PlaylistService) CreatePlaylist(playlist dao.Playlist) error {
	return s.playlistRepo.CreatePlaylist(playlist)
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
