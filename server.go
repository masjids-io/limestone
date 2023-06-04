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
	apb "github.com/mnadev/limestone/adhan_service/proto"
	"github.com/mnadev/limestone/event_service"
	epb "github.com/mnadev/limestone/event_service/proto"
	"github.com/mnadev/limestone/masjid_service"
	mpb "github.com/mnadev/limestone/masjid_service/proto"
	"github.com/mnadev/limestone/storage"
	"github.com/mnadev/limestone/user_service"
	upb "github.com/mnadev/limestone/user_service/proto"
)

var (
	grpcEndpoint = flag.String("grpc_endpoint", "localhost:8081", "gRPC server endpoint")
	httpEndpoint = flag.String("http_endpoint", "localhost:8080", "HTTP server endpoint")
)

func main() {
	lis, err := net.Listen("tcp", *grpcEndpoint)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	server := grpc.NewServer()

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
	apb.RegisterAdhanServiceServer(server, &adhan_service_server)
	event_server := event_service.EventServiceServer{
		SM: &storage.StorageManager{
			DB: DB,
		},
	}
	epb.RegisterEventServiceServer(server, &event_server)
	masjid_server := masjid_service.MasjidServiceServer{
		SM: &storage.StorageManager{
			DB: DB,
		},
	}
	mpb.RegisterMasjidServiceServer(server, &masjid_server)
	user_server := user_service.UserServiceServer{
		SM: &storage.StorageManager{
			DB: DB,
		},
	}
	upb.RegisterUserServiceServer(server, &user_server)

	go func() {
		if err := server.Serve(lis); err != nil {
			log.Fatalf("failed to serve gRPC traffic: %s", err)
		}
	}()

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	mux := runtime.NewServeMux()
	err = apb.RegisterAdhanServiceHandlerServer(ctx, mux, &adhan_service_server)
	if err != nil {
		log.Fatalf("failed to serve: %s", err)
	}
	err = epb.RegisterEventServiceHandlerServer(ctx, mux, &event_server)
	if err != nil {
		log.Fatalf("failed to serve: %s", err)
	}
	err = mpb.RegisterMasjidServiceHandlerServer(ctx, mux, &masjid_server)
	if err != nil {
		log.Fatalf("failed to serve: %s", err)
	}
	err = upb.RegisterUserServiceHandlerServer(ctx, mux, &user_server)
	if err != nil {
		log.Fatalf("failed to serve: %s", err)
	}

	if err = http.ListenAndServe(*httpEndpoint, mux); err != nil {
		log.Fatalf("failed to serve HTTP traffic: %s", err)
	}
}
