package service

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"newnewmedia.com/microservices/playlist/dao"
	"newnewmedia.com/microservices/playlist/repository"
)

type PlaylistService struct {
	playlistRepo *repository.PlaylistRepository
}

func NewPlaylistService(playlistRepo *repository.PlaylistRepository) *PlaylistService {
	return &PlaylistService{playlistRepo: playlistRepo}
}

func (s *PlaylistService) CreatePlaylist(playlist dao.Playlist) error {
	return repository.CreatePlaylist(playlist)
}

func (s *PlaylistService) GetPlaylistByID(id primitive.ObjectID) (dao.Playlist, error) {
	return repository.GetPlaylistByID(id)
}

func (s *PlaylistService) GetPlaylists() ([]dao.Playlist, error) {
	return repository.GetPlaylists()
}

func (s *PlaylistService) UpdatePlaylist(id primitive.ObjectID, playlist dao.Playlist) error {
	return repository.UpdatePlaylist(id, playlist)
}

func (s *PlaylistService) DeletePlaylist(id primitive.ObjectID) error {
	return repository.DeletePlaylist(id)
}
