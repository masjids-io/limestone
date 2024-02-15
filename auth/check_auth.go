package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v4"
	auth_grpc "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var jwtKey = []byte(os.Getenv("jwt-key")) // Use a secure way to store and retrieve this key

// AuthMiddleware checks for JWT in the authorization header and validates it (ONLY WHEN MAKING CALLS VIA GRPC PORT)
func AuthMiddleware(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	fmt.Println("Got inside of interceptor")
	fmt.Println(info.FullMethod)
	tokenStr, err := auth_grpc.AuthFromMD(ctx, "bearer")
	if err != nil {
		return nil, err
	}
	fmt.Println(tokenStr)
	// Verify the token
	parsedToken, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil || !parsedToken.Valid {
		return nil, status.Errorf(codes.Unauthenticated, "invalid auth token: %v", err)
	}

	return handler(ctx, req)
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func writeJSONError(w http.ResponseWriter, statusCode int, errMsg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(ErrorResponse{Error: errMsg})
}

// HTTPAuthMiddleware checks for JWT in the authorization header and validates it
// (ONLY WHEN MAKING CALLS VIA HTTP PORT, grpc interceptor doesnt function during http calls)
func HttpAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		uri := r.RequestURI
		if uri != "/users/create" && uri != "/users/login" && uri != "/users/verify" {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				writeJSONError(w, http.StatusUnauthorized, "Authorization header is required")
				return
			}

			tokenString := strings.TrimPrefix(authHeader, "Bearer ")
			token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
				return jwtKey, nil
			})

			if err != nil || !token.Valid {
				writeJSONError(w, http.StatusUnauthorized, "Invalid token")
				return
			}
		}
		next.ServeHTTP(w, r)
	})
}
