package helper

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	pb "github.com/mnadev/limestone/gen/go"
	"github.com/mnadev/limestone/internal/application/domain/entity"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func ToEntityRevertMatch(p *pb.RevertMatch) (*entity.RevertMatch, error) {
	if p == nil {
		return nil, errors.New("protobuf RevertMatch message is nil")
	}

	var matchID uuid.UUID
	if p.GetMatchId() != "" {
		parsedID, err := uuid.Parse(p.GetMatchId())
		if err != nil {
			return nil, fmt.Errorf("invalid MatchId format: %w", err)
		}
		matchID = parsedID
	}

	var status entity.RevertMatchStatus
	switch p.GetStatus() {
	case pb.RevertMatch_INITIATED:
		status = entity.RevertMatchStatus_INITIATED
	case pb.RevertMatch_ACCEPTED:
		status = entity.RevertMatchStatus_ACCEPTED
	case pb.RevertMatch_REJECTED:
		status = entity.RevertMatchStatus_REJECTED
	case pb.RevertMatch_ENDED:
		status = entity.RevertMatchStatus_ENDED
	case pb.RevertMatch_STATUS_UNSPECIFIED:
		status = entity.RevertMatchStatus_UNSPECIFIED
	default:
		return nil, fmt.Errorf("unknown RevertMatch status: %v", p.GetStatus())
	}

	return &entity.RevertMatch{
		ID:                 matchID,
		InitiatorProfileID: p.GetInitiatorProfileId(),
		ReceiverProfileID:  p.GetReceiverProfileId(),
		Status:             status,
		// CreatedAt dan UpdatedAt diisi oleh database
	}, nil
}

func ToProtoRevertMatch(e *entity.RevertMatch) *pb.RevertMatch {
	if e == nil {
		return nil
	}

	var status pb.RevertMatch_Status
	switch e.Status {
	case entity.RevertMatchStatus_INITIATED:
		status = pb.RevertMatch_INITIATED
	case entity.RevertMatchStatus_ACCEPTED:
		status = pb.RevertMatch_ACCEPTED
	case entity.RevertMatchStatus_REJECTED:
		status = pb.RevertMatch_REJECTED
	case entity.RevertMatchStatus_ENDED:
		status = pb.RevertMatch_ENDED
	case entity.RevertMatchStatus_UNSPECIFIED:
		status = pb.RevertMatch_STATUS_UNSPECIFIED
	default:
		status = pb.RevertMatch_STATUS_UNSPECIFIED
	}

	return &pb.RevertMatch{
		MatchId:            e.ID.String(),
		InitiatorProfileId: e.InitiatorProfileID,
		ReceiverProfileId:  e.ReceiverProfileID,
		Status:             status,
		CreateTime:         timestamppb.New(e.CreatedAt),
		UpdateTime:         timestamppb.New(e.UpdatedAt),
	}
}

func ToEntityRevertProfile(p *pb.RevertProfile) (*entity.RevertProfile, error) {
	if p == nil {
		return nil, errors.New("protobuf RevertProfile message is nil")
	}

	var profileID uuid.UUID
	if p.GetId() != "" {
		parsedID, err := uuid.Parse(p.GetId())
		if err != nil {
			return nil, fmt.Errorf("invalid RevertProfile ID format: %w", err)
		}
		profileID = parsedID
	}

	var gender entity.RevertProfileGender
	switch p.GetGender() {
	case pb.RevertProfile_MALE:
		gender = entity.RevertProfileGender_MALE
	case pb.RevertProfile_FEMALE:
		gender = entity.RevertProfileGender_FEMALE
	case pb.RevertProfile_GENDER_UNSPECIFIED:
		gender = entity.RevertProfileGender_UNSPECIFIED
	default:
		return nil, fmt.Errorf("unknown RevertProfile gender: %v", p.GetGender())
	}

	var birthDate entity.BirthDate
	if p.GetBirthDate() != nil {
		birthDate = entity.BirthDate{
			Year:  p.GetBirthDate().GetYear(),
			Month: entity.Month(int32(p.GetBirthDate().GetMonth())),
			Day:   int8(p.GetBirthDate().GetDay()),
		}
	}

	return &entity.RevertProfile{
		ID:        profileID,
		UserID:    p.GetUserId(),
		Name:      p.GetName(),
		Gender:    gender,
		BirthDate: birthDate,
		// CreatedAt dan UpdatedAt tidak diset di sini karena diisi oleh DB
	}, nil
}

func ToProtoRevertProfile(e *entity.RevertProfile) *pb.RevertProfile {
	if e == nil {
		return nil
	}

	var gender pb.RevertProfile_Gender
	switch e.Gender {
	case entity.RevertProfileGender_MALE:
		gender = pb.RevertProfile_MALE
	case entity.RevertProfileGender_FEMALE:
		gender = pb.RevertProfile_FEMALE
	case entity.RevertProfileGender_UNSPECIFIED:
		gender = pb.RevertProfile_GENDER_UNSPECIFIED
	default:
		gender = pb.RevertProfile_GENDER_UNSPECIFIED
	}

	protoBirthDate := &pb.RevertProfile_BirthDate{
		Year:  int32(e.BirthDate.Year),
		Month: pb.RevertProfile_BirthDate_Month(e.BirthDate.Month), // Konversi int ke proto enum Month
		Day:   int32(e.BirthDate.Day),
	}

	return &pb.RevertProfile{
		Id:         e.ID.String(),
		UserId:     e.UserID,
		Name:       e.Name,
		Gender:     gender,
		BirthDate:  protoBirthDate,
		CreateTime: timestamppb.New(e.CreatedAt),
		UpdateTime: timestamppb.New(e.UpdatedAt),
	}
}
