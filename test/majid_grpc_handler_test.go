package test

import (
	"context"
	"github.com/google/uuid"
	pb "github.com/mnadev/limestone/gen/go"
	"github.com/mnadev/limestone/internal/application/domain/entity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"time"
)

func (suite *GrpcHandlerTestSuite) TestCreateMasjid_Success() {
	ctx := context.Background()
	req := &pb.CreateMasjidRequest{
		Masjid: &pb.Masjid{
			Name:       "Test Masjid",
			IsVerified: false,
			Address: &pb.Masjid_Address{
				AddressLine_1: "Line 1",
				City:          "Test City",
				CountryCode:   "TC",
			},
			PhoneNumber: &pb.Masjid_PhoneNumber{
				Number: "1234567890",
			},
		},
	}

	resp, err := suite.MasjidHandler.CreateMasjid(ctx, req)
	require.NoError(suite.T(), err)
	require.NotNil(suite.T(), resp)
	assert.Equal(suite.T(), codes.OK.String(), resp.GetCode())
	assert.Equal(suite.T(), "success", resp.GetStatus())
	assert.Equal(suite.T(), "masjid created successfully", resp.Message)
	createdMasjid := resp.GetData().(*pb.StandardMasjidResponse_Masjid).Masjid
	assert.NotEmpty(suite.T(), createdMasjid.GetId())
	assert.Equal(suite.T(), "Test Masjid", createdMasjid.GetName())
	assert.False(suite.T(), createdMasjid.GetIsVerified())
	assert.Equal(suite.T(), "Line 1", createdMasjid.GetAddress().GetAddressLine_1())
	assert.Equal(suite.T(), "Test City", createdMasjid.GetAddress().GetCity())
	assert.Equal(suite.T(), "TC", createdMasjid.GetAddress().GetCountryCode())
	assert.Equal(suite.T(), "1234567890", createdMasjid.GetPhoneNumber().GetNumber())

	var dbMasjid entity.Masjid
	err = suite.DB.Where("id = ?", createdMasjid.GetId()).First(&dbMasjid).Error
	require.NoError(suite.T(), err)
	assert.Equal(suite.T(), "Test Masjid", dbMasjid.Name)

	err = suite.DB.Delete(&dbMasjid).Error
	require.NoError(suite.T(), err)
}

