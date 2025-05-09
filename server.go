package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/lpernett/godotenv"
	"github.com/mnadev/limestone/adhan_service"
	"github.com/mnadev/limestone/auth"
	"github.com/mnadev/limestone/event_service"
	"github.com/mnadev/limestone/masjid_service"
	pb "github.com/mnadev/limestone/proto"
	"github.com/mnadev/limestone/storage"
	"github.com/mnadev/limestone/user_service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"net"
	"net/http"
	"os"
)

var (
	grpcEndpoint = flag.String("grpc_endpoint", "localhost:8081", "gRPC server endpoint")
	httpEndpoint = flag.String("http_endpoint", "localhost:8080", "HTTP server endpoint")
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	lis, err := net.Listen("tcp", *grpcEndpoint)
	if err != nil {
		log.Fatalf("failed to listen for gRPC: %v", err)
	}
	log.Printf("gRPC server listening on %s", *grpcEndpoint)

	server := grpc.NewServer(
		grpc.UnaryInterceptor(auth.VerifyJWTInterceptor),
	)

	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")
	dbUser := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dsn := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable",
		host,
		port,
		dbUser,
		dbName,
		password,
	)
	DB, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	log.Println("Database connected successfully.")
	DB.AutoMigrate(storage.AdhanFile{})
	DB.AutoMigrate(storage.Event{})
	DB.AutoMigrate(storage.Masjid{})
	DB.AutoMigrate(storage.User{})

	adhan_service_server := adhan_service.AdhanServiceServer{
		Smgr: &storage.StorageManager{
			DB: DB,
		},
	}
	pb.RegisterAdhanServiceServer(server, &adhan_service_server)
	event_server := event_service.EventServiceServer{
		Smgr: &storage.StorageManager{
			DB: DB,
		},
	}
	pb.RegisterEventServiceServer(server, &event_server)
	masjid_server := masjid_service.MasjidServiceServer{
		Smgr: &storage.StorageManager{
			DB: DB,
		},
	}
	pb.RegisterMasjidServiceServer(server, &masjid_server)
	user_server := user_service.UserServiceServer{
		Smgr: &storage.StorageManager{
			DB: DB,
		},
	}
	pb.RegisterUserServiceServer(server, &user_server)

	reflection.Register(server) //debug grpc with grpcurl

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
		log.Fatalf("failed to register AdhanService handler: %s", err)
	}
	err = pb.RegisterEventServiceHandlerServer(ctx, mux, &event_server)
	if err != nil {
		log.Fatalf("failed to register EventService handler: %s", err)
	}
	err = pb.RegisterMasjidServiceHandlerServer(ctx, mux, &masjid_server)
	if err != nil {
		log.Fatalf("failed to register MasjidService handler: %s", err)
	}
	err = pb.RegisterUserServiceHandlerServer(ctx, mux, &user_server)
	if err != nil {
		log.Fatalf("failed to register UserService handler: %s", err)
	}
	log.Printf("HTTP server listening on %s", *httpEndpoint)
	if err = http.ListenAndServe(*httpEndpoint, mux); err != nil {
		log.Fatalf("failed to serve HTTP traffic: %s", err)
	}
}
