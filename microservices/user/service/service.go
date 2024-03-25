package service

import (
	"errors"
	"log"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/bson/primitive"
	utils "newnew.media/commons/utils"
	"newnew.media/microservices/user/dao"
	"newnew.media/microservices/user/dto"
	repository "newnew.media/microservices/user/repository"
)

type UserService struct {
	userRepo    *repository.UserRepository
	redisClient *redis.Client
	natsClient  *nats.Conn
}

func NewUserService(userRepo *repository.UserRepository, redisClient *redis.Client, natsClient *nats.Conn) *UserService {
	return &UserService{userRepo: userRepo, redisClient: redisClient, natsClient: natsClient}
}

func (s *UserService) CreateUser(user dto.CreateUserRequest) error {
	// Check if a user with the same email or Spotify ID already exists
	if err := s.CheckUserExists(user); err != nil {
		return err
	}
	newUser := dao.User{
		Password:      user.Password,
		City:          user.City,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
		EmailVerified: false, // Assuming EmailVerified should initially be false
	}

	// Set SpotifyID if provided
	if user.SpotifyID != "" {
		newUser.SpotifyID = user.SpotifyID
	}

	if user.Email != "" {
		newUser.Email = user.Email
		errChan := make(chan error, 1) // Create a channel to receive errors

		go func() {
			err := utils.SendNATSmessage(s.natsClient, "user-registered", []byte(user.Email))
			if err != nil {
				errChan <- err // Send the error to the channel
			} else {
				errChan <- nil // Send nil to indicate success
			}
		}()

		err := <-errChan // Wait for the result from the goroutine
		if err != nil {
			newUser.EmailSent = false
			log.Print("Failed to send msg to nats: ", err)
			return errors.New("failed to send msg to nats")
		} else {
			newUser.EmailSent = true
		}
	}
	// Create user in the repository
	if err := s.userRepo.CreateUser(newUser); err != nil {
		return err // Return error if user creation fails
	}
	return nil
}

func (s *UserService) GetUserByID(id primitive.ObjectID) (dao.User, error) {
	return s.userRepo.GetUserByID(id)
}

func (s *UserService) GetUserByEmail(email string) (dao.User, error) {
	return s.userRepo.GetUserByEmail(email)
}

// get all users
func (s *UserService) GetUsers() ([]dao.User, error) {
	return s.userRepo.GetUsers()
}

func (s *UserService) GetUserBySpotifyID(spotifyID string) (dao.User, error) {
	return s.userRepo.GetUserBySpotifyID(spotifyID)
}

func (s *UserService) GetUserByYouTubeID(youtubeID string) (dao.User, error) {
	return s.userRepo.GetUserByYouTubeID(youtubeID)
}

func (s *UserService) GetUserByFavoriteGenres(genres []primitive.ObjectID) ([]dao.User, error) {
	return s.userRepo.GetUserByFavoriteGenres(genres)
}

func (s *UserService) GetUserByFavoritePlaces(places []primitive.ObjectID) ([]dao.User, error) {
	return s.userRepo.GetUserByFavoritePlaces(places)
}

func (s *UserService) UpdateUser(id primitive.ObjectID, updates map[string]interface{}) error {
	return s.userRepo.UpdateUser(id, updates)
}

func (s *UserService) DeleteUser(id primitive.ObjectID) error {
	return s.userRepo.DeleteUser(id)
}

func (s *UserService) GetUserRoles(id primitive.ObjectID) ([]dao.UserRole, error) {
	userRoles, err := s.userRepo.GetUserRoles(id)
	if err != nil {
		return []dao.UserRole{}, errors.New("could not get user roles")
	}
	return userRoles, nil
}

func (s *UserService) AddRole(userID primitive.ObjectID, role dto.Role) error {

	// Validate the role
	if !isValidRole(role) {
		return errors.New("invalid Role")
	}

	userRole := dao.UserRole{
		UserID: userID,
		Role:   dao.Role(role),
	}
	// since we're checking userRoles on emailLogin,
	// we'll add the new role to user'd JWT on their next login.
	return s.userRepo.AddUserRole(userRole)
}

func (s *UserService) RemoveRole(userID primitive.ObjectID, role dto.Role) error {
	daoRole := dao.Role(role)
	return s.userRepo.RemoveUserRole(userID, daoRole)
}

// Function to validate role
func isValidRole(role dto.Role) bool {
	switch role {
	// SEC TODO: // remove dto.Admin later
	// case dto.Audience, dto.Artist, dto.Place, dto.Admin, dto.Crew:

	case dto.Audience, dto.Artist, dto.Place, dto.Admin, dto.Crew:
		return true
	default:
		return false
	}
}

func (s *UserService) CheckUserExists(user dto.CreateUserRequest) error {
	if user.Email != "" {
		if _, err := s.userRepo.GetUserByEmail(user.Email); err == nil {
			// User with the same email already exists
			return errors.New("user with the same email already exists")
		}
	}

	if user.SpotifyID != "" {
		if _, err := s.userRepo.GetUserBySpotifyID(user.SpotifyID); err == nil {
			// User with the same Spotify ID already exists
			return errors.New("user with the same Spotify ID already exists")
		}
	}

	return nil
}
