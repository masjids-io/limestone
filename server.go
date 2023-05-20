package main

import (
	"fmt"
	"log"
	"net"
	"os"

	"google.golang.org/grpc"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/mnadev/limestone/event_service"
	epb "github.com/mnadev/limestone/event_service/proto"
	"github.com/mnadev/limestone/masjid_service"
	mpb "github.com/mnadev/limestone/masjid_service/proto"
	"github.com/mnadev/limestone/storage"
	"github.com/mnadev/limestone/user_service"
	upb "github.com/mnadev/limestone/user_service/proto"
)

func main() {
	lis, err := net.Listen("tcp", ":10000")
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
	DB.AutoMigrate(storage.Event{})
	DB.AutoMigrate(storage.Masjid{})
	DB.AutoMigrate(storage.User{})

	epb.RegisterEventServiceServer(server, &event_service.EventServiceServer{
		SM: &storage.StorageManager{
			DB: DB,
		},
	})
	mpb.RegisterMasjidServiceServer(server, &masjid_service.MasjidServiceServer{
		SM: &storage.StorageManager{
			DB: DB,
		},
	})
	upb.RegisterUserServiceServer(server, &user_service.UserServiceServer{
		SM: &storage.StorageManager{
			DB: DB,
		},
	})

	if err := server.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %s", err)
	}
}
