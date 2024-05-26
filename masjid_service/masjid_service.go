package masjid_service

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/mnadev/limestone/proto"
	"github.com/mnadev/limestone/storage"
)

type MasjidServiceServer struct {
	smgr *storage.StorageManager
	pb.UnimplementedMasjidServiceServer
}

func (s *MasjidServiceServer) CreateMasjid(ctx context.Context, in *pb.CreateMasjidRequest) (*pb.Masjid, error) {
	masjid, err := s.smgr.CreateMasjid(in.GetMasjid())
	if err != nil {
		return nil, err
	}
	return masjid.ToProto(), status.Error(codes.OK, codes.OK.String())
}

func (s *MasjidServiceServer) GetMasjid(ctx context.Context, in *pb.GetMasjidRequest) (*pb.Masjid, error) {
	masjid, err := s.smgr.GetMasjid(in.GetMasjidId())
	if err != nil {
		return nil, err
	}
	return masjid.ToProto(), status.Error(codes.OK, codes.OK.String())
}

func (s *MasjidServiceServer) UpdateMasjid(ctx context.Context, in *pb.UpdateMasjidRequest) (*pb.Masjid, error) {
	masjid, err := s.smgr.UpdateMasjid(in.GetMasjid())
	if err != nil {
		return nil, err
	}
	return masjid.ToProto(), status.Error(codes.OK, codes.OK.String())
}

func (s *MasjidServiceServer) DeleteMasjid(ctx context.Context, in *pb.DeleteMasjidRequest) (*pb.DeleteMasjidResponse, error) {
	err := s.smgr.DeleteMasjid(in.GetMasjidId())
	if err != nil {
		return nil, err
	}
	return &pb.DeleteMasjidResponse{}, status.Error(codes.OK, codes.OK.String())
}
