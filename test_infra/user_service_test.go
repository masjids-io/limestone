package test_infra

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
	"google.golang.org/protobuf/testing/protocmp"

	pb "github.com/mnadev/limestone/proto"
)

func TestUserService(t *testing.T) {
	suite.Run(t, new(IntegrationTestSuite))
}

func (suite *IntegrationTestSuite) TestCreateUser_Success() {
	ctx := context.Background()
	out, err := suite.UserServiceClient.CreateUser(ctx, &pb.CreateUserRequest{
		User:     GetUserProto(UserEmail, Username),
		Password: Password,
	})

	suite.Nil(err)
	AssertProtoEqual(suite.T(), *GetUserProto(UserEmail, Username), *out,
		pb.User{}, protocmp.IgnoreFields(&pb.User{}, "create_time", "update_time"))
	AssertUserTimestampsCurrent(suite.T(), out)
}

func (suite *IntegrationTestSuite) TestCreateUser_PasswordTooShort() {
	ctx := context.Background()
	out, err := suite.UserServiceClient.CreateUser(ctx, &pb.CreateUserRequest{
		User:     GetUserProto(UserEmail, Username),
		Password: BadPassword,
	})

	suite.Nil(out)
	suite.Error(err)
}

func (suite *IntegrationTestSuite) TestUpdateUser_Success() {
	ctx := context.Background()
	out, err := suite.UserServiceClient.CreateUser(ctx, &pb.CreateUserRequest{
		User:     GetUserProto(UserEmail, Username),
		Password: Password,
	})

	suite.Nil(err)
	AssertProtoEqual(suite.T(), *GetUserProto(UserEmail, Username), *out,
		pb.User{}, protocmp.IgnoreFields(&pb.User{}, "create_time", "update_time"))
	AssertUserTimestampsCurrent(suite.T(), out)

	out, err = suite.UserServiceClient.UpdateUser(ctx, &pb.UpdateUserRequest{
		User:     GetUserProto(UserEmail, "new_name"),
		Password: Password,
	})

	suite.Nil(err)
	AssertProtoEqual(suite.T(), *GetUserProto(UserEmail, "new_name"), *out,
		pb.User{}, protocmp.IgnoreFields(&pb.User{}, "create_time", "update_time"))
	AssertUserTimestampsCurrent(suite.T(), out)
}

func (suite *IntegrationTestSuite) TestUpdateUser_BadPassword() {
	ctx := context.Background()
	out, err := suite.UserServiceClient.CreateUser(ctx, &pb.CreateUserRequest{
		User:     GetUserProto(UserEmail, Username),
		Password: Password,
	})

	suite.Nil(err)
	AssertProtoEqual(suite.T(), *GetUserProto(UserEmail, Username), *out,
		pb.User{}, protocmp.IgnoreFields(&pb.User{}, "create_time", "update_time"))
	AssertUserTimestampsCurrent(suite.T(), out)

	out, err = suite.UserServiceClient.UpdateUser(ctx, &pb.UpdateUserRequest{
		User:     GetUserProto(UserEmail, "new_name"),
		Password: BadPassword,
	})

	suite.Error(err)
	suite.Nil(out)
}

func (suite *IntegrationTestSuite) TestUpdateUser_NotFound() {
	ctx := context.Background()

	user := GetUserProto(UserEmail, Username)

	out, err := suite.UserServiceClient.UpdateUser(ctx, &pb.UpdateUserRequest{
		User:     user,
		Password: Password,
	})

	suite.Error(err)
	suite.Nil(out)
}

func (suite *IntegrationTestSuite) TestGetUserWithEmail_Success() {
	ctx := context.Background()
	out, err := suite.UserServiceClient.CreateUser(ctx, &pb.CreateUserRequest{
		User:     GetUserProto(UserEmail, Username),
		Password: Password,
	})

	suite.Nil(err)
	AssertProtoEqual(suite.T(), *GetUserProto(UserEmail, Username), *out,
		pb.User{}, protocmp.IgnoreFields(&pb.User{}, "create_time", "update_time"))
	AssertUserTimestampsCurrent(suite.T(), out)

	resp, err := suite.UserServiceClient.GetUser(ctx, &pb.GetUserRequest{
		Id: &pb.GetUserRequest_Email{
			Email: out.GetEmail(),
		},
		Password: Password,
	})

	suite.Nil(err)
	AssertProtoEqual(suite.T(), *GetUserProto(UserEmail, Username), *resp.GetUser(),
		pb.User{}, protocmp.IgnoreFields(&pb.User{}, "create_time", "update_time"))
	AssertUserTimestampsCurrent(suite.T(), resp.GetUser())
}

