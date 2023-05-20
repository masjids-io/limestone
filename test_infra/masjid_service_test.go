package test_infra

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
	"google.golang.org/protobuf/testing/protocmp"

	mpb "github.com/mnadev/limestone/masjid_service/proto"
)

func TestSuite(t *testing.T) {
	suite.Run(t, new(IntegrationTestSuite))
}

func (suite *IntegrationTestSuite) TestCreateMasjid_Success() {
	ctx := context.Background()
	out, err := suite.MasjidServiceClient.CreateMasjid(ctx, &mpb.CreateMasjidRequest{
		Masjid: GetMasjidProto(),
	})

	AssertProtoEqual(suite.T(), *GetMasjidProto(), *out,
		mpb.Masjid{}, protocmp.IgnoreFields(&mpb.Masjid{}, "create_time", "update_time"))
	AssertMasjidTimestampsCurrent(suite.T(), out)
	suite.Nil(err)
}

func (suite *IntegrationTestSuite) TestUpdateMasjid_Success() {
	ctx := context.Background()
	out, err := suite.MasjidServiceClient.CreateMasjid(ctx, &mpb.CreateMasjidRequest{
		Masjid: GetMasjidProto(),
	})

	AssertProtoEqual(suite.T(), *GetMasjidProto(), *out,
		mpb.Masjid{}, protocmp.IgnoreFields(&mpb.Masjid{}, "create_time", "update_time"))
	AssertMasjidTimestampsCurrent(suite.T(), out)
	suite.Nil(err)

	update := GetMasjidProto()
	update.Name = "Masjid 2"

	out, err = suite.MasjidServiceClient.UpdateMasjid(ctx, &mpb.UpdateMasjidRequest{
		Masjid: update,
	})

	AssertProtoEqual(suite.T(), *update, *out,
		mpb.Masjid{}, protocmp.IgnoreFields(&mpb.Masjid{}, "create_time", "update_time"))
	AssertMasjidTimestampsCurrent(suite.T(), out)
	suite.Nil(err)
}

func (suite *IntegrationTestSuite) TestUpdateMasjid_NotFound() {
	ctx := context.Background()

	out, err := suite.MasjidServiceClient.UpdateMasjid(ctx, &mpb.UpdateMasjidRequest{
		Masjid: GetMasjidProto(),
	})

	suite.Error(err)
	suite.Nil(out)
}

func (suite *IntegrationTestSuite) TestGetMasjid_Success() {
	ctx := context.Background()
	out, err := suite.MasjidServiceClient.CreateMasjid(ctx, &mpb.CreateMasjidRequest{
		Masjid: GetMasjidProto(),
	})

	AssertProtoEqual(suite.T(), *GetMasjidProto(), *out,
		mpb.Masjid{}, protocmp.IgnoreFields(&mpb.Masjid{}, "create_time", "update_time"))
	AssertMasjidTimestampsCurrent(suite.T(), out)
	suite.Nil(err)

	out, err = suite.MasjidServiceClient.GetMasjid(ctx, &mpb.GetMasjidRequest{
		MasjidId: DefaultId,
	})

	AssertProtoEqual(suite.T(), *GetMasjidProto(), *out,
		mpb.Masjid{}, protocmp.IgnoreFields(&mpb.Masjid{}, "create_time", "update_time"))
	AssertMasjidTimestampsCurrent(suite.T(), out)
	suite.Nil(err)
}

func (suite *IntegrationTestSuite) TestGetMasjid_NotFound() {
	ctx := context.Background()
	out, err := suite.MasjidServiceClient.GetMasjid(ctx, &mpb.GetMasjidRequest{
		MasjidId: DefaultId,
	})

	suite.Error(err)
	suite.Nil(out)
}

func (suite *IntegrationTestSuite) TestDeleteMasjid_Success() {
	ctx := context.Background()
	out, err := suite.MasjidServiceClient.CreateMasjid(ctx, &mpb.CreateMasjidRequest{
		Masjid: GetMasjidProto(),
	})

	AssertProtoEqual(suite.T(), *GetMasjidProto(), *out,
		mpb.Masjid{}, protocmp.IgnoreFields(&mpb.Masjid{}, "create_time", "update_time"))
	AssertMasjidTimestampsCurrent(suite.T(), out)
	suite.Nil(err)

	_, err = suite.MasjidServiceClient.DeleteMasjid(ctx, &mpb.DeleteMasjidRequest{
		MasjidId: DefaultId,
	})

	suite.Nil(err)
}

func (suite *IntegrationTestSuite) TestDeleteMasjid_NotFound() {
	ctx := context.Background()
	_, err := suite.MasjidServiceClient.DeleteMasjid(ctx, &mpb.DeleteMasjidRequest{
		MasjidId: DefaultId,
	})

	suite.Error(err)
}
