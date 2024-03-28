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

	// Convert DAO counter offers to DTO counter offers
	counterOffersDTO := make([]dto.CounterOffer, len(offerDAO.CounterOffers))
	for i, counterOfferDAO := range offerDAO.CounterOffers {
		counterOffersDTO[i] = dto.CounterOffer{
			CounterID: counterOfferDAO.CounterID.Hex(),
			Status:    string(counterOfferDAO.Status),
		}
	}

	offerDTO := &dto.Offer{
		// Populate the DTO fields from the DAO
		ID:            offerDAO.ID.Hex(),
		Songs:         utils.ConvertObjectIDsToString(offerDAO.Songs),
		ArtistID:      offerDAO.Artist.Hex(),
		PlaceID:       offerDAO.Place.Hex(),
		ValidUntil:    offerDAO.ValidUntil.Format(layoutString),
		OfferedAt:     offerDAO.OfferedAt.Format(layoutString),
		Preferences:   preferences,
		CounterOffers: counterOffersDTO,
	}
	return offerDTO, nil
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

// func (s *OfferServiceImpl) CreateCounterOffer(ctx context.Context, counter *dto.CreateCounterOffer) (*dto.CounterOffer, error) {
// 	// Use the repository to create a counter offer
// 	return s.repo.CreateCounterOffer(ctx, counter)
// }

func (s *OfferServiceImpl) CreateCounterOffer(ctx context.Context, counter *dto.CreateCounterOffer) (*dto.CounterOffer, error) {
	// Step 1: Validate the input
	if counter.OfferID == "" {
		return nil, errors.New("offer ID is required")
	}

	offerObjID, err := primitive.ObjectIDFromHex(counter.OfferID)
	if err != nil {
		log.Printf("invalid offer id : %s", counter.OfferID)

		return nil, err
	}

	// Step 2: Create a new counter offer
	offer, err := s.repo.GetOfferByID(ctx, offerObjID)
	if err != nil {
		log.Printf("error getting offer : %s", counter.OfferID)
		return nil, err
	}

	preferences, _ := convertStringsToPreferences(counter.Preferences)

	counterOffer := &dao.Counter{
		Offer:       offer.ID,
		OfferedAt:   time.Now(),
		ValidUntil:  time.Now().AddDate(0, 0, counter.ValidUntil), // Example: valid for 7 days
		Status:      dao.Pending,
		Preferences: preferences,
		ParentOffer: offer.ID,
	}

	// Step 3: Save the counter offer
	newCounter, err := s.repo.CreateCounterOffer(ctx, counterOffer)
	if err != nil {
		return nil, err
	}

	// Step 4: Update the original offer
	counterOfferDTO := &dto.CounterOffer{
		CounterID: newCounter.ID.Hex(),
		Status:    string(newCounter.Status),
	}

	counterObjID, _ := primitive.ObjectIDFromHex(counter.OfferID)
	if err := s.repo.UpdateOfferWithCounterOffer(ctx, counterObjID, &dao.CounterOffer{
		CounterID: counterOffer.ID,
		Status:    counterOffer.Status,
	}); err != nil {
		return nil, err
	}

	// Step 5: Return the counter offer
	return counterOfferDTO, nil
}

func (s *OfferServiceImpl) GetCounterOfferByID(ctx context.Context, id string) (*dto.Counter, error) {
	// Convert the string ID to a primitive.ObjectID
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	counterOffer, err := s.repo.GetCounterOfferByID(ctx, objectID)
	if err != nil {
		log.Printf("error getting counter offer: %s", id)
		return nil, err
	}

	preferences := convertPreferencesToStrings(counterOffer.Preferences)

	var counter = dto.Counter{
		ID:          counterOffer.ID.Hex(),
		OfferID:     counterOffer.Offer.Hex(),
		OfferedAt:   counterOffer.OfferedAt.Format(layoutString),
		ValidUntil:  counterOffer.ValidUntil.Format(layoutString),
		Status:      string(counterOffer.Status),
		Preferences: preferences,
		ParentOffer: counterOffer.ParentOffer.Hex(),
	}

	// Use the repository to get the counter offer by ID
	return &counter, nil
}

func (s *OfferServiceImpl) UpdateCounterOfferStatus(ctx context.Context, id string, status string) error {
	// Convert the string ID to a primitive.ObjectID
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	// Use the repository to update the counter offer status

	dtoStatus, err := getDTOStatusFromString(status)
	if err != nil {
		return err
	}
	daoStatus := dao.Status(dtoStatus)

	return s.repo.UpdateCounterOfferStatus(ctx, objectID, daoStatus)
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
	case string(dto.Countered):
		return dto.Countered, nil
	default:
		return "", errors.New("invalid status")
	}
}
