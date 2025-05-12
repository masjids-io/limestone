package main

import (
	pb "github.com/mnadev/limestone/gen/go"
	"github.com/mnadev/limestone/internal/application/services"
	"github.com/mnadev/limestone/internal/infrastructure/grpc/auth"
	"github.com/mnadev/limestone/internal/infrastructure/storage"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"gorm.io/gorm"
	"log"
	"net"
)

func setupGRPCServer(db *gorm.DB) (*grpc.Server, net.Listener) {
	lis, err := net.Listen("tcp", *grpcEndpoint)
	if err != nil {
		log.Fatalf("failed to listen for gRPC: %v", err)
	}
	log.Printf("gRPC server listening on %s", *grpcEndpoint)

	server := grpc.NewServer(
		grpc.UnaryInterceptor(auth.VerifyJWTInterceptor),
	)

	smgr := &storage.StorageManager{DB: db}
	//pb.RegisterAdhanServiceServer(server, &adhan_service.AdhanServiceServer{Smgr: smgr})
	//pb.RegisterEventServiceServer(server, &event_service.EventServiceServer{Smgr: smgr})
	//pb.RegisterMasjidServiceServer(server, &masjid_service.MasjidServiceServer{Smgr: smgr})
	pb.RegisterUserServiceServer(server, &user_service.UserServiceServer{Smgr: smgr})
	reflection.Register(server) //debug grpc with grpcurl

	return server, lis
}

func startGRPCServer(server *grpc.Server, lis net.Listener) {
	go func() {
		if err := server.Serve(lis); err != nil {
			log.Fatalf("failed to serve gRPC traffic: %s", err)
		}
	}()
}
