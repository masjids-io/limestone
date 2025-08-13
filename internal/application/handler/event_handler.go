package handler

import (
	"context"
	"errors"
	"github.com/google/uuid"
	pb "github.com/mnadev/limestone/gen/go"
	"github.com/mnadev/limestone/internal/application/domain/entity"
	"github.com/mnadev/limestone/internal/application/helper"
	"github.com/mnadev/limestone/internal/application/services"
	"github.com/mnadev/limestone/internal/infrastructure/auth"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"
	"time"
)

type EventGrpcHandler struct {
	pb.UnimplementedEventServiceServer
	Svc       *services.EventService
	MasjidSvc *services.MasjidService
}

func NewEventGrpcHandler(svc *services.EventService, masjidSvc *services.MasjidService) *EventGrpcHandler {
	return &EventGrpcHandler{
		Svc:       svc,
		MasjidSvc: masjidSvc,
	}
}

func (h *EventGrpcHandler) CreateEvent(ctx context.Context, req *pb.CreateEventRequest) (*pb.StandardEventResponse, error) {
	event := req.GetEvent()

	masjidIdStr := event.GetMasjidId()
	masjidId, err := uuid.Parse(masjidIdStr)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid masjid_id format: %v", err)
	}

	_, err = h.MasjidSvc.GetMasjid(ctx, masjidId.String())
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, status.Errorf(codes.NotFound, "masjid with id '%s' not found", masjidIdStr)
		}
		return nil, status.Errorf(codes.Internal, "failed to verify masjid: %v", err)
	}

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

	return helper.StandardEventResponse(codes.OK, "success", "event created successfully", createdEvent, nil, nil)
}

func (h *EventGrpcHandler) UpdateEvent(ctx context.Context, req *pb.UpdateEventRequest) (*pb.StandardEventResponse, error) {
	eventData := req.GetEvent()
	if eventData == nil {
		return nil, status.Errorf(codes.InvalidArgument, "event data is required")
	}

	eventID, err := uuid.Parse(eventData.GetId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid event ID format: %v", err)
	}

	// 1. READ: Ambil data event yang ada dari database terlebih dahulu.
	existingEvent, err := h.Svc.GetById(ctx, eventID.String()) // Anda perlu method Get(ctx, id) di service.
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, status.Errorf(codes.NotFound, "event with id %s not found", eventID)
		}
		return nil, status.Errorf(codes.Internal, "failed to retrieve event for update: %v", err)
	}

	// 2. MODIFY: Timpa field pada data yang ada dengan nilai dari request.
	// Hanya timpa jika nilai dari request bukan nilai default (zero-value).
	if eventData.GetName() != "" {
		existingEvent.Name = eventData.GetName()
	}
	if eventData.GetDescription() != "" {
		existingEvent.Description = eventData.GetDescription()
	}
	if eventData.GetMasjidId() != "" {
		existingEvent.MasjidId = eventData.GetMasjidId()
	}
	if eventData.GetStartTime() != nil && eventData.GetStartTime().IsValid() {
		existingEvent.StartTime = eventData.GetStartTime().AsTime()
	}
	if eventData.GetEndTime() != nil && eventData.GetEndTime().IsValid() {
		existingEvent.EndTime = eventData.GetEndTime().AsTime()
	}
	if eventData.GetMaxParticipants() > 0 {
		existingEvent.MaxParticipants = eventData.GetMaxParticipants()
	}
	if eventData.GetLivestreamLink() != "" {
		existingEvent.LivestreamLink = eventData.GetLivestreamLink()
	}

	// Catatan: Untuk boolean, pendekatan ini memiliki batasan.
	// Anda tidak bisa mengubah nilai dari 'true' ke 'false' karena 'false' adalah nilai default.
	// Untuk saat ini, kita asumsikan hanya bisa mengubah dari 'false' ke 'true'.
	if eventData.GetIsPaid() {
		existingEvent.IsPaid = true
	}
	if eventData.GetRequiresRsvp() {
		existingEvent.RequiresRsvp = true
	}

	existingEvent.UpdatedAt = time.Now()

	// 3. WRITE: Simpan kembali seluruh objek yang sudah diperbarui.
	// Pastikan service 'Update' Anda menerima dan menyimpan seluruh objek entity.
	updatedEvent, err := h.Svc.Update(ctx, existingEvent)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update event: %v", err)
	}

	return helper.StandardEventResponse(codes.OK, "success", "event updated successfully", updatedEvent, nil, nil)
}

func (h *EventGrpcHandler) GetEvent(ctx context.Context, req *pb.GetEventRequest) (*pb.StandardEventResponse, error) {
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
	return helper.StandardEventResponse(codes.OK, "success", "event retrieved successfully", event, nil, nil)
}

func (h *EventGrpcHandler) DeleteEvent(ctx context.Context, req *pb.DeleteEventRequest) (*pb.StandardEventResponse, error) {
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
	return helper.StandardEventResponse(codes.OK, "success", "event deleted successfully", nil, nil, &pb.DeleteEventResponse{})
}

func (h *EventGrpcHandler) ListEvents(ctx context.Context, req *pb.ListEventsRequest) (*pb.ListEventsResponse, error) {
	allowedRolesForAnyUser := []string{
		string(entity.MASJID_ADMIN),
		string(entity.MASJID_MEMBER),
		string(entity.MASJID_VOLUNTEER),
		string(entity.MASJID_IMAM),
	}
	if err := auth.RequireRole(ctx, allowedRolesForAnyUser, "ListEvents"); err != nil {
		return nil, err
	}

	page := req.GetPage()
	if page <= 0 {
		page = 1
	}
	limit := req.GetLimit()
	if limit <= 0 {
		limit = 10
	}
	searchQuery := req.GetSearch()

	events, totalItems, err := h.Svc.ListEvents(ctx, searchQuery, page, limit)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list events: %v", err)
	}

	var totalPages int32
	if totalItems > 0 {
		totalPages = int32((totalItems + int64(limit) - 1) / int64(limit))
	}

	protoEvents := make([]*pb.Event, 0, len(events))
	for _, event := range events {
		protoEvents = append(protoEvents, &pb.Event{
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
		})
	}

	return &pb.ListEventsResponse{
		Code:        codes.OK.String(),
		Status:      "success",
		Message:     "Events retrieved successfully",
		Data:        protoEvents,
		TotalItems:  totalItems,
		CurrentPage: page,
		TotalPages:  totalPages,
	}, nil
}
