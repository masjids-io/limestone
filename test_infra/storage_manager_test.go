package test_infra

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/testing/protocmp"

	pb "github.com/mnadev/limestone/proto"
	"github.com/mnadev/limestone/storage"
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

func (suite *UnitTestSuite) TestCreateNikkahProfileSuccess() {
	_, err := suite.StorageManager.CreateUser(GetUserProto(UserEmail, Username), Password)
	suite.Assert().Equal(status.Code(err), codes.OK)
	_, err = suite.StorageManager.CreateNikkahProfile(GetNikkahProfileProto())
	suite.Assert().Equal(status.Code(err), codes.OK)
}

func (suite *UnitTestSuite) TestCreateNikkahProfileWithNil() {
	_, err := suite.StorageManager.CreateUser(GetUserProto(UserEmail, Username), Password)
	suite.Assert().Equal(status.Code(err), codes.OK)
	_, err = suite.StorageManager.CreateNikkahProfile(nil)
	suite.Assert().Equal(status.Code(err), codes.InvalidArgument)
}

func (suite *UnitTestSuite) TestCreateNikkahProfileWithNoUserId() {
	profile := GetNikkahProfileProto()
	profile.UserId = ""
	_, err := suite.StorageManager.CreateNikkahProfile(profile)
	suite.Assert().Equal(status.Code(err), codes.InvalidArgument)
}

func (suite *UnitTestSuite) TestUpdateNikkahProfileSuccess() {
	_, err := suite.StorageManager.CreateUser(GetUserProto(UserEmail, Username), Password)
	suite.Assert().Equal(status.Code(err), codes.OK)
	profile, err := suite.StorageManager.CreateNikkahProfile(GetNikkahProfileProto())
	suite.Assert().Equal(status.Code(err), codes.OK)

	profile.Name = "xyz"
	_, err = suite.StorageManager.UpdateNikkahProfile(profile.ToProto())
	suite.Assert().Equal(status.Code(err), codes.OK)
}

func (suite *UnitTestSuite) TestUpdateNikkahProfileWithNil() {
	_, err := suite.StorageManager.UpdateNikkahProfile(nil)
	suite.Assert().Equal(status.Code(err), codes.InvalidArgument)
}

func (suite *UnitTestSuite) TestUpdateNikkahProfileWithNoId() {
	profile := GetNikkahProfileProto()
	profile.Id = ""
	_, err := suite.StorageManager.UpdateNikkahProfile(profile)
	suite.Assert().Equal(status.Code(err), codes.InvalidArgument)
}

func (suite *UnitTestSuite) TestUpdateNikkahProfileWithNoUserId() {
	_, err := suite.StorageManager.CreateUser(GetUserProto(UserEmail, Username), Password)
	suite.Assert().Equal(status.Code(err), codes.OK)
	profile, err := suite.StorageManager.CreateNikkahProfile(GetNikkahProfileProto())
	suite.Assert().Equal(status.Code(err), codes.OK)

	profile.UserID = ""
	_, err = suite.StorageManager.UpdateNikkahProfile(profile.ToProto())
	suite.Assert().Equal(status.Code(err), codes.InvalidArgument)
}

func (suite *UnitTestSuite) TestUpdateNikkahProfileNotFound() {
	profile := GetNikkahProfileProto()
	profile.Id = "01"
	_, err := suite.StorageManager.UpdateNikkahProfile(profile)
	suite.Assert().Equal(status.Code(err), codes.NotFound)
}

func (suite *UnitTestSuite) TestGetNikkahProfileSuccess() {
	_, err := suite.StorageManager.CreateUser(GetUserProto(UserEmail, Username), Password)
	suite.Assert().Equal(status.Code(err), codes.OK)
	_, err = suite.StorageManager.CreateNikkahProfile(GetNikkahProfileProto())
	suite.Assert().Equal(status.Code(err), codes.OK)

	profile, err := suite.StorageManager.GetNikkahProfile(DefaultId)
	suite.Assert().Equal(status.Code(err), codes.OK)
	AssertProtoEqual(suite.T(), *GetNikkahProfileProto(), *profile.ToProto(),
		pb.NikkahProfile{}, protocmp.IgnoreFields(&pb.NikkahProfile{}, "create_time", "update_time"))
}

