package test_infra

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/testing/protocmp"

	pb "github.com/mnadev/limestone/proto"
)

func TestStorageManager(t *testing.T) {
	suite.Run(t, new(UnitTestSuite))
}

func (suite *UnitTestSuite) TestCreateUserSuccess() {
	user, err := suite.StorageManager.CreateUser(GetUserProto(UserEmail, Username), Password)

	suite.Assert().Equal(status.Code(err), codes.OK)
	AssertProtoEqual(suite.T(), *GetUserProto(UserEmail, Username), *user.ToProto(),
		pb.User{}, protocmp.IgnoreFields(&pb.User{}, "create_time", "update_time"))
	AssertUserTimestampsCurrent(suite.T(), user.ToProto())
}

func (suite *UnitTestSuite) TestCreateUsePasswordTooShort() {
	user, err := suite.StorageManager.CreateUser(GetUserProto(UserEmail, Username), BadPassword)
	suite.Nil(user)
	suite.Assert().Equal(status.Code(err), codes.InvalidArgument)
}

func (suite *UnitTestSuite) TestUpdateUserWithEmailSuccess() {
	user, err := suite.StorageManager.CreateUser(GetUserProto(UserEmail, Username), Password)
	suite.Assert().Equal(status.Code(err), codes.OK)

	user.Email = "a@example.com"
	user, err = suite.StorageManager.UpdateUser(user.ToProto())

	suite.Assert().Equal(status.Code(err), codes.OK)
	AssertProtoEqual(suite.T(), *GetUserProto("a@example.com", Username), *user.ToProto(),
		pb.User{}, protocmp.IgnoreFields(&pb.User{}, "create_time", "update_time"))
	AssertUserTimestampsCurrent(suite.T(), user.ToProto())
}

func (suite *UnitTestSuite) TestUpdateUserWithEmailNotFound() {
	user, err := suite.StorageManager.UpdateUser(GetUserProto(UserEmail, Username))

	suite.Assert().Equal(status.Code(err), codes.NotFound)
	suite.Nil(user)
}

func (suite *UnitTestSuite) TestGetUserSuccess() {
	_, err := suite.StorageManager.CreateUser(GetUserProto(UserEmail, Username), Password)
	suite.Assert().Equal(status.Code(err), codes.OK)

	user, err := suite.StorageManager.GetUser(DefaultId)

	suite.Assert().Equal(status.Code(err), codes.OK)
	AssertProtoEqual(suite.T(), *GetUserProto(UserEmail, Username), *user.ToProto(),
		pb.User{}, protocmp.IgnoreFields(&pb.User{}, "create_time", "update_time"))
	AssertUserTimestampsCurrent(suite.T(), user.ToProto())
}

func (suite *UnitTestSuite) TestGetUserNotFound() {
	_, err := suite.StorageManager.CreateUser(GetUserProto(UserEmail, Username), Password)
	suite.Assert().Equal(status.Code(err), codes.OK)

	user, err := suite.StorageManager.GetUser("wrongid")
	suite.Assert().Equal(status.Code(err), codes.NotFound)
	suite.Nil(user)
}

func (suite *UnitTestSuite) TestDeleteUserSuccess() {
	_, err := suite.StorageManager.CreateUser(GetUserProto(UserEmail, Username), Password)
	suite.Assert().Equal(status.Code(err), codes.OK)

	err = suite.StorageManager.DeleteUser(DefaultId)
	suite.Assert().Equal(status.Code(err), codes.OK)
}

func (suite *UnitTestSuite) TestDeleteUserNotFound() {
	_, err := suite.StorageManager.CreateUser(GetUserProto(UserEmail, Username), Password)
	suite.Assert().Equal(status.Code(err), codes.OK)

	err = suite.StorageManager.DeleteUser("wrongid")
	suite.Assert().Equal(status.Code(err), codes.NotFound)
}

func (suite *UnitTestSuite) TestCreateMasjidSuccess() {
	masjid, err := suite.StorageManager.CreateMasjid(GetMasjidProto())

	suite.Assert().Equal(status.Code(err), codes.OK)
	AssertProtoEqual(suite.T(), *GetMasjidProto(), *masjid.ToProto(),
		pb.Masjid{}, protocmp.IgnoreFields(&pb.Masjid{}, "create_time", "update_time"))
	AssertMasjidTimestampsCurrent(suite.T(), masjid.ToProto())
}

