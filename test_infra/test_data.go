package test_infra

import (
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"

	pb "github.com/mnadev/limestone/proto"
)

const (
	DefaultId        = "00000000-0000-0000-0000-000000000000"
	UserEmail        = "example@example.com"
	Password         = "password"
	BadPassword      = "passwor"
	Username         = "coolguy1234"
	FirstName        = "John"
	LastName         = "Doe"
	PhoneNumber      = "+1234567890"
	MasjidName       = "Masjid 1"
	EventName        = "Event 1"
	EventDescription = "Some Event"
	LivestreamLink   = "http://example.com"
)

func GetUserProto(email string, username string) *pb.User {
	return &pb.User{
		Id:          DefaultId,
		Email:       email,
		Username:    username,
		FirstName:   FirstName,
		LastName:    LastName,
		PhoneNumber: PhoneNumber,
		Gender:      pb.User_FEMALE,
	}
}

func GetMasjidProto() *pb.Masjid {
	return &pb.Masjid{
		Id:         DefaultId,
		Name:       MasjidName,
		IsVerified: false,
		Address: &pb.Masjid_Address{
			AddressLine_1: "123 Maple Ave",
			ZoneCode:      "TX",
			PostalCode:    "12345",
			City:          "Springfield",
			CountryCode:   "US",
		},
		PhoneNumber: &pb.Masjid_PhoneNumber{
			CountryCode: "+1",
			Number:      "111-111-1111",
			Extension:   "1111",
		},
		PrayerConfig: &pb.PrayerTimesConfiguration{
			Method:           pb.PrayerTimesConfiguration_NORTH_AMERICA,
			FajrAngle:        15,
			IshaAngle:        10,
			IshaInterval:     0,
			AsrMethod:        pb.PrayerTimesConfiguration_SHAFI_HANBALI_MALIKI,
			HighLatitudeRule: pb.PrayerTimesConfiguration_NO_HIGH_LATITUDE_RULE,
			Adjustments: &pb.PrayerTimesConfiguration_PrayerAdjustments{
				FajrAdjustment:    -2,
				DhuhrAdjustment:   -1,
				AsrAdjustment:     0,
				MaghribAdjustment: 1,
				IshaAdjustment:    2,
			},
		},
	}
}

func GetEventProto() *pb.Event {
	return &pb.Event{
		Id:             DefaultId,
		OrganizationId: DefaultId,
		Name:           EventName,
		Description:    EventDescription,
		StartTime:      timestamppb.New(time.Date(2020, 10, 20, 20, 20, 20, 20, time.UTC)),
		EndTime:        timestamppb.New(time.Date(2020, 10, 20, 20, 20, 20, 20, time.UTC)),
		Types: []pb.Event_EventType{
			pb.Event_COMMUNITY,
			pb.Event_ATHLETIC,
		},
		GenderRestriction: pb.Event_MALE_ONLY,
		IsPaid:            true,
		RequiresRsvp:      true,
		MaxParticipants:   500,
		LivestreamLink:    LivestreamLink,
	}
}

func GetAdhanFileProto() *pb.AdhanFile {
	return &pb.AdhanFile{
		Id:       DefaultId,
		MasjidId: DefaultId,
		File:     []byte("SomeData"),
	}
}

func GetNikkahProfileProto() *pb.NikkahProfile {
	return &pb.NikkahProfile{
		Id:     DefaultId,
		UserId: DefaultId,
		Name:   "John Doe",
		Gender: pb.NikkahProfile_MALE,
		BirthDate: &pb.NikkahProfile_BirthDate{
			Year:  1990,
			Month: pb.NikkahProfile_BirthDate_JANUARY,
			Day:   1,
		},
	}
}

func GetNikkahLikeProto() *pb.NikkahLike {
	return &pb.NikkahLike{
		LikeId:         DefaultId,
		LikerProfileId: DefaultId,
		LikedProfileId: DefaultId,
		Status:         pb.NikkahLike_INITIATED,
	}
}

func GetNikkahMatchProto() *pb.NikkahMatch {
	return &pb.NikkahMatch{
		MatchId:            DefaultId,
		InitiatorProfileId: DefaultId,
		ReceiverProfileId:  DefaultId,
		Status:             pb.NikkahMatch_INITIATED,
	}
}

func GetRevertProfileProto() *pb.RevertProfile {
	return &pb.RevertProfile{
		Id:     DefaultId,
		UserId: DefaultId,
		Name:   "John Doe",
		Gender: pb.RevertProfile_MALE,
		BirthDate: &pb.RevertProfile_BirthDate{
			Year:  1990,
			Month: pb.RevertProfile_BirthDate_JANUARY,
			Day:   1,
		},
	}
}

func GetRevertMatchProto() *pb.RevertMatch {
	return &pb.RevertMatch{
		MatchId:            DefaultId,
		InitiatorProfileId: DefaultId,
		ReceiverProfileId:  DefaultId,
		Status:             pb.RevertMatch_INITIATED,
	}
}