func (suite *UnitTestSuite) TestGetNikkahProfileNotFound() {
	_, err := suite.StorageManager.GetNikkahProfile(DefaultId)
	suite.Assert().Equal(status.Code(err), codes.NotFound)
}

func (suite *UnitTestSuite) TestDeleteNikkahProfileSuccess() {
	_, err := suite.StorageManager.CreateUser(GetUserProto(UserEmail, Username), Password)
	suite.Assert().Equal(status.Code(err), codes.OK)
	_, err = suite.StorageManager.CreateNikkahProfile(GetNikkahProfileProto())
	suite.Assert().Equal(status.Code(err), codes.OK)

	err = suite.StorageManager.DeleteNikkahProfile(DefaultId)
	suite.Assert().Equal(status.Code(err), codes.OK)
}

func (suite *UnitTestSuite) TestDeleteNikkahProfileNotFound() {
	err := suite.StorageManager.DeleteNikkahProfile(DefaultId)
	suite.Assert().Equal(status.Code(err), codes.NotFound)
}

func (suite *UnitTestSuite) TestCreateNikkahLikeSuccess() {
	_, err := suite.StorageManager.CreateUser(GetUserProto(UserEmail, Username), Password)
	suite.Assert().Equal(status.Code(err), codes.OK)
	_, err = suite.StorageManager.CreateNikkahProfile(GetNikkahProfileProto())
	suite.Assert().Equal(status.Code(err), codes.OK)
	_, err = suite.StorageManager.CreateNikkahLike(GetNikkahLikeProto())
	suite.Assert().Equal(status.Code(err), codes.OK)
}

func (suite *UnitTestSuite) TestCreateNikkahLikeWithNil() {
	_, err := suite.StorageManager.CreateNikkahLike(nil)
	suite.Assert().Equal(status.Code(err), codes.InvalidArgument)
}

func (suite *UnitTestSuite) TestCreateNikkahLikeWithNoLikerProfileId() {
	like := GetNikkahLikeProto()
	like.LikerProfileId = ""
	_, err := suite.StorageManager.CreateNikkahLike(like)
	suite.Assert().Equal(status.Code(err), codes.InvalidArgument)
}

func (suite *UnitTestSuite) TestCreateNikkahLikeWithNoLikedProfileId() {
	_, err := suite.StorageManager.CreateUser(GetUserProto(UserEmail, Username), Password)
	suite.Assert().Equal(status.Code(err), codes.OK)
	_, err = suite.StorageManager.CreateNikkahProfile(GetNikkahProfileProto())
	suite.Assert().Equal(status.Code(err), codes.OK)
	like := GetNikkahLikeProto()
	like.LikedProfileId = ""
	_, err = suite.StorageManager.CreateNikkahLike(like)
	suite.Assert().Equal(status.Code(err), codes.InvalidArgument)
}

func (suite *UnitTestSuite) TestCreateNikkahProfileDoesntExist() {
	_, err := suite.StorageManager.CreateNikkahLike(GetNikkahLikeProto())
	suite.Assert().Equal(status.Code(err), codes.NotFound)
}

func (suite *UnitTestSuite) TestUpdateNikkahLikeSuccess() {
	_, err := suite.StorageManager.CreateUser(GetUserProto(UserEmail, Username), Password)
	suite.Assert().Equal(status.Code(err), codes.OK)
	_, err = suite.StorageManager.CreateNikkahProfile(GetNikkahProfileProto())
	suite.Assert().Equal(status.Code(err), codes.OK)
	like, err := suite.StorageManager.CreateNikkahLike(GetNikkahLikeProto())
	suite.Assert().Equal(status.Code(err), codes.OK)

	like.Status = storage.LikeStatus(pb.NikkahLike_INITIATED)
	_, err = suite.StorageManager.UpdateNikkahLike(like.ToProto())
	suite.Assert().Equal(status.Code(err), codes.OK)
}

