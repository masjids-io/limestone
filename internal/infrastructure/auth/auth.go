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
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

func GenerateJWT(userID string) (string, string, error) {
	accessSecret := os.Getenv("ACCESS_SECRET")
	refreshSecret := os.Getenv("REFRESH_SECRET")
	accessExpiration := os.Getenv("ACCESS_EXPIRATION")
	refreshExpiration := os.Getenv("REFRESH_EXPIRATION")

	accessExpMinutes, err := strconv.Atoi(accessExpiration)
	if err != nil {
		accessExpMinutes = 60
	}
	refreshExpHours, err := strconv.Atoi(refreshExpiration)
	if err != nil {
		refreshExpHours = 24 * 7
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

func RefreshToken(refreshTokenString string) (string, string, error) {
	refreshSecret := os.Getenv("REFRESH_SECRET")
	accessSecret := os.Getenv("ACCESS_SECRET")
	accessExpiration := os.Getenv("ACCESS_EXPIRATION")

	if refreshSecret == "" || accessSecret == "" || accessExpiration == "" {
		log.Println("REFRESH_SECRET, ACCESS_SECRET, or ACCESS_EXPIRATION not set")
		return "", "", status.Errorf(codes.Internal, "server configuration error: missing secrets or expiration")
	}

	token, err := jwt.Parse(refreshTokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(refreshSecret), nil
	})

	if err != nil {
		log.Printf("Error parsing refresh token: %v", err)
		return "", "", status.Errorf(codes.Unauthenticated, "invalid refresh token: %v", err)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		log.Println("Invalid refresh token claims")
		return "", "", status.Errorf(codes.Unauthenticated, "invalid refresh token claims")
	}

	userID, ok := claims["user_id"].(string)
	if !ok {
		log.Println("User ID not found in refresh token claims")
		return "", "", status.Errorf(codes.Unauthenticated, "invalid refresh token claims: user_id not found")
	}

	accessExpMinutes, err := strconv.Atoi(accessExpiration)
	if err != nil {
		accessExpMinutes = 60
	}

	now := time.Now()
	accessClaims := jwt.MapClaims{
		"user_id": userID,
		"exp":     now.Add(time.Minute * time.Duration(accessExpMinutes)).Unix(),
		"iat":     now.Unix(),
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessString, err := accessToken.SignedString([]byte(accessSecret))
	if err != nil {
		return "", "", fmt.Errorf("failed to sign new access token: %w", err)
	}

	refreshExpiration := os.Getenv("REFRESH_EXPIRATION")
	refreshExpHours, err := strconv.Atoi(refreshExpiration)
	if err != nil {
		refreshExpHours = 24 * 7
	}

	newRefreshTokenClaims := jwt.MapClaims{
		"user_id": userID,
		"exp":     now.Add(time.Hour * time.Duration(refreshExpHours)).Unix(),
		"iat":     now.Unix(),
	}
	newRefreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, newRefreshTokenClaims)
	newRefreshTokenString, err := newRefreshToken.SignedString([]byte(refreshSecret))
	if err != nil {
		return "", "", fmt.Errorf("failed to sign new refresh token: %w", err)
	}

	return accessString, newRefreshTokenString, nil
}

func VerifyJWTInterceptorRest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		route := UnprotectedRoute{Path: r.URL.Path, Method: r.Method}
		if _, isUnprotected := UnprotectedRoutesHTTP[route]; isUnprotected {
			next.ServeHTTP(w, r)
			return
		}

		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		bearerToken := strings.Split(authHeader, " ")
		if len(bearerToken) != 2 || strings.ToLower(bearerToken[0]) != "bearer" {
			http.Error(w, "Invalid authorization header", http.StatusUnauthorized)
			return
		}

		tokenString := bearerToken[1]
		accessSecret := os.Getenv("ACCESS_SECRET")
		if accessSecret == "" {
			log.Println("ACCESS_SECRET not set for HTTP")
			http.Error(w, "Server configuration error", http.StatusInternalServerError)
			return
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(accessSecret), nil
		})

		if err != nil {
			log.Printf("Error parsing token (HTTP): %v", err)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			userID, ok := claims["user_id"].(string)
			if !ok {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), UserIDContextKey, userID)
			next.ServeHTTP(w, r.WithContext(ctx))
		} else {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
		}
	})
}
