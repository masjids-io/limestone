package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"

	"google.golang.org/grpc"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"

	"github.com/mnadev/limestone/adhan_service"
	"github.com/mnadev/limestone/auth"
	"github.com/mnadev/limestone/event_service"
	"github.com/mnadev/limestone/masjid_service"
	pb "github.com/mnadev/limestone/proto"
	"github.com/mnadev/limestone/storage"
	"github.com/mnadev/limestone/user_service"
)

var (
	grpcEndpoint = flag.String("grpc_endpoint", ":8081", "gRPC server endpoint")
	httpEndpoint = flag.String("http_endpoint", ":8080", "HTTP server endpoint")
)

func main() {
	lis, err := net.Listen("tcp", *grpcEndpoint)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	server := grpc.NewServer(
		grpc.UnaryInterceptor(auth.AuthInterceptor),
	)

	host := os.Getenv("db-host")
	port := os.Getenv("db-port")
	dbName := os.Getenv("db-name")
	dbUser := os.Getenv("db-user")
	password := os.Getenv("db-password")
	dsn := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable",
		host,
		port,
		dbUser,
		dbName,
		password,
		// "localhost",
		// "5432",
		// "admin",
		// "postgres",
		// "password",
	)
	DB, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	DB.AutoMigrate(storage.AdhanFile{})
	DB.AutoMigrate(storage.Event{})
	DB.AutoMigrate(storage.Masjid{})
	DB.AutoMigrate(storage.User{})

	adhan_service_server := adhan_service.AdhanServiceServer{
		SM: &storage.StorageManager{
			DB: DB,
		},
	}
	pb.RegisterAdhanServiceServer(server, &adhan_service_server)
	event_server := event_service.EventServiceServer{
		SM: &storage.StorageManager{
			DB: DB,
		},
	}
	pb.RegisterEventServiceServer(server, &event_server)
	masjid_server := masjid_service.MasjidServiceServer{
		SM: &storage.StorageManager{
			DB: DB,
		},
	}
	pb.RegisterMasjidServiceServer(server, &masjid_server)
	user_server := user_service.UserServiceServer{
		SM: &storage.StorageManager{
			DB: DB,
		},
	}
	pb.RegisterUserServiceServer(server, &user_server)

	go func() {
		if err := server.Serve(lis); err != nil {
			log.Fatalf("failed to serve gRPC traffic: %s", err)
		}
	}()

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	mux := runtime.NewServeMux()
	err = pb.RegisterAdhanServiceHandlerServer(ctx, mux, &adhan_service_server)
	if err != nil {
		log.Fatalf("failed to serve: %s", err)
	}
	err = pb.RegisterEventServiceHandlerServer(ctx, mux, &event_server)
	if err != nil {
		log.Fatalf("failed to serve: %s", err)
	}
	err = pb.RegisterMasjidServiceHandlerServer(ctx, mux, &masjid_server)
	if err != nil {
		log.Fatalf("failed to serve: %s", err)
	}
	err = pb.RegisterUserServiceHandlerServer(ctx, mux, &user_server)
	if err != nil {
		log.Fatalf("failed to serve: %s", err)
	}

	if err = http.ListenAndServe(*httpEndpoint, mux); err != nil {
		log.Fatalf("failed to serve HTTP traffic: %s", err)
	}
}
