package auth

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

// HashPassword returns the hashed password
func HashPassword(p string) (string, error) {
	h, err := bcrypt.GenerateFromPassword([]byte(p), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}
	return string(h), nil
}

// CheckPassword checks if the provided plaintext password is equivalent to the hashed password
func CheckPassword(p string, hp string) error {
	return bcrypt.CompareHashAndPassword([]byte(hp), []byte(p))
}
