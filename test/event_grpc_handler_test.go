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
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

func (suite *GrpcHandlerTestSuite) TestCreateEvent_Success() {
	ctx := context.Background()
	masjidID := uuid.New().String()
	startTime := time.Now().Add(time.Hour)
	endTime := startTime.Add(2 * time.Hour)

	req := &pb.CreateEventRequest{
		Event: &pb.Event{
			MasjidId:          masjidID,
			Name:              "Test Event",
			Description:       "This is a test event",
			StartTime:         timestamppb.New(startTime),
			EndTime:           timestamppb.New(endTime),
			GenderRestriction: pb.Event_MALE_ONLY,
			IsPaid:            true,
			RequiresRsvp:      true,
			MaxParticipants:   100,
			LivestreamLink:    "http://example.com/live",
		},
	}

	resp, err := suite.EventHandler.CreateEvent(ctx, req)
	require.NoError(suite.T(), err)
	require.NotNil(suite.T(), resp)
	assert.Equal(suite.T(), codes.OK.String(), resp.GetCode())
	assert.Equal(suite.T(), "success", resp.GetStatus())
	assert.Equal(suite.T(), "event created successfully", resp.Message)
	createdEvent := resp.GetData().(*pb.StandardEventResponse_Event).Event
	assert.NotEmpty(suite.T(), createdEvent.GetId())
	assert.Equal(suite.T(), masjidID, createdEvent.GetMasjidId())
	assert.Equal(suite.T(), "Test Event", createdEvent.GetName())
	assert.Equal(suite.T(), pb.Event_MALE_ONLY, createdEvent.GetGenderRestriction())
	assert.True(suite.T(), createdEvent.GetIsPaid())
	assert.True(suite.T(), createdEvent.GetRequiresRsvp())
	assert.Equal(suite.T(), int32(100), createdEvent.GetMaxParticipants())
	assert.Equal(suite.T(), "http://example.com/live", createdEvent.GetLivestreamLink())
	assert.WithinDuration(suite.T(), startTime, createdEvent.GetStartTime().AsTime(), time.Second)
	assert.WithinDuration(suite.T(), endTime, createdEvent.GetEndTime().AsTime(), time.Second)

	var dbEvent entity.Event
	err = suite.DB.Where("id = ?", createdEvent.GetId()).First(&dbEvent).Error
	require.NoError(suite.T(), err)
	assert.Equal(suite.T(), "Test Event", dbEvent.Name)

	err = suite.DB.Delete(&dbEvent).Error
	require.NoError(suite.T(), err)
}

func (suite *GrpcHandlerTestSuite) TestGetEvent_Success() {
	ctx := context.Background()
	masjidID := uuid.New().String()
	startTime := time.Now().Add(time.Hour)
	endTime := startTime.Add(2 * time.Hour)

	existingEvent := &entity.Event{
		ID:                uuid.New(),
		MasjidId:          masjidID,
		Name:              "Existing Event",
		StartTime:         startTime,
		EndTime:           endTime,
		GenderRestriction: entity.GenderRestriction(pb.Event_FEMALE_ONLY),
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
	}
	err := suite.DB.Create(&existingEvent).Error
	require.NoError(suite.T(), err)

	req := &pb.GetEventRequest{
		Id: existingEvent.ID.String(),
	}

	resp, err := suite.EventHandler.GetEvent(ctx, req)
	require.NoError(suite.T(), err)
	require.NotNil(suite.T(), resp)
	assert.Equal(suite.T(), codes.OK.String(), resp.GetCode())
	assert.Equal(suite.T(), "success", resp.GetStatus())
	assert.Equal(suite.T(), "event retrieved successfully", resp.Message)
	retrievedEvent := resp.GetData().(*pb.StandardEventResponse_Event).Event
	assert.Equal(suite.T(), existingEvent.ID.String(), retrievedEvent.GetId())
	assert.Equal(suite.T(), "Existing Event", retrievedEvent.GetName())
	assert.Equal(suite.T(), pb.Event_FEMALE_ONLY, retrievedEvent.GetGenderRestriction())
	assert.WithinDuration(suite.T(), startTime, retrievedEvent.GetStartTime().AsTime(), time.Second)
	assert.WithinDuration(suite.T(), endTime, retrievedEvent.GetEndTime().AsTime(), time.Second)

	err = suite.DB.Delete(&existingEvent).Error
	require.NoError(suite.T(), err)
}

func (suite *GrpcHandlerTestSuite) TestGetEvent_NotFound() {
	ctx := context.Background()
	nonExistingID := uuid.New().String()
	req := &pb.GetEventRequest{
		Id: nonExistingID,
	}

	resp, err := suite.EventHandler.GetEvent(ctx, req)
	require.Error(suite.T(), err)
	st, ok := status.FromError(err)
	require.True(suite.T(), ok)
	assert.Equal(suite.T(), codes.NotFound, st.Code())
	assert.Equal(suite.T(), "event not found", st.Message())
	assert.Nil(suite.T(), resp)
}

