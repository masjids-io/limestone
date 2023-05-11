package masjid_service

import (
	"context"

	mpb "github.com/mnadev/limestone/masjid_service/proto"
	"github.com/mnadev/limestone/storage"
)

type MasjidServiceServer struct {
	SM *storage.StorageManager
	mpb.UnimplementedMasjidServiceServer
}

func (srvr *MasjidServiceServer) CreateMasjid(ctx context.Context, in *mpb.CreateMasjidRequest) (*mpb.Masjid, error) {
	masjid, err := srvr.SM.CreateMasjid(in.GetMasjid())
	if err != nil {
		return nil, err
	}
	return masjid.ToProto(), nil
}

func (srvr *MasjidServiceServer) GetMasjid(ctx context.Context, in *mpb.GetMasjidRequest) (*mpb.Masjid, error) {
	masjid, err := srvr.SM.GetMasjid(in.GetMasjidId())
	if err != nil {
		return nil, err
	}
	return masjid.ToProto(), nil
}

func (srvr *MasjidServiceServer) UpdateMasjid(ctx context.Context, in *mpb.UpdateMasjidRequest) (*mpb.Masjid, error) {
	masjid, err := srvr.SM.UpdateMasjid(in.GetMasjid())
	if err != nil {
		return nil, err
	}
	return masjid.ToProto(), nil
}

func (srvr *MasjidServiceServer) DeleteMasjid(ctx context.Context, in *mpb.DeleteMasjidRequest) (*mpb.DeleteMasjidResponse, error) {
	err := srvr.SM.DeleteMasjid(in.GetMasjidId())
	if err != nil {
		return nil, err
	}
	return &mpb.DeleteMasjidResponse{}, nil
}
