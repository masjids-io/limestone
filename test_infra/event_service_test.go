package test_infra

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
	"google.golang.org/protobuf/testing/protocmp"

	epb "github.com/mnadev/limestone/event_service/proto"
)

func TestSuite(t *testing.T) {
	suite.Run(t, new(IntegrationTestSuite))
}

func (suite *IntegrationTestSuite) TestCreateEvent_Success() {
	ctx := context.Background()
	out, err := suite.EventServiceClient.CreateEvent(ctx, &epb.CreateEventRequest{
		Event: GetEventProto(),
	})

	AssertProtoEqual(suite.T(), *GetEventProto(), *out,
		epb.Event{}, protocmp.IgnoreFields(&epb.Event{}, "create_time", "update_time"))
	AssertEventTimestampsCurrent(suite.T(), out)
	suite.Nil(err)
}

func (suite *IntegrationTestSuite) TestUpdateEvent_Success() {
	ctx := context.Background()
	out, err := suite.EventServiceClient.CreateEvent(ctx, &epb.CreateEventRequest{
		Event: GetEventProto(),
	})

	AssertProtoEqual(suite.T(), *GetEventProto(), *out,
		epb.Event{}, protocmp.IgnoreFields(&epb.Event{}, "create_time", "update_time"))
	AssertEventTimestampsCurrent(suite.T(), out)
	suite.Nil(err)

	update := GetEventProto()
	update.Name = "Event 2"

	out, err = suite.EventServiceClient.UpdateEvent(ctx, &epb.UpdateEventRequest{
		Event: update,
	})

	AssertProtoEqual(suite.T(), *update, *out,
		epb.Event{}, protocmp.IgnoreFields(&epb.Event{}, "create_time", "update_time"))
	AssertEventTimestampsCurrent(suite.T(), out)
	suite.Nil(err)
}

func (suite *IntegrationTestSuite) TestUpdateEvent_NotFound() {
	ctx := context.Background()

	out, err := suite.EventServiceClient.UpdateEvent(ctx, &epb.UpdateEventRequest{
		Event: GetEventProto(),
	})

	suite.Error(err)
	suite.Nil(out)
}

func (suite *IntegrationTestSuite) TestGetEvent_Success() {
	ctx := context.Background()
	out, err := suite.EventServiceClient.CreateEvent(ctx, &epb.CreateEventRequest{
		Event: GetEventProto(),
	})

	AssertProtoEqual(suite.T(), *GetEventProto(), *out,
		epb.Event{}, protocmp.IgnoreFields(&epb.Event{}, "create_time", "update_time"))
	AssertEventTimestampsCurrent(suite.T(), out)
	suite.Nil(err)

	out, err = suite.EventServiceClient.GetEvent(ctx, &epb.GetEventRequest{
		EventId: DefaultId,
	})

	AssertProtoEqual(suite.T(), *GetEventProto(), *out,
		epb.Event{}, protocmp.IgnoreFields(&epb.Event{}, "create_time", "update_time"))
	AssertEventTimestampsCurrent(suite.T(), out)
	suite.Nil(err)
}

func (suite *IntegrationTestSuite) TestGetEvent_NotFound() {
	ctx := context.Background()
	out, err := suite.EventServiceClient.GetEvent(ctx, &epb.GetEventRequest{
		EventId: DefaultId,
	})

	suite.Error(err)
	suite.Nil(out)
}

func (suite *IntegrationTestSuite) TestDeleteEvent_Success() {
	ctx := context.Background()
	out, err := suite.EventServiceClient.CreateEvent(ctx, &epb.CreateEventRequest{
		Event: GetEventProto(),
	})

	AssertProtoEqual(suite.T(), *GetEventProto(), *out,
		epb.Event{}, protocmp.IgnoreFields(&epb.Event{}, "create_time", "update_time"))
	AssertEventTimestampsCurrent(suite.T(), out)
	suite.Nil(err)

	_, err = suite.EventServiceClient.DeleteEvent(ctx, &epb.DeleteEventRequest{
		EventId: DefaultId,
	})

	suite.Nil(err)
}

func (suite *IntegrationTestSuite) TestDeleteEvent_NotFound() {
	ctx := context.Background()
	_, err := suite.EventServiceClient.DeleteEvent(ctx, &epb.DeleteEventRequest{
		EventId: DefaultId,
	})

	suite.Error(err)
}
