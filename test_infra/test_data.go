package test_infra

import (
	userservicepb "github.com/mnadev/limestone/user_service/proto"
	_ "google.golang.org/protobuf/types/known/timestamppb"
)

const (
	UserEmail   = "example@example.com"
	Password    = "password"
	BadPassword = "passwor"
	Username    = "coolguy1234"
	FirstName   = "John"
	LastName    = "Doe"
	PhoneNumber = "+1234567890"
)

func GetUserProto(email string, username string) *userservicepb.User {
	return &userservicepb.User{
		UserId:      "00000000-0000-0000-0000-000000000000",
		Email:       email,
		Username:    username,
		FirstName:   FirstName,
		LastName:    LastName,
		PhoneNumber: PhoneNumber,
		Gender:      userservicepb.User_FEMALE,
	}
}
