package handler

import (
	"context"
	"errors"
	"github.com/google/uuid"
	pb "github.com/mnadev/limestone/gen/go"
	"github.com/mnadev/limestone/internal/application/domain/entity"
	services "github.com/mnadev/limestone/internal/application/services"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"
	"time"
)

type MasjidGrpcHandler struct {
	pb.UnimplementedMasjidServiceServer
	Svc *services.MasjidService
}

func NewMasjidGrpcHandler(svc *services.MasjidService) *MasjidGrpcHandler {
	return &MasjidGrpcHandler{Svc: svc}
}

func (h *MasjidGrpcHandler) CreateMasjid(ctx context.Context, req *pb.CreateMasjidRequest) (*pb.Masjid, error) {
	masjid := req.GetMasjid()

	masjidEntity := &entity.Masjid{
		ID:         uuid.New(), // Generate a new UUID
		Name:       masjid.GetName(),
		IsVerified: masjid.GetIsVerified(),
		Address: entity.Address{
			AddressLine1: masjid.GetAddress().GetAddressLine_1(),
			AddressLine2: masjid.GetAddress().GetAddressLine_2(),
			ZoneCode:     masjid.GetAddress().GetZoneCode(),
			PostalCode:   masjid.GetAddress().GetPostalCode(),
			City:         masjid.GetAddress().GetCity(),
			CountryCode:  masjid.GetAddress().GetCountryCode(),
		},
		PhoneNumber: entity.PhoneNumber{
			PhoneCountryCode: masjid.GetPhoneNumber().GetCountryCode(),
			Number:           masjid.GetPhoneNumber().GetNumber(),
			Extension:        masjid.GetPhoneNumber().GetExtension(),
		},
		PrayerConfig: entity.PrayerTimesConfiguration{
			CalculationMethod: entity.CalculationMethod(masjid.GetPrayerConfig().GetMethod()),
			FajrAngle:         masjid.GetPrayerConfig().GetFajrAngle(),
			IshaAngle:         masjid.GetPrayerConfig().GetIshaAngle(),
			IshaInterval:      masjid.GetPrayerConfig().GetIshaInterval(),
			AsrMethod:         entity.AsrJuristicMethod(masjid.GetPrayerConfig().GetAsrMethod()),
			HighLatitudeRule:  entity.HighLatitudeRule(masjid.GetPrayerConfig().GetHighLatitudeRule()),
			Adjustments: entity.PrayerAdjustments{
				FajrAdjustment:    int32(int(masjid.GetPrayerConfig().GetAdjustments().GetFajrAdjustment())),
				DhuhrAdjustment:   int32(int(masjid.GetPrayerConfig().GetAdjustments().GetDhuhrAdjustment())),
				AsrAdjustment:     int32(int(masjid.GetPrayerConfig().GetAdjustments().GetAsrAdjustment())),
				MaghribAdjustment: int32(int(masjid.GetPrayerConfig().GetAdjustments().GetMaghribAdjustment())),
				IshaAdjustment:    int32(int(masjid.GetPrayerConfig().GetAdjustments().GetIshaAdjustment())),
			},
		},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	createdMasjid, err := h.Svc.CreateMasjid(ctx, masjidEntity)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create masjid: %v", err)
	}

	return convertMasjidEntityToProto(createdMasjid), nil
}

