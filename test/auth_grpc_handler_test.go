package test

import (
	"context"
	"github.com/google/uuid"
	pb "github.com/mnadev/limestone/gen/go"
	"github.com/mnadev/limestone/internal/application/domain/entity"
	"github.com/mnadev/limestone/internal/application/handler"
	"github.com/mnadev/limestone/internal/infrastructure/auth"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"os"
)

func (suite *GrpcHandlerTestSuite) TestAuthenticateUser_SuccessWithUsername() {
	ctx := context.Background()
	user := &entity.User{
		ID:             uuid.New(),
		Username:       "testuser10",
		HashedPassword: "$2a$10$GFHqUp5YmuaZQOPOpP5BLusG0IE8zMujioYZxmWYJOipOqhXXTS8O",
	}
	err := suite.DB.Create(&user).Error
	require.NoError(suite.T(), err, "Failed to create test user")

	authHandler := handler.NewAuthGrpcHandler(suite.AuthService)

	req := &pb.AuthenticateUserRequest{
		Identifier: &pb.AuthenticateUserRequest_Username{
			Username: "testuser10",
		},
		Password: "password",
	}

	resp, err := authHandler.AuthenticateUser(ctx, req)

	require.NoError(suite.T(), err)
	require.NotNil(suite.T(), resp)
	assert.Equal(suite.T(), codes.OK.String(), resp.GetCode())
	assert.Equal(suite.T(), "success", resp.GetStatus())
	assert.Equal(suite.T(), "Authentication successful", resp.GetMessage())
	assert.NotEmpty(suite.T(), resp.GetDatas().(*pb.StandardAuthResponse_AuthenticateUserData).AuthenticateUserData.GetAccessToken())
	assert.NotEmpty(suite.T(), resp.GetDatas().(*pb.StandardAuthResponse_AuthenticateUserData).AuthenticateUserData.GetRefreshToken())
	assert.Equal(suite.T(), user.ID.String(), resp.GetDatas().(*pb.StandardAuthResponse_AuthenticateUserData).AuthenticateUserData.GetUserId())

	// Clean up
	err = suite.DB.Delete(&user).Error
	require.NoError(suite.T(), err, "Failed to delete test user")
}

func (suite *GrpcHandlerTestSuite) TestAuthenticateUser_SuccessWithEmail() {
	ctx := context.Background()
	user := &entity.User{
		ID:             uuid.New(),
		Email:          "test10@example.com",
		HashedPassword: "$2a$10$GFHqUp5YmuaZQOPOpP5BLusG0IE8zMujioYZxmWYJOipOqhXXTS8O",
	}
	err := suite.DB.Create(&user).Error
	require.NoError(suite.T(), err, "Failed to create test user")

	authHandler := handler.NewAuthGrpcHandler(suite.AuthService)

	req := &pb.AuthenticateUserRequest{
		Password: "password",
		Identifier: &pb.AuthenticateUserRequest_Email{
			Email: "test10@example.com",
		},
	}

	resp, err := authHandler.AuthenticateUser(ctx, req)

	require.NoError(suite.T(), err)
	require.NotNil(suite.T(), resp)
	assert.Equal(suite.T(), codes.OK.String(), resp.GetCode())
	assert.Equal(suite.T(), "success", resp.GetStatus())
	assert.Equal(suite.T(), "Authentication successful", resp.GetMessage())
	assert.NotEmpty(suite.T(), resp.GetDatas().(*pb.StandardAuthResponse_AuthenticateUserData).AuthenticateUserData.GetAccessToken())
	assert.NotEmpty(suite.T(), resp.GetDatas().(*pb.StandardAuthResponse_AuthenticateUserData).AuthenticateUserData.GetRefreshToken())
	assert.Equal(suite.T(), user.ID.String(), resp.GetDatas().(*pb.StandardAuthResponse_AuthenticateUserData).AuthenticateUserData.GetUserId())

	// Clean up
	err = suite.DB.Delete(&user).Error
	require.NoError(suite.T(), err, "Failed to delete test user")
}

