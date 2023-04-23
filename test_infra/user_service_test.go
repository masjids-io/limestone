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
		User:     GetUserProto(UserEmail, Username),
		Password: BadPassword,
	})

	suite.Nil(out)
	suite.Error(err)
}

func (suite *IntegrationTestSuite) TestUpdateUser_Success() {
	ctx := context.Background()
	out, err := suite.UserServiceClient.CreateUser(ctx, &userservicepb.CreateUserRequest{
		User:     GetUserProto(UserEmail, Username),
		Password: Password,
	})

	AssertUserProtoEqual(suite.T(), GetUserProto(UserEmail, Username), out)
	suite.Nil(err)

	out, err = suite.UserServiceClient.UpdateUser(ctx, &userservicepb.UpdateUserRequest{
		User:     GetUserProto(UserEmail, "new_name"),
		Password: Password,
	})

	AssertUserProtoEqual(suite.T(), GetUserProto(UserEmail, "new_name"), out)
	suite.Nil(err)
}

func (suite *IntegrationTestSuite) TestUpdateUser_BadPassword() {
	ctx := context.Background()
	out, err := suite.UserServiceClient.CreateUser(ctx, &userservicepb.CreateUserRequest{
		User:     GetUserProto(UserEmail, Username),
		Password: Password,
	})

	AssertUserProtoEqual(suite.T(), GetUserProto(UserEmail, Username), out)
	suite.Nil(err)

	out, err = suite.UserServiceClient.UpdateUser(ctx, &userservicepb.UpdateUserRequest{
		User:     GetUserProto(UserEmail, "new_name"),
		Password: BadPassword,
	})

	suite.Error(err)
	suite.Nil(out)
}

func (suite *IntegrationTestSuite) TestUpdateUser_NotFound() {
	ctx := context.Background()

	user := GetUserProto(UserEmail, Username)

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
		User:     GetUserProto(UserEmail, Username),
		Password: Password,
	})

	AssertUserProtoEqual(suite.T(), GetUserProto(UserEmail, Username), out)
	suite.Nil(err)

	out, err = suite.UserServiceClient.GetUser(ctx, &userservicepb.GetUserRequest{
		Id: &userservicepb.GetUserRequest_Email{
			Email: out.GetEmail(),
		},
		Password: Password,
	})

	AssertUserProtoEqual(suite.T(), GetUserProto(UserEmail, Username), out)
	suite.Nil(err)
}

func (suite *IntegrationTestSuite) TestGetUserWithEmail_BadPassword() {
	ctx := context.Background()
	out, err := suite.UserServiceClient.CreateUser(ctx, &userservicepb.CreateUserRequest{
		User:     GetUserProto(UserEmail, Username),
		Password: Password,
	})

	AssertUserProtoEqual(suite.T(), GetUserProto(UserEmail, Username), out)
	suite.Nil(err)

	out, err = suite.UserServiceClient.GetUser(ctx, &userservicepb.GetUserRequest{
		Id: &userservicepb.GetUserRequest_Email{
			Email: out.GetEmail(),
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
			Email: "bad@example.com",
		},
		Password: Password,
	})

	suite.Error(err)
	suite.Nil(out)
}

func (suite *IntegrationTestSuite) TestGetUserWithUsername_Success() {
	ctx := context.Background()
	out, err := suite.UserServiceClient.CreateUser(ctx, &userservicepb.CreateUserRequest{
		User:     GetUserProto(UserEmail, Username),
		Password: Password,
	})

	AssertUserProtoEqual(suite.T(), GetUserProto(UserEmail, Username), out)
	suite.Nil(err)

	out, err = suite.UserServiceClient.GetUser(ctx, &userservicepb.GetUserRequest{
		Id: &userservicepb.GetUserRequest_Username{
			Username: out.GetUsername(),
		},
		Password: Password,
	})

	AssertUserProtoEqual(suite.T(), GetUserProto(UserEmail, Username), out)
	suite.Nil(err)
}

