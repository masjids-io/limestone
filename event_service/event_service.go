package event_service

import (
	"context"

	pb "github.com/mnadev/limestone/proto"
	"github.com/mnadev/limestone/storage"
)

type EventServiceServer struct {
	SM *storage.StorageManager
	pb.UnimplementedEventServiceServer
}

func (srvr *EventServiceServer) CreateEvent(ctx context.Context, in *pb.CreateEventRequest) (*pb.Event, error) {
	event, err := srvr.SM.CreateEvent(in.GetEvent())
	if err != nil {
		return nil, err
	}
	return event.ToProto(), nil
}

func (srvr *EventServiceServer) GetEvent(ctx context.Context, in *pb.GetEventRequest) (*pb.Event, error) {
	event, err := srvr.SM.GetEvent(in.GetEventId())
	if err != nil {
		return nil, err
	}
	return event.ToProto(), nil
}

func (srvr *EventServiceServer) UpdateEvent(ctx context.Context, in *pb.UpdateEventRequest) (*pb.Event, error) {
	event, err := srvr.SM.UpdateEvent(in.GetEvent())
	if err != nil {
		return nil, err
	}
	return event.ToProto(), nil
}

func (srvr *EventServiceServer) DeleteEvent(ctx context.Context, in *pb.DeleteEventRequest) (*pb.DeleteEventResponse, error) {
	err := srvr.SM.DeleteEvent(in.GetEventId())
	if err != nil {
		return nil, err
	}
	return &pb.DeleteEventResponse{}, nil
}
