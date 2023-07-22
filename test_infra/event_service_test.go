package test_infra

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/testing/protocmp"

	pb "github.com/mnadev/limestone/proto"
)

func TestEvenService(t *testing.T) {
	suite.Run(t, new(IntegrationTestSuite))
}

func (suite *IntegrationTestSuite) TestCreateEvent_Success() {
	ctx := context.Background()
	out, err := suite.EventServiceClient.CreateEvent(ctx, &pb.CreateEventRequest{
		Event: GetEventProto(),
	})

	suite.Assert().Equal(status.Code(err), codes.OK)
	AssertProtoEqual(suite.T(), *GetEventProto(), *out,
		pb.Event{}, protocmp.IgnoreFields(&pb.Event{}, "create_time", "update_time"))
	AssertEventTimestampsCurrent(suite.T(), out)
}

func (suite *IntegrationTestSuite) TestUpdateEvent_Success() {
	ctx := context.Background()
	out, err := suite.EventServiceClient.CreateEvent(ctx, &pb.CreateEventRequest{
		Event: GetEventProto(),
	})

	suite.Assert().Equal(status.Code(err), codes.OK)
	AssertProtoEqual(suite.T(), *GetEventProto(), *out,
		pb.Event{}, protocmp.IgnoreFields(&pb.Event{}, "create_time", "update_time"))
	AssertEventTimestampsCurrent(suite.T(), out)

	update := GetEventProto()
	update.Name = "Event 2"

	out, err = suite.EventServiceClient.UpdateEvent(ctx, &pb.UpdateEventRequest{
		Event: update,
	})

	suite.Assert().Equal(status.Code(err), codes.OK)
	AssertProtoEqual(suite.T(), *update, *out,
		pb.Event{}, protocmp.IgnoreFields(&pb.Event{}, "create_time", "update_time"))
	AssertEventTimestampsCurrent(suite.T(), out)
}

func (suite *IntegrationTestSuite) TestUpdateEvent_NotFound() {
	ctx := context.Background()

	out, err := suite.EventServiceClient.UpdateEvent(ctx, &pb.UpdateEventRequest{
		Event: GetEventProto(),
	})

	suite.Assert().Equal(status.Code(err), codes.NotFound)
	suite.Nil(out)
}

func (suite *IntegrationTestSuite) TestGetEvent_Success() {
	ctx := context.Background()
	out, err := suite.EventServiceClient.CreateEvent(ctx, &pb.CreateEventRequest{
		Event: GetEventProto(),
	})

	suite.Assert().Equal(status.Code(err), codes.OK)
	AssertProtoEqual(suite.T(), *GetEventProto(), *out,
		pb.Event{}, protocmp.IgnoreFields(&pb.Event{}, "create_time", "update_time"))
	AssertEventTimestampsCurrent(suite.T(), out)

	out, err = suite.EventServiceClient.GetEvent(ctx, &pb.GetEventRequest{
		EventId: DefaultId,
	})

	suite.Assert().Equal(status.Code(err), codes.OK)
	AssertProtoEqual(suite.T(), *GetEventProto(), *out,
		pb.Event{}, protocmp.IgnoreFields(&pb.Event{}, "create_time", "update_time"))
	AssertEventTimestampsCurrent(suite.T(), out)
}

func (suite *IntegrationTestSuite) TestGetEvent_NotFound() {
	ctx := context.Background()
	out, err := suite.EventServiceClient.GetEvent(ctx, &pb.GetEventRequest{
		EventId: DefaultId,
	})

	suite.Assert().Equal(status.Code(err), codes.NotFound)
	suite.Nil(out)
}

func (suite *IntegrationTestSuite) TestDeleteEvent_Success() {
	ctx := context.Background()
	out, err := suite.EventServiceClient.CreateEvent(ctx, &pb.CreateEventRequest{
		Event: GetEventProto(),
	})

	suite.Assert().Equal(status.Code(err), codes.OK)
	AssertProtoEqual(suite.T(), *GetEventProto(), *out,
		pb.Event{}, protocmp.IgnoreFields(&pb.Event{}, "create_time", "update_time"))
	AssertEventTimestampsCurrent(suite.T(), out)

	_, err = suite.EventServiceClient.DeleteEvent(ctx, &pb.DeleteEventRequest{
		EventId: DefaultId,
	})

	suite.Assert().Equal(status.Code(err), codes.OK)
}

func (suite *IntegrationTestSuite) TestDeleteEvent_NotFound() {
	ctx := context.Background()
	_, err := suite.EventServiceClient.DeleteEvent(ctx, &pb.DeleteEventRequest{
		EventId: DefaultId,
	})

	suite.Assert().Equal(status.Code(err), codes.NotFound)
}
