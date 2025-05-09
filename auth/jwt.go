package auth

import (
	"context"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	grpcauth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
	"os"
	"strconv"
	"time"
)

func GenerateJWT(userID string) (string, string, error) {
	accessSecret := os.Getenv("ACCESS_SECRET")
	refreshSecret := os.Getenv("REFRESH_SECRET")
	accessExpiration := os.Getenv("ACCESS_EXPIRATION")
	refreshExpiration := os.Getenv("REFRESH_EXPIRATION")

	accessExpMinutes, err := strconv.Atoi(accessExpiration)
	if err != nil {
		accessExpMinutes = 60 // Default to 30 minutes
	}
	refreshExpHours, err := strconv.Atoi(refreshExpiration)
	if err != nil {
		refreshExpHours = 24 * 7 // Default to 7 days
	}

	now := time.Now()
	accessClaims := jwt.MapClaims{
		"user_id": userID,
		"exp":     now.Add(time.Minute * time.Duration(accessExpMinutes)).Unix(),
		"iat":     now.Unix(),
	}

	refreshClaims := jwt.MapClaims{
		"user_id": userID,
		"exp":     now.Add(time.Hour * time.Duration(refreshExpHours)).Unix(),
		"iat":     now.Unix(),
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)

	accessString, err := accessToken.SignedString([]byte(accessSecret))
	if err != nil {
		return "", "", fmt.Errorf("failed to sign access token: %w", err)
	}

	refreshString, err := refreshToken.SignedString([]byte(refreshSecret))
	if err != nil {
		return "", "", fmt.Errorf("failed to sign refresh token: %w", err)
	}

	return accessString, refreshString, nil
}

type AuthContextKey string

const UserIDContextKey AuthContextKey = "userID"

func VerifyJWTInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	if _, isUnprotected := UnprotectedRoutes[info.FullMethod]; isUnprotected {
		return handler(ctx, req)
	}
	tokenString, err := grpcauth.AuthFromMD(ctx, "Bearer")
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "invalid authorization header: %v", err)
	}
	accessSecret := os.Getenv("ACCESS_SECRET")
	if accessSecret == "" {
		log.Println("ACCESS_SECRET not set")
		return nil, status.Errorf(codes.Internal, "server configuration error: ACCESS_SECRET not set")
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(accessSecret), nil
	})

	if err != nil {
		log.Printf("Error parsing token: %v", err)
		return nil, status.Errorf(codes.Unauthenticated, "invalid token: %v", err)
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userID, ok := claims["user_id"].(string)
		if !ok {
			return nil, status.Errorf(codes.Unauthenticated, "invalid token claims: user_id not found")
		}

		newCtx := context.WithValue(ctx, UserIDContextKey, userID)
		return handler(newCtx, req)
	}

	return nil, status.Errorf(codes.Unauthenticated, "invalid token")
}
