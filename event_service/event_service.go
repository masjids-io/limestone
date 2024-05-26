package event_service

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/mnadev/limestone/proto"
	"github.com/mnadev/limestone/storage"
)

type EventServiceServer struct {
	Smgr *storage.StorageManager
	pb.UnimplementedEventServiceServer
}

func (s *EventServiceServer) CreateEvent(ctx context.Context, in *pb.CreateEventRequest) (*pb.Event, error) {
	event, err := s.Smgr.CreateEvent(in.GetEvent())
	if err != nil {
		return nil, err
	}
	return event.ToProto(), status.Error(codes.OK, codes.OK.String())
}

func (s *EventServiceServer) GetEvent(ctx context.Context, in *pb.GetEventRequest) (*pb.Event, error) {
	event, err := s.Smgr.GetEvent(in.GetEventId())
	if err != nil {
		return nil, err
	}
	return event.ToProto(), status.Error(codes.OK, codes.OK.String())
}

func (s *EventServiceServer) UpdateEvent(ctx context.Context, in *pb.UpdateEventRequest) (*pb.Event, error) {
	event, err := s.Smgr.UpdateEvent(in.GetEvent())
	if err != nil {
		return nil, err
	}
	return event.ToProto(), status.Error(codes.OK, codes.OK.String())
}

func (s *EventServiceServer) DeleteEvent(ctx context.Context, in *pb.DeleteEventRequest) (*pb.DeleteEventResponse, error) {
	err := s.Smgr.DeleteEvent(in.GetEventId())
	if err != nil {
		return nil, err
	}
	return &pb.DeleteEventResponse{}, status.Error(codes.OK, codes.OK.String())
}