func (suite *GrpcHandlerTestSuite) TestDeleteEvent_Success() {
	ctx := context.Background()
	masjidID := uuid.New().String()
	startTime := time.Now().Add(time.Hour)
	endTime := startTime.Add(2 * time.Hour)

	existingEvent := &entity.Event{
		MasjidId:  masjidID,
		Name:      "Event to Delete",
		StartTime: startTime,
		EndTime:   endTime,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	err := suite.DB.Create(&existingEvent).Error
	require.NoError(suite.T(), err)

	req := &pb.DeleteEventRequest{
		Id: existingEvent.ID.String(),
	}

	resp, err := suite.EventHandler.DeleteEvent(ctx, req)
	require.NoError(suite.T(), err)
	require.NotNil(suite.T(), resp)
	assert.Equal(suite.T(), codes.OK.String(), resp.GetCode())
	assert.Equal(suite.T(), "success", resp.GetStatus())
	assert.Equal(suite.T(), "event deleted successfully", resp.Message)
	assert.NotNil(suite.T(), resp.GetData())
	_, ok := resp.GetData().(*pb.StandardEventResponse_DeleteEventResponse)
	assert.True(suite.T(), ok)

	var deletedEvent entity.Event
	err = suite.DB.Where("id = ?", existingEvent.ID).First(&deletedEvent).Error
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), err.Error(), "record not found")
}

func (suite *GrpcHandlerTestSuite) TestUpdateEvent_Success() {
	ctx := context.Background()
	masjidID := uuid.New().String()
	startTime := time.Now().Add(time.Hour)
	endTime := startTime.Add(2 * time.Hour)
	updatedEndTime := endTime.Add(time.Hour)

	existingEvent := &entity.Event{
		ID:                uuid.New(),
		MasjidId:          masjidID,
		Name:              "Old Event Name",
		Description:       "Old Description",
		StartTime:         startTime,
		EndTime:           endTime,
		GenderRestriction: entity.GenderRestriction(pb.Event_NO_RESTRICTION),
		IsPaid:            false,
		RequiresRsvp:      false,
		MaxParticipants:   50,
		LivestreamLink:    "old.example.com/live",
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
	}
	err := suite.DB.Create(&existingEvent).Error
	require.NoError(suite.T(), err)

	req := &pb.UpdateEventRequest{
		Event: &pb.Event{
			Id:                existingEvent.ID.String(),
			Name:              "New Event Name",
			Description:       "New Description",
			StartTime:         timestamppb.New(startTime),
			EndTime:           timestamppb.New(updatedEndTime),
			GenderRestriction: pb.Event_FEMALE_ONLY,
			IsPaid:            true,
			RequiresRsvp:      true,
			MaxParticipants:   100,
			LivestreamLink:    "new.example.com/live",
		},
	}

	resp, err := suite.EventHandler.UpdateEvent(ctx, req)
	require.NoError(suite.T(), err)
	require.NotNil(suite.T(), resp)
	assert.Equal(suite.T(), codes.OK.String(), resp.GetCode())
	assert.Equal(suite.T(), "success", resp.GetStatus())
	assert.Equal(suite.T(), "event updated successfully", resp.Message)
	updatedEventResp := resp.GetData().(*pb.StandardEventResponse_Event).Event
	assert.Equal(suite.T(), existingEvent.ID.String(), updatedEventResp.GetId())
	assert.Equal(suite.T(), "New Event Name", updatedEventResp.GetName())
	assert.Equal(suite.T(), "New Description", updatedEventResp.GetDescription())
	assert.WithinDuration(suite.T(), startTime, updatedEventResp.GetStartTime().AsTime(), time.Second)
	assert.WithinDuration(suite.T(), updatedEndTime, updatedEventResp.GetEndTime().AsTime(), time.Second)
	assert.Equal(suite.T(), pb.Event_FEMALE_ONLY, updatedEventResp.GetGenderRestriction())
	assert.True(suite.T(), updatedEventResp.GetIsPaid())
	assert.True(suite.T(), updatedEventResp.GetRequiresRsvp())
	assert.Equal(suite.T(), int32(100), updatedEventResp.GetMaxParticipants())
	assert.Equal(suite.T(), "new.example.com/live", updatedEventResp.GetLivestreamLink())

	var dbEvent entity.Event
	err = suite.DB.Where("id = ?", existingEvent.ID).First(&dbEvent).Error
	require.NoError(suite.T(), err)
	assert.Equal(suite.T(), "New Event Name", dbEvent.Name)
	assert.Equal(suite.T(), "New Description", dbEvent.Description)
	assert.WithinDuration(suite.T(), startTime, dbEvent.StartTime, time.Second)
	assert.WithinDuration(suite.T(), updatedEndTime, dbEvent.EndTime, time.Second)
	assert.Equal(suite.T(), entity.GenderRestriction(pb.Event_FEMALE_ONLY), dbEvent.GenderRestriction)
	assert.True(suite.T(), dbEvent.IsPaid)
	assert.True(suite.T(), dbEvent.RequiresRsvp)
	assert.Equal(suite.T(), int32(100), dbEvent.MaxParticipants)
	assert.Equal(suite.T(), "new.example.com/live", dbEvent.LivestreamLink)

	err = suite.DB.Delete(&dbEvent).Error
	require.NoError(suite.T(), err)
}

func (suite *GrpcHandlerTestSuite) TestUpdateEvent_NoEventID() {
	ctx := context.Background()
	startTime := time.Now().Add(time.Hour)
	endTime := startTime.Add(2 * time.Hour)

	req := &pb.UpdateEventRequest{
		Event: &pb.Event{
			Name:      "Event Without ID",
			StartTime: timestamppb.New(startTime),
			EndTime:   timestamppb.New(endTime),
		},
	}

	resp, err := suite.EventHandler.UpdateEvent(ctx, req)
	require.Error(suite.T(), err)
	st, ok := status.FromError(err)
	require.True(suite.T(), ok)
	assert.Equal(suite.T(), codes.InvalidArgument, st.Code())
	assert.Equal(suite.T(), "event ID is required for update", st.Message())
	assert.Nil(suite.T(), resp)
}
