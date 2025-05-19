package test

import (
	"context"
	"github.com/google/uuid"
	"google.golang.org/grpc/status"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"

	pb "github.com/mnadev/limestone/gen/go"
	"github.com/mnadev/limestone/internal/application/domain/entity"
	"gorm.io/gorm"
)

// create user
func (suite *GrpcHandlerTestSuite) TestCreateUser_Success() {
	ctx := context.Background()
	req := &pb.CreateUserRequest{
		Email:           "test@example.com",
		Username:        "testuser",
		Password:        "password123",
		FirstName:       "Test",
		LastName:        "User",
		PhoneNumber:     "1234567890",
		Gender:          pb.CreateUserRequest_MALE,
		IsEmailVerified: false,
	}

	resp, err := suite.UserHandler.CreateUser(ctx, req)

	require.NoError(suite.T(), err)
	require.NotNil(suite.T(), resp)
	assert.Equal(suite.T(), codes.OK.String(), resp.Code)
	assert.Equal(suite.T(), "success", resp.Status)
	assert.Equal(suite.T(), "user created successfully", resp.Message)
	assert.NotEmpty(suite.T(), resp.GetData().(*pb.StandardUserResponse_GetUserResponse).GetUserResponse.GetId())
	assert.Equal(suite.T(), req.Email, resp.GetData().(*pb.StandardUserResponse_GetUserResponse).GetUserResponse.GetEmail())
	assert.Equal(suite.T(), req.Username, resp.GetData().(*pb.StandardUserResponse_GetUserResponse).GetUserResponse.GetUsername())

	var createdUser entity.User
	err = suite.DB.Where("email = ?", req.Email).First(&createdUser).Error
	require.NoError(suite.T(), err)
	assert.Equal(suite.T(), req.FirstName, createdUser.FirstName)
}

func (suite *GrpcHandlerTestSuite) TestCreateUser_InvalidEmail() {
	ctx := context.Background()
	req := &pb.CreateUserRequest{
		Email:           "", // Invalid email
		Username:        "testuser",
		Password:        "password123",
		FirstName:       "Test",
		LastName:        "User",
		PhoneNumber:     "1234567890",
		Gender:          pb.CreateUserRequest_MALE,
		IsEmailVerified: false,
	}

	resp, err := suite.UserHandler.CreateUser(ctx, req)

	require.Error(suite.T(), err)
	st, ok := status.FromError(err)

	require.True(suite.T(), ok)
	assert.Equal(suite.T(), codes.InvalidArgument, st.Code())
	assert.Equal(suite.T(), "email is required", st.Message())
	assert.Nil(suite.T(), resp)
	//suite.T().Logf("Error Message: %s", st)
}

func (suite *GrpcHandlerTestSuite) TestCreateUser_DuplicateEmail() {
	ctx := context.Background()
	existingUser := &entity.User{
		Email:    "existing@example.com",
		Username: "existinguser",
	}
	err := suite.DB.Create(&existingUser).Error
	require.NoError(suite.T(), err)

	req := &pb.CreateUserRequest{
		Email:           "existing@example.com", // Duplicate email
		Username:        "newuser",
		Password:        "password123",
		FirstName:       "Test",
		LastName:        "User",
		PhoneNumber:     "1234567890",
		Gender:          pb.CreateUserRequest_MALE,
		IsEmailVerified: false,
	}

	resp, err := suite.UserHandler.CreateUser(ctx, req)

	require.Error(suite.T(), err)
	st, ok := status.FromError(err)
	require.True(suite.T(), ok)
	assert.Equal(suite.T(), codes.AlreadyExists, st.Code())
	assert.Equal(suite.T(), "email or username already exists", st.Message())
	assert.Nil(suite.T(), resp)
	//suite.T().Logf("Error: %s", st)
}

