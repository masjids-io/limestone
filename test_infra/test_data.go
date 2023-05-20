package test_infra

import (
	"time"

	epb "github.com/mnadev/limestone/event_service/proto"
	mpb "github.com/mnadev/limestone/masjid_service/proto"
	upb "github.com/mnadev/limestone/user_service/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
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

func GetUserProto(email string, username string) *upb.User {
	return &upb.User{
		UserId:      DefaultId,
		Email:       email,
		Username:    username,
		FirstName:   FirstName,
		LastName:    LastName,
		PhoneNumber: PhoneNumber,
		Gender:      upb.User_FEMALE,
	}
}

func GetMasjidProto() *mpb.Masjid {
	return &mpb.Masjid{
		Id:         DefaultId,
		Name:       MasjidName,
		IsVerified: false,
		Address: &mpb.Masjid_Address{
			AddressLine_1: "123 Maple Ave",
			ZoneCode:      "TX",
			PostalCode:    "12345",
			City:          "Springfield",
			CountryCode:   "US",
		},
		PhoneNumber: &mpb.Masjid_PhoneNumber{
			CountryCode: "+1",
			Number:      "111-111-1111",
			Extension:   "1111",
		},
		PrayerConfig: &mpb.PrayerTimesConfiguration{
			Method:           mpb.PrayerTimesConfiguration_NORTH_AMERICA,
			FajrAngle:        15,
			IshaAngle:        10,
			IshaInterval:     0,
			AsrMethod:        mpb.PrayerTimesConfiguration_SHAFI_HANBALI_MALIKI,
			HighLatitudeRule: mpb.PrayerTimesConfiguration_NO_HIGH_LATITUDE_RULE,
			Adjustments: &mpb.PrayerTimesConfiguration_PrayerAdjustments{
				FajrAdjustment:    -2,
				DhuhrAdjustment:   -1,
				AsrAdjustment:     0,
				MaghribAdjustment: 1,
				IshaAdjustment:    2,
			},
		},
	}
}

func GetEventProto() *epb.Event {
	return &epb.Event{
		Id:                DefaultId,
		Owner:             &epb.Event_MasjidId{MasjidId: DefaultId},
		Name:              EventName,
		Description:       EventDescription,
		StartTime:         timestamppb.New(time.Date(2020, 10, 20, 20, 20, 20, 20, time.UTC)),
		EndTime:           timestamppb.New(time.Date(2020, 10, 20, 20, 20, 20, 20, time.UTC)),
		GenderRestriction: epb.Event_MALE_ONLY,
		IsPaid:            true,
		RequiresRsvp:      true,
		MaxParticipants:   500,
		LivestreamLink:    LivestreamLink,
	}
}
