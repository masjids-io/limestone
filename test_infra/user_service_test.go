package test_infra

import (
	"context"
	"testing"

	userservicepb "github.com/mnadev/limestone/user_service/proto"
	"github.com/stretchr/testify/suite"
)

func TestSuite(t *testing.T) {
	suite.Run(t, new(IntegrationTestSuite))
}

func (suite *IntegrationTestSuite) TestCreateUser_Success() {
	ctx := context.Background()
	out, err := suite.UserServiceClient.CreateUser(ctx, &userservicepb.CreateUserRequest{
		User:     GetUserProto(UserEmail, Username),
		Password: Password,
	})

	AssertUserProtoEqual(suite.T(), GetUserProto(UserEmail, Username), out)
	suite.Nil(err)
}

func (suite *IntegrationTestSuite) TestCreateUser_PasswordTooShort() {
	ctx := context.Background()
	out, err := suite.UserServiceClient.CreateUser(ctx, &userservicepb.CreateUserRequest{
		User:     GetUserProto("a@example.com", "user1"),
		Password: BadPassword,
	})

	suite.Nil(out)
	suite.Error(err)
}

func (suite *IntegrationTestSuite) TestUpdateUserWithEmail_Success() {
	ctx := context.Background()
	out, err := suite.UserServiceClient.CreateUser(ctx, &userservicepb.CreateUserRequest{
		User:     GetUserProto("a@example.com", "user2"),
		Password: Password,
	})

	AssertUserProtoEqual(suite.T(), GetUserProto("a@example.com", "user2"), out)
	suite.Nil(err)

	out, err = suite.UserServiceClient.UpdateUser(ctx, &userservicepb.UpdateUserRequest{
		User:     GetUserProto("b@example.com", "user2"),
		Password: Password,
	})

	AssertUserProtoEqual(suite.T(), GetUserProto("b@example.com", "user2"), out)
	suite.Nil(err)
}

func (suite *IntegrationTestSuite) TestUpdateUserWithEmail_BadPassword() {
	ctx := context.Background()
	out, err := suite.UserServiceClient.CreateUser(ctx, &userservicepb.CreateUserRequest{
		User:     GetUserProto("c@example.com", "user3"),
		Password: Password,
	})

	AssertUserProtoEqual(suite.T(), GetUserProto("c@example.com", "user3"), out)
	suite.Nil(err)

	out, err = suite.UserServiceClient.UpdateUser(ctx, &userservicepb.UpdateUserRequest{
		User:     GetUserProto("d@example.com", "user3"),
		Password: BadPassword,
	})

	suite.Error(err)
	suite.Nil(out)
}

func (suite *IntegrationTestSuite) TestUpdateUserWithEmail_NotFound() {
	ctx := context.Background()

	user := GetUserProto("e@example.com", "user4")
	user.UserId = "00000000-0000-0000-0000-000000000001"

	out, err := suite.UserServiceClient.UpdateUser(ctx, &userservicepb.UpdateUserRequest{
		User:     user,
		Password: Password,
	})

	suite.Error(err)
	suite.Nil(out)
}

func (suite *IntegrationTestSuite) TestGetUserWithEmail_Success() {
	ctx := context.Background()
	out, err := suite.UserServiceClient.CreateUser(ctx, &userservicepb.CreateUserRequest{
		User:     GetUserProto("f@example.com", "user5"),
		Password: Password,
	})

	AssertUserProtoEqual(suite.T(), GetUserProto("f@example.com", "user5"), out)
	suite.Nil(err)

	out, err = suite.UserServiceClient.GetUser(ctx, &userservicepb.GetUserRequest{
		Id: &userservicepb.GetUserRequest_Email{
			Email: "f@example.com",
		},
		Password: Password,
	})

	AssertUserProtoEqual(suite.T(), GetUserProto("f@example.com", "user5"), out)
	suite.Nil(err)
}

func (suite *IntegrationTestSuite) TestGetUserWithEmail_BadPassword() {
	ctx := context.Background()
	out, err := suite.UserServiceClient.CreateUser(ctx, &userservicepb.CreateUserRequest{
		User:     GetUserProto("g@example.com", "user6"),
		Password: Password,
	})

	AssertUserProtoEqual(suite.T(), GetUserProto("g@example.com", "user6"), out)
	suite.Nil(err)

	out, err = suite.UserServiceClient.GetUser(ctx, &userservicepb.GetUserRequest{
		Id: &userservicepb.GetUserRequest_Email{
			Email: "g@example.com",
		},
		Password: BadPassword,
	})

	suite.Error(err)
	suite.Nil(out)
}

func (suite *IntegrationTestSuite) TestGetUserWithEmail_NotFound() {
	ctx := context.Background()
	out, err := suite.UserServiceClient.GetUser(ctx, &userservicepb.GetUserRequest{
		Id: &userservicepb.GetUserRequest_Email{
			Email: "h@example.com",
		},
		Password: Password,
	})

	suite.Error(err)
	suite.Nil(out)
}

