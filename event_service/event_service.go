package event_service

import (
	"context"

	epb "github.com/mnadev/limestone/event_service/proto"
	"github.com/mnadev/limestone/storage"
)

type EventServiceServer struct {
	SM *storage.StorageManager
	epb.UnimplementedEventServiceServer
}

func (srvr *EventServiceServer) CreateEvent(ctx context.Context, in *epb.CreateEventRequest) (*epb.Event, error) {
	event, err := srvr.SM.CreateEvent(in.GetEvent())
	if err != nil {
		return nil, err
	}
	return event.ToProto(), nil
}

func (srvr *EventServiceServer) GetEvent(ctx context.Context, in *epb.GetEventRequest) (*epb.Event, error) {
	event, err := srvr.SM.GetEvent(in.GetEventId())
	if err != nil {
		return nil, err
	}
	return event.ToProto(), nil
}

func (srvr *EventServiceServer) UpdateEvent(ctx context.Context, in *epb.UpdateEventRequest) (*epb.Event, error) {
	event, err := srvr.SM.UpdateEvent(in.GetEvent())
	if err != nil {
		return nil, err
	}
	return event.ToProto(), nil
}

func (srvr *EventServiceServer) DeleteEvent(ctx context.Context, in *epb.DeleteEventRequest) (*epb.DeleteEventResponse, error) {
	err := srvr.SM.DeleteEvent(in.GetEventId())
	if err != nil {
		return nil, err
	}
	return &epb.DeleteEventResponse{}, nil
}
