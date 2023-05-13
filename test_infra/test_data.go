package test_infra

import (
	mpb "github.com/mnadev/limestone/masjid_service/proto"
	upb "github.com/mnadev/limestone/user_service/proto"
)

const (
	DefaultId   = "00000000-0000-0000-0000-000000000000"
	UserEmail   = "example@example.com"
	Password    = "password"
	BadPassword = "passwor"
	Username    = "coolguy1234"
	FirstName   = "John"
	LastName    = "Doe"
	PhoneNumber = "+1234567890"
	MasjidName  = "Masjid 1"
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
