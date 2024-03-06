package service

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"newnewmedia.com/microservices/user/dao"
	"newnewmedia.com/microservices/user/dto"

	repository "newnewmedia.com/microservices/user/repository"
)

type UserService struct {
	userRepo *repository.UserRepository
}

func NewUserService(userRepo *repository.UserRepository) *UserService {
	return &UserService{userRepo: userRepo}
}

func (s *UserService) CreateUser(user dao.User) error {
	return s.userRepo.CreateUser(user)
}

func (s *UserService) GetUserByID(id primitive.ObjectID) (dao.User, error) {
	return s.userRepo.GetUserByID(id)
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
