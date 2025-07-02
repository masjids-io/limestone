package test

import (
	"bytes"
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

const maxAdhanFileSizeMB = 5

func createValidAudioContent() []byte {
	return []byte{0xFF, 0xF3}
}

func (suite *GrpcHandlerTestSuite) TestCreateAdhan_Success() {
	ctx := context.Background()
	masjidID := uuid.New().String()
	validAudioContent := createValidAudioContent()

	req := &pb.CreateAdhanFileRequest{
		AdhanFile: &pb.AdhanFile{
			MasjidId: masjidID,
			File:     validAudioContent,
		},
	}

	resp, err := suite.AdhanHandler.CreateAdhan(ctx, req)

	require.NoError(suite.T(), err)
	require.NotNil(suite.T(), resp)
	assert.Equal(suite.T(), codes.OK.String(), resp.GetCode())
	assert.Equal(suite.T(), "success", resp.GetStatus())
	assert.Equal(suite.T(), "adhan file created successfully", resp.Message)
	retrievedAdhan := resp.GetData().(*pb.StandardAdhanResponse_AdhanFile).AdhanFile
	assert.NotEmpty(suite.T(), retrievedAdhan.GetId())
	assert.Equal(suite.T(), masjidID, retrievedAdhan.GetMasjidId())
	assert.Equal(suite.T(), validAudioContent, retrievedAdhan.GetFile())

	var createdAdhan entity.Adhan
	err = suite.DB.Where("id = ?", retrievedAdhan.GetId()).First(&createdAdhan).Error
	require.NoError(suite.T(), err)
	assert.Equal(suite.T(), masjidID, createdAdhan.MasjidId)
	assert.Equal(suite.T(), validAudioContent, createdAdhan.File)

	err = suite.DB.Delete(&createdAdhan).Error
	require.NoError(suite.T(), err)
}

func (suite *GrpcHandlerTestSuite) TestCreateAdhan_NoAdhanFile() {
	ctx := context.Background()
	req := &pb.CreateAdhanFileRequest{}

	resp, err := suite.AdhanHandler.CreateAdhan(ctx, req)

	require.Error(suite.T(), err)
	st, ok := status.FromError(err)
	require.True(suite.T(), ok)
	assert.Equal(suite.T(), codes.InvalidArgument, st.Code())
	assert.Equal(suite.T(), "adhan file data is required", st.Message())
	assert.Nil(suite.T(), resp)
}

func (suite *GrpcHandlerTestSuite) TestCreateAdhan_NoMasjidID() {
	ctx := context.Background()
	validAudioContent := createValidAudioContent()
	req := &pb.CreateAdhanFileRequest{
		AdhanFile: &pb.AdhanFile{
			File: validAudioContent,
		},
	}

	resp, err := suite.AdhanHandler.CreateAdhan(ctx, req)

	require.Error(suite.T(), err)
	st, ok := status.FromError(err)
	require.True(suite.T(), ok)
	assert.Equal(suite.T(), codes.InvalidArgument, st.Code())
	assert.Equal(suite.T(), "masjid ID is required", st.Message())
	assert.Nil(suite.T(), resp)
}

func (suite *GrpcHandlerTestSuite) TestCreateAdhan_NoFileContent() {
	ctx := context.Background()
	masjidID := uuid.New().String()
	req := &pb.CreateAdhanFileRequest{
		AdhanFile: &pb.AdhanFile{
			MasjidId: masjidID,
			File:     []byte{},
		},
	}

	resp, err := suite.AdhanHandler.CreateAdhan(ctx, req)

	require.Error(suite.T(), err)
	st, ok := status.FromError(err)
	require.True(suite.T(), ok)
	assert.Equal(suite.T(), codes.InvalidArgument, st.Code())
	assert.Equal(suite.T(), "adhan file content is required", st.Message())
	assert.Nil(suite.T(), resp)
}

func (suite *GrpcHandlerTestSuite) TestCreateAdhan_FileSizeExceedsLimit() {
	ctx := context.Background()
	masjidID := uuid.New().String()
	exceedingSize := bytes.Repeat([]byte{0x01}, int(maxAdhanFileSizeMB*1024*1024+1))

	req := &pb.CreateAdhanFileRequest{
		AdhanFile: &pb.AdhanFile{
			MasjidId: masjidID,
			File:     exceedingSize,
		},
	}

	resp, err := suite.AdhanHandler.CreateAdhan(ctx, req)

	require.Error(suite.T(), err)
	st, ok := status.FromError(err)
	require.True(suite.T(), ok)
	assert.Equal(suite.T(), codes.InvalidArgument, st.Code())
	assert.Contains(suite.T(), st.Message(), "adhan file size exceeds maximum allowed size")
	assert.Nil(suite.T(), resp)
}

func (suite *GrpcHandlerTestSuite) TestCreateAdhan_InvalidFileType() {
	ctx := context.Background()
	masjidID := uuid.New().String()
	invalidAudioContent := []byte{0x00, 0x00, 0x00}

	req := &pb.CreateAdhanFileRequest{
		AdhanFile: &pb.AdhanFile{
			MasjidId: masjidID,
			File:     invalidAudioContent,
		},
	}

	resp, err := suite.AdhanHandler.CreateAdhan(ctx, req)

	require.Error(suite.T(), err)
	st, ok := status.FromError(err)
	require.True(suite.T(), ok)
	assert.Equal(suite.T(), codes.InvalidArgument, st.Code())
	assert.Equal(suite.T(), "invalid adhan file type. Only MP3 and WAV are supported for now.", st.Message())
	assert.Nil(suite.T(), resp)
}

func (suite *GrpcHandlerTestSuite) TestUpdateAdhan_Success() {
	ctx := context.Background()
	masjidID := uuid.New().String()
	initialAudioContent := createValidAudioContent()

	initialAdhan := &entity.Adhan{
		ID:        uuid.New(),
		MasjidId:  masjidID,
		File:      initialAudioContent,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	err := suite.DB.Create(&initialAdhan).Error
	require.NoError(suite.T(), err, "Failed to create initial adhan file")

	req := &pb.UpdateAdhanFileRequest{
		Id: initialAdhan.ID.String(),
		AdhanFile: &pb.AdhanFile{
			MasjidId: masjidID,
			File:     initialAudioContent,
		},
	}

	resp, err := suite.AdhanHandler.UpdateAdhan(ctx, req)

	require.NoError(suite.T(), err)
	require.NotNil(suite.T(), resp)
	assert.Equal(suite.T(), codes.OK.String(), resp.GetCode())
	assert.Equal(suite.T(), "success", resp.GetStatus())
	assert.Equal(suite.T(), "adhan file updated successfully", resp.Message)
	updatedAdhanResp := resp.GetData().(*pb.StandardAdhanResponse_AdhanFile).AdhanFile
	assert.Equal(suite.T(), initialAdhan.ID.String(), updatedAdhanResp.GetId())
	assert.Equal(suite.T(), masjidID, updatedAdhanResp.GetMasjidId())
	assert.Equal(suite.T(), initialAudioContent, updatedAdhanResp.GetFile())

	var updatedAdhanEntity entity.Adhan
	err = suite.DB.Where("id = ?", initialAdhan.ID).First(&updatedAdhanEntity).Error
	require.NoError(suite.T(), err)
	assert.Equal(suite.T(), masjidID, updatedAdhanEntity.MasjidId)
	assert.Equal(suite.T(), initialAudioContent, updatedAdhanEntity.File)

	err = suite.DB.Delete(&updatedAdhanEntity).Error
	require.NoError(suite.T(), err)
	//suite.T().Logf("Error Message: %s", err)
}

func (suite *GrpcHandlerTestSuite) TestUpdateAdhan_NoAdhanFileData() {
	ctx := context.Background()
	adhanID := uuid.New().String()
	req := &pb.UpdateAdhanFileRequest{
		Id: adhanID,
	}

	resp, err := suite.AdhanHandler.UpdateAdhan(ctx, req)

	require.Error(suite.T(), err)
	st, ok := status.FromError(err)
	require.True(suite.T(), ok)
	assert.Equal(suite.T(), codes.InvalidArgument, st.Code())
	assert.Equal(suite.T(), "adhan file data is required for update", st.Message())
	assert.Nil(suite.T(), resp)
}

func (suite *GrpcHandlerTestSuite) TestUpdateAdhan_NoAdhanID() {
	ctx := context.Background()
	validAudioContent := createValidAudioContent()
	req := &pb.UpdateAdhanFileRequest{
		AdhanFile: &pb.AdhanFile{
			MasjidId: uuid.New().String(),
			File:     validAudioContent,
		},
	}

	resp, err := suite.AdhanHandler.UpdateAdhan(ctx, req)

	require.Error(suite.T(), err)
	st, ok := status.FromError(err)
	require.True(suite.T(), ok)
	assert.Equal(suite.T(), codes.InvalidArgument, st.Code())
	assert.Equal(suite.T(), "adhan file ID is required for update", st.Message())
	assert.Nil(suite.T(), resp)
}

func (suite *GrpcHandlerTestSuite) TestGetAdhanById_Success() {
	ctx := context.Background()
	masjidID := uuid.New().String()
	validAudioContent := createValidAudioContent()

	existingAdhan := &entity.Adhan{
		ID:        uuid.New(),
		MasjidId:  masjidID,
		File:      validAudioContent,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	err := suite.DB.Create(&existingAdhan).Error
	require.NoError(suite.T(), err, "Failed to create existing adhan file")

	req := &pb.GetAdhanFileRequest{
		Id: existingAdhan.ID.String(),
	}

	resp, err := suite.AdhanHandler.GetAdhanById(ctx, req)

	require.NoError(suite.T(), err)
	require.NotNil(suite.T(), resp)
	assert.Equal(suite.T(), codes.OK.String(), resp.GetCode())
	assert.Equal(suite.T(), "success", resp.GetStatus())
	assert.Equal(suite.T(), "adhan file retrieved successfully", resp.Message)
	retrievedAdhan := resp.GetData().(*pb.StandardAdhanResponse_AdhanFile).AdhanFile
	assert.Equal(suite.T(), existingAdhan.ID.String(), retrievedAdhan.GetId())
	assert.Equal(suite.T(), masjidID, retrievedAdhan.GetMasjidId())
	assert.Equal(suite.T(), validAudioContent, retrievedAdhan.GetFile())

	err = suite.DB.Delete(&existingAdhan).Error
	require.NoError(suite.T(), err)
}

func (suite *GrpcHandlerTestSuite) TestGetAdhanById_AdhanNotFound() {
	ctx := context.Background()
	nonExistentAdhanID := uuid.New().String()
	req := &pb.GetAdhanFileRequest{
		Id: nonExistentAdhanID,
	}

	resp, err := suite.AdhanHandler.GetAdhanById(ctx, req)

	require.Error(suite.T(), err)
	st, ok := status.FromError(err)
	require.True(suite.T(), ok)
	assert.Equal(suite.T(), codes.NotFound, st.Code())
	assert.Contains(suite.T(), st.Message(), "adhan file with ID")
	assert.Contains(suite.T(), st.Message(), "not found")
	assert.Nil(suite.T(), resp)
	//suite.T().Logf("Error Message: %s", err)
}

func (suite *GrpcHandlerTestSuite) TestDeleteAdhan_Success() {
	ctx := context.Background()
	masjidID := uuid.New().String()
	validAudioContent := createValidAudioContent()

	existingAdhan := &entity.Adhan{
		ID:        uuid.New(),
		MasjidId:  masjidID,
		File:      validAudioContent,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	err := suite.DB.Create(&existingAdhan).Error
	require.NoError(suite.T(), err, "Failed to create existing adhan file")

	req := &pb.DeleteAdhanFileRequest{
		Id: existingAdhan.ID.String(),
	}

	resp, err := suite.AdhanHandler.DeleteAdhan(ctx, req)

	require.NoError(suite.T(), err)
	require.NotNil(suite.T(), resp)
	assert.Equal(suite.T(), codes.OK.String(), resp.GetCode())
	assert.Equal(suite.T(), "success", resp.GetStatus())
	assert.Equal(suite.T(), "adhan file deleted successfully", resp.Message)
	assert.NotNil(suite.T(), resp.GetData())
	_, ok := resp.GetData().(*pb.StandardAdhanResponse_DeleteAdhanFileResponse)
	assert.True(suite.T(), ok)

	var deletedAdhan entity.Adhan
	err = suite.DB.Where("id = ?", existingAdhan.ID).First(&deletedAdhan).Error
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), err.Error(), "record not found")
}

func (suite *GrpcHandlerTestSuite) TestDeleteAdhan_NoAdhanID() {
	ctx := context.Background()
	req := &pb.DeleteAdhanFileRequest{}

	resp, err := suite.AdhanHandler.DeleteAdhan(ctx, req)

	require.Error(suite.T(), err)
	st, ok := status.FromError(err)
	require.True(suite.T(), ok)
	assert.Equal(suite.T(), codes.InvalidArgument, st.Code())
	assert.Equal(suite.T(), "id is required", st.Message())
	assert.Nil(suite.T(), resp)
}
