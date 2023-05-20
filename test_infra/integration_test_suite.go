package test_infra

import (
	"context"
	"database/sql"
	"net"

	_ "github.com/proullon/ramsql/driver"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/mnadev/limestone/event_service"
	epb "github.com/mnadev/limestone/event_service/proto"
	"github.com/mnadev/limestone/masjid_service"
	mpb "github.com/mnadev/limestone/masjid_service/proto"
	"github.com/mnadev/limestone/storage"
	user_service "github.com/mnadev/limestone/user_service"
	upb "github.com/mnadev/limestone/user_service/proto"
)

type IntegrationTestSuite struct {
	suite.Suite
	DB                  *gorm.DB
	Server              *grpc.Server
	EventServiceClient  epb.EventServiceClient
	MasjidServiceClient mpb.MasjidServiceClient
	UserServiceClient   upb.UserServiceClient
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
	suite.DB.AutoMigrate(&storage.Masjid{})
	suite.DB.AutoMigrate(&storage.Event{})
	epb.RegisterEventServiceServer(suite.Server, &event_service.EventServiceServer{
		SM: &storage.StorageManager{
			DB: suite.DB,
		},
	})
	mpb.RegisterMasjidServiceServer(suite.Server, &masjid_service.MasjidServiceServer{
		SM: &storage.StorageManager{
			DB: suite.DB,
		},
	})
	upb.RegisterUserServiceServer(suite.Server, &user_service.UserServiceServer{
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

	suite.EventServiceClient = epb.NewEventServiceClient(conn)
	suite.MasjidServiceClient = mpb.NewMasjidServiceClient(conn)
	suite.UserServiceClient = upb.NewUserServiceClient(conn)
}

func (suite *IntegrationTestSuite) AfterTest(suiteName, testName string) {
	suite.Server.Stop()
}
