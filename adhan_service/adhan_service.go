package adhan_service

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/mnadev/limestone/proto"
	"github.com/mnadev/limestone/storage"
)

type AdhanServiceServer struct {
	SM *storage.StorageManager
	pb.UnimplementedAdhanServiceServer
}

func (srvr *AdhanServiceServer) CreateAdhanFile(ctx context.Context, in *pb.CreateAdhanFileRequest) (*pb.AdhanFile, error) {
	adhan_file, err := srvr.SM.CreateAdhanFile(in.GetAdhanFile())
	if err != nil {
		return nil, err
	}
	return adhan_file.ToProto(), status.Error(codes.OK, codes.OK.String())
}

func (srvr *AdhanServiceServer) GetAdhanFile(ctx context.Context, in *pb.GetAdhanFileRequest) (*pb.AdhanFile, error) {
	adhan_file, err := srvr.SM.GetAdhanFile(in.GetMasjidId())
	if err != nil {
		return nil, err
	}
	return adhan_file.ToProto(), status.Error(codes.OK, codes.OK.String())
}

func (srvr *AdhanServiceServer) UpdateAdhanFile(ctx context.Context, in *pb.UpdateAdhanFileRequest) (*pb.AdhanFile, error) {
	adhan_file, err := srvr.SM.UpdateAdhanFile(in.GetAdhanFile())
	if err != nil {
		return nil, err
	}
	return adhan_file.ToProto(), status.Error(codes.OK, codes.OK.String())
}

func (srvr *AdhanServiceServer) DeleteAdhanFile(ctx context.Context, in *pb.DeleteAdhanFileRequest) (*pb.DeleteAdhanFileResponse, error) {
	err := srvr.SM.DeleteAdhanFile(in.GetId())
	if err != nil {
		return nil, err
	}
	return &pb.DeleteAdhanFileResponse{}, status.Error(codes.OK, codes.OK.String())
}
