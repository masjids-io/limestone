package masjid_service

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/mnadev/limestone/proto"
	"github.com/mnadev/limestone/storage"
)

type MasjidServiceServer struct {
	SM *storage.StorageManager
	pb.UnimplementedMasjidServiceServer
}

func (srvr *MasjidServiceServer) CreateMasjid(ctx context.Context, in *pb.CreateMasjidRequest) (*pb.Masjid, error) {
	masjid, err := srvr.SM.CreateMasjid(in.GetMasjid())
	if err != nil {
		return nil, err
	}
	return masjid.ToProto(), status.Error(codes.OK, codes.OK.String())
}

func (srvr *MasjidServiceServer) GetMasjid(ctx context.Context, in *pb.GetMasjidRequest) (*pb.Masjid, error) {
	masjid, err := srvr.SM.GetMasjid(in.GetMasjidId())
	if err != nil {
		return nil, err
	}
	return masjid.ToProto(), status.Error(codes.OK, codes.OK.String())
}

func (srvr *MasjidServiceServer) UpdateMasjid(ctx context.Context, in *pb.UpdateMasjidRequest) (*pb.Masjid, error) {
	masjid, err := srvr.SM.UpdateMasjid(in.GetMasjid())
	if err != nil {
		return nil, err
	}
	return masjid.ToProto(), status.Error(codes.OK, codes.OK.String())
}

func (srvr *MasjidServiceServer) DeleteMasjid(ctx context.Context, in *pb.DeleteMasjidRequest) (*pb.DeleteMasjidResponse, error) {
	err := srvr.SM.DeleteMasjid(in.GetMasjidId())
	if err != nil {
		return nil, err
	}
	return &pb.DeleteMasjidResponse{}, status.Error(codes.OK, codes.OK.String())
}