func (suite *UnitTestSuite) TestUpdateNikkahLikeNotFound() {
	like, err := suite.StorageManager.GetNikkahLike(DefaultId)
	suite.Assert().Equal(status.Code(err), codes.NotFound)
	suite.Nil(like)
}

func (suite *UnitTestSuite) TestUpdateNikkahLikeWithLikerProfileIdMismatch() {
	_, err := suite.StorageManager.CreateUser(GetUserProto(UserEmail, Username), Password)
	suite.Assert().Equal(status.Code(err), codes.OK)
	_, err = suite.StorageManager.CreateNikkahProfile(GetNikkahProfileProto())
	suite.Assert().Equal(status.Code(err), codes.OK)
	like, err := suite.StorageManager.CreateNikkahLike(GetNikkahLikeProto())
	suite.Assert().Equal(status.Code(err), codes.OK)

	like.LikerProfileID = "invalid"
	_, err = suite.StorageManager.UpdateNikkahLike(like.ToProto())
	suite.Assert().Equal(status.Code(err), codes.InvalidArgument)
}

func (suite *UnitTestSuite) TestUpdateNikkahLikeWithLikedProfileIdMismatch() {
	_, err := suite.StorageManager.CreateUser(GetUserProto(UserEmail, Username), Password)
	suite.Assert().Equal(status.Code(err), codes.OK)
	_, err = suite.StorageManager.CreateNikkahProfile(GetNikkahProfileProto())
	suite.Assert().Equal(status.Code(err), codes.OK)
	like, err := suite.StorageManager.CreateNikkahLike(GetNikkahLikeProto())
	suite.Assert().Equal(status.Code(err), codes.OK)

	like.LikedProfileID = "invalid"
	_, err = suite.StorageManager.UpdateNikkahLike(like.ToProto())
	suite.Assert().Equal(status.Code(err), codes.InvalidArgument)
}

func (suite *UnitTestSuite) TestGetNikkahLikeSuccess() {
	_, err := suite.StorageManager.CreateUser(GetUserProto(UserEmail, Username), Password)
	suite.Assert().Equal(status.Code(err), codes.OK)
	_, err = suite.StorageManager.CreateNikkahProfile(GetNikkahProfileProto())
	suite.Assert().Equal(status.Code(err), codes.OK)
	_, err = suite.StorageManager.CreateNikkahLike(GetNikkahLikeProto())
	suite.Assert().Equal(status.Code(err), codes.OK)

	like, err := suite.StorageManager.GetNikkahLike(DefaultId)
	suite.Assert().Equal(status.Code(err), codes.OK)
	AssertProtoEqual(suite.T(), *GetNikkahLikeProto(), like.ToProto(),
		pb.NikkahLike{}, protocmp.IgnoreFields(&pb.NikkahLike{}, "create_time", "update_time"))
}

func (suite *UnitTestSuite) TestGetNikkahLikeNotFound() {
	_, err := suite.StorageManager.GetNikkahLike(DefaultId)
	suite.Assert().Equal(status.Code(err), codes.NotFound)
}

func (suite *UnitTestSuite) TestCreateNikkahMatchSuccess() {
	_, err := suite.StorageManager.CreateUser(GetUserProto(UserEmail, Username), Password)
	suite.Assert().Equal(status.Code(err), codes.OK)
	_, err = suite.StorageManager.CreateNikkahProfile(GetNikkahProfileProto())
	suite.Assert().Equal(status.Code(err), codes.OK)
	_, err = suite.StorageManager.CreateNikkahMatch(GetNikkahMatchProto())
	suite.Assert().Equal(status.Code(err), codes.OK)
}

func (suite *UnitTestSuite) TestCreateNikkahMatchWithNil() {
	_, err := suite.StorageManager.CreateNikkahMatch(nil)
	suite.Assert().Equal(status.Code(err), codes.InvalidArgument)
}