func (suite *IntegrationTestSuite) TestGetUserWithEmail_BadPassword() {
	ctx := context.Background()
	out, err := suite.UserServiceClient.CreateUser(ctx, &pb.CreateUserRequest{
		User:     GetUserProto(UserEmail, Username),
		Password: Password,
	})

	suite.Nil(err)
	AssertProtoEqual(suite.T(), *GetUserProto(UserEmail, Username), *out,
		pb.User{}, protocmp.IgnoreFields(&pb.User{}, "create_time", "update_time"))
	AssertUserTimestampsCurrent(suite.T(), out)

	resp, err := suite.UserServiceClient.GetUser(ctx, &pb.GetUserRequest{
		Id: &pb.GetUserRequest_Email{
			Email: out.GetEmail(),
		},
		Password: BadPassword,
	})

	suite.Error(err)
	suite.Nil(resp)
}

func (suite *IntegrationTestSuite) TestGetUserWithEmail_NotFound() {
	ctx := context.Background()
	resp, err := suite.UserServiceClient.GetUser(ctx, &pb.GetUserRequest{
		Id: &pb.GetUserRequest_Email{
			Email: "bad@example.com",
		},
		Password: Password,
	})

	suite.Error(err)
	suite.Nil(resp)
}

func (suite *IntegrationTestSuite) TestGetUserWithUsername_Success() {
	ctx := context.Background()
	out, err := suite.UserServiceClient.CreateUser(ctx, &pb.CreateUserRequest{
		User:     GetUserProto(UserEmail, Username),
		Password: Password,
	})

	suite.Nil(err)
	AssertProtoEqual(suite.T(), *GetUserProto(UserEmail, Username), *out,
		pb.User{}, protocmp.IgnoreFields(&pb.User{}, "create_time", "update_time"))
	AssertUserTimestampsCurrent(suite.T(), out)

	resp, err := suite.UserServiceClient.GetUser(ctx, &pb.GetUserRequest{
		Id: &pb.GetUserRequest_Username{
			Username: out.GetUsername(),
		},
		Password: Password,
	})

	suite.Nil(err)
	AssertProtoEqual(suite.T(), *GetUserProto(UserEmail, Username), *resp.GetUser(),
		pb.User{}, protocmp.IgnoreFields(&pb.User{}, "create_time", "update_time"))
	AssertUserTimestampsCurrent(suite.T(), resp.GetUser())
}

func (suite *IntegrationTestSuite) TestGetUserWithUsername_BadPassword() {
	ctx := context.Background()
	out, err := suite.UserServiceClient.CreateUser(ctx, &pb.CreateUserRequest{
		User:     GetUserProto(UserEmail, Username),
		Password: Password,
	})

	suite.Nil(err)
	AssertProtoEqual(suite.T(), *GetUserProto(UserEmail, Username), *out,
		pb.User{}, protocmp.IgnoreFields(&pb.User{}, "create_time", "update_time"))
	AssertUserTimestampsCurrent(suite.T(), out)

	resp, err := suite.UserServiceClient.GetUser(ctx, &pb.GetUserRequest{
		Id: &pb.GetUserRequest_Username{
			Username: out.GetUsername(),
		},
		Password: BadPassword,
	})

	suite.Error(err)
	suite.Nil(resp)
}

func (suite *IntegrationTestSuite) TestGetUserWithUsername_NotFound() {
	ctx := context.Background()
	resp, err := suite.UserServiceClient.GetUser(ctx, &pb.GetUserRequest{
		Id: &pb.GetUserRequest_Username{
			Username: "bad_user",
		},
		Password: Password,
	})

	suite.Error(err)
	suite.Nil(resp)
}

