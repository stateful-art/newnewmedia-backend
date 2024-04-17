package service

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	dao "newnew.media/microservices/offer/dao"
	dto "newnew.media/microservices/offer/dto"

	utils "newnew.media/commons/utils"
	repository "newnew.media/microservices/offer/repository"
)

var layoutString string = os.Getenv("TIME_LAYOUT_STRING")

type OfferServiceImpl struct {
	repo repository.OfferRepository
}

func NewOfferService(repo repository.OfferRepository) *OfferServiceImpl {
	return &OfferServiceImpl{repo: repo}
}

func (s *OfferServiceImpl) CreateOffer(ctx context.Context, offer *dto.CreateOffer) (*dto.CreateOfferResponse, error) {

	placeObjID, err := primitive.ObjectIDFromHex(offer.PlaceID)
	if err != nil {
		fmt.Println("Error getting  placeObjID:", err)
	}
	artistObjID, err := primitive.ObjectIDFromHex(offer.ArtistID)
	if err != nil {
		fmt.Println("Error getting  artistObjID:", err)
	}

	songs, err := utils.ConvertStringsToObjectIDs(offer.Songs)
	if err != nil {
		log.Print("could not convert song ids to objectIDs")
	}

	preferences, err := convertStringsToPreferences(offer.Preferences)
	if err != nil {
		fmt.Println(err.Error())
	}

	validUntil := time.Now().AddDate(0, 0, offer.ValidUntil)
	offerDAO := dao.Offer{
		Songs:       songs,
		Artist:      artistObjID,
		Place:       placeObjID,
		OfferedAt:   time.Now(),
		ValidUntil:  validUntil,
		Status:      dao.Pending,
		Preferences: preferences,
	}

	newOffer, err := s.repo.CreateOffer(ctx, &offerDAO)
	if err != nil {
		log.Print(err.Error())
	}

	validUntilStr := validUntil.Format(layoutString)
	createdOffer := dto.CreateOfferResponse{
		ID:         newOffer.ID.Hex(),
		Songs:      offer.Songs,
		PlaceID:    offer.PlaceID,
		OfferedAt:  newOffer.OfferedAt.String(),
		ValidUntil: validUntilStr,
		Status:     string(newOffer.Status),
	}

	return &createdOffer, nil
}

// func (s *OfferServiceImpl) GetOfferByID(ctx context.Context, id string) (*dto.Offer, error) {
// 	// Convert the string ID to a primitive.ObjectID
// 	objectID, err := primitive.ObjectIDFromHex(id)
// 	if err != nil {
// 		return nil, err
// 	}
// 	// Use the repository to get the offer by ID
// 	return s.repo.GetOfferByID(ctx, objectID)
// }

func (s *OfferServiceImpl) GetOfferByID(ctx context.Context, id string) (*dto.Offer, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	offerDAO, err := s.repo.GetOfferByID(ctx, objectID)
	if err != nil {
		return nil, err
	}

	preferences := convertPreferencesToStrings(offerDAO.Preferences)

	offerDTO := &dto.Offer{
		// Populate the DTO fields from the DAO
		ID:          offerDAO.ID.Hex(),
		Songs:       utils.ConvertObjectIDsToString(offerDAO.Songs),
		ArtistID:    offerDAO.Artist.Hex(),
		PlaceID:     offerDAO.Place.Hex(),
		ValidUntil:  offerDAO.ValidUntil.Format(layoutString),
		OfferedAt:   offerDAO.OfferedAt.Format(layoutString),
		Preferences: preferences,
	}
	return offerDTO, nil
}

func (s *OfferServiceImpl) GetOffersByPlace(ctx context.Context, id string) ([]*dto.Offer, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	offersDAO, err := s.repo.GetOffersByPlace(ctx, objectID)
	if err != nil {
		return nil, err
	}

	// Transform DAO offers to DTO offers
	dtoOffers := make([]*dto.Offer, len(offersDAO))
	for i, o := range offersDAO {
		dtoOffers[i] = convertOfferDAOToDTO(o)
	}

	return dtoOffers, nil
}

func (s *OfferServiceImpl) GetOffersByArtist(ctx context.Context, id string) ([]*dto.Offer, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	offersDAO, err := s.repo.GetOffersByArtist(ctx, objectID)
	if err != nil {
		return nil, err
	}

	// Transform DAO offers to DTO offers
	dtoOffers := make([]*dto.Offer, len(offersDAO))
	for i, o := range offersDAO {
		dtoOffers[i] = convertOfferDAOToDTO(o)
	}

	return dtoOffers, nil
}

func (s *OfferServiceImpl) UpdateOfferStatus(ctx context.Context, id string, status string) error {
	// Convert the string ID to a primitive.ObjectID
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	dtoStatus, err := getDTOStatusFromString(status)
	if err != nil {
		return err
	}
	daoStatus := dao.Status(dtoStatus)

	// Use the repository to update the offer status
	return s.repo.UpdateOfferStatus(ctx, objectID, daoStatus)
}

func (s *OfferServiceImpl) DeleteOffer(ctx context.Context, id string) error {
	// Convert the string ID to a primitive.ObjectID
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	// Use the repository to delete the offer
	return s.repo.DeleteOffer(ctx, objectID)
}

func convertStringsToPreferences(strings []string) ([]dao.Preference, error) {
	var preferences []dao.Preference

	for _, str := range strings {
		switch dao.Preference(str) {
		case dao.Public:
			preferences = append(preferences, dao.Public)
		case dao.Private:
			preferences = append(preferences, dao.Private)
		case dao.Collective:
			preferences = append(preferences, dao.Collective)
		case dao.Individual:
			preferences = append(preferences, dao.Individual)
		default:
			return nil, fmt.Errorf("invalid preference: %s", str)
		}
	}

	return preferences, nil
}

func convertPreferencesToStrings(preferences []dao.Preference) []string {
	var strings []string

	for _, preference := range preferences {
		strings = append(strings, string(preference))
	}

	return strings
}

// stringToStatus converts a string to a dto.Status.
func getDTOStatusFromString(status string) (dto.Status, error) {
	switch status {
	case string(dto.Pending):
		return dto.Pending, nil
	case string(dto.Accepted):
		return dto.Accepted, nil
	case string(dto.Rejected):
		return dto.Rejected, nil

	default:
		return "", errors.New("invalid status")
	}
}

func convertOfferDAOToDTO(offer *dao.Offer) *dto.Offer {
	preferences := convertPreferencesToStrings(offer.Preferences)

	return &dto.Offer{
		ID:          offer.ID.Hex(),
		Songs:       utils.ConvertObjectIDsToString(offer.Songs),
		ArtistID:    offer.Artist.Hex(),
		PlaceID:     offer.Place.Hex(),
		ValidUntil:  offer.ValidUntil.Format(layoutString),
		OfferedAt:   offer.OfferedAt.Format(layoutString),
		Preferences: preferences,
	}
}
