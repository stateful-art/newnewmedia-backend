package service

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
	dao "newnew.media/microservices/user/dao"
)

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

func GenerateJWTToken(user *dao.User) (string, error) {
	// Define the claims of the token
	claims := jwt.MapClaims{
		"ID":    user.ID,
		"email": user.Email,
		// Add any other claims you need
		"exp": time.Now().Add(time.Hour * 24).Unix(), // Token expires after 24 hours
	}

	// Create a new token object, specifying signing method and the claims
	// TODO: switch to SigningMethodES256 later for more security.
	// ECDSA approach is more secure than using HS256 because it uses asymmetric encryption,
	// which means the private key is used to sign the token, and the public key is used to verify it.
	// This eliminates the need to securely share a secret key between the server and the client.
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token with the secret key
	tokenString, err := token.SignedString([]byte("xyz"))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
