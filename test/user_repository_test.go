package test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/mnadev/limestone/gen/go"
	"github.com/mnadev/limestone/internal/application/domain/entity"
	grpc_handler "github.com/mnadev/limestone/internal/application/handler"
	"github.com/mnadev/limestone/internal/application/services"
	"github.com/mnadev/limestone/internal/infrastructure/auth" // Import auth package
	"github.com/mnadev/limestone/test/mocks"
)

type GrpcHandlerTestSuite struct {
	suite.Suite
	MockUserRepo *mocks.MockUserRepository
	UserHandler  *grpc_handler.UserGrpcHandler
}

func (suite *GrpcHandlerTestSuite) SetupTest() {
	suite.MockUserRepo = new(mocks.MockUserRepository)
	userService := services.NewUserService(suite.MockUserRepo)
	suite.UserHandler = grpc_handler.NewUserGrpcHandler(userService)

	auth.ResetRequireRole()
}

func (suite *GrpcHandlerTestSuite) TestCreateUser_Success() {
	ctx := context.Background()
	req := &pb.CreateUserRequest{
		Email:           "test@example.com",
		Username:        "testuser",
		Password:        "password12345",
		FirstName:       "Test",
		LastName:        "User",
		PhoneNumber:     "1234567890",
		Gender:          pb.CreateUserRequest_MALE,
		Role:            pb.CreateUserRequest_MASJID_MEMBER,
		IsEmailVerified: true,
	}

	returnedUser := &entity.User{
		ID:          uuid.New(),
		Email:       req.Email,
		Username:    req.Username,
		FirstName:   req.FirstName,
		LastName:    req.LastName,
		PhoneNumber: req.PhoneNumber,
		Gender:      entity.Gender(req.GetGender().String()),
		Role:        entity.Role(req.GetRole().String()),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	suite.MockUserRepo.On("Create", mock.Anything, mock.AnythingOfType("*entity.User")).Return(returnedUser, nil).Once()

	resp, err := suite.UserHandler.CreateUser(ctx, req)

	require.NoError(suite.T(), err)
	require.NotNil(suite.T(), resp)
	assert.Equal(suite.T(), codes.OK.String(), resp.Code)
	assert.Equal(suite.T(), "success", resp.Status)
	assert.Equal(suite.T(), "user created successfully", resp.Message)
	assert.Equal(suite.T(), returnedUser.ID.String(), resp.GetData().(*pb.StandardUserResponse_GetUserResponse).GetUserResponse.GetId())
	assert.Equal(suite.T(), req.Email, resp.GetData().(*pb.StandardUserResponse_GetUserResponse).GetUserResponse.GetEmail())
	assert.Equal(suite.T(), req.Username, resp.GetData().(*pb.StandardUserResponse_GetUserResponse).GetUserResponse.GetUsername())

	suite.MockUserRepo.AssertExpectations(suite.T())
}

func (suite *GrpcHandlerTestSuite) TestCreateUser_MissingRequiredFields() {
	testCases := []struct {
		name          string
		req           *pb.CreateUserRequest
		expectedCode  codes.Code
		expectedError string
	}{
		{
			name: "Empty Email",
			req: &pb.CreateUserRequest{
				Username:    "user",
				Password:    "password12345",
				FirstName:   "F",
				LastName:    "L",
				PhoneNumber: "1",
				Role:        pb.CreateUserRequest_MASJID_ADMIN,
			},
			expectedCode:  codes.InvalidArgument,
			expectedError: "email is required",
		},
		{
			name: "Empty Username",
			req: &pb.CreateUserRequest{
				Email:       "e@e.com",
				Password:    "password12345",
				FirstName:   "F",
				LastName:    "L",
				PhoneNumber: "1",
				Role:        pb.CreateUserRequest_MASJID_ADMIN,
			},
			expectedCode:  codes.InvalidArgument,
			expectedError: "username is required",
		},
		{
			name: "Empty Password",
			req: &pb.CreateUserRequest{
				Email:       "e@e.com",
				Username:    "user",
				FirstName:   "F",
				LastName:    "L",
				PhoneNumber: "1",
				Role:        pb.CreateUserRequest_MASJID_ADMIN,
			},
			expectedCode:  codes.InvalidArgument,
			expectedError: "password is required",
		},
		{
			name: "Password Too Short",
			req: &pb.CreateUserRequest{
				Email:       "e@e.com",
				Username:    "user",
				Password:    "short",
				FirstName:   "F",
				LastName:    "L",
				PhoneNumber: "1",
				Role:        pb.CreateUserRequest_MASJID_ADMIN,
			},
			expectedCode:  codes.InvalidArgument,
			expectedError: "password must be at least 8 characters",
		},
		{
			name: "Empty First Name",
			req: &pb.CreateUserRequest{
				Email:       "e@e.com",
				Username:    "user",
				Password:    "password12345",
				LastName:    "L",
				PhoneNumber: "1",
				Role:        pb.CreateUserRequest_MASJID_ADMIN,
			},
			expectedCode:  codes.InvalidArgument,
			expectedError: "first name is required",
		},
		{
			name: "Empty Last Name",
			req: &pb.CreateUserRequest{
				Email:       "e@e.com",
				Username:    "user",
				Password:    "password12345",
				FirstName:   "F",
				PhoneNumber: "1",
				Role:        pb.CreateUserRequest_MASJID_ADMIN,
			},
			expectedCode:  codes.InvalidArgument,
			expectedError: "last name is required",
		},
		{
			name: "Empty Phone Number",
			req: &pb.CreateUserRequest{
				Email:     "e@e.com",
				Username:  "user",
				Password:  "password12345",
				FirstName: "F",
				LastName:  "L",
				Role:      pb.CreateUserRequest_MASJID_ADMIN,
			},
			expectedCode:  codes.InvalidArgument,
			expectedError: "phone number is required",
		},
	}

	for _, tc := range testCases {
		suite.T().Run(tc.name, func(t *testing.T) {
			ctx := context.Background()
			resp, err := suite.UserHandler.CreateUser(ctx, tc.req)

			require.Error(t, err)
			st, ok := status.FromError(err)
			require.True(t, ok)
			assert.Equal(t, tc.expectedCode, st.Code())
			assert.Equal(t, tc.expectedError, st.Message())
			assert.Nil(t, resp)

			suite.MockUserRepo.AssertNotCalled(t, "Create", mock.Anything, mock.Anything) // Ensure Create is NOT called for invalid inputs
			suite.MockUserRepo.AssertExpectations(t)
		})
	}
}

func (suite *GrpcHandlerTestSuite) TestCreateUser_DuplicateEmailOrUsername() {
	ctx := context.Background()
	req := &pb.CreateUserRequest{
		Email:           "existing@example.com",
		Username:        "existinguser",
		Password:        "password12345",
		FirstName:       "Test",
		LastName:        "User",
		PhoneNumber:     "1234567890",
		Gender:          pb.CreateUserRequest_MALE,
		Role:            pb.CreateUserRequest_MASJID_ADMIN,
		IsEmailVerified: false,
	}

	suite.MockUserRepo.On("Create", mock.Anything, mock.AnythingOfType("*entity.User")).Return(nil, errors.New("pq: duplicate key value violates unique constraint \"users_email_key\"")).Once()

	resp, err := suite.UserHandler.CreateUser(ctx, req)

	require.Error(suite.T(), err)
	st, ok := status.FromError(err)
	require.True(suite.T(), ok)
	assert.Equal(suite.T(), codes.AlreadyExists, st.Code())
	assert.Equal(suite.T(), "email or username already exists", st.Message())
	assert.Nil(suite.T(), resp)

	suite.MockUserRepo.AssertExpectations(suite.T())
}

func (suite *GrpcHandlerTestSuite) TestCreateUser_ServiceInternalError() {
	ctx := context.Background()
	req := &pb.CreateUserRequest{
		Email:           "valid@example.com",
		Username:        "validuser",
		Password:        "password12345",
		FirstName:       "Test",
		LastName:        "User",
		PhoneNumber:     "1234567890",
		Gender:          pb.CreateUserRequest_MALE,
		Role:            pb.CreateUserRequest_MASJID_ADMIN,
		IsEmailVerified: false,
	}

	suite.MockUserRepo.On("Create", mock.Anything, mock.AnythingOfType("*entity.User")).Return(nil, errors.New("database connection lost")).Once()

	resp, err := suite.UserHandler.CreateUser(ctx, req)

	require.Error(suite.T(), err)
	st, ok := status.FromError(err)
	require.True(suite.T(), ok)
	assert.Equal(suite.T(), codes.Internal, st.Code())
	assert.Equal(suite.T(), "database connection lost", st.Message()) // Memastikan error message diteruskan
	assert.Nil(suite.T(), resp)

	suite.MockUserRepo.AssertExpectations(suite.T())
}

func TestGrpcHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(GrpcHandlerTestSuite))
}
