package main

import (
	"fmt"
	"log"
	"net"
	"os"

	"github.com/mnadev/limestone/storage"
	userservicepb "github.com/mnadev/limestone/user_service/proto"
	"google.golang.org/grpc"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
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
	DB.AutoMigrate(storage.User{})
	userservicepb.RegisterUserServiceServer(server, &userservice.UserServiceServer{
		SM: &storage.StorageManager{
			DB: DB,
		},
	})

	if err := server.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %s", err)
	}
}
