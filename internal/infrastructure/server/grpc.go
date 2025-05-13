package server

import (
	pb "github.com/mnadev/limestone/gen/go"
	services "github.com/mnadev/limestone/internal/application/services"
	"log"
	"net"

	"github.com/mnadev/limestone/internal/application/handler"
	"github.com/mnadev/limestone/internal/infrastructure/auth"
	"github.com/mnadev/limestone/internal/infrastructure/storage"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"gorm.io/gorm"
)

func SetupGRPCServer(db *gorm.DB, grpcEndpoint string) (*grpc.Server, net.Listener) {
	listener, err := net.Listen("tcp", grpcEndpoint) // Use a constant or config
	if err != nil {
		log.Printf("failed to listen for gRPC: %s", err)
	}
	log.Printf("gRPC server listening on %s", grpcEndpoint)

	server := grpc.NewServer(
		grpc.UnaryInterceptor(auth.VerifyJWTInterceptor), // Apply auth globally
	)

	// Initialize repositories and services
	userRepo := storage.NewGormUserRepository(db)
	userService := services.NewUserService(userRepo)
	authService := services.NewAuthService(userRepo) // Assuming AuthService

	// Initialize handlers
	userHandler := handler.NewUserGrpcHandler(userService)
	authHandler := handler.NewAuthGrpcHandler(authService)

	// Register services with their handlers
	pb.RegisterUserServiceServer(server, userHandler)
	pb.RegisterAuthServiceServer(server, authHandler) // Assuming AuthService

	reflection.Register(server)

	return server, listener
}

func StartGRPCServer(server *grpc.Server, lis net.Listener) {
	go func() {
		if err := server.Serve(lis); err != nil {
			log.Fatalf("failed to serve gRPC traffic: %v", err)
		}
	}()
}