func (suite *GrpcHandlerTestSuite) TestGetMasjid_Success() {
	ctx := context.Background()
	existingMasjid := &entity.Masjid{
		ID:   uuid.New(),
		Name: "Existing Masjid",
		Address: entity.Address{
			AddressLine1: "Existing Line 1",
			City:         "Existing City",
			CountryCode:  "EC",
		},
		PhoneNumber: entity.PhoneNumber{
			Number: "9876543210",
		},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	err := suite.DB.Create(&existingMasjid).Error
	require.NoError(suite.T(), err)

	req := &pb.GetMasjidRequest{
		Id: existingMasjid.ID.String(),
	}

	resp, err := suite.MasjidHandler.GetMasjid(ctx, req)
	require.NoError(suite.T(), err)
	require.NotNil(suite.T(), resp)
	assert.Equal(suite.T(), codes.OK.String(), resp.GetCode())
	assert.Equal(suite.T(), "success", resp.GetStatus())
	assert.Equal(suite.T(), "masjid retrieved successfully", resp.Message)
	retrievedMasjid := resp.GetData().(*pb.StandardMasjidResponse_Masjid).Masjid
	assert.Equal(suite.T(), existingMasjid.ID.String(), retrievedMasjid.GetId())
	assert.Equal(suite.T(), "Existing Masjid", retrievedMasjid.GetName())
	assert.Equal(suite.T(), "Existing Line 1", retrievedMasjid.GetAddress().GetAddressLine_1())
	assert.Equal(suite.T(), "Existing City", retrievedMasjid.GetAddress().GetCity())
	assert.Equal(suite.T(), "EC", retrievedMasjid.GetAddress().GetCountryCode())
	assert.Equal(suite.T(), "9876543210", retrievedMasjid.GetPhoneNumber().GetNumber())

	err = suite.DB.Delete(&existingMasjid).Error
	require.NoError(suite.T(), err)
}

func (suite *GrpcHandlerTestSuite) TestGetMasjid_NotFound() {
	ctx := context.Background()
	nonExistingID := uuid.New().String()
	req := &pb.GetMasjidRequest{
		Id: nonExistingID,
	}

	resp, err := suite.MasjidHandler.GetMasjid(ctx, req)
	require.Error(suite.T(), err)
	st, ok := status.FromError(err)
	require.True(suite.T(), ok)
	assert.Equal(suite.T(), codes.NotFound, st.Code())
	assert.Equal(suite.T(), "masjid not found", st.Message())
	assert.Nil(suite.T(), resp)
}

func (suite *GrpcHandlerTestSuite) TestDeleteMasjid_Success() {
	ctx := context.Background()

	existingMasjid := &entity.Masjid{
		ID:        uuid.New(),
		Name:      "Masjid to Delete",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	err := suite.DB.Create(&existingMasjid).Error
	require.NoError(suite.T(), err)

	req := &pb.DeleteMasjidRequest{
		Id: existingMasjid.ID.String(),
	}

	resp, err := suite.MasjidHandler.DeleteMasjid(ctx, req)
	require.NoError(suite.T(), err)
	require.NotNil(suite.T(), resp)
	assert.Equal(suite.T(), codes.OK.String(), resp.GetCode())
	assert.Equal(suite.T(), "success", resp.GetStatus())
	assert.Equal(suite.T(), "masjid deleted successfully", resp.Message)
	assert.NotNil(suite.T(), resp.GetData())
	_, ok := resp.GetData().(*pb.StandardMasjidResponse_DeleteMasjidResponse)
	assert.True(suite.T(), ok)

	var deletedMasjid entity.Masjid
	err = suite.DB.Where("id = ?", existingMasjid.ID).First(&deletedMasjid).Error
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), err.Error(), "record not found")
}

func (suite *GrpcHandlerTestSuite) TestUpdateMasjid_Success() {
	ctx := context.Background()

	existingMasjid := &entity.Masjid{
		ID:   uuid.New(),
		Name: "Old Masjid Name",
		Address: entity.Address{
			AddressLine1: "Old Line 1",
			City:         "Old City",
			CountryCode:  "OC",
		},
		PhoneNumber: entity.PhoneNumber{
			Number: "1112223333",
		},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	err := suite.DB.Create(&existingMasjid).Error
	require.NoError(suite.T(), err)

	updatedName := "New Masjid Name"
	updatedCity := "New City"

	req := &pb.UpdateMasjidRequest{
		Masjid: &pb.Masjid{
			Id:   existingMasjid.ID.String(),
			Name: updatedName,
			Address: &pb.Masjid_Address{
				City: updatedCity,
			},
		},
	}

	resp, err := suite.MasjidHandler.UpdateMasjid(ctx, req)
	require.NoError(suite.T(), err)
	require.NotNil(suite.T(), resp)
	assert.Equal(suite.T(), codes.OK.String(), resp.GetCode())
	assert.Equal(suite.T(), "success", resp.GetStatus())
	assert.Equal(suite.T(), "masjid updated successfully", resp.Message)
	updatedMasjidResp := resp.GetData().(*pb.StandardMasjidResponse_Masjid).Masjid
	assert.Equal(suite.T(), existingMasjid.ID.String(), updatedMasjidResp.GetId())
	assert.Equal(suite.T(), updatedName, updatedMasjidResp.GetName())
	assert.Equal(suite.T(), updatedCity, updatedMasjidResp.GetAddress().GetCity())

	var dbMasjid entity.Masjid
	err = suite.DB.Where("id = ?", existingMasjid.ID).First(&dbMasjid).Error
	require.NoError(suite.T(), err)
	assert.Equal(suite.T(), updatedName, dbMasjid.Name)
	assert.Equal(suite.T(), updatedCity, dbMasjid.Address.City)

	err = suite.DB.Delete(&dbMasjid).Error
	require.NoError(suite.T(), err)
	//suite.T().Logf("Error Message: %s", updatedMasjidResp)
}
