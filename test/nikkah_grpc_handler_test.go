package test

import (
	"context"
	"github.com/google/uuid"
	pb "github.com/mnadev/limestone/gen/go"
	"github.com/mnadev/limestone/internal/application/domain/entity"
	"github.com/mnadev/limestone/internal/infrastructure/auth"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"time"
)

func (suite *GrpcHandlerTestSuite) TestCreateNikkahProfile_Success() {

	userCtx := context.Background()
	userReq := &pb.CreateUserRequest{
		Email:           "testprofile@example.com",
		Username:        "testuser_profile",
		Password:        "password123",
		FirstName:       "Profile",
		LastName:        "User",
		PhoneNumber:     "1234567890",
		Gender:          pb.CreateUserRequest_MALE,
		IsEmailVerified: false,
	}

	userResp, userErr := suite.UserHandler.CreateUser(userCtx, userReq)
	require.NoError(suite.T(), userErr, "Failed to create test user for profile creation via handler")
	require.NotNil(suite.T(), userResp, "User creation response was nil")

	createdUserProto, ok := userResp.GetData().(*pb.StandardUserResponse_GetUserResponse)
	require.True(suite.T(), ok, "Failed to cast user response data to GetUserResponse")
	testUserID := createdUserProto.GetUserResponse.GetId()
	require.NotEmpty(suite.T(), testUserID, "Created user ID was empty")

	nikkahCtx := context.WithValue(context.Background(), auth.UserIDContextKey, testUserID)

	birthDate := time.Date(1990, 5, 15, 0, 0, 0, 0, time.UTC)
	req := &pb.CreateNikkahProfileRequest{
		Profile: &pb.NikkahProfile{
			Name:   "Test Nikkah Profile",
			Gender: pb.NikkahProfile_FEMALE,
			BirthDate: &pb.NikkahProfile_BirthDate{
				Year:  int32(birthDate.Year()),
				Month: pb.NikkahProfile_BirthDate_Month(int32(birthDate.Month())),
				Day:   int32(birthDate.Day()),
			},
			UserId: testUserID,
		},
	}

	resp, err := suite.NikkahHandler.CreateNikkahProfile(nikkahCtx, req)

	require.NoError(suite.T(), err, "Expected no error for successful profile creation")
	require.NotNil(suite.T(), resp, "Expected non-nil response for successful profile creation")

	assert.Equal(suite.T(), codes.OK.String(), resp.GetCode(), "Expected code to be OK")
	assert.Equal(suite.T(), "success", resp.GetStatus(), "Expected status to be 'success'")
	assert.Equal(suite.T(), "Nikkah profile created successfully", resp.GetMessage(), "Expected correct success message")

	profileData := resp.GetNikkahProfile()
	require.NotNil(suite.T(), profileData, "Expected profile data in response")

	assert.NotEmpty(suite.T(), profileData.GetId(), "Expected profile ID to be generated")
	assert.Equal(suite.T(), req.GetProfile().GetName(), profileData.GetName(), "Expected name to match")
	assert.Equal(suite.T(), req.GetProfile().GetGender(), profileData.GetGender(), "Expected gender to match")
	assert.Equal(suite.T(), req.GetProfile().GetBirthDate().GetYear(), profileData.GetBirthDate().GetYear(), "Expected birth year to match")
	assert.Equal(suite.T(), req.GetProfile().GetBirthDate().GetMonth(), profileData.GetBirthDate().GetMonth(), "Expected birth month to match")
	assert.Equal(suite.T(), req.GetProfile().GetBirthDate().GetDay(), profileData.GetBirthDate().GetDay(), "Expected birth day to match")
	assert.Equal(suite.T(), req.GetProfile().GetUserId(), profileData.GetUserId(), "Expected user ID to match")
	assert.NotNil(suite.T(), profileData.GetCreateTime(), "Expected CreatedAt to be set")
	assert.NotNil(suite.T(), profileData.GetUpdateTime(), "Expected UpdatedAt to be set")

	var createdProfile entity.NikkahProfile
	err = suite.DB.Where("user_id = ?", testUserID).First(&createdProfile).Error
	require.NoError(suite.T(), err, "Failed to find created profile in database")
	assert.Equal(suite.T(), profileData.GetId(), createdProfile.ID.String(), "DB profile ID should match response ID")
	assert.Equal(suite.T(), req.GetProfile().GetName(), createdProfile.Name, "DB profile name should match")
	assert.Equal(suite.T(), testUserID, createdProfile.UserID, "DB profile user ID should match")
}

