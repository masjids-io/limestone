package test_infra

import (
	"context"
	"testing"

	upb "github.com/mnadev/limestone/user_service/proto"
	"github.com/stretchr/testify/suite"
	"google.golang.org/protobuf/testing/protocmp"
)

func TestSuite(t *testing.T) {
	suite.Run(t, new(IntegrationTestSuite))
}

func (suite *IntegrationTestSuite) TestCreateUser_Success() {
	ctx := context.Background()
	out, err := suite.UserServiceClient.CreateUser(ctx, &upb.CreateUserRequest{
		User:     GetUserProto(UserEmail, Username),
		Password: Password,
	})

	AssertProtoEqual(suite.T(), *GetUserProto(UserEmail, Username), *out,
		upb.User{}, protocmp.IgnoreFields(&upb.User{}, "create_time", "update_time"))
	AssertUserTimestampsCurrent(suite.T(), out)
	suite.Nil(err)
}

func (suite *IntegrationTestSuite) TestCreateUser_PasswordTooShort() {
	ctx := context.Background()
	out, err := suite.UserServiceClient.CreateUser(ctx, &upb.CreateUserRequest{
		User:     GetUserProto(UserEmail, Username),
		Password: BadPassword,
	})

	suite.Nil(out)
	suite.Error(err)
}

func (suite *IntegrationTestSuite) TestUpdateUser_Success() {
	ctx := context.Background()
	out, err := suite.UserServiceClient.CreateUser(ctx, &upb.CreateUserRequest{
		User:     GetUserProto(UserEmail, Username),
		Password: Password,
	})

	AssertProtoEqual(suite.T(), *GetUserProto(UserEmail, Username), *out,
		upb.User{}, protocmp.IgnoreFields(&upb.User{}, "create_time", "update_time"))
	AssertUserTimestampsCurrent(suite.T(), out)
	suite.Nil(err)

	out, err = suite.UserServiceClient.UpdateUser(ctx, &upb.UpdateUserRequest{
		User:     GetUserProto(UserEmail, "new_name"),
		Password: Password,
	})

	AssertProtoEqual(suite.T(), *GetUserProto(UserEmail, "new_name"), *out,
		upb.User{}, protocmp.IgnoreFields(&upb.User{}, "create_time", "update_time"))
	AssertUserTimestampsCurrent(suite.T(), out)
	suite.Nil(err)
}

func (suite *IntegrationTestSuite) TestUpdateUser_BadPassword() {
	ctx := context.Background()
	out, err := suite.UserServiceClient.CreateUser(ctx, &upb.CreateUserRequest{
		User:     GetUserProto(UserEmail, Username),
		Password: Password,
	})

	AssertProtoEqual(suite.T(), *GetUserProto(UserEmail, Username), *out,
		upb.User{}, protocmp.IgnoreFields(&upb.User{}, "create_time", "update_time"))
	AssertUserTimestampsCurrent(suite.T(), out)
	suite.Nil(err)

	out, err = suite.UserServiceClient.UpdateUser(ctx, &upb.UpdateUserRequest{
		User:     GetUserProto(UserEmail, "new_name"),
		Password: BadPassword,
	})

	suite.Error(err)
	suite.Nil(out)
}

func (suite *IntegrationTestSuite) TestUpdateUser_NotFound() {
	ctx := context.Background()

	user := GetUserProto(UserEmail, Username)

	out, err := suite.UserServiceClient.UpdateUser(ctx, &upb.UpdateUserRequest{
		User:     user,
		Password: Password,
	})

	suite.Error(err)
	suite.Nil(out)
}

func (suite *IntegrationTestSuite) TestGetUserWithEmail_Success() {
	ctx := context.Background()
	out, err := suite.UserServiceClient.CreateUser(ctx, &upb.CreateUserRequest{
		User:     GetUserProto(UserEmail, Username),
		Password: Password,
	})

	AssertProtoEqual(suite.T(), *GetUserProto(UserEmail, Username), *out,
		upb.User{}, protocmp.IgnoreFields(&upb.User{}, "create_time", "update_time"))
	AssertUserTimestampsCurrent(suite.T(), out)
	suite.Nil(err)

	out, err = suite.UserServiceClient.GetUser(ctx, &upb.GetUserRequest{
		Id: &upb.GetUserRequest_Email{
			Email: out.GetEmail(),
		},
		Password: Password,
	})

	AssertProtoEqual(suite.T(), *GetUserProto(UserEmail, Username), *out,
		upb.User{}, protocmp.IgnoreFields(&upb.User{}, "create_time", "update_time"))
	AssertUserTimestampsCurrent(suite.T(), out)
	suite.Nil(err)
}