func (suite *GrpcHandlerTestSuite) TestAuthenticateUser_NoIdentifier() {
	ctx := context.Background()
	authHandler := handler.NewAuthGrpcHandler(suite.AuthService)

	req := &pb.AuthenticateUserRequest{
		Password: "password",
	}

	resp, err := authHandler.AuthenticateUser(ctx, req)

	require.Error(suite.T(), err)
	st, ok := status.FromError(err)
	require.True(suite.T(), ok)
	assert.Equal(suite.T(), codes.Canceled, st.Code())
	assert.Equal(suite.T(), "username or email must be provided", st.Message())
	assert.Nil(suite.T(), resp)
}

func (suite *GrpcHandlerTestSuite) TestAuthenticateUser_InvalidCredentials() {
	ctx := context.Background()
	authHandler := handler.NewAuthGrpcHandler(suite.AuthService)

	req := &pb.AuthenticateUserRequest{
		Identifier: &pb.AuthenticateUserRequest_Username{
			Username: "nonexistentuse",
		},
		Password: "wrongpassword",
	}

	resp, err := authHandler.AuthenticateUser(ctx, req)

	require.Error(suite.T(), err)
	st, ok := status.FromError(err)
	require.True(suite.T(), ok)
	assert.Equal(suite.T(), codes.Canceled, st.Code())
	assert.Equal(suite.T(), "invalid username/email or password", st.Message())
	assert.Nil(suite.T(), resp)
}

func (suite *GrpcHandlerTestSuite) TestRefreshToken_Success() {
	require.NotEmpty(suite.T(), os.Getenv("REFRESH_SECRET"), "REFRESH_SECRET environment variable must be set for this test")
	require.NotEmpty(suite.T(), os.Getenv("ACCESS_SECRET"), "ACCESS_SECRET environment variable must be set for this test")
	require.NotEmpty(suite.T(), os.Getenv("ACCESS_EXPIRATION"), "ACCESS_EXPIRATION environment variable must be set for this test")
	require.NotEmpty(suite.T(), os.Getenv("REFRESH_EXPIRATION"), "REFRESH_EXPIRATION environment variable must be set for this test")

	ctx := context.Background()

	user := &entity.User{
		ID:             uuid.New(),
		Username:       "testuser11",
		HashedPassword: "$2a$10$GFHqUp5YmuaZQOPOpP5BLusG0IE8zMujioYZxmWYJOipOqhXXTS8O",
	}
	err := suite.DB.Create(&user).Error
	require.NoError(suite.T(), err, "Failed to create test user")
	userID := user.ID.String()

	_, refreshToken, err := auth.GenerateJWT(userID)
	require.NoError(suite.T(), err, "Failed to generate initial JWT tokens")

	authHandler := handler.NewAuthGrpcHandler(suite.AuthService)

	req := &pb.RefreshTokenRequest{
		RefreshToken: refreshToken,
	}

	resp, err := authHandler.RefreshToken(ctx, req)

	require.NoError(suite.T(), err)
	require.NotNil(suite.T(), resp)
	assert.Equal(suite.T(), codes.OK.String(), resp.GetCode())
	assert.Equal(suite.T(), "success", resp.GetStatus())
	assert.Equal(suite.T(), "Token refreshed", resp.GetMessage())
	assert.NotEmpty(suite.T(), resp.GetDatas().(*pb.StandardAuthResponse_RefreshTokenData).RefreshTokenData.GetAccessToken())
	assert.NotEmpty(suite.T(), resp.GetDatas().(*pb.StandardAuthResponse_RefreshTokenData).RefreshTokenData.GetRefreshToken())

	err = suite.DB.Delete(&user).Error
	require.NoError(suite.T(), err, "Failed to delete test user")
}

func (suite *GrpcHandlerTestSuite) TestRefreshToken_InvalidToken() {
	ctx := context.Background()
	authHandler := handler.NewAuthGrpcHandler(suite.AuthService)

	req := &pb.RefreshTokenRequest{
		RefreshToken: "invalid-refresh-token",
	}

	resp, err := authHandler.RefreshToken(ctx, req)

	require.Error(suite.T(), err)
	st, ok := status.FromError(err)
	require.True(suite.T(), ok)
	assert.Equal(suite.T(), codes.Internal, st.Code())
	assert.Contains(suite.T(), st.Message(), "failed to refresh access token")
	assert.Nil(suite.T(), resp)
}
