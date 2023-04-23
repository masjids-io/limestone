package test_infra

import (
	"context"
	"database/sql"
	"net"

	"github.com/mnadev/limestone/storage"
	user_service "github.com/mnadev/limestone/user_service"
	userservicepb "github.com/mnadev/limestone/user_service/proto"
	_ "github.com/proullon/ramsql/driver"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type IntegrationTestSuite struct {
	suite.Suite
	DB                *gorm.DB
	Server            *grpc.Server
	UserServiceClient userservicepb.UserServiceClient
}

func (suite *IntegrationTestSuite) BeforeTest(suiteName, testName string) {
	ctx := context.Background()

	buffer := 1024 * 1024
	listener := bufconn.Listen(buffer)

	suite.Server = grpc.NewServer()

	sqlDB, err := sql.Open("ramsql", "Test"+testName)

	if err != nil {
		panic(err)
	}

	suite.DB, err = gorm.Open(postgres.New(postgres.Config{
		Conn: sqlDB,
	}), &gorm.Config{})

	if err != nil {
		panic(err)
	}

	suite.DB.AutoMigrate(&storage.User{})
	userservicepb.RegisterUserServiceServer(suite.Server, &user_service.UserServiceServer{
		SM: &storage.StorageManager{
			DB: suite.DB,
		},
	})

	go func() {
		if err := suite.Server.Serve(listener); err != nil {
			panic(err)
		}
	}()

	conn, _ := grpc.DialContext(ctx, "", grpc.WithContextDialer(func(context.Context, string) (net.Conn, error) {
		return listener.Dial()
	}), grpc.WithInsecure(), grpc.WithBlock())

	suite.UserServiceClient = userservicepb.NewUserServiceClient(conn)
}

func (suite *IntegrationTestSuite) AfterTest(suiteName, testName string) {
	suite.Server.Stop()
}
