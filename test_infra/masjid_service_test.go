package test_infra

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
	"google.golang.org/protobuf/testing/protocmp"

	pb "github.com/mnadev/limestone/proto"
)

func TestMasjidService(t *testing.T) {
	suite.Run(t, new(IntegrationTestSuite))
}

func (suite *IntegrationTestSuite) TestCreateMasjid_Success() {
	ctx := context.Background()
	out, err := suite.MasjidServiceClient.CreateMasjid(ctx, &pb.CreateMasjidRequest{
		Masjid: GetMasjidProto(),
	})

	AssertProtoEqual(suite.T(), *GetMasjidProto(), *out,
		pb.Masjid{}, protocmp.IgnoreFields(&pb.Masjid{}, "create_time", "update_time"))
	AssertMasjidTimestampsCurrent(suite.T(), out)
	suite.Nil(err)
}

func (suite *IntegrationTestSuite) TestUpdateMasjid_Success() {
	ctx := context.Background()
	out, err := suite.MasjidServiceClient.CreateMasjid(ctx, &pb.CreateMasjidRequest{
		Masjid: GetMasjidProto(),
	})

	AssertProtoEqual(suite.T(), *GetMasjidProto(), *out,
		pb.Masjid{}, protocmp.IgnoreFields(&pb.Masjid{}, "create_time", "update_time"))
	AssertMasjidTimestampsCurrent(suite.T(), out)
	suite.Nil(err)

	update := GetMasjidProto()
	update.Name = "Masjid 2"

	out, err = suite.MasjidServiceClient.UpdateMasjid(ctx, &pb.UpdateMasjidRequest{
		Masjid: update,
	})

	AssertProtoEqual(suite.T(), *update, *out,
		pb.Masjid{}, protocmp.IgnoreFields(&pb.Masjid{}, "create_time", "update_time"))
	AssertMasjidTimestampsCurrent(suite.T(), out)
	suite.Nil(err)
}

func (suite *IntegrationTestSuite) TestUpdateMasjid_NotFound() {
	ctx := context.Background()

	out, err := suite.MasjidServiceClient.UpdateMasjid(ctx, &pb.UpdateMasjidRequest{
		Masjid: GetMasjidProto(),
	})

	suite.Error(err)
	suite.Nil(out)
}

func (suite *IntegrationTestSuite) TestGetMasjid_Success() {
	ctx := context.Background()
	out, err := suite.MasjidServiceClient.CreateMasjid(ctx, &pb.CreateMasjidRequest{
		Masjid: GetMasjidProto(),
	})

	AssertProtoEqual(suite.T(), *GetMasjidProto(), *out,
		pb.Masjid{}, protocmp.IgnoreFields(&pb.Masjid{}, "create_time", "update_time"))
	AssertMasjidTimestampsCurrent(suite.T(), out)
	suite.Nil(err)

	out, err = suite.MasjidServiceClient.GetMasjid(ctx, &pb.GetMasjidRequest{
		MasjidId: DefaultId,
	})

	AssertProtoEqual(suite.T(), *GetMasjidProto(), *out,
		pb.Masjid{}, protocmp.IgnoreFields(&pb.Masjid{}, "create_time", "update_time"))
	AssertMasjidTimestampsCurrent(suite.T(), out)
	suite.Nil(err)
}

func (suite *IntegrationTestSuite) TestGetMasjid_NotFound() {
	ctx := context.Background()
	out, err := suite.MasjidServiceClient.GetMasjid(ctx, &pb.GetMasjidRequest{
		MasjidId: DefaultId,
	})

	suite.Error(err)
	suite.Nil(out)
}

func (suite *IntegrationTestSuite) TestDeleteMasjid_Success() {
	ctx := context.Background()
	out, err := suite.MasjidServiceClient.CreateMasjid(ctx, &pb.CreateMasjidRequest{
		Masjid: GetMasjidProto(),
	})

	AssertProtoEqual(suite.T(), *GetMasjidProto(), *out,
		pb.Masjid{}, protocmp.IgnoreFields(&pb.Masjid{}, "create_time", "update_time"))
	AssertMasjidTimestampsCurrent(suite.T(), out)
	suite.Nil(err)

	_, err = suite.MasjidServiceClient.DeleteMasjid(ctx, &pb.DeleteMasjidRequest{
		MasjidId: DefaultId,
	})

	suite.Nil(err)
}

func (suite *IntegrationTestSuite) TestDeleteMasjid_NotFound() {
	ctx := context.Background()
	_, err := suite.MasjidServiceClient.DeleteMasjid(ctx, &pb.DeleteMasjidRequest{
		MasjidId: DefaultId,
	})

	suite.Error(err)
}