func (suite *IntegrationTestSuite) TestGetUserWithEmail_BadPassword() {
	ctx := context.Background()
	out, err := suite.UserServiceClient.CreateUser(ctx, &upb.CreateUserRequest{
		User:     GetUserProto(UserEmail, Username),
		Password: Password,
	})

	AssertProtoEqual(suite.T(), *GetUserProto(UserEmail, Username), *out,
		upb.User{}, protocmp.IgnoreFields(&upb.User{}, "create_time", "update_time"))
	AssertUserTimestampsCurrent(suite.T(), out)
	suite.Nil(err)

	out, err = suite.UserServiceClient.GetUser(ctx, &upb.GetUserRequest{
		Id: &upb.GetUserRequest_Email{
			Email: out.GetEmail(),
		},
		Password: BadPassword,
	})

	suite.Error(err)
	suite.Nil(out)
}

func (suite *IntegrationTestSuite) TestGetUserWithEmail_NotFound() {
	ctx := context.Background()
	out, err := suite.UserServiceClient.GetUser(ctx, &upb.GetUserRequest{
		Id: &upb.GetUserRequest_Email{
			Email: "bad@example.com",
		},
		Password: Password,
	})

	suite.Error(err)
	suite.Nil(out)
}

func (suite *IntegrationTestSuite) TestGetUserWithUsername_Success() {
	ctx := context.Background()
	out, err := suite.UserServiceClient.CreateUser(ctx, &upb.CreateUserRequest{
		User:     GetUserProto(UserEmail, Username),
		Password: Password,
	})

	AssertProtoEqual(suite.T(), *GetUserProto(UserEmail, Username), *out,
		upb.User{}, protocmp.IgnoreFields(&upb.User{}, "create_time", "update_time"))
	AssertUserTimestampsCurrent(suite.T(), out)
	suite.Nil(err)

	out, err = suite.UserServiceClient.GetUser(ctx, &upb.GetUserRequest{
		Id: &upb.GetUserRequest_Username{
			Username: out.GetUsername(),
		},
		Password: Password,
	})

	AssertProtoEqual(suite.T(), *GetUserProto(UserEmail, Username), *out,
		upb.User{}, protocmp.IgnoreFields(&upb.User{}, "create_time", "update_time"))
	AssertUserTimestampsCurrent(suite.T(), out)
	suite.Nil(err)
}

func (suite *IntegrationTestSuite) TestGetUserWithUsername_BadPassword() {
	ctx := context.Background()
	out, err := suite.UserServiceClient.CreateUser(ctx, &upb.CreateUserRequest{
		User:     GetUserProto(UserEmail, Username),
		Password: Password,
	})

	AssertProtoEqual(suite.T(), *GetUserProto(UserEmail, Username), *out,
		upb.User{}, protocmp.IgnoreFields(&upb.User{}, "create_time", "update_time"))
	AssertUserTimestampsCurrent(suite.T(), out)
	suite.Nil(err)

	out, err = suite.UserServiceClient.GetUser(ctx, &upb.GetUserRequest{
		Id: &upb.GetUserRequest_Username{
			Username: out.GetUsername(),
		},
		Password: BadPassword,
	})

	suite.Error(err)
	suite.Nil(out)
}

func (suite *IntegrationTestSuite) TestGetUserWithUsername_NotFound() {
	ctx := context.Background()
	out, err := suite.UserServiceClient.GetUser(ctx, &upb.GetUserRequest{
		Id: &upb.GetUserRequest_Username{
			Username: "bad_user",
		},
		Password: Password,
	})

	suite.Error(err)
	suite.Nil(out)
}

