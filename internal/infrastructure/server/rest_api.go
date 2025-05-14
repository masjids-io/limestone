package server

import (
	"context"
	"encoding/json"
	"fmt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	pb "github.com/mnadev/limestone/gen/go"
	"github.com/mnadev/limestone/internal/application/handler"
	services "github.com/mnadev/limestone/internal/application/services"
	"github.com/mnadev/limestone/internal/infrastructure/storage"
	"gorm.io/gorm"
)

func SetupRESTGateway(ctx context.Context, db *gorm.DB, grpcEndpoint string) *runtime.ServeMux {
	mux := runtime.NewServeMux(
		runtime.WithErrorHandler(customErrorHandler),
	)

	// User Service
	userRepo := storage.NewGormUserRepository(db)
	userService := services.NewUserService(userRepo)
	userHandler := handler.NewUserGrpcHandler(userService)
	if err := pb.RegisterUserServiceHandlerServer(ctx, mux, userHandler); err != nil {
		log.Fatalf("failed to register UserService handler: %s", err)
	}

	// Auth Service
	authService := services.NewAuthService(userRepo)
	authHandler := handler.NewAuthGrpcHandler(authService)
	if err := pb.RegisterAuthServiceHandlerServer(ctx, mux, authHandler); err != nil {
		log.Fatalf("failed to register AuthService handler: %s", err)
	}

	//Masjid Service
	masjidRepo := storage.NewGormMasjidRepository(db)
	masjidService := services.NewMasjidService(masjidRepo)
	masjidHandler := handler.NewMasjidGrpcHandler(masjidService)
	if err := pb.RegisterMasjidServiceHandlerServer(ctx, mux, masjidHandler); err != nil {
		log.Fatalf("failed to register MasjidService handler: %s", err)
	}

	return mux
}

func StartRESTGateway(ctx context.Context, handler http.Handler, httpEndpoint string) {
	log.Printf("HTTP server listening on %s", httpEndpoint)
	if err := http.ListenAndServe(httpEndpoint, handler); err != nil {
		log.Fatalf("failed to serve HTTP traffic: %s", err)
	}
}

func customErrorHandler(ctx context.Context, mux *runtime.ServeMux, marshaler runtime.Marshaler, w http.ResponseWriter, req *http.Request, err error) {
	s, ok := status.FromError(err)
	if !ok {
		s = status.New(codes.Internal, "internal server error")
	}

	httpStatus := runtime.HTTPStatusFromCode(s.Code())
	fmt.Println(httpStatus)

	errorResponse := map[string]interface{}{
		"code":    string(s.Code().String()),
		"status":  "error",
		"message": s.Message(),
	}

	jsonResponse, err := json.Marshal(errorResponse)
	if err != nil {
		log.Printf("failed to marshal error message: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpStatus)
	w.Write(jsonResponse)
}
