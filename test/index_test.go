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
	DB          *gorm.DB
	UserHandler *handler.UserGrpcHandler
	UserService *services.UserService
	AuthService *services.AuthService
	AuthHandler *handler.AuthGrpcHandler
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
