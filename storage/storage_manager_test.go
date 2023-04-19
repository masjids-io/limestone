package storage

import (
	"database/sql"
	"testing"
	"time"

	userpb "github.com/mnadev/limestone/user/proto"
	_ "github.com/proullon/ramsql/driver"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	_ "google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
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

func GetUserProto(email string) *userpb.User {
	return &userpb.User{
		UserId:      "00000000-0000-0000-0000-000000000000",
		Email:       email,
		Username:    Username,
		FirstName:   FirstName,
		LastName:    LastName,
		PhoneNumber: PhoneNumber,
		Gender:      userpb.User_FEMALE,
	}
}

func AssertUserProtoEqual(t *testing.T, expected, actual *userpb.User) bool {
	actual_without_timestamp := userpb.User{
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

func InitStorageManager(testName string) (*StorageManager, error) {
	sqlDB, err := sql.Open("ramsql", "Test"+testName)

	if err != nil {
		return nil, err
	}

	db, err := gorm.Open(postgres.New(postgres.Config{
		Conn: sqlDB,
	}), &gorm.Config{})

	if err != nil {
		return nil, err
	}

	db.AutoMigrate(&User{})

	return &StorageManager{
		DB: db,
	}, nil
}

func ClearUsersTable(s *StorageManager) {
	s.DB.Where("id = *").Delete(&User{})
}

func TestCreateUser_Success(t *testing.T) {
	s, err := InitStorageManager("CreateUser_Success")
	require.Nil(t, err)

	user, err := s.CreateUser(GetUserProto(UserEmail), Password)
	AssertUserProtoEqual(t, GetUserProto(UserEmail), user.ToProto())
	assert.Nil(t, err)

	ClearUsersTable(s)
}

func TestCreateUser_PasswordTooShort(t *testing.T) {
	s, err := InitStorageManager("CreateUser_PasswordTooShort")
	require.Nil(t, err)

	user, err := s.CreateUser(GetUserProto(UserEmail), BadPassword)
	assert.Nil(t, user)
	assert.Error(t, err)

	ClearUsersTable(s)
}

func TestUpdateUserWithEmail_Success(t *testing.T) {
	s, err := InitStorageManager("UpdateUserWithEmail_Success")
	require.Nil(t, err)

	user, err := s.CreateUser(GetUserProto(UserEmail), Password)
	assert.Nil(t, err)

	user.Email = "a@example.com"
	user, err = s.UpdateUser(user.ToProto(), Password)
	AssertUserProtoEqual(t, GetUserProto("a@example.com"), user.ToProto())
	assert.Nil(t, err)

	ClearUsersTable(s)
}

func TestUpdateUserWithEmail_BadPassword(t *testing.T) {
	s, err := InitStorageManager("UpdateUserWithEmail_BadPassword")
	require.Nil(t, err)

	_, err = s.CreateUser(GetUserProto(UserEmail), Password)
	assert.Nil(t, err)

	user, err := s.UpdateUser(GetUserProto("a@example.com"), BadPassword)
	assert.Error(t, err)
	assert.Nil(t, user)

	ClearUsersTable(s)
}

func TestUpdateUserWithEmail_NotFound(t *testing.T) {
	s, err := InitStorageManager("UpdateUserWithEmail_NotFound")
	require.Nil(t, err)

	user, err := s.UpdateUser(GetUserProto(UserEmail), Password)
	assert.Error(t, err)
	assert.Nil(t, user)

	ClearUsersTable(s)
}

func TestGetUserWithEmail_Success(t *testing.T) {
	s, err := InitStorageManager("GetUserWithEmail_Success")
	require.Nil(t, err)

	_, err = s.CreateUser(GetUserProto(UserEmail), Password)
	assert.Nil(t, err)

	user, err := s.GetUserWithEmail(UserEmail, Password)
	AssertUserProtoEqual(t, GetUserProto(UserEmail), user.ToProto())
	assert.Nil(t, err)

	ClearUsersTable(s)
}

func TestGetUserWithEmail_BadPassword(t *testing.T) {
	s, err := InitStorageManager("GetUserWithEmail_BadPassword")
	require.Nil(t, err)

	_, err = s.CreateUser(GetUserProto(UserEmail), Password)
	assert.Nil(t, err)

	user, err := s.GetUserWithEmail(UserEmail, BadPassword)
	assert.Error(t, err)
	assert.Nil(t, user)

	ClearUsersTable(s)
}

func TestGetUserWithEmail_NotFound(t *testing.T) {
	s, err := InitStorageManager("GetUserWithEmail_NotFound")
	require.Nil(t, err)

	_, err = s.CreateUser(GetUserProto(UserEmail), Password)
	assert.Nil(t, err)

	user, err := s.GetUserWithEmail("a@example.com", Password)
	assert.Error(t, err)
	assert.Nil(t, user)

	ClearUsersTable(s)
}

func TestGetUserWithUsername_Success(t *testing.T) {
	s, err := InitStorageManager("GetUserWithUsername_Success")
	require.Nil(t, err)

	_, err = s.CreateUser(GetUserProto(UserEmail), Password)
	assert.Nil(t, err)

	user, err := s.GetUserWithUsername(Username, Password)
	AssertUserProtoEqual(t, GetUserProto(UserEmail), user.ToProto())
	assert.Nil(t, err)

	ClearUsersTable(s)
}

func TestGetUserWithUsername_BadPassword(t *testing.T) {
	s, err := InitStorageManager("GetUserWithUsername_BadPassword")
	require.Nil(t, err)

	_, err = s.CreateUser(GetUserProto(UserEmail), Password)
	assert.Nil(t, err)

	user, err := s.GetUserWithUsername(Username, BadPassword)
	assert.Error(t, err)
	assert.Nil(t, user)

	ClearUsersTable(s)
}

func TestGetUserWithUsername_NotFound(t *testing.T) {
	s, err := InitStorageManager("GetUserWithUsername_NotFound")
	require.Nil(t, err)

	_, err = s.CreateUser(GetUserProto(UserEmail), Password)
	assert.Nil(t, err)

	user, err := s.GetUserWithUsername("notcoolguy1234", Password)
	assert.Error(t, err)
	assert.Nil(t, user)

	ClearUsersTable(s)
}

func TestDeleteUserWithEmail_Success(t *testing.T) {
	s, err := InitStorageManager("DeleteUserWithEmail_Success")
	require.Nil(t, err)

	_, err = s.CreateUser(GetUserProto(UserEmail), Password)
	assert.Nil(t, err)

	err = s.DeleteUserWithEmail(UserEmail, Password)
	assert.Nil(t, err)

	ClearUsersTable(s)
}

func TestDeleteUserWithEmail_BadPassword(t *testing.T) {
	s, err := InitStorageManager("DeleteUserWithEmail_BadPassword")
	require.Nil(t, err)

	_, err = s.CreateUser(GetUserProto(UserEmail), Password)
	assert.Nil(t, err)

	err = s.DeleteUserWithEmail(UserEmail, BadPassword)
	assert.Error(t, err)

	ClearUsersTable(s)
}

func TestDeleteUserWithEmail_NotFound(t *testing.T) {
	s, err := InitStorageManager("DeleteUserWithEmail_NotFound")
	require.Nil(t, err)

	_, err = s.CreateUser(GetUserProto(UserEmail), Password)
	assert.Nil(t, err)

	err = s.DeleteUserWithEmail("a@example.com", Password)
	assert.Error(t, err)

	ClearUsersTable(s)
}

func TestDeleteUserWithUsername_Success(t *testing.T) {
	s, err := InitStorageManager("DeleteUserWithUsername_Success")
	require.Nil(t, err)

	_, err = s.CreateUser(GetUserProto(UserEmail), Password)
	assert.Nil(t, err)

	err = s.DeleteUserWithUsername(Username, Password)
	assert.Nil(t, err)

	ClearUsersTable(s)
}

func TestDeleteUserWithUsername_BadPassword(t *testing.T) {
	s, err := InitStorageManager("DeleteUserWithUsername_BadPassword")
	require.Nil(t, err)

	_, err = s.CreateUser(GetUserProto(UserEmail), Password)
	assert.Nil(t, err)

	err = s.DeleteUserWithUsername(Username, BadPassword)
	assert.Error(t, err)

	ClearUsersTable(s)
}

func TestDeleteUserWithUsername_NotFound(t *testing.T) {
	s, err := InitStorageManager("DeleteUserWithUsername_NotFound")
	require.Nil(t, err)

	_, err = s.CreateUser(GetUserProto(UserEmail), Password)
	assert.Nil(t, err)

	err = s.DeleteUserWithUsername("notcoolguy1234", Password)
	assert.Error(t, err)

	ClearUsersTable(s)
}