func (suite *GrpcHandlerTestSuite) TestCreateNikkahProfile_Unauthenticated() {
	ctx := context.Background()
	birthDate := time.Date(1990, 5, 15, 0, 0, 0, 0, time.UTC)
	req := &pb.CreateNikkahProfileRequest{
		Profile: &pb.NikkahProfile{
			Name:   "Test Nikkah Profile",
			Gender: pb.NikkahProfile_FEMALE,
			BirthDate: &pb.NikkahProfile_BirthDate{
				Year:  int32(birthDate.Year()),
				Month: pb.NikkahProfile_BirthDate_Month(int32(birthDate.Month())),
				Day:   int32(birthDate.Day()),
			},
			UserId: "some_user_id",
		},
	}

	resp, err := suite.NikkahHandler.CreateNikkahProfile(ctx, req)

	require.Error(suite.T(), err)
	st, ok := status.FromError(err)
	require.True(suite.T(), ok)
	assert.Equal(suite.T(), codes.Unauthenticated, st.Code())
	assert.Equal(suite.T(), "Unauthenticated: user ID not found in context", st.Message())
	assert.Nil(suite.T(), resp)
}

func (suite *GrpcHandlerTestSuite) TestGetSelfNikkahProfile_Success() {
	// 1. Arrange: Buat user menggunakan UserHandler
	userCtx := context.Background()
	userReq := &pb.CreateUserRequest{
		Email:           "getselfprofile@example.com", // Ubah email agar unik
		Username:        "getselfuser",
		Password:        "password123",
		FirstName:       "GetSelf",
		LastName:        "User",
		PhoneNumber:     "1234567890",
		Gender:          pb.CreateUserRequest_MALE,
		IsEmailVerified: false,
	}

	userResp, userErr := suite.UserHandler.CreateUser(userCtx, userReq)
	require.NoError(suite.T(), userErr, "Failed to create test user via handler")
	require.NotNil(suite.T(), userResp, "User creation response was nil")

	// Validasi respons pembuatan user (opsional, tapi baik untuk keandalan)
	assert.Equal(suite.T(), codes.OK.String(), userResp.GetCode(), "Expected user creation code to be OK")
	assert.Equal(suite.T(), "success", userResp.GetStatus(), "Expected user creation status to be 'success'")

	createdUserProto, ok := userResp.GetData().(*pb.StandardUserResponse_GetUserResponse)
	require.True(suite.T(), ok, "Failed to cast user response data to GetUserResponse")
	testUserID := createdUserProto.GetUserResponse.GetId()
	require.NotEmpty(suite.T(), testUserID, "Created user ID was empty")

	// 2. Sekarang, buat Nikkah profile untuk user yang baru dibuat
	//    Kali ini, kita buat langsung di DB karena GetSelfNikkahProfile
	//    tidak membuat profile, melainkan mengambilnya.
	profileID := uuid.New()
	birthDate := time.Date(1992, time.June, 10, 0, 0, 0, 0, time.UTC)

	expectedProfile := &entity.NikkahProfile{
		ID:     profileID,
		UserID: testUserID, // Menggunakan ID user yang baru dibuat
		Name:   "Retrieve Self Profile",
		Gender: entity.Gender(pb.NikkahProfile_FEMALE),
		BirthDate: entity.BirthDate{ // Pastikan `entity.BirthDateEntity` adalah tipe yang benar
			Year:  int32(birthDate.Year()),
			Month: entity.Month(birthDate.Month()),
			Day:   int8(birthDate.Day()),
		},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	err := suite.DB.Create(expectedProfile).Error
	require.NoError(suite.T(), err, "Failed to create test nikkah profile directly in DB")

	// 3. Set up the context dengan User ID dari user yang baru dibuat
	ctx := context.WithValue(context.Background(), auth.UserIDContextKey, testUserID)
	req := &pb.GetSelfNikkahProfileRequest{}

	// 4. Act: Panggil handler GetSelfNikkahProfile
	resp, err := suite.NikkahHandler.GetSelfNikkahProfile(ctx, req)

	// 5. Assert: Periksa respons
	require.NoError(suite.T(), err, "Expected no error for successful profile retrieval")
	require.NotNil(suite.T(), resp, "Expected non-nil response for successful profile retrieval")

	assert.Equal(suite.T(), codes.OK.String(), resp.GetCode(), "Expected code to be OK")
	assert.Equal(suite.T(), "success", resp.GetStatus(), "Expected status to be 'success'")
	assert.Equal(suite.T(), "Self nikkah profile retrieved successfully", resp.GetMessage(), "Expected correct success message")

	retrievedProfileProto := resp.GetNikkahProfile()
	require.NotNil(suite.T(), retrievedProfileProto, "Expected Nikkah profile data in response")

	assert.Equal(suite.T(), expectedProfile.ID.String(), retrievedProfileProto.GetId(), "Expected profile ID to match")
	assert.Equal(suite.T(), expectedProfile.Name, retrievedProfileProto.GetName(), "Expected name to match")
	assert.Equal(suite.T(), pb.NikkahProfile_Gender(pb.NikkahProfile_Gender_value[string(expectedProfile.Gender)]), retrievedProfileProto.GetGender(), "Expected gender to match")
	assert.Equal(suite.T(), int32(expectedProfile.BirthDate.Year), retrievedProfileProto.GetBirthDate().GetYear(), "Expected birth year to match")
	assert.Equal(suite.T(), pb.NikkahProfile_BirthDate_Month(int32(expectedProfile.BirthDate.Month)), retrievedProfileProto.GetBirthDate().GetMonth(), "Expected birth month to match")
	assert.Equal(suite.T(), int32(expectedProfile.BirthDate.Day), retrievedProfileProto.GetBirthDate().GetDay(), "Expected birth day to match")
	assert.Equal(suite.T(), expectedProfile.UserID, retrievedProfileProto.GetUserId(), "Expected user ID to match")
	assert.NotNil(suite.T(), retrievedProfileProto.GetCreateTime(), "Expected CreateTime to be set")
	assert.NotNil(suite.T(), retrievedProfileProto.GetUpdateTime(), "Expected UpdateTime to be set")
}

func (suite *GrpcHandlerTestSuite) TestGetSelfNikkahProfile_ProfileNotFound() {
	userCtx := context.Background()
	userReq := &pb.CreateUserRequest{
		Email:           "nouserprofile@example.com", // Email unik
		Username:        "nouserprofile",
		Password:        "password123",
		FirstName:       "NoProfile",
		LastName:        "User",
		PhoneNumber:     "1111111111",
		Gender:          pb.CreateUserRequest_FEMALE,
		IsEmailVerified: false,
	}

	userResp, userErr := suite.UserHandler.CreateUser(userCtx, userReq)
	require.NoError(suite.T(), userErr, "Failed to create test user for ProfileNotFound scenario")
	require.NotNil(suite.T(), userResp, "User creation response was nil")

	createdUserProto, ok := userResp.GetData().(*pb.StandardUserResponse_GetUserResponse)
	require.True(suite.T(), ok, "Failed to cast user response data to GetUserResponse")
	testUserID := createdUserProto.GetUserResponse.GetId()
	require.NotEmpty(suite.T(), testUserID, "Created user ID was empty")
	
	ctx := context.WithValue(context.Background(), auth.UserIDContextKey, testUserID)
	req := &pb.GetSelfNikkahProfileRequest{}

	// 2. Act: Panggil handler
	resp, err := suite.NikkahHandler.GetSelfNikkahProfile(ctx, req)

	// 3. Assert: Periksa error
	require.Error(suite.T(), err, "Expected error for profile not found")
	st, ok := status.FromError(err)
	require.True(suite.T(), ok, "Expected gRPC status error")
	assert.Equal(suite.T(), codes.NotFound, st.Code(), "Expected NotFound error code")
	assert.Equal(suite.T(), "Nikkah profile not found for this user", st.Message(), "Expected correct error message")
	assert.Nil(suite.T(), resp, "Expected nil response for error")
}