func (suite *IntegrationTestSuite) TestDeleteUserWithEmail_Success() {
	ctx := context.Background()
	out, err := suite.UserServiceClient.CreateUser(ctx, &upb.CreateUserRequest{
		User:     GetUserProto(UserEmail, Username),
		Password: Password,
	})

	AssertProtoEqual(suite.T(), *GetUserProto(UserEmail, Username), *out,
		upb.User{}, protocmp.IgnoreFields(&upb.User{}, "create_time", "update_time"))
	AssertUserTimestampsCurrent(suite.T(), out)
	suite.Nil(err)

	_, err = suite.UserServiceClient.DeleteUser(ctx, &upb.DeleteUserRequest{
		Id: &upb.DeleteUserRequest_Email{
			Email: out.GetEmail(),
		},
		Password: Password,
	})

	suite.Nil(err)
}

func (suite *IntegrationTestSuite) TestDeleteUserWithEmail_BadPassword() {
	ctx := context.Background()
	out, err := suite.UserServiceClient.CreateUser(ctx, &upb.CreateUserRequest{
		User:     GetUserProto(UserEmail, Username),
		Password: Password,
	})

	AssertProtoEqual(suite.T(), *GetUserProto(UserEmail, Username), *out,
		upb.User{}, protocmp.IgnoreFields(&upb.User{}, "create_time", "update_time"))
	AssertUserTimestampsCurrent(suite.T(), out)
	suite.Nil(err)

	_, err = suite.UserServiceClient.DeleteUser(ctx, &upb.DeleteUserRequest{
		Id: &upb.DeleteUserRequest_Email{
			Email: out.GetEmail(),
		},
		Password: BadPassword,
	})

	suite.Error(err)
}

func (suite *IntegrationTestSuite) TestDeleteUserWithEmail_NotFound() {
	ctx := context.Background()
	_, err := suite.UserServiceClient.DeleteUser(ctx, &upb.DeleteUserRequest{
		Id: &upb.DeleteUserRequest_Email{
			Email: "bad@example.com",
		},
		Password: Password,
	})

	suite.Error(err)
}

func (suite *IntegrationTestSuite) TestDeleteUserWithUsername_Success() {
	ctx := context.Background()
	out, err := suite.UserServiceClient.CreateUser(ctx, &upb.CreateUserRequest{
		User:     GetUserProto(UserEmail, Username),
		Password: Password,
	})

	AssertProtoEqual(suite.T(), *GetUserProto(UserEmail, Username), *out,
		upb.User{}, protocmp.IgnoreFields(&upb.User{}, "create_time", "update_time"))
	AssertUserTimestampsCurrent(suite.T(), out)
	suite.Nil(err)

	_, err = suite.UserServiceClient.DeleteUser(ctx, &upb.DeleteUserRequest{
		Id: &upb.DeleteUserRequest_Username{
			Username: out.GetUsername(),
		},
		Password: Password,
	})

	suite.Nil(err)
}

func (suite *IntegrationTestSuite) TestDeleteUserWithUsername_BadPassword() {
	ctx := context.Background()
	out, err := suite.UserServiceClient.CreateUser(ctx, &upb.CreateUserRequest{
		User:     GetUserProto(UserEmail, Username),
		Password: Password,
	})

	AssertProtoEqual(suite.T(), *GetUserProto(UserEmail, Username), *out,
		upb.User{}, protocmp.IgnoreFields(&upb.User{}, "create_time", "update_time"))
	AssertUserTimestampsCurrent(suite.T(), out)
	suite.Nil(err)

	_, err = suite.UserServiceClient.DeleteUser(ctx, &upb.DeleteUserRequest{
		Id: &upb.DeleteUserRequest_Username{
			Username: out.GetUsername(),
		},
		Password: BadPassword,
	})

	suite.Error(err)
}

func (suite *IntegrationTestSuite) TestDeleteUserWithUsername_NotFound() {
	ctx := context.Background()
	_, err := suite.UserServiceClient.DeleteUser(ctx, &upb.DeleteUserRequest{
		Id: &upb.DeleteUserRequest_Username{
			Username: "bad_user",
		},
		Password: Password,
	})

	suite.Error(err)
}
