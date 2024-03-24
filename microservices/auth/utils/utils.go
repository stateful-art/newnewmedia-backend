package utils

import (
	"regexp"
	"strings"
	"unicode"
)

func IsValidEmail(email string) bool {
	// Regular expression for basic email validation
	// This regex pattern is simple and may not cover all edge cases
	pattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	matched, err := regexp.MatchString(pattern, email)
	return err == nil && matched
}

func TrimInput(input string) string {
	return strings.TrimSpace(input)
}

func IsValidPassword(password string) bool {
	// Ensure minimum length of 8 characters
	if len(password) < 8 {
		return false
	}

	// Check for at least one uppercase letter, one lowercase letter,
	// one digit, and one special character
	hasUpper := false
	hasLower := false
	hasDigit := false
	hasSpecial := false

	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsDigit(char):
			hasDigit = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}

	return hasUpper && hasLower && hasDigit && hasSpecial
}
