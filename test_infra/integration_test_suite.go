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

	"github.com/mnadev/limestone/adhan_service"
	"github.com/mnadev/limestone/event_service"
	"github.com/mnadev/limestone/masjid_service"
	pb "github.com/mnadev/limestone/proto"
	"github.com/mnadev/limestone/storage"
	"github.com/mnadev/limestone/user_service"
)

type IntegrationTestSuite struct {
	suite.Suite
	DB                  *gorm.DB
	Server              *grpc.Server
	AdhanServiceClient  pb.AdhanServiceClient
	EventServiceClient  pb.EventServiceClient
	MasjidServiceClient pb.MasjidServiceClient
	UserServiceClient   pb.UserServiceClient
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

	suite.DB.AutoMigrate(&storage.AdhanFile{})
	suite.DB.AutoMigrate(&storage.Event{})
	suite.DB.AutoMigrate(&storage.Masjid{})
	suite.DB.AutoMigrate(&storage.User{})

	pb.RegisterAdhanServiceServer(suite.Server, &adhan_service.AdhanServiceServer{
		Smgr: &storage.StorageManager{
			DB: suite.DB,
		},
	})
	pb.RegisterEventServiceServer(suite.Server, &event_service.EventServiceServer{
		Smgr: &storage.StorageManager{
			DB: suite.DB,
		},
	})
	pb.RegisterMasjidServiceServer(suite.Server, &masjid_service.MasjidServiceServer{
		Smgr: &storage.StorageManager{
			DB: suite.DB,
		},
	})
	pb.RegisterUserServiceServer(suite.Server, &user_service.UserServiceServer{
		Smgr: &storage.StorageManager{
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

	suite.AdhanServiceClient = pb.NewAdhanServiceClient(conn)
	suite.EventServiceClient = pb.NewEventServiceClient(conn)
	suite.MasjidServiceClient = pb.NewMasjidServiceClient(conn)
	suite.UserServiceClient = pb.NewUserServiceClient(conn)
}

func (suite *IntegrationTestSuite) AfterTest(suiteName, testName string) {
	suite.Server.Stop()
}
