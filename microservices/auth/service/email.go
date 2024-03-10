package service

import (
	"errors"
	// Import your user model and database connection
	dto "newnew.media/microservices/user/dto"
	userService "newnew.media/microservices/user/service"
)

type EmailAuthService struct {
	userService *userService.UserService
}

func NewEmailAuthService(userService *userService.UserService) *EmailAuthService {
	return &EmailAuthService{userService: userService}
}

// RegisterUser registers a new user.
func (eas *EmailAuthService) RegisterUser(user dto.CreateUserRequest) (bool, error) {
	// Check if the user already exists
	err := eas.userService.CheckUserExists(user)
	if err != nil {
		return false, errors.New("user already exists")
	}

	// Hash the password
	hashedPassword, err := HashPassword(user.Password)
	if err != nil {
		return false, err
	}

	// // Create a new user with the hashed password
	user.Password = hashedPassword

	error := eas.userService.CreateUser(user)

	return error == nil, error
}

func (eas *EmailAuthService) LoginUser(email, password string) (string, error) {
	// Retrieve the user from the database
	user, err := eas.userService.GetUserByEmail(email)
	if err != nil {
		return "", err
	}

	// Verify the password
	if !CheckPasswordHash(password, user.Password) {
		return "", errors.New("invalid password")
	}

	// Generate a JWT token for the user
	token, err := GenerateJWTToken(&user)
	if err != nil {
		return "", err
	}

	return token, nil
}
