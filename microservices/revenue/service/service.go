package service

import (
	"errors"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"newnew.media/microservices/revenue/dao"
	"newnew.media/microservices/revenue/dto"

	"newnew.media/microservices/revenue/repository"
)

type RevenueService struct {
	revenueRepo *repository.RevenueRepository
}

func NewRevenueService(revenueRepo *repository.RevenueRepository) *RevenueService {
	return &RevenueService{revenueRepo: revenueRepo}
}

// CreateRevenue inserts a new revenue entry into the database.
func (s *RevenueService) CreateRevenue(revenue dao.Revenue) error {

	return repository.CreateRevenue(revenue)
}

// GetRevenueByID retrieves revenue by its ID.
func (s *RevenueService) GetRevenueByID(id primitive.ObjectID) (dao.Revenue, error) {
	return repository.GetRevenueByID(id)
}

// GetRevenueByArtistID retrieves revenue entries by artist ID.
func (s *RevenueService) GetRevenueByArtistID(artistID primitive.ObjectID) ([]dao.Revenue, error) {
	return repository.GetRevenueByArtistID(artistID)
}

// GetRevenueByPlaceID retrieves revenue entries by place ID.
func (s *RevenueService) GetRevenueByPlaceID(placeID primitive.ObjectID) ([]dao.Revenue, error) {
	return repository.GetRevenueByPlaceID(placeID)
}

// GetRevenueByPlaylistID retrieves revenue entries by playlist ID.
func (s *RevenueService) GetRevenueByPlaylistID(playlistID primitive.ObjectID) ([]dao.Revenue, error) {
	return repository.GetRevenueByPlaylistID(playlistID)
}

func (s *RevenueService) CalculateCollectiveRevenueSplit(playlist dto.Playlist, totalRevenue float64) (map[primitive.ObjectID]float64, error) {
	if playlist.Type != dto.Private || playlist.RevenueSharingModel != dto.CollectiveSharing {
		return nil, errors.New("playlist must be private with collective revenue sharing model")
	}

	// Deduct owner's share based on revenue cut percentage
	ownerShare := totalRevenue * playlist.RevenueCutPercentage / 100.0
	remainingRevenue := totalRevenue - ownerShare

	// Count the number of unique artists in the playlist
	uniqueArtistsCount := 0
	uniqueArtists := make(map[primitive.ObjectID]bool)
	for _, song := range playlist.Songs {
		if _, found := uniqueArtists[song.ArtistID]; !found {
			uniqueArtists[song.ArtistID] = true
			uniqueArtistsCount++
		}
	}

	// Calculate share for each artist equally
	artistShare := remainingRevenue / float64(uniqueArtistsCount)

	// Prepare map to store artist shares
	artistShares := make(map[primitive.ObjectID]float64)

	// Assign equal share to each artist
	for artist := range uniqueArtists {
		artistShares[artist] = artistShare
	}

	return artistShares, nil
}

func (s *RevenueService) CalculateIndividualRevenueSplit(playlist dto.Playlist, totalRevenue float64) (map[primitive.ObjectID]float64, error) {
	if playlist.Type != dto.Private || playlist.RevenueSharingModel != dto.IndividualSharing {
		return nil, errors.New("playlist must be private with individual revenue sharing model")
	}

	// Deduct owner's share based on revenue cut percentage
	ownerShare := totalRevenue * playlist.RevenueCutPercentage / 100.0
	remainingRevenue := totalRevenue - ownerShare

	// Prepare map to store individual artist shares
	artistShares := make(map[primitive.ObjectID]float64)

	// Calculate total play count for all songs in the playlist
	totalPlayCount := float64(0)
	for _, song := range playlist.Songs {
		totalPlayCount += float64(song.PlayCount)
	}

	// Calculate share for each artist based on their song's play count
	for _, song := range playlist.Songs {
		artistShare := remainingRevenue * (float64(song.PlayCount) / totalPlayCount)
		artistShares[song.ArtistID] += artistShare
	}

	return artistShares, nil
}
