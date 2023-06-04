package adhan_service

import (
	"context"

	apb "github.com/mnadev/limestone/adhan_service/proto"
	"github.com/mnadev/limestone/storage"
)

type AdhanServiceServer struct {
	SM *storage.StorageManager
	apb.UnimplementedAdhanServiceServer
}

func (srvr *AdhanServiceServer) CreateAdhanFile(ctx context.Context, in *apb.CreateAdhanFileRequest) (*apb.AdhanFile, error) {
	adhan_file, err := srvr.SM.CreateAdhanFile(in.GetAdhanFile())
	if err != nil {
		return nil, err
	}
	return adhan_file.ToProto(), nil
}

func (srvr *AdhanServiceServer) GetAdhanFile(ctx context.Context, in *apb.GetAdhanFileRequest) (*apb.AdhanFile, error) {
	adhan_file, err := srvr.SM.GetAdhanFile(in.GetMasjidId())
	if err != nil {
		return nil, err
	}
	return adhan_file.ToProto(), nil
}

func (srvr *AdhanServiceServer) UpdateAdhanFile(ctx context.Context, in *apb.UpdateAdhanFileRequest) (*apb.AdhanFile, error) {
	adhan_file, err := srvr.SM.UpdateAdhanFile(in.GetAdhanFile())
	if err != nil {
		return nil, err
	}
	return adhan_file.ToProto(), nil
}

func (srvr *AdhanServiceServer) DeleteAdhanFile(ctx context.Context, in *apb.DeleteAdhanFileRequest) (*apb.DeleteAdhanFileResponse, error) {
	err := srvr.SM.DeleteAdhanFile(in.GetId())
	if err != nil {
		return nil, err
	}
	return &apb.DeleteAdhanFileResponse{}, nil
}