func (suite *UnitTestSuite) TestUpdateMasjidSuccess() {
	masjid, err := suite.StorageManager.CreateMasjid(GetMasjidProto())
	suite.Assert().Equal(status.Code(err), codes.OK)

	masjid.Name = "Masjid 2"
	masjid, err = suite.StorageManager.UpdateMasjid(masjid.ToProto())

	want := GetMasjidProto()
	want.Name = "Masjid 2"

	suite.Assert().Equal(status.Code(err), codes.OK)
	AssertProtoEqual(suite.T(), *want, *masjid.ToProto(), pb.Masjid{},
		protocmp.IgnoreFields(&pb.Masjid{}, "create_time", "update_time"))
	AssertMasjidTimestampsCurrent(suite.T(), masjid.ToProto())
}

func (suite *UnitTestSuite) TestUpdateMasjidNotFound() {
	masjid, err := suite.StorageManager.UpdateMasjid(GetMasjidProto())

	suite.Assert().Equal(status.Code(err), codes.NotFound)
	suite.Nil(masjid)
}

func (suite *UnitTestSuite) TestGetMasjidSuccess() {
	_, err := suite.StorageManager.CreateMasjid(GetMasjidProto())
	suite.Assert().Equal(status.Code(err), codes.OK)

	masjid, err := suite.StorageManager.GetMasjid(DefaultId)

	suite.Assert().Equal(status.Code(err), codes.OK)
	AssertProtoEqual(suite.T(), *GetMasjidProto(), *masjid.ToProto(),
		pb.Masjid{}, protocmp.IgnoreFields(&pb.Masjid{}, "create_time", "update_time"))
	AssertMasjidTimestampsCurrent(suite.T(), masjid.ToProto())
}

func (suite *UnitTestSuite) TestGetMasjidNotFound() {
	masjid, err := suite.StorageManager.GetMasjid(DefaultId)

	suite.Assert().Equal(status.Code(err), codes.NotFound)
	suite.Nil(masjid)
}

func (suite *UnitTestSuite) TestDeleteMasjidSuccess() {
	_, err := suite.StorageManager.CreateMasjid(GetMasjidProto())
	suite.Assert().Equal(status.Code(err), codes.OK)

	err = suite.StorageManager.DeleteMasjid(DefaultId)
	suite.Assert().Equal(status.Code(err), codes.OK)
}

func (suite *UnitTestSuite) TestDeleteMasjidNotFound() {
	err := suite.StorageManager.DeleteMasjid(DefaultId)
	suite.Assert().Equal(status.Code(err), codes.NotFound)
}

func (suite *UnitTestSuite) TestCreateEventSuccess() {
	event, err := suite.StorageManager.CreateEvent(GetEventProto())

	suite.Assert().Equal(status.Code(err), codes.OK)
	AssertProtoEqual(suite.T(), *GetEventProto(), *event.ToProto(),
		pb.Event{}, protocmp.IgnoreFields(&pb.Event{}, "create_time", "update_time"))
	AssertEventTimestampsCurrent(suite.T(), event.ToProto())
}

func (suite *UnitTestSuite) TestUpdateEventSuccess() {
	event, err := suite.StorageManager.CreateEvent(GetEventProto())
	suite.Assert().Equal(status.Code(err), codes.OK)

	event.Name = "Event 2"
	event, err = suite.StorageManager.UpdateEvent(event.ToProto())

	want := GetEventProto()
	want.Name = "Event 2"

	suite.Assert().Equal(status.Code(err), codes.OK)
	AssertProtoEqual(suite.T(), *want, *event.ToProto(), pb.Event{},
		protocmp.IgnoreFields(&pb.Event{}, "create_time", "update_time"))
	AssertEventTimestampsCurrent(suite.T(), event.ToProto())
}

func (suite *UnitTestSuite) TestUpdateEventNotFound() {
	event, err := suite.StorageManager.UpdateEvent(GetEventProto())

	suite.Assert().Equal(status.Code(err), codes.NotFound)
	suite.Nil(event)
}