func (h *MasjidGrpcHandler) UpdateMasjid(ctx context.Context, req *pb.UpdateMasjidRequest) (*pb.Masjid, error) {
	masjid := req.GetMasjid()
	if masjid == nil {
		return nil, status.Errorf(codes.InvalidArgument, "masjid data is required")
	}

	masjidIDStr := masjid.GetId()
	if masjidIDStr == "" {
		return nil, status.Errorf(codes.InvalidArgument, "masjid ID is required for update")
	}

	masjidID, err := uuid.Parse(masjidIDStr)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid masjid ID format")
	}

	masjidEntity := &entity.Masjid{
		ID: masjidID,
	}

	if masjid.GetName() != "" {
		masjidEntity.Name = masjid.GetName()
	}
	// Add similar checks and assignments for other fields you want to support updating
	if masjid.GetAddress() != nil {
		masjidEntity.Address = entity.Address{
			AddressLine1: masjid.GetAddress().GetAddressLine_1(),
			AddressLine2: masjid.GetAddress().GetAddressLine_2(),
			ZoneCode:     masjid.GetAddress().GetZoneCode(),
			PostalCode:   masjid.GetAddress().GetPostalCode(),
			City:         masjid.GetAddress().GetCity(),
			CountryCode:  masjid.GetAddress().GetCountryCode(),
		}
	}

	if masjid.GetPhoneNumber() != nil {
		masjidEntity.PhoneNumber = entity.PhoneNumber{
			PhoneCountryCode: masjid.GetPhoneNumber().GetCountryCode(),
			Number:           masjid.GetPhoneNumber().GetNumber(),
			Extension:        masjid.GetPhoneNumber().GetExtension(),
		}
	}

	if masjid.GetPrayerConfig() != nil {
		masjidEntity.PrayerConfig = entity.PrayerTimesConfiguration{
			CalculationMethod: entity.CalculationMethod(masjid.GetPrayerConfig().GetMethod()),
			FajrAngle:         masjid.GetPrayerConfig().GetFajrAngle(),
			IshaAngle:         masjid.GetPrayerConfig().GetIshaAngle(),
			IshaInterval:      masjid.GetPrayerConfig().GetIshaInterval(),
			AsrMethod:         entity.AsrJuristicMethod(masjid.GetPrayerConfig().GetAsrMethod()),
			HighLatitudeRule:  entity.HighLatitudeRule(masjid.GetPrayerConfig().GetHighLatitudeRule()),
			Adjustments: entity.PrayerAdjustments{
				FajrAdjustment:    int32(int(masjid.GetPrayerConfig().GetAdjustments().GetFajrAdjustment())),
				DhuhrAdjustment:   int32(int(masjid.GetPrayerConfig().GetAdjustments().GetDhuhrAdjustment())),
				AsrAdjustment:     int32(int(masjid.GetPrayerConfig().GetAdjustments().GetAsrAdjustment())),
				MaghribAdjustment: int32(int(masjid.GetPrayerConfig().GetAdjustments().GetMaghribAdjustment())),
				IshaAdjustment:    int32(int(masjid.GetPrayerConfig().GetAdjustments().GetIshaAdjustment())),
			},
		}
	}

	masjidEntity.UpdatedAt = time.Now()

	updatedMasjid, err := h.Svc.UpdateMasjid(ctx, masjidEntity)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update masjid: %v", err)
	}

	return convertMasjidEntityToProto(updatedMasjid), nil
}

func (h *MasjidGrpcHandler) DeleteMasjid(ctx context.Context, req *pb.DeleteMasjidRequest) (*pb.DeleteMasjidResponse, error) {
	err := h.Svc.DeleteMasjid(ctx, req.GetId())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to delete masjid: %v", err)
	}
	// Return an empty DeleteMasjidResponse on success.
	return &pb.DeleteMasjidResponse{}, nil
}

func (h *MasjidGrpcHandler) GetMasjid(ctx context.Context, req *pb.GetMasjidRequest) (*pb.Masjid, error) {
	masjid, err := h.Svc.GetMasjid(ctx, req.GetId())
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, status.Errorf(codes.NotFound, "masjid not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to get masjid: %v", err)
	}
	return convertMasjidEntityToProto(masjid), nil
}

func (h *MasjidGrpcHandler) ListMasjids(ctx context.Context, req *pb.ListMasjidsRequest) (*pb.ListMasjidsResponse, error) {
	masjids, err := h.Svc.ListMasjids(ctx, req.GetPageSize(), req.GetPageToken())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list masjids: %v", err)
	}

	response := &pb.ListMasjidsResponse{}
	for _, masjid := range masjids {
		response.Masjids = append(response.Masjids, convertMasjidEntityToProto(masjid))
	}
	return response, nil
}

// Helper function to convert between entity.Masjid and pb.Masjid
func convertMasjidEntityToProto(masjid *entity.Masjid) *pb.Masjid {
	return &pb.Masjid{
		Id:         masjid.ID.String(),
		Name:       masjid.Name,
		IsVerified: masjid.IsVerified,
		Address: &pb.Masjid_Address{
			AddressLine_1: masjid.Address.AddressLine1,
			AddressLine_2: masjid.Address.AddressLine2,
			ZoneCode:      masjid.Address.ZoneCode,
			PostalCode:    masjid.Address.PostalCode,
			City:          masjid.Address.City,
			CountryCode:   masjid.Address.CountryCode,
		},
		PhoneNumber: &pb.Masjid_PhoneNumber{
			CountryCode: masjid.PhoneNumber.PhoneCountryCode,
			Number:      masjid.PhoneNumber.Number,
			Extension:   masjid.PhoneNumber.Extension,
		},

		CreateTime: timestamppb.New(masjid.CreatedAt),
		UpdateTime: timestamppb.New(masjid.UpdatedAt),
	}
}