func (suite *IntegrationTestSuite) TestGetUserWithUsername_BadPassword() {
	ctx := context.Background()
	out, err := suite.UserServiceClient.CreateUser(ctx, &userservicepb.CreateUserRequest{
		User:     GetUserProto(UserEmail, Username),
		Password: Password,
	})

	AssertUserProtoEqual(suite.T(), GetUserProto(UserEmail, Username), out)
	suite.Nil(err)

	out, err = suite.UserServiceClient.GetUser(ctx, &userservicepb.GetUserRequest{
		Id: &userservicepb.GetUserRequest_Username{
			Username: out.GetUsername(),
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
			Username: "bad_user",
		},
		Password: Password,
	})

	suite.Error(err)
	suite.Nil(out)
}

func (suite *IntegrationTestSuite) TestDeleteUserWithEmail_Success() {
	ctx := context.Background()
	out, err := suite.UserServiceClient.CreateUser(ctx, &userservicepb.CreateUserRequest{
		User:     GetUserProto(UserEmail, Username),
		Password: Password,
	})

	AssertUserProtoEqual(suite.T(), GetUserProto(UserEmail, Username), out)
	suite.Nil(err)

	_, err = suite.UserServiceClient.DeleteUser(ctx, &userservicepb.DeleteUserRequest{
		Id: &userservicepb.DeleteUserRequest_Email{
			Email: out.GetEmail(),
		},
		Password: Password,
	})

	suite.Nil(err)
}

func (suite *IntegrationTestSuite) TestDeleteUserWithEmail_BadPassword() {
	ctx := context.Background()
	out, err := suite.UserServiceClient.CreateUser(ctx, &userservicepb.CreateUserRequest{
		User:     GetUserProto(UserEmail, Username),
		Password: Password,
	})

	AssertUserProtoEqual(suite.T(), GetUserProto(UserEmail, Username), out)
	suite.Nil(err)

	_, err = suite.UserServiceClient.DeleteUser(ctx, &userservicepb.DeleteUserRequest{
		Id: &userservicepb.DeleteUserRequest_Email{
			Email: out.GetEmail(),
		},
		Password: BadPassword,
	})

	suite.Error(err)
}

func (suite *IntegrationTestSuite) TestDeleteUserWithEmail_NotFound() {
	ctx := context.Background()
	_, err := suite.UserServiceClient.DeleteUser(ctx, &userservicepb.DeleteUserRequest{
		Id: &userservicepb.DeleteUserRequest_Email{
			Email: "bad@example.com",
		},
		Password: Password,
	})

	suite.Error(err)
}

func (suite *IntegrationTestSuite) TestDeleteUserWithUsername_Success() {
	ctx := context.Background()
	out, err := suite.UserServiceClient.CreateUser(ctx, &userservicepb.CreateUserRequest{
		User:     GetUserProto(UserEmail, Username),
		Password: Password,
	})

	AssertUserProtoEqual(suite.T(), GetUserProto(UserEmail, Username), out)
	suite.Nil(err)

	_, err = suite.UserServiceClient.DeleteUser(ctx, &userservicepb.DeleteUserRequest{
		Id: &userservicepb.DeleteUserRequest_Username{
			Username: out.GetUsername(),
		},
		Password: Password,
	})

	suite.Nil(err)
}

func (suite *IntegrationTestSuite) TestDeleteUserWithUsername_BadPassword() {
	ctx := context.Background()
	out, err := suite.UserServiceClient.CreateUser(ctx, &userservicepb.CreateUserRequest{
		User:     GetUserProto(UserEmail, Username),
		Password: Password,
	})

	AssertUserProtoEqual(suite.T(), GetUserProto(UserEmail, Username), out)
	suite.Nil(err)

	_, err = suite.UserServiceClient.DeleteUser(ctx, &userservicepb.DeleteUserRequest{
		Id: &userservicepb.DeleteUserRequest_Username{
			Username: out.GetUsername(),
		},
		Password: BadPassword,
	})

	suite.Error(err)
}

func (suite *IntegrationTestSuite) TestDeleteUserWithUsername_NotFound() {
	ctx := context.Background()
	_, err := suite.UserServiceClient.DeleteUser(ctx, &userservicepb.DeleteUserRequest{
		Id: &userservicepb.DeleteUserRequest_Username{
			Username: "bad_user",
		},
		Password: Password,
	})

	suite.Error(err)
}
