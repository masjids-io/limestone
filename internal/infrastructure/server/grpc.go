package server

import (
	pb "github.com/mnadev/limestone/gen/go"
	services "github.com/mnadev/limestone/internal/application/services"
	"log"
	"net"

	"github.com/mnadev/limestone/internal/application/handler"
	"github.com/mnadev/limestone/internal/infrastructure/grpc/auth"
	"github.com/mnadev/limestone/internal/infrastructure/storage"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"gorm.io/gorm"
)

func SetupGRPCServer(db *gorm.DB, grpcEndpoint string) (*grpc.Server, net.Listener) {
	listener, err := net.Listen("tcp", grpcEndpoint)
	if err != nil {
		log.Fatalf("failed to listen for gRPC: %v", err)
	}
	log.Printf("gRPC server listening on %s", grpcEndpoint)

	repo := storage.NewGormUserRepository(db)
	svc := services.NewUserService(repo)
	handler := handler.NewUserGrpcHandler(svc)

	server := grpc.NewServer(
		grpc.UnaryInterceptor(auth.VerifyJWTInterceptor),
	)
	pb.RegisterUserServiceServer(server, handler)
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
