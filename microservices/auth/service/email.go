package service

import (
	"errors"
	// Import your user model and database connection

	userDTO "newnew.media/microservices/user/dto"
	userService "newnew.media/microservices/user/service"
)

type EmailAuthService struct {
	userService *userService.UserService
}

func NewEmailAuthService(userService *userService.UserService) *EmailAuthService {
	return &EmailAuthService{userService: userService}
}

// RegisterUser registers a new user.
func (eas *EmailAuthService) RegisterUser(user userDTO.CreateUserRequest) (bool, error) {
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

	if !user.EmailVerified {
		return "", errors.New("please verify your email")
	}
	// Verify the password
	if !CheckPasswordHash(password, user.Password) {
		return "", errors.New("invalid password")
	}

	userRoles, err := eas.userService.GetUserRoles(user.ID)
	if err != nil {
		return "", err
	}
	// var roles []userDAO.Role
	// roles = append(roles, userDAO.Audience)
	// Generate a JWT token for the user

	var roles []string
	for _, userRole := range userRoles {
		roles = append(roles, string(userRole.Role))
	}

	token, err := GenerateJWTToken(&user, roles)
	if err != nil {
		return "", err
	}

	return token, nil
}
