package storage

import (
	"database/sql"
	"testing"

	userpb "github.com/mnadev/limestone/user/proto"
	_ "github.com/proullon/ramsql/driver"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
		Email:       email,
		Username:    Username,
		FirstName:   FirstName,
		LastName:    LastName,
		PhoneNumber: PhoneNumber,
		// Gender:      userpb.User_MALE,
	}
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

	var user *User = nil
	tx := db.Table("users").Take(user)
	tx.Commit()

	if user == nil {
		db.AutoMigrate(&User{})
	}

	return &StorageManager{
		DB: db,
	}, nil
}

func TestCreateUser_Success(t *testing.T) {
	s, err := InitStorageManager("CreateUser_Success")
	require.Nil(t, err)

	err = s.CreateUser(GetUserProto(UserEmail), Password)
	assert.Nil(t, err)
}

func TestCreateUser_PasswordTooShort(t *testing.T) {
	s, err := InitStorageManager("CreateUser_PasswordTooShort")
	require.Nil(t, err)

	err = s.CreateUser(GetUserProto(UserEmail), BadPassword)
	assert.Error(t, err)
}

func TestGetUserWithEmail_Success(t *testing.T) {
	s, err := InitStorageManager("GetUserWithEmail_Success")
	require.Nil(t, err)

	err = s.CreateUser(GetUserProto(UserEmail), Password)
	assert.Nil(t, err)

	user, err := s.GetUserWithEmail(UserEmail, Password)
	assert.Equal(t, *GetUserProto(UserEmail), *user.ToProto())
	assert.Nil(t, err)
}

func TestGetUserWithEmail_BadPassword(t *testing.T) {
	s, err := InitStorageManager("GetUserWithEmail_BadPassword")
	require.Nil(t, err)

	err = s.CreateUser(GetUserProto(UserEmail), Password)
	assert.Nil(t, err)

	user, err := s.GetUserWithEmail(UserEmail, BadPassword)
	assert.Error(t, err)
	assert.Nil(t, user)
}

func TestGetUserWithEmail_NotFound(t *testing.T) {
	s, err := InitStorageManager("GetUserWithEmail_NotFound")
	require.Nil(t, err)

	err = s.CreateUser(GetUserProto(UserEmail), Password)
	assert.Nil(t, err)

	user, err := s.GetUserWithEmail("a@example.com", Password)
	assert.Error(t, err)
	assert.Nil(t, user)
}

func TestGetUserWithUsername_Success(t *testing.T) {
	s, err := InitStorageManager("GetUserWithUsername_Success")
	require.Nil(t, err)

	err = s.CreateUser(GetUserProto(UserEmail), Password)
	assert.Nil(t, err)

	user, err := s.GetUserWithUsername(Username, Password)
	assert.Equal(t, *GetUserProto(UserEmail), *user.ToProto())
	assert.Nil(t, err)
}

func TestGetUserWithUsername_BadPassword(t *testing.T) {
	s, err := InitStorageManager("GetUserWithUsername_BadPassword")
	require.Nil(t, err)

	err = s.CreateUser(GetUserProto(UserEmail), Password)
	assert.Nil(t, err)

	user, err := s.GetUserWithUsername(Username, BadPassword)
	assert.Error(t, err)
	assert.Nil(t, user)
}

func TestGetUserWithUsername_NotFound(t *testing.T) {
	s, err := InitStorageManager("GetUserWithUsername_NotFound")
	require.Nil(t, err)

	err = s.CreateUser(GetUserProto(UserEmail), Password)
	assert.Nil(t, err)

	user, err := s.GetUserWithUsername("notcoolguy1234", Password)
	assert.Error(t, err)
	assert.Nil(t, user)
}
