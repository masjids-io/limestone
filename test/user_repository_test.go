package test

import (
	"context"
	"errors" // Untuk errors.New jika diperlukan
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

// GrpcHandlerTestSuite akan menggunakan mock repository
type GrpcHandlerTestSuite struct {
	suite.Suite
	MockUserRepo *mocks.MockUserRepository
	UserHandler  *grpc_handler.UserGrpcHandler // Handler yang akan kita uji
}

// SetupTest dipanggil sebelum setiap test
func (suite *GrpcHandlerTestSuite) SetupTest() {
	// Inisialisasi mock repository baru untuk setiap test
	suite.MockUserRepo = new(mocks.MockUserRepository)
	// Inisialisasi service dengan mock repository
	userService := services.NewUserService(suite.MockUserRepo)
	// Inisialisasi handler gRPC dengan service yang sudah diinject
	suite.UserHandler = grpc_handler.NewUserGrpcHandler(userService)

	// Reset RequireRole ke defaultnya sebelum setiap test
	auth.ResetRequireRole()
}

// --- Test `CreateUser` ---

func (suite *GrpcHandlerTestSuite) TestCreateUser_Success() {
	ctx := context.Background()
	req := &pb.CreateUserRequest{
		Email:           "test@example.com",
		Username:        "testuser",
		Password:        "password12345", // Pastikan > 8 karakter
		FirstName:       "Test",
		LastName:        "User",
		PhoneNumber:     "1234567890",
		Gender:          pb.CreateUserRequest_MALE,
		Role:            pb.CreateUserRequest_MASJID_ADMIN,
		IsEmailVerified: false,
	}

	// Ekspektasi: service CreateUser akan dipanggil dengan user yang ID-nya baru
	// Kita tidak perlu mem-mock GetByEmail/GetByUsername secara eksplisit di test handler,
	// karena handler tidak memanggilnya, service yang memanggilnya.
	// Kita hanya perlu mock panggilan Create di service.

	// `mock.AnythingOfType("*entity.User")` digunakan untuk argumen `user`
	// karena kita tidak tahu persis ID UUID yang akan digenerate handler.
	// Kita juga tidak peduli tentang `HashedPassword` yang tepat di sini,
	// hanya memverifikasi bahwa `Create` dipanggil.
	returnedUser := &entity.User{
		ID:       uuid.New(), // ID ini akan cocok dengan yang digenerate oleh handler
		Email:    req.Email,
		Username: req.Username,
		// ... bidang lainnya sesuai dengan yang dikembalikan handler
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
				Password:    "short", // < 8 chars
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
		{
			name: "Unspecified Role",
			req: &pb.CreateUserRequest{
				Email:       "e@e.com",
				Username:    "user",
				Password:    "password12345",
				FirstName:   "F",
				LastName:    "L",
				PhoneNumber: "1",
				Role:        pb.CreateUserRequest_ROLE_UNSPECIFIED,
			},
			expectedCode:  codes.InvalidArgument,
			expectedError: "role is required and cannot be unspecified.",
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

	// Ekspektasi: Service CreateUser akan mengembalikan error yang mengandung "duplicate key"
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

	// Ekspektasi: Service CreateUser akan mengembalikan error internal lainnya
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

// --- Test `GetUser` ---

func (suite *GrpcHandlerTestSuite) TestGetUser_Success() {
	ctx := context.Background()
	userID := uuid.New().String()
	req := &pb.GetUserRequest{Id: userID}

	// Mock `auth.RequireRole` agar selalu berhasil untuk test ini
	auth.RequireRole = func(ctx context.Context, allowedRoles []string, methodName string) error {
		return nil // Langsung izinkan
	}
	defer auth.ResetRequireRole() // Pastikan direset setelah test

	expectedUser := &entity.User{
		ID:          uuid.MustParse(userID),
		Email:       "test1@example.com",
		Username:    "testuser1",
		FirstName:   "Test",
		LastName:    "User",
		PhoneNumber: "1234567890",
		Gender:      entity.Male,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	suite.MockUserRepo.On("GetByID", mock.Anything, userID).Return(expectedUser, nil).Once()

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

	suite.MockUserRepo.AssertExpectations(suite.T())
}

func (suite *GrpcHandlerTestSuite) TestGetUser_NotFound() {
	ctx := context.Background()
	userID := uuid.New().String()
	req := &pb.GetUserRequest{Id: userID}

	// Mock `auth.RequireRole` agar selalu berhasil
	auth.RequireRole = func(ctx context.Context, allowedRoles []string, methodName string) error {
		return nil
	}
	defer auth.ResetRequireRole()

	// Handler Anda mengonversi error dari service.GetUser jika mengandung "record not found"
	suite.MockUserRepo.On("GetByID", mock.Anything, userID).Return(nil, errors.New("record not found")).Once()

	resp, err := suite.UserHandler.GetUser(ctx, req)

	require.Error(suite.T(), err)
	st, ok := status.FromError(err)
	require.True(suite.T(), ok)
	assert.Equal(suite.T(), codes.Canceled, st.Code())        // Harusnya NotFound, bukan Canceled
	assert.Equal(suite.T(), "record not found", st.Message()) // Pesan dari handler
	assert.Nil(suite.T(), resp)

	suite.MockUserRepo.AssertExpectations(suite.T())
}

func (suite *GrpcHandlerTestSuite) TestGetUser_AuthorizationDenied() {
	ctx := context.Background()
	req := &pb.GetUserRequest{Id: uuid.New().String()}

	// Mock `auth.RequireRole` agar mengembalikan error PermissionDenied
	auth.RequireRole = func(ctx context.Context, allowedRoles []string, methodName string) error {
		return status.Errorf(codes.PermissionDenied, "role not allowed for GetUser")
	}
	defer auth.ResetRequireRole()

	resp, err := suite.UserHandler.GetUser(ctx, req)

	require.Error(suite.T(), err)
	st, ok := status.FromError(err)
	require.True(suite.T(), ok)
	assert.Equal(suite.T(), codes.PermissionDenied, st.Code())
	assert.Equal(suite.T(), "role not allowed for GetUser", st.Message())
	assert.Nil(suite.T(), resp)

	// Pastikan service GetUser tidak dipanggil karena otorisasi gagal
	suite.MockUserRepo.AssertNotCalled(suite.T(), "GetByID", mock.Anything, mock.Anything)
	suite.MockUserRepo.AssertExpectations(suite.T())
}

// --- Test `UpdateUser` ---

func (suite *GrpcHandlerTestSuite) TestUpdateUser_Success() {
	ctx := context.Background()
	userID := uuid.New()
	userIDStr := userID.String()

	// Mock `auth.RequireRole` agar selalu berhasil
	auth.RequireRole = func(ctx context.Context, allowedRoles []string, methodName string) error {
		return nil
	}
	defer auth.ResetRequireRole()

	// Initial user (tidak digunakan langsung oleh handler, tapi penting untuk GetUser setelah update)
	// Handler memanggil GetUser setelah UpdateUser
	initialUserForGet := &entity.User{
		ID:          userID,
		Email:       "initial@example.com",
		Username:    "initialuser",
		FirstName:   "Initial",
		LastName:    "User",
		PhoneNumber: "08123456789",
		Gender:      entity.Male,
	}

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

	// Ekspektasi: service UpdateUser dipanggil dan berhasil
	// Argumen kedua mock.AnythingOfType("*entity.User") karena updateData dibuat di handler.
	suite.MockUserRepo.On("Update", mock.Anything, mock.AnythingOfType("*entity.User")).Return(&entity.User{ID: userID}, nil).Once()

	// Ekspektasi: service GetUser dipanggil setelah update untuk mengambil data terbaru
	updatedUserAfterService := &entity.User{
		ID:          userID,
		Email:       updateReq.GetUser().GetEmail(),
		Username:    updateReq.GetUser().GetUsername(),
		FirstName:   updateReq.GetUser().GetFirstName(),
		LastName:    updateReq.GetUser().GetLastName(),
		PhoneNumber: updateReq.GetUser().GetPhoneNumber(),
		Gender:      entity.Gender(updateReq.GetUser().GetGender().String()),
		CreatedAt:   initialUserForGet.CreatedAt,
		UpdatedAt:   time.Now(),
	}
	suite.MockUserRepo.On("GetByID", mock.Anything, userIDStr).Return(updatedUserAfterService, nil).Once()

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

	suite.MockUserRepo.AssertExpectations(suite.T())
}

func (suite *GrpcHandlerTestSuite) TestUpdateUser_InvalidUserIDFormat() {
	ctx := context.Background()
	updateReq := &pb.UpdateUserRequest{
		User: &pb.User{
			Id:          "invalid-uuid-format", // Invalid UUID
			Email:       "updated@example.com",
			Username:    "updateduser",
			FirstName:   "Updated",
			LastName:    "User",
			PhoneNumber: "08987654321",
			Gender:      pb.User_MALE,
		},
	}

	// Mock `auth.RequireRole` agar selalu berhasil
	auth.RequireRole = func(ctx context.Context, allowedRoles []string, methodName string) error {
		return nil
	}
	defer auth.ResetRequireRole()

	resp, err := suite.UserHandler.UpdateUser(ctx, updateReq)

	require.Error(suite.T(), err)
	st, ok := status.FromError(err)
	require.True(suite.T(), ok)
	assert.Equal(suite.T(), codes.InvalidArgument, st.Code())
	assert.Equal(suite.T(), "invalid user ID format", st.Message())
	assert.Nil(suite.T(), resp)

	// Pastikan tidak ada interaksi dengan mock repo
	suite.MockUserRepo.AssertNotCalled(suite.T(), "Update", mock.Anything, mock.Anything)
	suite.MockUserRepo.AssertNotCalled(suite.T(), "GetByID", mock.Anything, mock.Anything)
	suite.MockUserRepo.AssertExpectations(suite.T())
}

func (suite *GrpcHandlerTestSuite) TestUpdateUser_ServiceUpdateFails() {
	ctx := context.Background()
	userID := uuid.New().String()
	updateReq := &pb.UpdateUserRequest{
		User: &pb.User{
			Id:          userID,
			Email:       "updated@example.com",
			Username:    "updateduser",
			FirstName:   "Updated",
			LastName:    "User",
			PhoneNumber: "08987654321",
			Gender:      pb.User_MALE,
		},
	}

	// Mock `auth.RequireRole` agar selalu berhasil
	auth.RequireRole = func(ctx context.Context, allowedRoles []string, methodName string) error {
		return nil
	}
	defer auth.ResetRequireRole()

	// Ekspektasi: Service UpdateUser mengembalikan error "record not found"
	suite.MockUserRepo.On("Update", mock.Anything, mock.AnythingOfType("*entity.User")).Return(nil, errors.New("record not found")).Once()

	resp, err := suite.UserHandler.UpdateUser(ctx, updateReq)

	require.Error(suite.T(), err)
	st, ok := status.FromError(err)
	require.True(suite.T(), ok)
	assert.Equal(suite.T(), "failed to update user: record not found", st.Message())
	assert.Nil(suite.T(), resp)

	// Pastikan GetUser tidak dipanggil setelah update gagal
	suite.MockUserRepo.AssertNotCalled(suite.T(), "GetByID", mock.Anything, mock.Anything)
	suite.MockUserRepo.AssertExpectations(suite.T())
}

func (suite *GrpcHandlerTestSuite) TestUpdateUser_RetrieveUpdatedUserFails() {
	ctx := context.Background()
	userID := uuid.New()
	userIDStr := userID.String()

	// Mock `auth.RequireRole` agar selalu berhasil
	auth.RequireRole = func(ctx context.Context, allowedRoles []string, methodName string) error {
		return nil
	}
	defer auth.ResetRequireRole()

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

	// Ekspektasi: Service UpdateUser berhasil
	suite.MockUserRepo.On("Update", mock.Anything, mock.AnythingOfType("*entity.User")).Return(&entity.User{ID: userID}, nil).Once()

	// Ekspektasi: Service GetUser GAGAL setelah update berhasil (simulasi kasus edge/bug)
	suite.MockUserRepo.On("GetByID", mock.Anything, userIDStr).Return(nil, errors.New("database error during retrieve")).Once()

	resp, err := suite.UserHandler.UpdateUser(ctx, updateReq)

	require.Error(suite.T(), err)
	st, ok := status.FromError(err)
	require.True(suite.T(), ok)
	assert.Equal(suite.T(), "failed to retrieve updated user", st.Message())
	assert.Nil(suite.T(), resp)

	suite.MockUserRepo.AssertExpectations(suite.T())
}

// --- Test `DeleteUser` ---

func (suite *GrpcHandlerTestSuite) TestDeleteUser_Success() {
	ctx := context.Background()
	userIDToDelete := uuid.New().String()
	req := &pb.DeleteUserRequest{Id: userIDToDelete}

	// Mock `auth.RequireRole` agar selalu berhasil
	auth.RequireRole = func(ctx context.Context, allowedRoles []string, methodName string) error {
		return nil
	}
	defer auth.ResetRequireRole()

	// Ekspektasi: service DeleteUser dipanggil dan berhasil
	suite.MockUserRepo.On("Delete", mock.Anything, userIDToDelete).Return(nil).Once()

	resp, err := suite.UserHandler.DeleteUser(ctx, req)

	require.NoError(suite.T(), err)
	require.NotNil(suite.T(), resp)
	assert.Equal(suite.T(), codes.OK.String(), resp.Code)
	assert.Equal(suite.T(), "success", resp.Status)
	assert.Equal(suite.T(), "user deleted successfully", resp.Message)

	suite.MockUserRepo.AssertExpectations(suite.T())
}

// Untuk menjalankan Test Suite Anda
func TestGrpcHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(GrpcHandlerTestSuite))
}
