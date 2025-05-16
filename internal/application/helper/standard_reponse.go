package helper

import (
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

func StandardMasjidResponse(code codes.Code, statusMessage string, message string, masjidEntity *entity.Masjid, listResponse *pb.ListMasjidsResponse, deleteResponse *pb.DeleteMasjidResponse) (*pb.StandardMasjidResponse, error) {
	resp := &pb.StandardMasjidResponse{
		Code:    code.String(),
		Status:  statusMessage,
		Message: message,
	}

	if masjidEntity != nil {
		protoMasjid := &pb.Masjid{
			Id:         masjidEntity.ID.String(),
			Name:       masjidEntity.Name,
			IsVerified: masjidEntity.IsVerified,
			Address: &pb.Masjid_Address{
				AddressLine_1: masjidEntity.Address.AddressLine1,
				AddressLine_2: masjidEntity.Address.AddressLine2,
				ZoneCode:      masjidEntity.Address.ZoneCode,
				PostalCode:    masjidEntity.Address.PostalCode,
				City:          masjidEntity.Address.City,
				CountryCode:   masjidEntity.Address.CountryCode,
			},
			PhoneNumber: &pb.Masjid_PhoneNumber{
				CountryCode: masjidEntity.PhoneNumber.PhoneCountryCode,
				Number:      masjidEntity.PhoneNumber.Number,
				Extension:   masjidEntity.PhoneNumber.Extension,
			},
			CreateTime: timestamppb.New(masjidEntity.CreatedAt),
			UpdateTime: timestamppb.New(masjidEntity.UpdatedAt),
		}
		resp.Data = &pb.StandardMasjidResponse_Masjid{Masjid: protoMasjid}
	} else if listResponse != nil {
		resp.Data = &pb.StandardMasjidResponse_ListMasjidResponse{ListMasjidResponse: listResponse}
	} else if deleteResponse != nil {
		resp.Data = &pb.StandardMasjidResponse_DeleteMasjidResponse{DeleteMasjidResponse: deleteResponse}
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
