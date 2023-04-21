package test_infra

import (
	"context"
	"testing"
	"time"

	userservicepb "github.com/mnadev/limestone/userservice/proto"
	"github.com/stretchr/testify/assert"
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

func AssertUserProtoEqual(t *testing.T, expected, actual *userservicepb.User) bool {
	actual_without_timestamp := userservicepb.User{
		UserId:      actual.UserId,
		Email:       actual.Email,
		Username:    actual.Username,
		FirstName:   actual.FirstName,
		LastName:    actual.LastName,
		PhoneNumber: actual.PhoneNumber,
		Gender:      actual.Gender,
	}

	assert.Equal(t, *expected, actual_without_timestamp)
	assert.LessOrEqual(t, time.Now().Unix()-actual.CreateTime.GetSeconds(), int64(1))
	assert.LessOrEqual(t, time.Now().Unix()-actual.UpdateTime.GetSeconds(), int64(1))

	return true
}

func TestCreateUser_Success(t *testing.T) {
	ctx := context.Background()
	client, closer := UserServiceClient(ctx)
	defer closer()

	out, err := client.CreateUser(ctx, &userservicepb.CreateUserRequest{
		User:     GetUserProto(UserEmail, Username),
		Password: Password,
	})

	AssertUserProtoEqual(t, GetUserProto(UserEmail, Username), out)
	assert.Nil(t, err)
}

func TestCreateUser_PasswordTooShort(t *testing.T) {
	ctx := context.Background()
	client, closer := UserServiceClient(ctx)
	defer closer()

	out, err := client.CreateUser(ctx, &userservicepb.CreateUserRequest{
		User:     GetUserProto("a@example.com", "user1"),
		Password: BadPassword,
	})

	assert.Nil(t, out)
	assert.Error(t, err)
}

func TestUpdateUserWithEmail_Success(t *testing.T) {
	ctx := context.Background()
	client, closer := UserServiceClient(ctx)
	defer closer()

	out, err := client.CreateUser(ctx, &userservicepb.CreateUserRequest{
		User:     GetUserProto("a@example.com", "user2"),
		Password: Password,
	})

	AssertUserProtoEqual(t, GetUserProto("a@example.com", "user2"), out)
	assert.Nil(t, err)

	out, err = client.UpdateUser(ctx, &userservicepb.UpdateUserRequest{
		User:     GetUserProto("b@example.com", "user2"),
		Password: Password,
	})

	AssertUserProtoEqual(t, GetUserProto("b@example.com", "user2"), out)
	assert.Nil(t, err)
}

func TestUpdateUserWithEmail_BadPassword(t *testing.T) {
	ctx := context.Background()
	client, closer := UserServiceClient(ctx)
	defer closer()

	out, err := client.CreateUser(ctx, &userservicepb.CreateUserRequest{
		User:     GetUserProto("c@example.com", "user3"),
		Password: Password,
	})

	AssertUserProtoEqual(t, GetUserProto("c@example.com", "user3"), out)
	assert.Nil(t, err)

	out, err = client.UpdateUser(ctx, &userservicepb.UpdateUserRequest{
		User:     GetUserProto("d@example.com", "user3"),
		Password: BadPassword,
	})

	assert.Error(t, err)
	assert.Nil(t, out)
}

func TestUpdateUserWithEmail_NotFound(t *testing.T) {
	ctx := context.Background()
	client, closer := UserServiceClient(ctx)
	defer closer()

	user := GetUserProto("e@example.com", "user4")
	user.UserId = "00000000-0000-0000-0000-000000000001"

	out, err := client.UpdateUser(ctx, &userservicepb.UpdateUserRequest{
		User:     user,
		Password: Password,
	})

	assert.Error(t, err)
	assert.Nil(t, out)
}

func TestGetUserWithEmail_Success(t *testing.T) {
	ctx := context.Background()
	client, closer := UserServiceClient(ctx)
	defer closer()

	out, err := client.CreateUser(ctx, &userservicepb.CreateUserRequest{
		User:     GetUserProto("f@example.com", "user5"),
		Password: Password,
	})

	AssertUserProtoEqual(t, GetUserProto("f@example.com", "user5"), out)
	assert.Nil(t, err)

	out, err = client.GetUser(ctx, &userservicepb.GetUserRequest{
		Id: &userservicepb.GetUserRequest_Email{
			Email: "f@example.com",
		},
		Password: Password,
	})

	AssertUserProtoEqual(t, GetUserProto("f@example.com", "user5"), out)
	assert.Nil(t, err)
}

func TestGetUserWithEmail_BadPassword(t *testing.T) {
	ctx := context.Background()
	client, closer := UserServiceClient(ctx)
	defer closer()

	out, err := client.CreateUser(ctx, &userservicepb.CreateUserRequest{
		User:     GetUserProto("g@example.com", "user6"),
		Password: Password,
	})

	AssertUserProtoEqual(t, GetUserProto("g@example.com", "user6"), out)
	assert.Nil(t, err)

	out, err = client.GetUser(ctx, &userservicepb.GetUserRequest{
		Id: &userservicepb.GetUserRequest_Email{
			Email: "g@example.com",
		},
		Password: BadPassword,
	})

	assert.Error(t, err)
	assert.Nil(t, out)
}

func TestGetUserWithEmail_NotFound(t *testing.T) {
	ctx := context.Background()
	client, closer := UserServiceClient(ctx)
	defer closer()

	out, err := client.GetUser(ctx, &userservicepb.GetUserRequest{
		Id: &userservicepb.GetUserRequest_Email{
			Email: "h@example.com",
		},
		Password: Password,
	})

	assert.Error(t, err)
	assert.Nil(t, out)
}

func TestGetUserWithUsername_Success(t *testing.T) {
	ctx := context.Background()
	client, closer := UserServiceClient(ctx)
	defer closer()

	out, err := client.CreateUser(ctx, &userservicepb.CreateUserRequest{
		User:     GetUserProto("i@example.com", "user7"),
		Password: Password,
	})

	AssertUserProtoEqual(t, GetUserProto("i@example.com", "user7"), out)
	assert.Nil(t, err)

	out, err = client.GetUser(ctx, &userservicepb.GetUserRequest{
		Id: &userservicepb.GetUserRequest_Username{
			Username: "user7",
		},
		Password: Password,
	})

	AssertUserProtoEqual(t, GetUserProto("i@example.com", "user7"), out)
	assert.Nil(t, err)
}

func TestGetUserWithUsername_BadPassword(t *testing.T) {
	ctx := context.Background()
	client, closer := UserServiceClient(ctx)
	defer closer()

	out, err := client.CreateUser(ctx, &userservicepb.CreateUserRequest{
		User:     GetUserProto("j@example.com", "user8"),
		Password: Password,
	})

	AssertUserProtoEqual(t, GetUserProto("j@example.com", "user8"), out)
	assert.Nil(t, err)

	out, err = client.GetUser(ctx, &userservicepb.GetUserRequest{
		Id: &userservicepb.GetUserRequest_Username{
			Username: "user8",
		},
		Password: BadPassword,
	})

	assert.Error(t, err)
	assert.Nil(t, out)
}

func TestGetUserWithUsername_NotFound(t *testing.T) {
	ctx := context.Background()
	client, closer := UserServiceClient(ctx)
	defer closer()

	out, err := client.GetUser(ctx, &userservicepb.GetUserRequest{
		Id: &userservicepb.GetUserRequest_Username{
			Username: "user9",
		},
		Password: Password,
	})

	assert.Error(t, err)
	assert.Nil(t, out)
}

func TestDeleteUserWithEmail_Success(t *testing.T) {
	ctx := context.Background()
	client, closer := UserServiceClient(ctx)
	defer closer()

	out, err := client.CreateUser(ctx, &userservicepb.CreateUserRequest{
		User:     GetUserProto("k@example.com", "user10"),
		Password: Password,
	})

	AssertUserProtoEqual(t, GetUserProto("k@example.com", "user10"), out)
	assert.Nil(t, err)

	_, err = client.DeleteUser(ctx, &userservicepb.DeleteUserRequest{
		Id: &userservicepb.DeleteUserRequest_Email{
			Email: "k@example.com",
		},
		Password: Password,
	})

	assert.Nil(t, err)
}

func TestDeleteUserWithEmail_BadPassword(t *testing.T) {
	ctx := context.Background()
	client, closer := UserServiceClient(ctx)
	defer closer()

	out, err := client.CreateUser(ctx, &userservicepb.CreateUserRequest{
		User:     GetUserProto("l@example.com", "user11"),
		Password: Password,
	})

	AssertUserProtoEqual(t, GetUserProto("l@example.com", "user11"), out)
	assert.Nil(t, err)

	_, err = client.DeleteUser(ctx, &userservicepb.DeleteUserRequest{
		Id: &userservicepb.DeleteUserRequest_Email{
			Email: "l@example.com",
		},
		Password: BadPassword,
	})

	assert.Error(t, err)
}

func TestDeleteUserWithEmail_NotFound(t *testing.T) {
	ctx := context.Background()
	client, closer := UserServiceClient(ctx)
	defer closer()

	_, err := client.DeleteUser(ctx, &userservicepb.DeleteUserRequest{
		Id: &userservicepb.DeleteUserRequest_Email{
			Email: "m@example.com",
		},
		Password: Password,
	})

	assert.Error(t, err)
}

func TestDeleteUserWithUsername_Success(t *testing.T) {
	ctx := context.Background()
	client, closer := UserServiceClient(ctx)
	defer closer()

	out, err := client.CreateUser(ctx, &userservicepb.CreateUserRequest{
		User:     GetUserProto("n@example.com", "user12"),
		Password: Password,
	})

	AssertUserProtoEqual(t, GetUserProto("n@example.com", "user12"), out)
	assert.Nil(t, err)

	_, err = client.DeleteUser(ctx, &userservicepb.DeleteUserRequest{
		Id: &userservicepb.DeleteUserRequest_Username{
			Username: "user12",
		},
		Password: Password,
	})

	assert.Nil(t, err)
}

func TestDeleteUserWithUsername_BadPassword(t *testing.T) {
	ctx := context.Background()
	client, closer := UserServiceClient(ctx)
	defer closer()

	out, err := client.CreateUser(ctx, &userservicepb.CreateUserRequest{
		User:     GetUserProto("o@example.com", "user13"),
		Password: Password,
	})

	AssertUserProtoEqual(t, GetUserProto("o@example.com", "user13"), out)
	assert.Nil(t, err)

	_, err = client.DeleteUser(ctx, &userservicepb.DeleteUserRequest{
		Id: &userservicepb.DeleteUserRequest_Username{
			Username: "user13",
		},
		Password: BadPassword,
	})

	assert.Error(t, err)
}

func TestDeleteUserWithUsername_NotFound(t *testing.T) {
	ctx := context.Background()
	client, closer := UserServiceClient(ctx)
	defer closer()

	_, err := client.DeleteUser(ctx, &userservicepb.DeleteUserRequest{
		Id: &userservicepb.DeleteUserRequest_Username{
			Username: "user14",
		},
		Password: Password,
	})

	assert.Error(t, err)
}
