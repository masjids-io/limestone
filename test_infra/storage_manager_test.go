package test_infra

import (
	"testing"

	epb "github.com/mnadev/limestone/event_service/proto"
	mpb "github.com/mnadev/limestone/masjid_service/proto"
	upb "github.com/mnadev/limestone/user_service/proto"
	"github.com/stretchr/testify/suite"
	"google.golang.org/protobuf/testing/protocmp"
)

func TestSuite(t *testing.T) {
	suite.Run(t, new(UnitTestSuite))
}

func (suite *UnitTestSuite) TestCreateUser_Success() {
	user, err := suite.StorageManager.CreateUser(GetUserProto(UserEmail, Username), Password)
	AssertProtoEqual(suite.T(), *GetUserProto(UserEmail, Username), *user.ToProto(),
		upb.User{}, protocmp.IgnoreFields(&upb.User{}, "create_time", "update_time"))
	AssertUserTimestampsCurrent(suite.T(), user.ToProto())
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
	AssertProtoEqual(suite.T(), *GetUserProto("a@example.com", Username), *user.ToProto(),
		upb.User{}, protocmp.IgnoreFields(&upb.User{}, "create_time", "update_time"))
	AssertUserTimestampsCurrent(suite.T(), user.ToProto())
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
	AssertProtoEqual(suite.T(), *GetUserProto(UserEmail, Username), *user.ToProto(),
		upb.User{}, protocmp.IgnoreFields(&upb.User{}, "create_time", "update_time"))
	AssertUserTimestampsCurrent(suite.T(), user.ToProto())
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
	AssertProtoEqual(suite.T(), *GetUserProto(UserEmail, Username), *user.ToProto(),
		upb.User{}, protocmp.IgnoreFields(&upb.User{}, "create_time", "update_time"))
	AssertUserTimestampsCurrent(suite.T(), user.ToProto())
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

func (suite *UnitTestSuite) TestCreateMasjid_Success() {
	masjid, err := suite.StorageManager.CreateMasjid(GetMasjidProto())
	AssertProtoEqual(suite.T(), *GetMasjidProto(), *masjid.ToProto(),
		mpb.Masjid{}, protocmp.IgnoreFields(&mpb.Masjid{}, "create_time", "update_time"))
	AssertMasjidTimestampsCurrent(suite.T(), masjid.ToProto())
	suite.Nil(err)
}

func (suite *UnitTestSuite) TestUpdateMasjid_Success() {
	masjid, err := suite.StorageManager.CreateMasjid(GetMasjidProto())
	suite.Nil(err)

	masjid.Name = "Masjid 2"
	masjid, err = suite.StorageManager.UpdateMasjid(masjid.ToProto())

	want := GetMasjidProto()
	want.Name = "Masjid 2"

	AssertProtoEqual(suite.T(), *want, *masjid.ToProto(), mpb.Masjid{},
		protocmp.IgnoreFields(&mpb.Masjid{}, "create_time", "update_time"))
	AssertMasjidTimestampsCurrent(suite.T(), masjid.ToProto())
	suite.Nil(err)
}

func (suite *UnitTestSuite) TestUpdateMasjid_NotFound() {
	masjid, err := suite.StorageManager.UpdateMasjid(GetMasjidProto())
	suite.Error(err)
	suite.Nil(masjid)
}

func (suite *UnitTestSuite) TestGetMasjid_Success() {
	_, err := suite.StorageManager.CreateMasjid(GetMasjidProto())
	suite.Nil(err)

	masjid, err := suite.StorageManager.GetMasjid(DefaultId)
	AssertProtoEqual(suite.T(), *GetMasjidProto(), *masjid.ToProto(),
		mpb.Masjid{}, protocmp.IgnoreFields(&mpb.Masjid{}, "create_time", "update_time"))
	AssertMasjidTimestampsCurrent(suite.T(), masjid.ToProto())
	suite.Nil(err)
}

func (suite *UnitTestSuite) TestGetMasjid_NotFound() {
	masjid, err := suite.StorageManager.GetMasjid(DefaultId)
	suite.Error(err)
	suite.Nil(masjid)
}

func (suite *UnitTestSuite) TestDeleteMasjid_Success() {
	_, err := suite.StorageManager.CreateMasjid(GetMasjidProto())
	suite.Nil(err)

	err = suite.StorageManager.DeleteMasjid(DefaultId)
	suite.Nil(err)
}

func (suite *UnitTestSuite) TestDeleteMasjid_NotFound() {
	err := suite.StorageManager.DeleteMasjid(DefaultId)
	suite.Error(err)
}

func (suite *UnitTestSuite) TestCreateEvent_Success() {
	event, err := suite.StorageManager.CreateEvent(GetEventProto())
	AssertProtoEqual(suite.T(), *GetEventProto(), *event.ToProto(),
		epb.Event{}, protocmp.IgnoreFields(&epb.Event{}, "create_time", "update_time"))
	AssertEventTimestampsCurrent(suite.T(), event.ToProto())
	suite.Nil(err)
}

func (suite *UnitTestSuite) TestUpdateEvent_Success() {
	event, err := suite.StorageManager.CreateEvent(GetEventProto())
	suite.Nil(err)

	event.Name = "Event 2"
	event, err = suite.StorageManager.UpdateEvent(event.ToProto())

	want := GetEventProto()
	want.Name = "Event 2"

	AssertProtoEqual(suite.T(), *want, *event.ToProto(), epb.Event{},
		protocmp.IgnoreFields(&epb.Event{}, "create_time", "update_time"))
	AssertEventTimestampsCurrent(suite.T(), event.ToProto())
	suite.Nil(err)
}

func (suite *UnitTestSuite) TestUpdateEvent_NotFound() {
	event, err := suite.StorageManager.UpdateEvent(GetEventProto())
	suite.Error(err)
	suite.Nil(event)
}

func (suite *UnitTestSuite) TestGetEvent_Success() {
	_, err := suite.StorageManager.CreateEvent(GetEventProto())
	suite.Nil(err)

	event, err := suite.StorageManager.GetEvent(DefaultId)
	AssertProtoEqual(suite.T(), *GetEventProto(), *event.ToProto(),
		epb.Event{}, protocmp.IgnoreFields(&epb.Event{}, "create_time", "update_time"))
	AssertEventTimestampsCurrent(suite.T(), event.ToProto())
	suite.Nil(err)
}

func (suite *UnitTestSuite) TestGetEvent_NotFound() {
	event, err := suite.StorageManager.GetEvent(DefaultId)
	suite.Error(err)
	suite.Nil(event)
}

func (suite *UnitTestSuite) TestDeleteEvent_Success() {
	_, err := suite.StorageManager.CreateEvent(GetEventProto())
	suite.Nil(err)

	err = suite.StorageManager.DeleteEvent(DefaultId)
	suite.Nil(err)
}

func (suite *UnitTestSuite) TestDeleteEvent_NotFound() {
	err := suite.StorageManager.DeleteEvent(DefaultId)
	suite.Error(err)
}
