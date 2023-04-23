package test_infra

import (
	"testing"

	_ "github.com/proullon/ramsql/driver"
	"github.com/stretchr/testify/suite"
	_ "google.golang.org/protobuf/types/known/timestamppb"
)

func TestSuite(t *testing.T) {
	suite.Run(t, new(UnitTestSuite))
}

func (suite *UnitTestSuite) TestCreateUser_Success() {
	user, err := suite.StorageManager.CreateUser(GetUserProto(UserEmail, Username), Password)
	AssertUserProtoEqual(suite.T(), GetUserProto(UserEmail, Username), user.ToProto())
	suite.Nil(err)
}

func (suite *UnitTestSuite) TestCreateUser_PasswordTooShort() {
	user, err := suite.StorageManager.CreateUser(GetUserProto(UserEmail, Username), BadPassword)
	suite.Nil(user)
	suite.Error(err)
}

func (suite *UnitTestSuite) TestUpdateUserWithEmail_Success() {
	user, err := suite.StorageManager.CreateUser(GetUserProto(UserEmail, Username), Password)
	suite.Nil(err)

	user.Email = "a@example.com"
	user, err = suite.StorageManager.UpdateUser(user.ToProto(), Password)
	AssertUserProtoEqual(suite.T(), GetUserProto("a@example.com", Username), user.ToProto())
	suite.Nil(err)
}

func (suite *UnitTestSuite) TestUpdateUserWithEmail_BadPassword() {
	_, err := suite.StorageManager.CreateUser(GetUserProto(UserEmail, Username), Password)
	suite.Nil(err)

	user, err := suite.StorageManager.UpdateUser(GetUserProto("a@example.com", Username), BadPassword)
	suite.Error(err)
	suite.Nil(user)
}

func (suite *UnitTestSuite) TestUpdateUserWithEmail_NotFound() {
	user, err := suite.StorageManager.UpdateUser(GetUserProto(UserEmail, Username), Password)
	suite.Error(err)
	suite.Nil(user)
}

func (suite *UnitTestSuite) TestGetUserWithEmail_Success() {
	_, err := suite.StorageManager.CreateUser(GetUserProto(UserEmail, Username), Password)
	suite.Nil(err)

	user, err := suite.StorageManager.GetUserWithEmail(UserEmail, Password)
	AssertUserProtoEqual(suite.T(), GetUserProto(UserEmail, Username), user.ToProto())
	suite.Nil(err)
}

func (suite *UnitTestSuite) TestGetUserWithEmail_BadPassword() {
	_, err := suite.StorageManager.CreateUser(GetUserProto(UserEmail, Username), Password)
	suite.Nil(err)

	user, err := suite.StorageManager.GetUserWithEmail(UserEmail, BadPassword)
	suite.Error(err)
	suite.Nil(user)
}

func (suite *UnitTestSuite) TestGetUserWithEmail_NotFound() {
	_, err := suite.StorageManager.CreateUser(GetUserProto(UserEmail, Username), Password)
	suite.Nil(err)

	user, err := suite.StorageManager.GetUserWithEmail("a@example.com", Password)
	suite.Error(err)
	suite.Nil(user)
}

func (suite *UnitTestSuite) TestGetUserWithUsername_Success() {
	_, err := suite.StorageManager.CreateUser(GetUserProto(UserEmail, Username), Password)
	suite.Nil(err)

	user, err := suite.StorageManager.GetUserWithUsername(Username, Password)
	AssertUserProtoEqual(suite.T(), GetUserProto(UserEmail, Username), user.ToProto())
	suite.Nil(err)
}

func (suite *UnitTestSuite) TestGetUserWithUsername_BadPassword() {
	_, err := suite.StorageManager.CreateUser(GetUserProto(UserEmail, Username), Password)
	suite.Nil(err)

	user, err := suite.StorageManager.GetUserWithUsername(Username, BadPassword)
	suite.Error(err)
	suite.Nil(user)
}

func (suite *UnitTestSuite) TestGetUserWithUsername_NotFound() {
	_, err := suite.StorageManager.CreateUser(GetUserProto(UserEmail, Username), Password)
	suite.Nil(err)

	user, err := suite.StorageManager.GetUserWithUsername("notcoolguy1234", Password)
	suite.Error(err)
	suite.Nil(user)
}

func (suite *UnitTestSuite) TestDeleteUserWithEmail_Success() {
	_, err := suite.StorageManager.CreateUser(GetUserProto(UserEmail, Username), Password)
	suite.Nil(err)

	err = suite.StorageManager.DeleteUserWithEmail(UserEmail, Password)
	suite.Nil(err)
}

func (suite *UnitTestSuite) TestDeleteUserWithEmail_BadPassword() {
	_, err := suite.StorageManager.CreateUser(GetUserProto(UserEmail, Username), Password)
	suite.Nil(err)

	err = suite.StorageManager.DeleteUserWithEmail(UserEmail, BadPassword)
	suite.Error(err)
}

func (suite *UnitTestSuite) TestDeleteUserWithEmail_NotFound() {
	_, err := suite.StorageManager.CreateUser(GetUserProto(UserEmail, Username), Password)
	suite.Nil(err)

	err = suite.StorageManager.DeleteUserWithEmail("a@example.com", Password)
	suite.Error(err)
}

func (suite *UnitTestSuite) TestDeleteUserWithUsername_Success() {
	_, err := suite.StorageManager.CreateUser(GetUserProto(UserEmail, Username), Password)
	suite.Nil(err)

	err = suite.StorageManager.DeleteUserWithUsername(Username, Password)
	suite.Nil(err)
}

func (suite *UnitTestSuite) TestDeleteUserWithUsername_BadPassword() {
	_, err := suite.StorageManager.CreateUser(GetUserProto(UserEmail, Username), Password)
	suite.Nil(err)

	err = suite.StorageManager.DeleteUserWithUsername(Username, BadPassword)
	suite.Error(err)
}

func (suite *UnitTestSuite) TestDeleteUserWithUsername_NotFound() {
	_, err := suite.StorageManager.CreateUser(GetUserProto(UserEmail, Username), Password)
	suite.Nil(err)

	err = suite.StorageManager.DeleteUserWithUsername("notcoolguy1234", Password)
	suite.Error(err)
}
