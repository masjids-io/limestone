package server

import (
	pb "github.com/mnadev/limestone/gen/go"
	"github.com/mnadev/limestone/internal/application/services"
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
	listener, err := net.Listen("tcp", grpcEndpoint)
	if err != nil {
		log.Printf("failed to listen for gRPC: %s", err)
	}
	log.Printf("gRPC server listening on %s", grpcEndpoint)

	server := grpc.NewServer(
		grpc.UnaryInterceptor(auth.VerifyJWTInterceptor),
	)

	// Initialize repositories and services
	userRepo := storage.NewGormUserRepository(db)
	userService := services.NewUserService(userRepo)
	authService := services.NewAuthService(userRepo)
	//masjid service
	masjidRepo := storage.NewGormMasjidRepository(db)
	masjidService := services.NewMasjidService(masjidRepo)
	//adhan service
	adhanRepo := storage.NewGormAdhanRepository(db)
	adhanService := services.NewAdhanService(adhanRepo)
	//event service
	eventRepo := storage.NewGormEventRepository(db)
	eventService := services.NewEventService(eventRepo)
	//nikkah service
	nikkahRepo := storage.NewGormNikkahRepository(db)
	nikkahService := services.NewNikkahService(nikkahRepo)

	// Initialize handlers
	userHandler := handler.NewUserGrpcHandler(userService)
	authHandler := handler.NewAuthGrpcHandler(authService)
	masjidHandler := handler.NewMasjidGrpcHandler(masjidService)
	adhanHandler := handler.NewAdhanGrpcHandler(adhanService)
	eventHandler := handler.NewEventGrpcHandler(eventService)
	nikkahHandler := handler.NewNikkahIoGrpcHandler(nikkahService)

	// Register services with their handlers
	pb.RegisterUserServiceServer(server, userHandler)
	pb.RegisterAuthServiceServer(server, authHandler)
	pb.RegisterMasjidServiceServer(server, masjidHandler)
	pb.RegisterAdhanServiceServer(server, adhanHandler)
	pb.RegisterEventServiceServer(server, eventHandler)
	pb.RegisterNikkahIoServiceServer(server, nikkahHandler)

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
