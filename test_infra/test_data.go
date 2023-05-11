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
		Id:          DefaultId,
		Name:        MasjidName,
		IsVerified:  false,
		Address:     &mpb.Masjid_Address{},
		PhoneNumber: &mpb.Masjid_PhoneNumber{},
	}
}
