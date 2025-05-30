package helper

import (
	"fmt"
	pb "github.com/mnadev/limestone/gen/go"
	"github.com/mnadev/limestone/internal/application/domain/entity"
	"google.golang.org/grpc/codes"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func StandardUserResponse(code codes.Code, statusMessage string, message string, userEntity *entity.User, deleteResponse *pb.DeleteUserResponse) (*pb.StandardUserResponse, error) {
	resp := &pb.StandardUserResponse{
		Code:    code.String(),
		Status:  statusMessage,
		Message: message,
	}

	if userEntity != nil {
		protoUser := &pb.User{
			Id:              userEntity.ID.String(),
			Email:           userEntity.Email,
			Username:        userEntity.Username,
			IsEmailVerified: userEntity.IsVerified,
			FirstName:       userEntity.FirstName,
			LastName:        userEntity.LastName,
			PhoneNumber:     userEntity.PhoneNumber,
			Gender:          pb.User_Gender(pb.User_Gender_value[userEntity.Gender.String()]),
			CreateTime:      timestamppb.New(userEntity.CreatedAt),
			UpdateTime:      timestamppb.New(userEntity.UpdatedAt),
		}
		resp.Data = &pb.StandardUserResponse_GetUserResponse{GetUserResponse: protoUser}
	} else if deleteResponse != nil {
		resp.Data = &pb.StandardUserResponse_DeleteUserResponse{DeleteUserResponse: deleteResponse}
	}

	return resp, nil
}

func StandardMasjidResponse(code codes.Code, status string, message string, masjid *entity.Masjid, listMasjidsResponse *pb.ListMasjidsResponse, deleteMasjidResponse *pb.DeleteMasjidResponse) (*pb.StandardMasjidResponse, error) {
	resp := &pb.StandardMasjidResponse{
		Code:    code.String(),
		Status:  status,
		Message: message,
	}

	if masjid != nil {
		resp.Data = &pb.StandardMasjidResponse_Masjid{
			Masjid: &pb.Masjid{
				Id:         masjid.ID.String(),
				Name:       masjid.Name,
				IsVerified: masjid.IsVerified,
				Location:   masjid.Location,
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
			},
		}
	} else if listMasjidsResponse != nil {
		resp.Data = &pb.StandardMasjidResponse_ListMasjidResponse{
			ListMasjidResponse: listMasjidsResponse,
		}
	} else if deleteMasjidResponse != nil {
		resp.Data = &pb.StandardMasjidResponse_DeleteMasjidResponse{
			DeleteMasjidResponse: deleteMasjidResponse,
		}
	}

	return resp, nil
}

func StandardEventResponse(code codes.Code, statusMessage string, message string, eventEntity *entity.Event, listResponse *pb.ListEventsResponse, deleteResponse *pb.DeleteEventResponse) (*pb.StandardEventResponse, error) {
	resp := &pb.StandardEventResponse{
		Code:    code.String(),
		Status:  statusMessage,
		Message: message,
	}

	if eventEntity != nil {
		protoEvent := &pb.Event{
			Id:                eventEntity.ID.String(),
			MasjidId:          eventEntity.MasjidId,
			Name:              eventEntity.Name,
			Description:       eventEntity.Description,
			StartTime:         timestamppb.New(eventEntity.StartTime),
			EndTime:           timestamppb.New(eventEntity.EndTime),
			GenderRestriction: pb.Event_GenderRestriction(eventEntity.GenderRestriction),
			IsPaid:            eventEntity.IsPaid,
			RequiresRsvp:      eventEntity.RequiresRsvp,
			MaxParticipants:   eventEntity.MaxParticipants,
			LivestreamLink:    eventEntity.LivestreamLink,
			CreateTime:        timestamppb.New(eventEntity.CreatedAt),
			UpdateTime:        timestamppb.New(eventEntity.UpdatedAt),
		}
		resp.Data = &pb.StandardEventResponse_Event{Event: protoEvent}
	} else if listResponse != nil {
		resp.Data = &pb.StandardEventResponse_ListEventResponse{ListEventResponse: listResponse}
	} else if deleteResponse != nil {
		resp.Data = &pb.StandardEventResponse_DeleteEventResponse{DeleteEventResponse: deleteResponse}
	}

	return resp, nil
}

func StandardAdhanResponse(code codes.Code, statusMessage string, message string, adhanEntity *entity.Adhan, deleteResponse *pb.DeleteAdhanFileResponse) (*pb.StandardAdhanResponse, error) {
	resp := &pb.StandardAdhanResponse{
		Code:    code.String(),
		Status:  statusMessage,
		Message: message,
	}

	if adhanEntity != nil {
		protoAdhan := &pb.AdhanFile{
			Id:         adhanEntity.ID.String(),
			MasjidId:   adhanEntity.MasjidId,
			File:       adhanEntity.File,
			CreateTime: timestamppb.New(adhanEntity.CreatedAt),
			UpdateTime: timestamppb.New(adhanEntity.UpdatedAt),
		}
		resp.Data = &pb.StandardAdhanResponse_AdhanFile{AdhanFile: protoAdhan}
	} else if deleteResponse != nil {
		resp.Data = &pb.StandardAdhanResponse_DeleteAdhanFileResponse{DeleteAdhanFileResponse: deleteResponse}
	}

	return resp, nil
}

func StandardNikkahResponse(code codes.Code, statusMessage string, message string, entityData interface{}) (*pb.StandardNikkahResponse, error) {
	resp := &pb.StandardNikkahResponse{
		Code:    code.String(),
		Status:  statusMessage,
		Message: message,
	}

	if entityData != nil {
		switch data := entityData.(type) {
		case *entity.NikkahProfile:
			resp.Data = &pb.StandardNikkahResponse_NikkahProfile{
				NikkahProfile: ToProtoNikkahProfile(data),
			}
		case *entity.NikkahLike:
			protoLike := &pb.NikkahLike{
				LikeId:         data.ID.String(),
				LikerProfileId: data.LikerProfileID,
				LikedProfileId: data.LikedProfileID,
				CreateTime:     timestamppb.New(data.CreatedAt),
				UpdateTime:     timestamppb.New(data.UpdatedAt),
			}
			switch data.Status {
			case entity.LikeStatusInitiated:
				protoLike.Status = pb.NikkahLike_INITIATED
			case entity.LikeStatusCompleted:
				protoLike.Status = pb.NikkahLike_COMPLETED
			case entity.LikeStatusCancelled:
				protoLike.Status = pb.NikkahLike_CANCELLED
			default:
				protoLike.Status = pb.NikkahLike_STATUS_UNSPECIFIED
			}
			resp.Data = &pb.StandardNikkahResponse_Like{Like: protoLike}

		case *entity.NikkahMatch:
			protoMatch := &pb.NikkahMatch{
				MatchId:            data.ID.String(),
				InitiatorProfileId: data.InitiatorProfileID.String(),
				ReceiverProfileId:  data.ReceiverProfileID.String(),
				CreateTime:         timestamppb.New(data.CreatedAt),
				UpdateTime:         timestamppb.New(data.UpdatedAt),
			}
			switch data.Status {
			case entity.MatchStatusInitiated:
				protoMatch.Status = pb.NikkahMatch_INITIATED
			case entity.MatchStatusAccepted:
				protoMatch.Status = pb.NikkahMatch_ACCEPTED
			case entity.MatchStatusRejected:
				protoMatch.Status = pb.NikkahMatch_REJECTED
			case entity.MatchStatusEnded:
				protoMatch.Status = pb.NikkahMatch_ENDED
			default:
				protoMatch.Status = pb.NikkahMatch_STATUS_UNSPECIFIED
			}
			resp.Data = &pb.StandardNikkahResponse_Match{Match: protoMatch}

		default:
			fmt.Printf("Warning: Unknown entity type passed to StandardNikkahResponse: %T\n", data)
		}
	}

	return resp, nil
}

func StandardRevertResponse(code codes.Code, statusStr string, message string, data interface{}) (*pb.StandardRevertResponse, error) {
	resp := &pb.StandardRevertResponse{
		Code:    code.String(),
		Status:  statusStr,
		Message: message,
	}

	if data != nil {
		switch d := data.(type) {
		case *entity.RevertProfile:
			resp.Data = &pb.StandardRevertResponse_RevertProfile{
				RevertProfile: ToProtoRevertProfile(d),
			}
		case *entity.RevertMatch:
			resp.Data = &pb.StandardRevertResponse_RevertMatch{
				RevertMatch: ToProtoRevertMatch(d),
			}
		case *pb.ListRevertProfilesResponse:
			resp.Data = &pb.StandardRevertResponse_ListRevertProfilesResponse{
				ListRevertProfilesResponse: d,
			}
		default:
			return nil, fmt.Errorf("unsupported data type for StandardRevertResponse: %T", d)
		}
	}
	return resp, nil
}
