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

	// Convert proto location to entity location
	var location entity.Location
	if protoProfile.GetLocation() != nil {
		location = entity.Location{
			Country:   protoProfile.GetLocation().GetCountry(),
			City:      protoProfile.GetLocation().GetCity(),
			State:     protoProfile.GetLocation().GetState(),
			ZipCode:   protoProfile.GetLocation().GetZipCode(),
			Latitude:  protoProfile.GetLocation().GetLatitude(),
			Longitude: protoProfile.GetLocation().GetLongitude(),
		}
	}

	// Convert proto education to entity education
	var education entity.Education
	if protoProfile.GetEducation() != pb.Education_EDUCATION_UNSPECIFIED {
		education = entity.Education(protoProfile.GetEducation())
	}

	// Convert proto height to entity height
	var height entity.Height
	if protoProfile.GetHeight() != nil {
		height = entity.Height{
			Cm: protoProfile.GetHeight().GetCm(),
		}
	}

	// Convert proto sect to entity sect
	var sect entity.Sect
	if protoProfile.GetSect() != pb.Sect_SECT_UNSPECIFIED {
		sect = entity.Sect(protoProfile.GetSect())
	}

	// Convert proto pictures to entity pictures
	var pictures []entity.Picture
	for _, pic := range protoProfile.GetPictures() {
		pictures = append(pictures, entity.Picture{
			Image:    pic.GetImage(),
			MimeType: pic.GetMimeType(),
		})
	}

	// Convert proto hobbies to entity hobbies
	var hobbies []entity.Hobbies
	for _, hobby := range protoProfile.GetHobbies() {
		if hobby != pb.Hobbies_HOBBIES_UNSPECIFIED {
			hobbies = append(hobbies, entity.Hobbies(hobby))
		}
	}

	return &entity.NikkahProfile{
		ID:         profileID,
		UserID:     protoProfile.GetUserId(),
		Name:       protoProfile.GetName(),
		BirthDate:  birthDate,
		Gender:     gender,
		Location:   location,
		Education:  education,
		Occupation: protoProfile.GetOccupation(),
		Height:     height,
		Sect:       sect,
		Pictures:   pictures,
		Hobbies:    hobbies,
		CreatedAt:  createdAt,
		UpdatedAt:  updatedAt,
	}, nil
}

func ToProtoNikkahProfile(eProfile *entity.NikkahProfile) *pb.NikkahProfile {
	if eProfile == nil {
		return nil
	}

	// Convert entity location to proto location
	var protoLocation *pb.Location
	if eProfile.Location != (entity.Location{}) {
		protoLocation = &pb.Location{
			Country:   eProfile.Location.Country,
			City:      eProfile.Location.City,
			State:     eProfile.Location.State,
			ZipCode:   eProfile.Location.ZipCode,
			Latitude:  eProfile.Location.Latitude,
			Longitude: eProfile.Location.Longitude,
		}
	}

	// Convert entity education to proto education
	var protoEducation pb.Education
	if eProfile.Education != entity.EducationUnspecified {
		protoEducation = pb.Education(eProfile.Education)
	} else {
		protoEducation = pb.Education_EDUCATION_UNSPECIFIED
	}

	// Convert entity height to proto height
	var protoHeight *pb.Height
	if eProfile.Height != (entity.Height{}) {
		protoHeight = &pb.Height{
			Cm: eProfile.Height.Cm,
		}
	}

	// Convert entity sect to proto sect
	var protoSect pb.Sect
	if eProfile.Sect != entity.SectUnspecified {
		protoSect = pb.Sect(eProfile.Sect)
	} else {
		protoSect = pb.Sect_SECT_UNSPECIFIED
	}

	// Convert entity pictures to proto pictures
	var protoPictures []*pb.Picture
	for _, pic := range eProfile.Pictures {
		protoPictures = append(protoPictures, &pb.Picture{
			Image:    pic.Image,
			MimeType: pic.MimeType,
		})
	}

	// Convert entity hobbies to proto hobbies
	var protoHobbies []pb.Hobbies
	for _, hobby := range eProfile.Hobbies {
		if hobby != entity.HobbiesUnspecified {
			protoHobbies = append(protoHobbies, pb.Hobbies(hobby))
		}
	}

	protoProfile := &pb.NikkahProfile{
		Id:         eProfile.ID.String(),
		UserId:     eProfile.UserID,
		Name:       eProfile.Name,
		Location:   protoLocation,
		Education:  protoEducation,
		Occupation: eProfile.Occupation,
		Height:     protoHeight,
		Sect:       protoSect,
		Pictures:   protoPictures,
		Hobbies:    protoHobbies,
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
