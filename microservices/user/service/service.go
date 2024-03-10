package service

import (
	"errors"
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
		Email:         user.Email,
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
		err := utils.SendNATSmessage(s.natsClient, "user-registered", user.Email)

		if err != nil {
			newUser.EmailSent = false
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

func (s *UserService) UpdateUser(id primitive.ObjectID, user dao.User) error {
	return s.userRepo.UpdateUser(id, user)
}

func (s *UserService) DeleteUser(id primitive.ObjectID) error {
	return s.userRepo.DeleteUser(id)
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
	return s.userRepo.AddUserRole(userRole)
}

func (s *UserService) RemoveRole(userID primitive.ObjectID, role dto.Role) error {
	daoRole := dao.Role(role)
	return s.userRepo.RemoveUserRole(userID, daoRole)
}

// Function to validate role
func isValidRole(role dto.Role) bool {
	switch role {
	case dto.Audience, dto.Artist, dto.Place, dto.Admin, dto.Crew:
		return true
	default:
		return false
	}
}

func (s *UserService) CheckUserExists(user dto.CreateUserRequest) error {
	if user.Email != "" {
		_, err := s.userRepo.GetUserByEmail(user.Email)
		if err == nil {
			// User with the same email already exists
			return errors.New("user with the same email already exists")
		}
	}

	if user.SpotifyID != "" {
		_, err := s.userRepo.GetUserBySpotifyID(user.SpotifyID)
		if err == nil {
			// User with the same Spotify ID already exists
			return errors.New("user with the same Spotify ID already exists")
		}
	}

	return nil
}
