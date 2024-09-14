package test_infra

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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

	suite.Assert().Equal(status.Code(err), codes.OK)
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
	suite.Assert().Equal(status.Code(err), codes.InvalidArgument)
}

func (suite *IntegrationTestSuite) TestUpdateUserSuccess() {
	ctx := context.Background()
	out, err := suite.UserServiceClient.CreateUser(ctx, &pb.CreateUserRequest{
		User:     GetUserProto(UserEmail, Username),
		Password: Password,
	})

	suite.Assert().Equal(status.Code(err), codes.OK)
	AssertProtoEqual(suite.T(), *GetUserProto(UserEmail, Username), *out,
		pb.User{}, protocmp.IgnoreFields(&pb.User{}, "create_time", "update_time"))
	AssertUserTimestampsCurrent(suite.T(), out)

	out, err = suite.UserServiceClient.UpdateUser(ctx, &pb.UpdateUserRequest{
		User: GetUserProto(UserEmail, "new_name"),
	})

	suite.Assert().Equal(status.Code(err), codes.OK)
	AssertProtoEqual(suite.T(), *GetUserProto(UserEmail, "new_name"), *out,
		pb.User{}, protocmp.IgnoreFields(&pb.User{}, "create_time", "update_time"))
	AssertUserTimestampsCurrent(suite.T(), out)
}

func (suite *IntegrationTestSuite) TestUpdateUser_NotFound() {
	ctx := context.Background()

	user := GetUserProto(UserEmail, Username)

	out, err := suite.UserServiceClient.UpdateUser(ctx, &pb.UpdateUserRequest{
		User: user,
	})

	suite.Assert().Equal(status.Code(err), codes.NotFound)
	suite.Nil(out)
}

func (suite *IntegrationTestSuite) TestGetUserSuccess() {
	ctx := context.Background()
	out, err := suite.UserServiceClient.CreateUser(ctx, &pb.CreateUserRequest{
		User:     GetUserProto(UserEmail, Username),
		Password: Password,
	})

	suite.Assert().Equal(status.Code(err), codes.OK)
	AssertProtoEqual(suite.T(), *GetUserProto(UserEmail, Username), *out,
		pb.User{}, protocmp.IgnoreFields(&pb.User{}, "create_time", "update_time"))
	AssertUserTimestampsCurrent(suite.T(), out)

	resp, err := suite.UserServiceClient.GetUser(ctx, &pb.GetUserRequest{
		Id: out.GetId(),
	})

	suite.Assert().Equal(status.Code(err), codes.OK)
	AssertProtoEqual(suite.T(), *GetUserProto(UserEmail, Username), *resp.GetUser(),
		pb.User{}, protocmp.IgnoreFields(&pb.User{}, "create_time", "update_time"))
	AssertUserTimestampsCurrent(suite.T(), resp.GetUser())
}

func (suite *IntegrationTestSuite) TestGetUserNotFound() {
	ctx := context.Background()
	resp, err := suite.UserServiceClient.GetUser(ctx, &pb.GetUserRequest{
		Id: "wrongid",
	})

	suite.Assert().Equal(status.Code(err), codes.NotFound)
	suite.Nil(resp)
}

func (suite *IntegrationTestSuite) TestDeleteUserSuccess() {
	ctx := context.Background()
	out, err := suite.UserServiceClient.CreateUser(ctx, &pb.CreateUserRequest{
		User:     GetUserProto(UserEmail, Username),
		Password: Password,
	})

	suite.Assert().Equal(status.Code(err), codes.OK)
	AssertProtoEqual(suite.T(), *GetUserProto(UserEmail, Username), *out,
		pb.User{}, protocmp.IgnoreFields(&pb.User{}, "create_time", "update_time"))
	AssertUserTimestampsCurrent(suite.T(), out)

	_, err = suite.UserServiceClient.DeleteUser(ctx, &pb.DeleteUserRequest{
		Id: out.GetId(),
	})

	suite.Assert().Equal(status.Code(err), codes.OK)
}

func (suite *IntegrationTestSuite) TestDeleteUserNotFound() {
	ctx := context.Background()
	_, err := suite.UserServiceClient.DeleteUser(ctx, &pb.DeleteUserRequest{
		Id: "wrongid",
	})

	suite.Assert().Equal(status.Code(err), codes.NotFound)
}
