package handler

import (
	"context"
	"github.com/google/uuid"
	pb "github.com/mnadev/limestone/gen/go"
	"github.com/mnadev/limestone/internal/application/domain/entity"
	"github.com/mnadev/limestone/internal/application/helper"
	"github.com/mnadev/limestone/internal/application/services"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"time"
)

type AdhanGrpcHandler struct {
	pb.UnimplementedAdhanServiceServer
	Svc *services.AdhanService
}

func NewAdhanGrpcHandler(svc *services.AdhanService) *AdhanGrpcHandler {
	return &AdhanGrpcHandler{Svc: svc}
}

const maxAdhanFileSizeMB = 5

func (h *AdhanGrpcHandler) CreateAdhan(ctx context.Context, req *pb.CreateAdhanFileRequest) (*pb.StandardAdhanResponse, error) {
	adhanFile := req.GetAdhanFile()
	if adhanFile == nil {
		return nil, status.Errorf(codes.InvalidArgument, "adhan file data is required")
	}

	masjidID := adhanFile.GetMasjidId()
	if masjidID == "" {
		return nil, status.Errorf(codes.InvalidArgument, "masjid ID is required")
	}

	audioBytes := adhanFile.GetFile()

	if len(audioBytes) == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "adhan file content is required")
	}

	if len(audioBytes) > int(maxAdhanFileSizeMB*1024*1024) {
		return nil, status.Errorf(codes.InvalidArgument, "adhan file size exceeds maximum allowed size (%d MB)", maxAdhanFileSizeMB)
	}

	if len(audioBytes) > int(maxAdhanFileSizeMB*1024*1024) {
		return nil, status.Errorf(codes.InvalidArgument, "adhan file size exceeds maximum allowed size (%d MB)", maxAdhanFileSizeMB)
	}

	if !helper.IsAudioFile(audioBytes) {
		return nil, status.Errorf(codes.InvalidArgument, "invalid adhan file type. Only MP3 and WAV are supported for now.")
	}

	adhanEntity := &entity.Adhan{
		ID:        uuid.New(),
		MasjidId:  masjidID,
		File:      audioBytes,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	createdAdhan, err := h.Svc.CreateAdhan(ctx, adhanEntity) // Asumsi nama service method CreateAdhanFile
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create adhan file: %v", err)
	}

	return helper.StandardAdhanResponse(codes.OK, "success", "adhan file created successfully", createdAdhan, nil)
}

func (h *AdhanGrpcHandler) UpdateAdhan(ctx context.Context, req *pb.UpdateAdhanFileRequest) (*pb.StandardAdhanResponse, error) {
	adhanFile := req.GetAdhanFile()
	if adhanFile == nil {
		return nil, status.Errorf(codes.InvalidArgument, "adhan file data is required for update")
	}

	idStr := req.GetId()
	if idStr == "" {
		return nil, status.Errorf(codes.InvalidArgument, "adhan file ID is required for update")
	}

	id, err := uuid.Parse(idStr)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid adhan file ID format: %v", err)
	}

	masjidID := adhanFile.GetMasjidId()
	audioBytes := adhanFile.GetFile()

	if len(audioBytes) == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "adhan file content is required")
	}

	if len(audioBytes) > int(maxAdhanFileSizeMB*1024*1024) {
		return nil, status.Errorf(codes.InvalidArgument, "adhan file size exceeds maximum allowed size (%d MB)", maxAdhanFileSizeMB)
	}

	if len(audioBytes) > int(maxAdhanFileSizeMB*1024*1024) {
		return nil, status.Errorf(codes.InvalidArgument, "adhan file size exceeds maximum allowed size (%d MB)", maxAdhanFileSizeMB)
	}

	if !helper.IsAudioFile(audioBytes) {
		return nil, status.Errorf(codes.InvalidArgument, "invalid adhan file type. Only MP3 and WAV are supported for now.")
	}

	existingAdhan, err := h.Svc.GetAdhanByID(ctx, idStr)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return nil, status.Errorf(codes.NotFound, "adhan file with ID %s not found", id)
		}
		return nil, status.Errorf(codes.Internal, "failed to get existing adhan file: %v", err)
	}

	updatedAdhanEntity := &entity.Adhan{
		ID:        id,
		MasjidId:  masjidID,
		File:      audioBytes,
		UpdatedAt: time.Now(),
	}
	if len(audioBytes) == 0 {
		updatedAdhanEntity.File = existingAdhan.File
	}

	updatedAdhan, err := h.Svc.UpdateAdhan(ctx, updatedAdhanEntity)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update adhan file: %v", err)
	}

	return helper.StandardAdhanResponse(codes.OK, "success", "adhan file updated successfully", updatedAdhan, nil)
}

func (h *AdhanGrpcHandler) GetAdhanById(ctx context.Context, req *pb.GetAdhanFileRequest) (*pb.StandardAdhanResponse, error) {
	idStr := req.GetId()
	if idStr == "" {
		return nil, status.Errorf(codes.InvalidArgument, "adhan file ID is required")
	}

	adhan, err := h.Svc.GetAdhanByID(ctx, idStr)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return nil, status.Errorf(codes.NotFound, "adhan file with ID %s not found", idStr)
		}
		return nil, status.Errorf(codes.Internal, "failed to get adhan file: %v", err)
	}

	return helper.StandardAdhanResponse(codes.OK, "success", "adhan file retrieved successfully", adhan, nil)
}

func (h *AdhanGrpcHandler) DeleteAdhan(ctx context.Context, req *pb.DeleteAdhanFileRequest) (*pb.StandardAdhanResponse, error) {

	userIDStr := req.GetId()
	if userIDStr == "" {
		return nil, status.Errorf(codes.InvalidArgument, "id is required")
	}

	_, err := uuid.Parse(userIDStr)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid user ID format")
	}

	err = h.Svc.DeleteAdhan(ctx, userIDStr)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "failed to retrieve Delete user")
	}

	return helper.StandardAdhanResponse(codes.OK, "success", "adhan file deleted successfully", nil, &pb.DeleteAdhanFileResponse{})
}
