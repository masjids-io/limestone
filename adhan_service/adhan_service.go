package adhan_service

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/mnadev/limestone/proto"
	"github.com/mnadev/limestone/storage"
)

type AdhanServiceServer struct {
	Smgr *storage.StorageManager
	pb.UnimplementedAdhanServiceServer
}

func (s *AdhanServiceServer) CreateAdhanFile(ctx context.Context, in *pb.CreateAdhanFileRequest) (*pb.AdhanFile, error) {
	adhan_file, err := s.Smgr.CreateAdhanFile(in.GetAdhanFile())
	if err != nil {
		return nil, err
	}
	return adhan_file.ToProto(), status.Error(codes.OK, codes.OK.String())
}

func (s *AdhanServiceServer) GetAdhanFile(ctx context.Context, in *pb.GetAdhanFileRequest) (*pb.AdhanFile, error) {
	adhan_file, err := s.Smgr.GetAdhanFile(in.GetId())
	if err != nil {
		return nil, err
	}
	return adhan_file.ToProto(), status.Error(codes.OK, codes.OK.String())
}

func (s *AdhanServiceServer) UpdateAdhanFile(ctx context.Context, in *pb.UpdateAdhanFileRequest) (*pb.AdhanFile, error) {
	adhan_file, err := s.Smgr.UpdateAdhanFile(in.GetAdhanFile())
	if err != nil {
		return nil, err
	}
	return adhan_file.ToProto(), status.Error(codes.OK, codes.OK.String())
}

func (s *AdhanServiceServer) DeleteAdhanFile(ctx context.Context, in *pb.DeleteAdhanFileRequest) (*pb.DeleteAdhanFileResponse, error) {
	err := s.Smgr.DeleteAdhanFile(in.GetId())
	if err != nil {
		return nil, err
	}
	return &pb.DeleteAdhanFileResponse{}, status.Error(codes.OK, codes.OK.String())
}
