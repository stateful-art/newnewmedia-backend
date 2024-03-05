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

// func (s *PlaylistService) CalculateCollectiveRevenueSplit(playlist dao.Playlist, totalRevenue float64) (map[primitive.ObjectID]float64, error) {
// 	if playlist.Type != dao.Private || playlist.RevenueSharingModel != dao.CollectiveSharing {
// 		return nil, errors.New("playlist must be private with collective revenue sharing model")
// 	}

// 	// Deduct owner's share based on revenue cut percentage
// 	ownerShare := totalRevenue * playlist.RevenueCutPercentage / 100.0
// 	remainingRevenue := totalRevenue - ownerShare

// 	// Count the number of unique artists in the playlist
// 	uniqueArtistsCount := 0
// 	uniqueArtists := make(map[primitive.ObjectID]bool)
// 	for _, song := range playlist.Songs {
// 		if _, found := uniqueArtists[song.ArtistID]; !found {
// 			uniqueArtists[song.ArtistID] = true
// 			uniqueArtistsCount++
// 		}
// 	}

// 	// Calculate share for each artist equally
// 	artistShare := remainingRevenue / float64(uniqueArtistsCount)

// 	// Prepare map to store artist shares
// 	artistShares := make(map[primitive.ObjectID]float64)

// 	// Assign equal share to each artist
// 	for artist := range uniqueArtists {
// 		artistShares[artist] = artistShare
// 	}

// 	return artistShares, nil
// }

// func (s *PlaylistService) CalculateIndividualRevenueSplit(playlist dao.Playlist, totalRevenue float64) (map[primitive.ObjectID]float64, error) {
// 	if playlist.Type != dao.Private || playlist.RevenueSharingModel != dao.IndividualSharing {
// 		return nil, errors.New("playlist must be private with individual revenue sharing model")
// 	}

// 	// Deduct owner's share based on revenue cut percentage
// 	ownerShare := totalRevenue * playlist.RevenueCutPercentage / 100.0
// 	remainingRevenue := totalRevenue - ownerShare

// 	// Prepare map to store individual artist shares
// 	artistShares := make(map[primitive.ObjectID]float64)

// 	// Calculate total play count for all songs in the playlist
// 	totalPlayCount := float64(0)
// 	for _, song := range playlist.Songs {
// 		totalPlayCount += float64(song.PlayCount)
// 	}

// 	// Calculate share for each artist based on their song's play count
// 	for _, song := range playlist.Songs {
// 		artistShare := remainingRevenue * (float64(song.PlayCount) / totalPlayCount)
// 		artistShares[song.ArtistID] += artistShare
// 	}

// 	return artistShares, nil
// }