func (suite *IntegrationTestSuite) TestDeleteUserWithEmail_Success() {
	ctx := context.Background()
	out, err := suite.UserServiceClient.CreateUser(ctx, &pb.CreateUserRequest{
		User:     GetUserProto(UserEmail, Username),
		Password: Password,
	})

	suite.Nil(err)
	AssertProtoEqual(suite.T(), *GetUserProto(UserEmail, Username), *out,
		pb.User{}, protocmp.IgnoreFields(&pb.User{}, "create_time", "update_time"))
	AssertUserTimestampsCurrent(suite.T(), out)

	_, err = suite.UserServiceClient.DeleteUser(ctx, &pb.DeleteUserRequest{
		Id: &pb.DeleteUserRequest_Email{
			Email: out.GetEmail(),
		},
		Password: Password,
	})

	suite.Nil(err)
}

func (suite *IntegrationTestSuite) TestDeleteUserWithEmail_BadPassword() {
	ctx := context.Background()
	out, err := suite.UserServiceClient.CreateUser(ctx, &pb.CreateUserRequest{
		User:     GetUserProto(UserEmail, Username),
		Password: Password,
	})

	suite.Nil(err)
	AssertProtoEqual(suite.T(), *GetUserProto(UserEmail, Username), *out,
		pb.User{}, protocmp.IgnoreFields(&pb.User{}, "create_time", "update_time"))
	AssertUserTimestampsCurrent(suite.T(), out)

	_, err = suite.UserServiceClient.DeleteUser(ctx, &pb.DeleteUserRequest{
		Id: &pb.DeleteUserRequest_Email{
			Email: out.GetEmail(),
		},
		Password: BadPassword,
	})

	suite.Error(err)
}

func (suite *IntegrationTestSuite) TestDeleteUserWithEmail_NotFound() {
	ctx := context.Background()
	_, err := suite.UserServiceClient.DeleteUser(ctx, &pb.DeleteUserRequest{
		Id: &pb.DeleteUserRequest_Email{
			Email: "bad@example.com",
		},
		Password: Password,
	})

	suite.Error(err)
}

func (suite *IntegrationTestSuite) TestDeleteUserWithUsername_Success() {
	ctx := context.Background()
	out, err := suite.UserServiceClient.CreateUser(ctx, &pb.CreateUserRequest{
		User:     GetUserProto(UserEmail, Username),
		Password: Password,
	})

	suite.Nil(err)
	AssertProtoEqual(suite.T(), *GetUserProto(UserEmail, Username), *out,
		pb.User{}, protocmp.IgnoreFields(&pb.User{}, "create_time", "update_time"))
	AssertUserTimestampsCurrent(suite.T(), out)

	_, err = suite.UserServiceClient.DeleteUser(ctx, &pb.DeleteUserRequest{
		Id: &pb.DeleteUserRequest_Username{
			Username: out.GetUsername(),
		},
		Password: Password,
	})

	suite.Nil(err)
}

func (suite *IntegrationTestSuite) TestDeleteUserWithUsername_BadPassword() {
	ctx := context.Background()
	out, err := suite.UserServiceClient.CreateUser(ctx, &pb.CreateUserRequest{
		User:     GetUserProto(UserEmail, Username),
		Password: Password,
	})

	suite.Nil(err)
	AssertProtoEqual(suite.T(), *GetUserProto(UserEmail, Username), *out,
		pb.User{}, protocmp.IgnoreFields(&pb.User{}, "create_time", "update_time"))
	AssertUserTimestampsCurrent(suite.T(), out)

	_, err = suite.UserServiceClient.DeleteUser(ctx, &pb.DeleteUserRequest{
		Id: &pb.DeleteUserRequest_Username{
			Username: out.GetUsername(),
		},
		Password: BadPassword,
	})

	suite.Error(err)
}

func (suite *IntegrationTestSuite) TestDeleteUserWithUsername_NotFound() {
	ctx := context.Background()
	_, err := suite.UserServiceClient.DeleteUser(ctx, &pb.DeleteUserRequest{
		Id: &pb.DeleteUserRequest_Username{
			Username: "bad_user",
		},
		Password: Password,
	})

	suite.Error(err)
}
