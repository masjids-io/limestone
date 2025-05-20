package test

import (
	"github.com/mnadev/limestone/internal/application/domain/entity"
	"testing"

	"github.com/lpernett/godotenv"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"

	"github.com/mnadev/limestone/internal/application/handler"
	"github.com/mnadev/limestone/internal/application/services"
	"github.com/mnadev/limestone/internal/infrastructure/database"
	"github.com/mnadev/limestone/internal/infrastructure/storage"
)

type GrpcHandlerTestSuite struct {
	suite.Suite
	DB            *gorm.DB
	UserHandler   *handler.UserGrpcHandler
	UserService   *services.UserService
	AuthService   *services.AuthService
	AuthHandler   *handler.AuthGrpcHandler
	AdhanService  *services.AdhanService
	AdhanHandler  *handler.AdhanGrpcHandler
	MasjidService *services.MasjidService
	MasjidHandler *handler.MasjidGrpcHandler
	EventService  *services.EventService
	EventHandler  *handler.EventGrpcHandler
}

func (suite *GrpcHandlerTestSuite) SetupSuite() {
	err := godotenv.Load("../.env")
	require.NoError(suite.T(), err, "Failed to load .env file")

	suite.DB = database.SetupDatabaseTesting()
	require.NotNil(suite.T(), suite.DB, "Failed to setup test database")

	//user service
	userRepo := storage.NewGormUserRepository(suite.DB)
	suite.UserService = services.NewUserService(userRepo)
	suite.UserHandler = handler.NewUserGrpcHandler(suite.UserService)

	//auth service
	authService := services.NewAuthService(userRepo)
	suite.AuthService = authService
	suite.AuthHandler = handler.NewAuthGrpcHandler(suite.AuthService)

	//adhan service
	adhanRepo := storage.NewGormAdhanRepository(suite.DB)
	suite.AdhanService = services.NewAdhanService(adhanRepo)
	suite.AdhanHandler = handler.NewAdhanGrpcHandler(suite.AdhanService)

	//masjid service
	masjidrepo := storage.NewGormMasjidRepository(suite.DB)
	suite.MasjidService = services.NewMasjidService(masjidrepo)
	suite.MasjidHandler = handler.NewMasjidGrpcHandler(suite.MasjidService)

	//event service
	eventRepo := storage.NewGormEventRepository(suite.DB)
	suite.EventService = services.NewEventService(eventRepo)
	suite.EventHandler = handler.NewEventGrpcHandler(suite.EventService)
}

func (suite *GrpcHandlerTestSuite) TearDownSuite() {
	if suite.DB != nil {
		err := suite.DB.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&entity.User{}).Error
		require.NoError(suite.T(), err, "Failed to clean up users table")

		db, err := suite.DB.DB()
		require.NoError(suite.T(), err, "Failed to get underlying DB connection")
		err = db.Close()
		require.NoError(suite.T(), err, "Failed to close database connection")
	}
}

func Test(t *testing.T) {
	suite.Run(t, new(GrpcHandlerTestSuite))
}
