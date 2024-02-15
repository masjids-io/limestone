package auth

import (
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type TokenType int64

var (
	access  TokenType = 0
	refresh TokenType = 1
)

// Claims represents the JWT claims structure
type Claims struct {
	Uuid  string    `json:"uuid"`
	Email string    `json:"email"`
	Type  TokenType `json:"type"` //To decide whether it is access(0) or refresh(1) token
	jwt.RegisteredClaims
}

func CreateJWT(uuid string, email string, t TokenType) (string, error) {
	// Set the expiration time to 30 days from now
	fmt.Println("Generating your token ......")
	var expirationTime time.Time
	if t == access {
		// 10 days expiry of access token
		expirationTime = time.Now().Add(10 * 24 * time.Hour)
	} else if t == refresh {
		// 30 days expiry of refresh token
		expirationTime = time.Now().Add(30 * 24 * time.Hour)
	}
	// Create the claims
	claims := Claims{
		Uuid:  uuid,
		Email: email,
		Type:  t,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}
	// Create the token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	var jwtKey = []byte(os.Getenv("jwt-key"))
	// Sign the token with the secret key
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}
	fmt.Printf("Here is your token: %s", tokenString)
	return tokenString, nil
}