func (suite *UnitTestSuite) TestCreateNikkahMatchWithNoInitiatorProfileId() {
	match := GetNikkahMatchProto()
	match.InitiatorProfileId = ""
	_, err := suite.StorageManager.CreateNikkahMatch(match)
	suite.Assert().Equal(status.Code(err), codes.InvalidArgument)
}

func (suite *UnitTestSuite) TestCreateNikkahMatchWithNoReceiverProfileId() {
	_, err := suite.StorageManager.CreateUser(GetUserProto(UserEmail, Username), Password)
	suite.Assert().Equal(status.Code(err), codes.OK)
	_, err = suite.StorageManager.CreateNikkahProfile(GetNikkahProfileProto())
	suite.Assert().Equal(status.Code(err), codes.OK)
	match := GetNikkahMatchProto()
	match.ReceiverProfileId = ""
	_, err = suite.StorageManager.CreateNikkahMatch(match)
	suite.Assert().Equal(status.Code(err), codes.InvalidArgument)
}

func (suite *UnitTestSuite) TestUpdateNikkahMatchSuccess() {
	_, err := suite.StorageManager.CreateUser(GetUserProto(UserEmail, Username), Password)
	suite.Assert().Equal(status.Code(err), codes.OK)
	_, err = suite.StorageManager.CreateNikkahProfile(GetNikkahProfileProto())
	suite.Assert().Equal(status.Code(err), codes.OK)
	match, err := suite.StorageManager.CreateNikkahMatch(GetNikkahMatchProto())
	suite.Assert().Equal(status.Code(err), codes.OK)

	match.Status = storage.MatchStatusEnded
	match, err = suite.StorageManager.UpdateNikkahMatch(match.ToProto())
	suite.Assert().Equal(status.Code(err), codes.OK)
	expected := GetNikkahMatchProto()
	expected.Status = pb.NikkahMatch_ENDED
	AssertProtoEqual(suite.T(), *expected, *match.ToProto(),
		pb.NikkahMatch{}, protocmp.IgnoreFields(&pb.NikkahMatch{}, "create_time", "update_time"))
}

func (suite *UnitTestSuite) TestUpdateNikkahMatchInitiatorProfileIdMismatch() {
	_, err := suite.StorageManager.CreateUser(GetUserProto(UserEmail, Username), Password)
	suite.Assert().Equal(status.Code(err), codes.OK)
	_, err = suite.StorageManager.CreateNikkahProfile(GetNikkahProfileProto())
	suite.Assert().Equal(status.Code(err), codes.OK)
	match, err := suite.StorageManager.CreateNikkahMatch(GetNikkahMatchProto())
	suite.Assert().Equal(status.Code(err), codes.OK)

	match.InitiatorProfileID, err = uuid.NewRandom()
	suite.Assert().Equal(status.Code(err), codes.OK)
	_, err = suite.StorageManager.UpdateNikkahMatch(match.ToProto())
	suite.Assert().Equal(status.Code(err), codes.InvalidArgument)
}

func (suite *UnitTestSuite) TestUpdateNikkahMatchReceiverProfileIdMismatch() {
	_, err := suite.StorageManager.CreateUser(GetUserProto(UserEmail, Username), Password)
	suite.Assert().Equal(status.Code(err), codes.OK)
	_, err = suite.StorageManager.CreateNikkahProfile(GetNikkahProfileProto())
	suite.Assert().Equal(status.Code(err), codes.OK)
	match, err := suite.StorageManager.CreateNikkahMatch(GetNikkahMatchProto())
	suite.Assert().Equal(status.Code(err), codes.OK)

	match.ReceiverProfileID, err = uuid.NewRandom()
	suite.Assert().Equal(status.Code(err), codes.OK)
	_, err = suite.StorageManager.UpdateNikkahMatch(match.ToProto())
	suite.Assert().Equal(status.Code(err), codes.InvalidArgument)
}

func (suite *UnitTestSuite) TestUpdateNikkahMatchNotFound() {
	like, err := suite.StorageManager.GetNikkahMatch(DefaultId)
	suite.Assert().Equal(status.Code(err), codes.NotFound)
	suite.Nil(like)
}
