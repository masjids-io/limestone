package test_infra

import (
	"context"
	"database/sql"
	"net"

	"github.com/mnadev/limestone/storage"
	user_service "github.com/mnadev/limestone/user_service"
	userservicepb "github.com/mnadev/limestone/userservice/proto"
	_ "github.com/proullon/ramsql/driver"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func UserServiceClient(ctx context.Context) (userservicepb.UserServiceClient, func()) {
	buffer := 1024 * 1024
	listener := bufconn.Listen(buffer)

	server := grpc.NewServer()

	sqlDB, err := sql.Open("ramsql", "LimestoneIntegrationTesting")

	if err != nil {
		panic(err)
	}

	db, err := gorm.Open(postgres.New(postgres.Config{
		Conn: sqlDB,
	}), &gorm.Config{})

	if err != nil {
		panic(err)
	}

	db.AutoMigrate(&storage.User{})
	userservicepb.RegisterUserServiceServer(server, &user_service.UserServiceServer{
		SM: &storage.StorageManager{
			DB: db,
		},
	})

	go func() {
		if err := server.Serve(listener); err != nil {
			panic(err)
		}
	}()

	conn, _ := grpc.DialContext(ctx, "", grpc.WithContextDialer(func(context.Context, string) (net.Conn, error) {
		return listener.Dial()
	}), grpc.WithInsecure(), grpc.WithBlock())

	client := userservicepb.NewUserServiceClient(conn)

	return client, server.Stop
}
