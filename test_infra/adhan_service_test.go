package test_infra

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
	"google.golang.org/protobuf/testing/protocmp"

	pb "github.com/mnadev/limestone/proto"
)

func TestAdhanService(t *testing.T) {
	suite.Run(t, new(IntegrationTestSuite))
}

func (suite *IntegrationTestSuite) TestCreateAdhanFile_Success() {
	ctx := context.Background()
	out, err := suite.AdhanServiceClient.CreateAdhanFile(ctx, &pb.CreateAdhanFileRequest{
		AdhanFile: GetAdhanFileProto(),
	})

	suite.Nil(err)
	AssertProtoEqual(suite.T(), *GetAdhanFileProto(), *out,
		pb.AdhanFile{}, protocmp.IgnoreFields(&pb.AdhanFile{}, "create_time", "update_time"))
	AssertAdhanFileTimestampsCurrent(suite.T(), out)
}

func (suite *IntegrationTestSuite) TestUpdateAdhanFile_Success() {
	ctx := context.Background()
	out, err := suite.AdhanServiceClient.CreateAdhanFile(ctx, &pb.CreateAdhanFileRequest{
		AdhanFile: GetAdhanFileProto(),
	})

	AssertProtoEqual(suite.T(), *GetAdhanFileProto(), *out,
		pb.AdhanFile{}, protocmp.IgnoreFields(&pb.AdhanFile{}, "create_time", "update_time"))
	AssertAdhanFileTimestampsCurrent(suite.T(), out)
	suite.Nil(err)

	update := GetAdhanFileProto()
	update.File = []byte("xyz")

	out, err = suite.AdhanServiceClient.UpdateAdhanFile(ctx, &pb.UpdateAdhanFileRequest{
		AdhanFile: update,
	})

	AssertProtoEqual(suite.T(), *update, *out,
		pb.AdhanFile{}, protocmp.IgnoreFields(&pb.AdhanFile{}, "create_time", "update_time"))
	AssertAdhanFileTimestampsCurrent(suite.T(), out)
	suite.Nil(err)
}

func (suite *IntegrationTestSuite) TestUpdateAdhanFile_NotFound() {
	ctx := context.Background()

	out, err := suite.AdhanServiceClient.UpdateAdhanFile(ctx, &pb.UpdateAdhanFileRequest{
		AdhanFile: GetAdhanFileProto(),
	})

	suite.Error(err)
	suite.Nil(out)
}

func (suite *IntegrationTestSuite) TestGetAdhanFile_Success() {
	ctx := context.Background()
	out, err := suite.AdhanServiceClient.CreateAdhanFile(ctx, &pb.CreateAdhanFileRequest{
		AdhanFile: GetAdhanFileProto(),
	})

	AssertProtoEqual(suite.T(), *GetAdhanFileProto(), *out,
		pb.AdhanFile{}, protocmp.IgnoreFields(&pb.AdhanFile{}, "create_time", "update_time"))
	AssertAdhanFileTimestampsCurrent(suite.T(), out)
	suite.Nil(err)

	out, err = suite.AdhanServiceClient.GetAdhanFile(ctx, &pb.GetAdhanFileRequest{
		MasjidId: DefaultId,
	})

	AssertProtoEqual(suite.T(), *GetAdhanFileProto(), *out,
		pb.AdhanFile{}, protocmp.IgnoreFields(&pb.AdhanFile{}, "create_time", "update_time"))
	AssertAdhanFileTimestampsCurrent(suite.T(), out)
	suite.Nil(err)
}

func (suite *IntegrationTestSuite) TestGetAdhanFile_NotFound() {
	ctx := context.Background()
	out, err := suite.AdhanServiceClient.GetAdhanFile(ctx, &pb.GetAdhanFileRequest{
		MasjidId: DefaultId,
	})

	suite.Error(err)
	suite.Nil(out)
}

func (suite *IntegrationTestSuite) TestDeleteAdhanFile_Success() {
	ctx := context.Background()
	out, err := suite.AdhanServiceClient.CreateAdhanFile(ctx, &pb.CreateAdhanFileRequest{
		AdhanFile: GetAdhanFileProto(),
	})

	AssertProtoEqual(suite.T(), *GetAdhanFileProto(), *out,
		pb.AdhanFile{}, protocmp.IgnoreFields(&pb.AdhanFile{}, "create_time", "update_time"))
	AssertAdhanFileTimestampsCurrent(suite.T(), out)
	suite.Nil(err)

	_, err = suite.AdhanServiceClient.DeleteAdhanFile(ctx, &pb.DeleteAdhanFileRequest{
		Id: DefaultId,
	})

	suite.Nil(err)
}

func (suite *IntegrationTestSuite) TestDeleteAdhanFile_NotFound() {
	ctx := context.Background()
	_, err := suite.AdhanServiceClient.DeleteAdhanFile(ctx, &pb.DeleteAdhanFileRequest{
		Id: DefaultId,
	})

	suite.Error(err)
}
