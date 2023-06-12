package test_infra

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/suite"
	"google.golang.org/protobuf/testing/protocmp"

	apb "github.com/mnadev/limestone/adhan_service/proto"
)

func TestSuite(t *testing.T) {
	suite.Run(t, new(IntegrationTestSuite))
}

func (suite *IntegrationTestSuite) TestCreateAdhanFile_Success() {
	ctx := context.Background()
	out, err := suite.AdhanServiceClient.CreateAdhanFile(ctx, &apb.CreateAdhanFileRequest{
		AdhanFile: GetAdhanFileProto(),
	})

	suite.Nil(err)
	AssertProtoEqual(suite.T(), *GetAdhanFileProto(), *out,
		apb.AdhanFile{}, protocmp.IgnoreFields(&apb.AdhanFile{}, "create_time", "update_time"))
	AssertAdhanFileTimestampsCurrent(suite.T(), out)
}

func (suite *IntegrationTestSuite) TestUpdateAdhanFile_Success() {
	ctx := context.Background()
	out, err := suite.AdhanServiceClient.CreateAdhanFile(ctx, &apb.CreateAdhanFileRequest{
		AdhanFile: GetAdhanFileProto(),
	})

	AssertProtoEqual(suite.T(), *GetAdhanFileProto(), *out,
		apb.AdhanFile{}, protocmp.IgnoreFields(&apb.AdhanFile{}, "create_time", "update_time"))
	AssertAdhanFileTimestampsCurrent(suite.T(), out)
	suite.Nil(err)

	update := GetAdhanFileProto()
	update.File = []byte("xyz")

	out, err = suite.AdhanServiceClient.UpdateAdhanFile(ctx, &apb.UpdateAdhanFileRequest{
		AdhanFile: update,
	})

	AssertProtoEqual(suite.T(), *update, *out,
		apb.AdhanFile{}, protocmp.IgnoreFields(&apb.AdhanFile{}, "create_time", "update_time"))
	AssertAdhanFileTimestampsCurrent(suite.T(), out)
	suite.Nil(err)
}

func (suite *IntegrationTestSuite) TestUpdateAdhanFile_NotFound() {
	ctx := context.Background()

	out, err := suite.AdhanServiceClient.UpdateAdhanFile(ctx, &apb.UpdateAdhanFileRequest{
		AdhanFile: GetAdhanFileProto(),
	})

	suite.Error(err)
	suite.Nil(out)
}

func (suite *IntegrationTestSuite) TestGetAdhanFile_Success() {
	ctx := context.Background()
	out, err := suite.AdhanServiceClient.CreateAdhanFile(ctx, &apb.CreateAdhanFileRequest{
		AdhanFile: GetAdhanFileProto(),
	})

	AssertProtoEqual(suite.T(), *GetAdhanFileProto(), *out,
		apb.AdhanFile{}, protocmp.IgnoreFields(&apb.AdhanFile{}, "create_time", "update_time"))
	AssertAdhanFileTimestampsCurrent(suite.T(), out)
	suite.Nil(err)

	out, err = suite.AdhanServiceClient.GetAdhanFile(ctx, &apb.GetAdhanFileRequest{
		MasjidId: DefaultId,
	})

	AssertProtoEqual(suite.T(), *GetAdhanFileProto(), *out,
		apb.AdhanFile{}, protocmp.IgnoreFields(&apb.AdhanFile{}, "create_time", "update_time"))
	AssertAdhanFileTimestampsCurrent(suite.T(), out)
	suite.Nil(err)
}

func (suite *IntegrationTestSuite) TestGetAdhanFile_NotFound() {
	ctx := context.Background()
	out, err := suite.AdhanServiceClient.GetAdhanFile(ctx, &apb.GetAdhanFileRequest{
		MasjidId: DefaultId,
	})

	suite.Error(err)
	suite.Nil(out)
}

func (suite *IntegrationTestSuite) TestDeleteAdhanFile_Success() {
	ctx := context.Background()
	out, err := suite.AdhanServiceClient.CreateAdhanFile(ctx, &apb.CreateAdhanFileRequest{
		AdhanFile: GetAdhanFileProto(),
	})

	AssertProtoEqual(suite.T(), *GetAdhanFileProto(), *out,
		apb.AdhanFile{}, protocmp.IgnoreFields(&apb.AdhanFile{}, "create_time", "update_time"))
	AssertAdhanFileTimestampsCurrent(suite.T(), out)
	suite.Nil(err)

	_, err = suite.AdhanServiceClient.DeleteAdhanFile(ctx, &apb.DeleteAdhanFileRequest{
		Id: DefaultId,
	})

	suite.Nil(err)
}

func (suite *IntegrationTestSuite) TestDeleteAdhanFile_NotFound() {
	ctx := context.Background()
	_, err := suite.AdhanServiceClient.DeleteAdhanFile(ctx, &apb.DeleteAdhanFileRequest{
		Id: DefaultId,
	})

	suite.Error(err)
}
