package handler

import (
	"context"
	"errors"
	"github.com/google/uuid"
	pb "github.com/mnadev/limestone/gen/go"
	"github.com/mnadev/limestone/internal/application/domain/entity"
	services "github.com/mnadev/limestone/internal/application/services"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"
	"time"
)

type EventGrpcHandler struct {
	pb.UnimplementedEventServiceServer
	Svc *services.EventService
}

func NewEventGrpcHandler(svc *services.EventService) *EventGrpcHandler {
	return &EventGrpcHandler{Svc: svc}
}

func (h *EventGrpcHandler) CreateEvent(ctx context.Context, req *pb.CreateEventRequest) (*pb.Event, error) {
	event := req.GetEvent()

	eventEntity := &entity.Event{
		ID:                uuid.New(),
		MasjidId:          event.GetMasjidId(),
		Name:              event.GetName(),
		Description:       event.GetDescription(),
		StartTime:         event.StartTime.AsTime(),
		EndTime:           event.EndTime.AsTime(),
		GenderRestriction: entity.GenderRestriction(event.GetGenderRestriction()),
		IsPaid:            event.GetIsPaid(),
		RequiresRsvp:      event.GetRequiresRsvp(),
		MaxParticipants:   event.GetMaxParticipants(),
		LivestreamLink:    event.GetLivestreamLink(),
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
	}

	createdEvent, err := h.Svc.Create(ctx, eventEntity)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create event: %v", err)
	}

	return convertEventEntityToProto(createdEvent), nil
}

func (h *EventGrpcHandler) UpdateEvent(ctx context.Context, req *pb.UpdateEventRequest) (*pb.Event, error) {
	event := req.GetEvent()
	if event == nil {
		return nil, status.Errorf(codes.InvalidArgument, "event data is required")
	}

	eventIDStr := event.GetId()
	if eventIDStr == "" {
		return nil, status.Errorf(codes.InvalidArgument, "event ID is required for update")
	}

	eventID, err := uuid.Parse(eventIDStr)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid event ID format: %v", err)
	}

	eventEntity := &entity.Event{
		ID:        eventID,
		UpdatedAt: time.Now(),
	}

	if event.GetMasjidId() != "" {
		eventEntity.MasjidId = event.GetMasjidId()
	}
	if event.GetName() != "" {
		eventEntity.Name = event.GetName()
	}
	if event.GetDescription() != "" {
		eventEntity.Description = event.GetDescription()
	}
	if event.GetStartTime() != nil {
		eventEntity.StartTime = event.GetStartTime().AsTime()
	}
	if event.GetEndTime() != nil {
		eventEntity.EndTime = event.GetEndTime().AsTime()
	}
	if event.GetGenderRestriction() != pb.Event_NO_RESTRICTION {
		eventEntity.GenderRestriction = entity.GenderRestriction(event.GetGenderRestriction())
	}
	if event.GetIsPaid() != false {
		eventEntity.IsPaid = event.GetIsPaid()
	}
	if event.GetRequiresRsvp() != false {
		eventEntity.RequiresRsvp = event.GetRequiresRsvp()
	}
	if event.GetMaxParticipants() != 0 {
		eventEntity.MaxParticipants = event.GetMaxParticipants()
	}
	if event.GetLivestreamLink() != "" {
		eventEntity.LivestreamLink = event.GetLivestreamLink()
	}

	updatedEvent, err := h.Svc.Update(ctx, eventEntity)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update event: %v", err)
	}

	return convertEventEntityToProto(updatedEvent), nil
}

func (h *EventGrpcHandler) GetEvent(ctx context.Context, req *pb.GetEventRequest) (*pb.Event, error) {
	eventIDStr := req.GetId()
	if eventIDStr == "" {
		return nil, status.Errorf(codes.InvalidArgument, "event ID is required")
	}

	eventID, err := uuid.Parse(eventIDStr)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid event ID format: %v", err)
	}

	event, err := h.Svc.GetById(ctx, eventID.String())
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, status.Errorf(codes.NotFound, "event not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to get event: %v", err)
	}
	return convertEventEntityToProto(event), nil
}

func (h *EventGrpcHandler) DeleteEvent(ctx context.Context, req *pb.DeleteEventRequest) (*pb.DeleteEventResponse, error) {
	eventIDStr := req.GetId()
	if eventIDStr == "" {
		return nil, status.Errorf(codes.InvalidArgument, "event ID is required")
	}

	eventID, err := uuid.Parse(eventIDStr)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid event ID format: %v", err)
	}

	err = h.Svc.Delete(ctx, eventID.String())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to delete event: %v", err)
	}
	return &pb.DeleteEventResponse{}, nil
}

func (h *EventGrpcHandler) ListEvents(ctx context.Context, req *pb.ListEventsRequest) (*pb.ListEventsResponse, error) {
	events, err := h.Svc.ListEvent(ctx, req.GetPageSize(), req.GetPageToken())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list masjids: %v", err)
	}

	response := &pb.ListEventsResponse{}
	for _, event := range events {
		response.Events = append(response.Events, convertEventEntityToProto(event))
	}
	return response, nil
}

func convertEventEntityToProto(event *entity.Event) *pb.Event {
	return &pb.Event{
		Id:                event.ID.String(),
		MasjidId:          event.MasjidId,
		Name:              event.Name,
		Description:       event.Description,
		StartTime:         timestamppb.New(event.StartTime),
		EndTime:           timestamppb.New(event.EndTime),
		GenderRestriction: pb.Event_GenderRestriction(event.GenderRestriction),
		IsPaid:            event.IsPaid,
		RequiresRsvp:      event.RequiresRsvp,
		MaxParticipants:   event.MaxParticipants,
		LivestreamLink:    event.LivestreamLink,
		CreateTime:        timestamppb.New(event.CreatedAt),
		UpdateTime:        timestamppb.New(event.UpdatedAt),
	}
}