func (suite *IntegrationTestSuite) TestGetUserWithUsername_Success() {
	ctx := context.Background()
	out, err := suite.UserServiceClient.CreateUser(ctx, &userservicepb.CreateUserRequest{
		User:     GetUserProto("i@example.com", "user7"),
		Password: Password,
	})

	AssertUserProtoEqual(suite.T(), GetUserProto("i@example.com", "user7"), out)
	suite.Nil(err)

	out, err = suite.UserServiceClient.GetUser(ctx, &userservicepb.GetUserRequest{
		Id: &userservicepb.GetUserRequest_Username{
			Username: "user7",
		},
		Password: Password,
	})

	AssertUserProtoEqual(suite.T(), GetUserProto("i@example.com", "user7"), out)
	suite.Nil(err)
}

func (suite *IntegrationTestSuite) TestGetUserWithUsername_BadPassword() {
	ctx := context.Background()
	out, err := suite.UserServiceClient.CreateUser(ctx, &userservicepb.CreateUserRequest{
		User:     GetUserProto("j@example.com", "user8"),
		Password: Password,
	})

	AssertUserProtoEqual(suite.T(), GetUserProto("j@example.com", "user8"), out)
	suite.Nil(err)

	out, err = suite.UserServiceClient.GetUser(ctx, &userservicepb.GetUserRequest{
		Id: &userservicepb.GetUserRequest_Username{
			Username: "user8",
		},
		Password: BadPassword,
	})

	suite.Error(err)
	suite.Nil(out)
}

func (suite *IntegrationTestSuite) TestGetUserWithUsername_NotFound() {
	ctx := context.Background()
	out, err := suite.UserServiceClient.GetUser(ctx, &userservicepb.GetUserRequest{
		Id: &userservicepb.GetUserRequest_Username{
			Username: "user9",
		},
		Password: Password,
	})

	suite.Error(err)
	suite.Nil(out)
}

func (suite *IntegrationTestSuite) TestDeleteUserWithEmail_Success() {
	ctx := context.Background()
	out, err := suite.UserServiceClient.CreateUser(ctx, &userservicepb.CreateUserRequest{
		User:     GetUserProto("k@example.com", "user10"),
		Password: Password,
	})

	AssertUserProtoEqual(suite.T(), GetUserProto("k@example.com", "user10"), out)
	suite.Nil(err)

	_, err = suite.UserServiceClient.DeleteUser(ctx, &userservicepb.DeleteUserRequest{
		Id: &userservicepb.DeleteUserRequest_Email{
			Email: "k@example.com",
		},
		Password: Password,
	})

	suite.Nil(err)
}

func (suite *IntegrationTestSuite) TestDeleteUserWithEmail_BadPassword() {
	ctx := context.Background()
	out, err := suite.UserServiceClient.CreateUser(ctx, &userservicepb.CreateUserRequest{
		User:     GetUserProto("l@example.com", "user11"),
		Password: Password,
	})

	AssertUserProtoEqual(suite.T(), GetUserProto("l@example.com", "user11"), out)
	suite.Nil(err)

	_, err = suite.UserServiceClient.DeleteUser(ctx, &userservicepb.DeleteUserRequest{
		Id: &userservicepb.DeleteUserRequest_Email{
			Email: "l@example.com",
		},
		Password: BadPassword,
	})

	suite.Error(err)
}

func (suite *IntegrationTestSuite) TestDeleteUserWithEmail_NotFound() {
	ctx := context.Background()
	_, err := suite.UserServiceClient.DeleteUser(ctx, &userservicepb.DeleteUserRequest{
		Id: &userservicepb.DeleteUserRequest_Email{
			Email: "m@example.com",
		},
		Password: Password,
	})

	suite.Error(err)
}

func (suite *IntegrationTestSuite) TestDeleteUserWithUsername_Success() {
	ctx := context.Background()
	out, err := suite.UserServiceClient.CreateUser(ctx, &userservicepb.CreateUserRequest{
		User:     GetUserProto("n@example.com", "user12"),
		Password: Password,
	})

	AssertUserProtoEqual(suite.T(), GetUserProto("n@example.com", "user12"), out)
	suite.Nil(err)

	_, err = suite.UserServiceClient.DeleteUser(ctx, &userservicepb.DeleteUserRequest{
		Id: &userservicepb.DeleteUserRequest_Username{
			Username: "user12",
		},
		Password: Password,
	})

	suite.Nil(err)
}

func (suite *IntegrationTestSuite) TestDeleteUserWithUsername_BadPassword() {
	ctx := context.Background()
	out, err := suite.UserServiceClient.CreateUser(ctx, &userservicepb.CreateUserRequest{
		User:     GetUserProto("o@example.com", "user13"),
		Password: Password,
	})

	AssertUserProtoEqual(suite.T(), GetUserProto("o@example.com", "user13"), out)
	suite.Nil(err)

	_, err = suite.UserServiceClient.DeleteUser(ctx, &userservicepb.DeleteUserRequest{
		Id: &userservicepb.DeleteUserRequest_Username{
			Username: "user13",
		},
		Password: BadPassword,
	})

	suite.Error(err)
}

func (suite *IntegrationTestSuite) TestDeleteUserWithUsername_NotFound() {
	ctx := context.Background()
	_, err := suite.UserServiceClient.DeleteUser(ctx, &userservicepb.DeleteUserRequest{
		Id: &userservicepb.DeleteUserRequest_Username{
			Username: "user14",
		},
		Password: Password,
	})

	suite.Error(err)
}
