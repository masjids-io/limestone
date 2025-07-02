package helper

import (
	"fmt"
	"github.com/google/uuid"
	pb "github.com/mnadev/limestone/gen/go"
	"github.com/mnadev/limestone/internal/application/domain/entity"
	"google.golang.org/protobuf/types/known/timestamppb"
	"time"
)

func ToEntityNikkahProfile(protoProfile *pb.NikkahProfile) (*entity.NikkahProfile, error) {
	if protoProfile == nil {
		return nil, fmt.Errorf("proto profile cannot be nil")
	}

	profileID, err := uuid.Parse(protoProfile.GetId())
	if err != nil && protoProfile.GetId() != "" {
		return nil, fmt.Errorf("invalid profile ID format: %w", err)
	}
	if protoProfile.GetId() == "" {
		profileID = uuid.Nil
	}

	var gender entity.Gender
	switch protoProfile.GetGender() {
	case pb.NikkahProfile_MALE:
		gender = entity.Male
	case pb.NikkahProfile_FEMALE:
		gender = entity.Female
	default:
		gender = entity.GenderUnspecified
	}

	var birthDate entity.BirthDate
	if protoProfile.GetBirthDate() != nil {
		birthDate = entity.BirthDate{
			Year:  protoProfile.GetBirthDate().GetYear(),
			Month: entity.Month(int32(protoProfile.GetBirthDate().GetMonth())),
			Day:   int8(protoProfile.GetBirthDate().GetDay()),
		}
	}

	var createdAt, updatedAt time.Time
	if protoProfile.GetCreateTime() != nil {
		createdAt = protoProfile.GetCreateTime().AsTime()
	}
	if protoProfile.GetUpdateTime() != nil {
		updatedAt = protoProfile.GetUpdateTime().AsTime()
	}

	return &entity.NikkahProfile{
		ID:        profileID,
		UserID:    protoProfile.GetUserId(),
		Name:      protoProfile.GetName(),
		BirthDate: birthDate,
		Gender:    gender,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}, nil
}

func ToProtoNikkahProfile(eProfile *entity.NikkahProfile) *pb.NikkahProfile {
	if eProfile == nil {
		return nil
	}

	protoProfile := &pb.NikkahProfile{
		Id:     eProfile.ID.String(),
		UserId: eProfile.UserID,
		Name:   eProfile.Name,
		BirthDate: &pb.NikkahProfile_BirthDate{
			Year:  eProfile.BirthDate.Year,
			Month: pb.NikkahProfile_BirthDate_Month(eProfile.BirthDate.Month),
			Day:   int32(eProfile.BirthDate.Day),
		},
		CreateTime: timestamppb.New(eProfile.CreatedAt),
		UpdateTime: timestamppb.New(eProfile.UpdatedAt),
	}

	switch eProfile.Gender {
	case entity.Male:
		protoProfile.Gender = pb.NikkahProfile_MALE
	case entity.Female:
		protoProfile.Gender = pb.NikkahProfile_FEMALE
	default:
		protoProfile.Gender = pb.NikkahProfile_GENDER_UNSPECIFIED
	}
	return protoProfile
}

func ToProtoNikkahLike(e *entity.NikkahLike) *pb.NikkahLike {
	if e == nil {
		return nil
	}

	return &pb.NikkahLike{
		LikeId:         e.ID.String(),
		LikerProfileId: e.LikerProfileID,
		LikedProfileId: e.LikedProfileID,
		Status:         pb.NikkahLike_Status(e.Status),
		CreateTime:     timestamppb.New(e.CreatedAt),
		UpdateTime:     timestamppb.New(e.UpdatedAt),
	}
}

func ToEntityNikkahLike(p *pb.NikkahLike) (*entity.NikkahLike, error) {
	if p == nil {
		return nil, nil
	}

	id, err := uuid.Parse(p.GetLikeId())
	if err != nil && p.GetLikeId() != "" {
		return nil, fmt.Errorf("invalid like ID format: %w", err)
	}

	likedProfileID, err := uuid.Parse(p.GetLikedProfileId())
	if err != nil && p.GetLikedProfileId() != "" {
		return nil, fmt.Errorf("invalid liked profile ID format: %w", err)
	}

	return &entity.NikkahLike{
		ID:             id,
		LikerProfileID: p.GetLikerProfileId(),
		LikedProfileID: likedProfileID.String(),
		Status:         entity.LikeStatus(p.GetStatus()),
		CreatedAt:      p.GetCreateTime().AsTime(),
		UpdatedAt:      p.GetUpdateTime().AsTime(),
	}, nil
}

func ToProtoNikkahMatch(e *entity.NikkahMatch) *pb.NikkahMatch {
	if e == nil {
		return nil
	}

	var protoStatus pb.NikkahMatch_Status
	switch e.Status {
	case entity.MatchStatusUnspecified:
		protoStatus = pb.NikkahMatch_STATUS_UNSPECIFIED
	case entity.MatchStatusInitiated:
		protoStatus = pb.NikkahMatch_INITIATED
	case entity.MatchStatusAccepted:
		protoStatus = pb.NikkahMatch_ACCEPTED
	case entity.MatchStatusRejected:
		protoStatus = pb.NikkahMatch_REJECTED
	case entity.MatchStatusEnded:
		protoStatus = pb.NikkahMatch_ENDED
	default:
		protoStatus = pb.NikkahMatch_STATUS_UNSPECIFIED
	}

	return &pb.NikkahMatch{
		MatchId:            e.ID.String(),
		InitiatorProfileId: e.InitiatorProfileID.String(),
		ReceiverProfileId:  e.ReceiverProfileID.String(),
		Status:             protoStatus,
		CreateTime:         timestamppb.New(e.CreatedAt),
		UpdateTime:         timestamppb.New(e.UpdatedAt),
	}
}
