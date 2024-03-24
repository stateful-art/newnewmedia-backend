package service

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
	dao "newnew.media/microservices/user/dao"
)

var jwtSecret = os.Getenv("JWT_SECRET")

// HashPassword hashes a password using bcrypt.
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// CheckPasswordHash verifies a password against a hash.
// returns true if check pass, false if cannot.
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// func GenerateJWTToken(user *dao.User) (string, error) {
// 	// Define the claims of the token
// 	claims := jwt.MapClaims{
// 		"ID":    user.ID,
// 		"email": user.Email,
// 		// Add any other claims you need
// 		"exp": time.Now().Add(time.Hour * 24).Unix(), // Token expires after 24 hours
// 	}
// 	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

// 	// Sign the token with the secret key
// 	tokenString, err := token.SignedString([]byte("xyz"))
// 	if err != nil {
// 		return "", err
// 	}

// 	return tokenString, nil
// }

func GenerateJWTToken(user *dao.User, roles []string) (string, error) {
	// Define the claims of the token
	claims := jwt.MapClaims{
		"ID":    user.ID,
		"email": user.Email,
		"roles": roles,                                 // Include the roles in the token claims
		"exp":   time.Now().Add(time.Hour * 24).Unix(), // Token expires after 24 hours
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token with the secret key
	tokenString, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func GenerateSpotifyJWTToken(user *dao.User, roles []string) (string, error) {
	// Define the claims of the token
	claims := jwt.MapClaims{
		"ID":         user.ID,
		"spotify_id": user.SpotifyID,
		"roles":      roles,                                 // Include the roles in the token claims
		"exp":        time.Now().Add(time.Hour * 24).Unix(), // Token expires after 24 hours
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token with the secret key
	tokenString, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