// get user by id
func (suite *GrpcHandlerTestSuite) TestGetUser_Success() {
	ctx := context.Background()
	userID := uuid.New().String()
	req := &pb.GetUserRequest{Id: userID}

	expectedUser := &entity.User{
		ID:        uuid.MustParse(userID),
		Email:     "test1@example.com",
		Username:  "testuser1",
		FirstName: "Test",
		LastName:  "User",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err := suite.DB.Create(&expectedUser).Error
	require.NoError(suite.T(), err, "Failed to create test user in database")

	resp, err := suite.UserHandler.GetUser(ctx, req)

	require.NoError(suite.T(), err)
	require.NotNil(suite.T(), resp)
	assert.Equal(suite.T(), codes.OK.String(), resp.Code)
	assert.Equal(suite.T(), "success", resp.Status)
	assert.Equal(suite.T(), "user retrieved successfully", resp.Message)
	actualUser := resp.GetData().(*pb.StandardUserResponse_GetUserResponse).GetUserResponse
	assert.Equal(suite.T(), expectedUser.ID.String(), actualUser.GetId())
	assert.Equal(suite.T(), expectedUser.Email, actualUser.GetEmail())
	assert.Equal(suite.T(), expectedUser.Username, actualUser.GetUsername())
	assert.Equal(suite.T(), expectedUser.FirstName, actualUser.GetFirstName())
}

func (suite *GrpcHandlerTestSuite) TestGetUser_NotFound() {
	ctx := context.Background()
	userID := uuid.New().String()
	req := &pb.GetUserRequest{Id: userID}

	resp, err := suite.UserHandler.GetUser(ctx, req)

	require.Error(suite.T(), err)
	st, ok := status.FromError(err)
	require.True(suite.T(), ok)
	assert.Equal(suite.T(), codes.Canceled, st.Code())
	assert.Equal(suite.T(), "record not found", st.Message())
	assert.Nil(suite.T(), resp)
}

// update user
func (suite *GrpcHandlerTestSuite) TestUpdateUser_Success() {
	ctx := context.Background()
	userID := uuid.New()
	userIDStr := userID.String()

	initialUser := &entity.User{
		ID:          userID,
		Email:       "initial@example.com",
		Username:    "initialuser",
		FirstName:   "Initial",
		LastName:    "User",
		PhoneNumber: "08123456789",
		Gender:      entity.Male,
	}
	err := suite.DB.Create(&initialUser).Error
	require.NoError(suite.T(), err, "Failed to create initial user")

	updateReq := &pb.UpdateUserRequest{
		User: &pb.User{
			Id:          userIDStr,
			Email:       "updated@example.com",
			Username:    "updateduser",
			FirstName:   "Updated",
			LastName:    "User",
			PhoneNumber: "08987654321",
			Gender:      pb.User_FEMALE,
		},
	}

	resp, err := suite.UserHandler.UpdateUser(ctx, updateReq)

	require.NoError(suite.T(), err)
	require.NotNil(suite.T(), resp)
	assert.Equal(suite.T(), codes.OK.String(), resp.Code)
	assert.Equal(suite.T(), "success", resp.Status)
	assert.Equal(suite.T(), "user updated successfully", resp.Message)

	updatedUserResp := resp.GetData().(*pb.StandardUserResponse_GetUserResponse).GetUserResponse
	assert.Equal(suite.T(), userIDStr, updatedUserResp.GetId())
	assert.Equal(suite.T(), "updated@example.com", updatedUserResp.GetEmail())
	assert.Equal(suite.T(), "updateduser", updatedUserResp.GetUsername())
	assert.Equal(suite.T(), "Updated", updatedUserResp.GetFirstName())
	assert.Equal(suite.T(), "User", updatedUserResp.GetLastName())
	assert.Equal(suite.T(), "08987654321", updatedUserResp.GetPhoneNumber())
	assert.Equal(suite.T(), pb.User_FEMALE, updatedUserResp.GetGender())

	var updatedUser entity.User
	err = suite.DB.First(&updatedUser, "id = ?", userID).Error
	require.NoError(suite.T(), err, "Failed to retrieve updated user from database")
	assert.Equal(suite.T(), "updated@example.com", updatedUser.Email)
	assert.Equal(suite.T(), "updateduser", updatedUser.Username)
	assert.Equal(suite.T(), "Updated", updatedUser.FirstName)
	assert.Equal(suite.T(), "User", updatedUser.LastName)
	assert.Equal(suite.T(), "08987654321", updatedUser.PhoneNumber)
	assert.Equal(suite.T(), entity.Female, updatedUser.Gender)
}

func (suite *GrpcHandlerTestSuite) TestUpdateUser_IDNotFound() {
	ctx := context.Background()
	nonExistentUserID := uuid.New().String()
	updateReq := &pb.UpdateUserRequest{
		User: &pb.User{
			Id:          nonExistentUserID,
			Email:       "updated@example.com",
			Username:    "updateduser",
			FirstName:   "Updated",
			LastName:    "User",
			PhoneNumber: "08987654321",
			Gender:      pb.User_MALE,
		},
	}

	resp, err := suite.UserHandler.UpdateUser(ctx, updateReq)

	require.Error(suite.T(), err)
	st, ok := status.FromError(err)
	require.True(suite.T(), ok)
	assert.Equal(suite.T(), codes.InvalidArgument, st.Code())
	assert.Contains(suite.T(), st.Message(), "failed to retrieve updated user")
	assert.Nil(suite.T(), resp)

	var user entity.User
	err = suite.DB.Where("id = ?", nonExistentUserID).First(&user).Error
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), gorm.ErrRecordNotFound, err)

	//suite.T().Logf("Error Message: %s", st)
}