func (suite *UnitTestSuite) TestGetEventSuccess() {
	_, err := suite.StorageManager.CreateEvent(GetEventProto())
	suite.Assert().Equal(status.Code(err), codes.OK)

	event, err := suite.StorageManager.GetEvent(DefaultId)

	suite.Assert().Equal(status.Code(err), codes.OK)
	AssertProtoEqual(suite.T(), *GetEventProto(), *event.ToProto(),
		pb.Event{}, protocmp.IgnoreFields(&pb.Event{}, "create_time", "update_time"))
	AssertEventTimestampsCurrent(suite.T(), event.ToProto())
}

func (suite *UnitTestSuite) TestGetEventNotFound() {
	event, err := suite.StorageManager.GetEvent(DefaultId)

	suite.Assert().Equal(status.Code(err), codes.NotFound)
	suite.Nil(event)
}

func (suite *UnitTestSuite) TestDeleteEventSuccess() {
	_, err := suite.StorageManager.CreateEvent(GetEventProto())
	suite.Assert().Equal(status.Code(err), codes.OK)

	err = suite.StorageManager.DeleteEvent(DefaultId)
	suite.Assert().Equal(status.Code(err), codes.OK)
}

func (suite *UnitTestSuite) TestDeleteEventNotFound() {
	err := suite.StorageManager.DeleteEvent(DefaultId)
	suite.Assert().Equal(status.Code(err), codes.NotFound)
}

func (suite *UnitTestSuite) TestCreateAdhanFileSuccess() {
	file, err := suite.StorageManager.CreateAdhanFile(GetAdhanFileProto())

	suite.Assert().Equal(status.Code(err), codes.OK)
	AssertProtoEqual(suite.T(), *GetAdhanFileProto(), *file.ToProto(),
		pb.AdhanFile{}, protocmp.IgnoreFields(&pb.AdhanFile{}, "create_time", "update_time"))
	AssertAdhanFileTimestampsCurrent(suite.T(), file.ToProto())
}

func (suite *UnitTestSuite) TestUpdateAdhanFileSuccess() {
	file, err := suite.StorageManager.CreateAdhanFile(GetAdhanFileProto())
	suite.Assert().Equal(status.Code(err), codes.OK)

	file.File = []byte("xyz")
	got, err := suite.StorageManager.UpdateAdhanFile(file.ToProto())

	want := GetAdhanFileProto()
	want.File = []byte("xyz")

	suite.Assert().Equal(status.Code(err), codes.OK)
	AssertProtoEqual(suite.T(), *want, *got.ToProto(), pb.AdhanFile{},
		protocmp.IgnoreFields(&pb.AdhanFile{}, "create_time", "update_time"))
	AssertAdhanFileTimestampsCurrent(suite.T(), file.ToProto())
}

func (suite *UnitTestSuite) TestUpdateAdhanFileNotFound() {
	file, err := suite.StorageManager.UpdateAdhanFile(GetAdhanFileProto())

	suite.Assert().Equal(status.Code(err), codes.NotFound)
	suite.Nil(file)
}

func (suite *UnitTestSuite) TestGetAdhanFileSuccess() {
	_, err := suite.StorageManager.CreateAdhanFile(GetAdhanFileProto())
	suite.Assert().Equal(status.Code(err), codes.OK)

	file, err := suite.StorageManager.GetAdhanFile(DefaultId)

	suite.Assert().Equal(status.Code(err), codes.OK)
	AssertProtoEqual(suite.T(), *GetAdhanFileProto(), *file.ToProto(),
		pb.AdhanFile{}, protocmp.IgnoreFields(&pb.AdhanFile{}, "create_time", "update_time"))
	AssertAdhanFileTimestampsCurrent(suite.T(), file.ToProto())
}

func (suite *UnitTestSuite) TestGetAdhanFileNotFound() {
	file, err := suite.StorageManager.GetAdhanFile(DefaultId)

	suite.Assert().Equal(status.Code(err), codes.NotFound)
	suite.Nil(file)
}

func (suite *UnitTestSuite) TestDeleteAdhanFileSuccess() {
	_, err := suite.StorageManager.CreateAdhanFile(GetAdhanFileProto())
	suite.Assert().Equal(status.Code(err), codes.OK)

	err = suite.StorageManager.DeleteAdhanFile(DefaultId)
	suite.Assert().Equal(status.Code(err), codes.OK)
}

func (suite *UnitTestSuite) TestDeleteAdhanFileNotFound() {
	err := suite.StorageManager.DeleteAdhanFile(DefaultId)
	suite.Assert().Equal(status.Code(err), codes.NotFound)
}
