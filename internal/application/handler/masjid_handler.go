package handler

import (
	"context"
	"errors"
	"github.com/google/uuid"
	pb "github.com/mnadev/limestone/gen/go"
	"github.com/mnadev/limestone/internal/application/domain/entity"
	"github.com/mnadev/limestone/internal/application/helper"
	"github.com/mnadev/limestone/internal/application/services"
	"github.com/mnadev/limestone/internal/infrastructure/auth"
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

func (h *MasjidGrpcHandler) CreateMasjid(ctx context.Context, req *pb.CreateMasjidRequest) (*pb.StandardMasjidResponse, error) {
	// --- Start Authorization (Coarse-Grained) ---
	allowedRolesForAnyUser := []string{string(entity.MASJID_ADMIN)}
	if err := auth.RequireRole(ctx, allowedRolesForAnyUser, "Create Masjid"); err != nil {
		return nil, err
	}
	// --- End Authorization (Coarse-Grained) ---

	masjid := req.GetMasjid()

	masjidEntity := &entity.Masjid{
		ID:         uuid.New(),
		Name:       masjid.GetName(),
		Location:   masjid.GetLocation(),
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

	cm, err := h.Svc.CreateMasjid(ctx, masjidEntity)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create masjid: %v", err)
	}

	return helper.StandardMasjidResponse(codes.OK, "success", "masjid created successfully", cm, nil, nil)
}

func (h *MasjidGrpcHandler) UpdateMasjid(ctx context.Context, req *pb.UpdateMasjidRequest) (*pb.StandardMasjidResponse, error) {
	// --- Start Authorization (Coarse-Grained) ---
	allowedRolesForAnyUser := []string{
		string(entity.MASJID_ADMIN),
		string(entity.MASJID_VOLUNTEER),
	}
	if err := auth.RequireRole(ctx, allowedRolesForAnyUser, "UpdateUser"); err != nil {
		return nil, err
	}
	// --- End Authorization (Coarse-Grained) ---

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

	um, err := h.Svc.UpdateMasjid(ctx, masjidEntity)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update masjid: %v", err)
	}

	return helper.StandardMasjidResponse(codes.OK, "success", "masjid updated successfully", um, nil, nil)
}

func (h *MasjidGrpcHandler) DeleteMasjid(ctx context.Context, req *pb.DeleteMasjidRequest) (*pb.StandardMasjidResponse, error) {
	// --- Start Authorization (Coarse-Grained) ---
	allowedRolesForAnyUser := []string{
		string(entity.MASJID_ADMIN),
	}
	if err := auth.RequireRole(ctx, allowedRolesForAnyUser, "UpdateUser"); err != nil {
		return nil, err
	}
	// --- End Authorization (Coarse-Grained) ---
	err := h.Svc.DeleteMasjid(ctx, req.GetId())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to delete masjid: %v", err)
	}
	return helper.StandardMasjidResponse(codes.OK, "success", "masjid deleted successfully", nil, nil, &pb.DeleteMasjidResponse{})
}

func (h *MasjidGrpcHandler) GetMasjid(ctx context.Context, req *pb.GetMasjidRequest) (*pb.StandardMasjidResponse, error) {
	// --- Start Authorization (Coarse-Grained) ---
	allowedRolesForAnyUser := []string{
		string(entity.MASJID_ADMIN),
		string(entity.MASJID_VOLUNTEER),
		string(entity.MASJID_MEMBER),
		string(entity.MASJID_IMAM),
	}
	if err := auth.RequireRole(ctx, allowedRolesForAnyUser, "UpdateUser"); err != nil {
		return nil, err
	}
	// --- End Authorization (Coarse-Grained) ---
	masjid, err := h.Svc.GetMasjid(ctx, req.GetId())
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, status.Errorf(codes.NotFound, "masjid not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to get masjid: %v", err)
	}
	return helper.StandardMasjidResponse(codes.OK, "success", "masjid retrieved successfully", masjid, nil, nil)
}

func (h *MasjidGrpcHandler) ListMasjids(ctx context.Context, req *pb.ListMasjidsRequest) (*pb.StandardMasjidResponse, error) {
	// --- Start Authorization (Coarse-Grained) ---
	//allowedRolesForAnyUser := []string{
	//	string(entity.MASJID_ADMIN),
	//	string(entity.MASJID_VOLUNTEER),
	//	string(entity.MASJID_MEMBER),
	//	string(entity.MASJID_IMAM),
	//}
	//if err := auth.RequireRole(ctx, allowedRolesForAnyUser, "ListMajids"); err != nil {
	//	return nil, err
	//}
	// --- End Authorization (Coarse-Grained) ---
	params := &entity.ListMasjidsQueryParams{
		Start:    req.GetStart(),
		Limit:    req.GetLimit(),
		Page:     req.GetPage(),
		Name:     req.GetName(),
		Location: req.GetLocation(),
	}

	masjids, totalCount, err := h.Svc.ListMasjids(ctx, params)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list masjids: %v", err)
	}

	protoMasjidsList := make([]*pb.Masjid, len(masjids))
	for i, masjid := range masjids {
		protoMasjidsList[i] = &pb.Masjid{
			Id:         masjid.ID.String(),
			Name:       masjid.Name,
			Location:   masjid.Location,
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
			PrayerConfig: &pb.PrayerTimesConfiguration{
				Method:           pb.PrayerTimesConfiguration_CalculationMethod(int32(masjid.PrayerConfig.CalculationMethod)),
				FajrAngle:        masjid.PrayerConfig.FajrAngle,
				IshaAngle:        masjid.PrayerConfig.IshaAngle,
				IshaInterval:     masjid.PrayerConfig.IshaInterval,
				AsrMethod:        pb.PrayerTimesConfiguration_AsrJuristicMethod(int32(masjid.PrayerConfig.AsrMethod)),
				HighLatitudeRule: pb.PrayerTimesConfiguration_HighLatitudeRule(int32(masjid.PrayerConfig.HighLatitudeRule)),
				Adjustments: &pb.PrayerTimesConfiguration_PrayerAdjustments{
					FajrAdjustment:    masjid.PrayerConfig.Adjustments.FajrAdjustment,
					DhuhrAdjustment:   masjid.PrayerConfig.Adjustments.DhuhrAdjustment,
					AsrAdjustment:     masjid.PrayerConfig.Adjustments.AsrAdjustment,
					MaghribAdjustment: masjid.PrayerConfig.Adjustments.MaghribAdjustment,
					IshaAdjustment:    masjid.PrayerConfig.Adjustments.IshaAdjustment,
				},
			},
			CreateTime: timestamppb.New(masjid.CreatedAt),
			UpdateTime: timestamppb.New(masjid.UpdatedAt),
		}
	}

	pageSize := params.Limit
	if pageSize <= 0 {
		pageSize = 10
	}

	totalPages := int32(0)
	if pageSize > 0 {
		totalPages = int32((totalCount + int32(pageSize) - 1) / int32(pageSize))
	}

	currentPage := params.Page
	if currentPage == 0 {
		if params.Start > 0 && pageSize > 0 {
			currentPage = (params.Start / pageSize) + 1
		} else {
			currentPage = 1
		}
	}

	listMasjidsResponse := &pb.ListMasjidsResponse{
		Masjids:     protoMasjidsList,
		TotalCount:  totalCount,
		CurrentPage: currentPage,
		TotalPages:  totalPages,
	}

	return helper.StandardMasjidResponse(codes.OK, "success", "masjids retrieved successfully", nil, listMasjidsResponse, nil)
}
