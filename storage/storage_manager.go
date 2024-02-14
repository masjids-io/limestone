package storage

import (
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
	"regexp"
	"strconv"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"

	"github.com/go-redis/redis"
	"github.com/mnadev/limestone/auth"
	pb "github.com/mnadev/limestone/proto"
)

type StorageManager struct {
	DB    *gorm.DB
	Cache *redis.Client
}

func gormToGrpcError(err error) error {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return status.Error(codes.NotFound, "requested entity was not found.")
	}
	if errors.Is(err, gorm.ErrDuplicatedKey) {
		return status.Error(codes.AlreadyExists, "entity already exists")
	}
	return status.Error(codes.Internal, "an internal error occurred")
}
func (s *StorageManager) SetCode(email string) (string, error) {
	randomNumber := GenerateRandomNumber()
	if err := Set(s.Cache, email, randomNumber, 10*time.Minute); err != nil {
		// handle error
		fmt.Println(err)
		return "", err
	}
	return randomNumber, nil
}

func GenerateRandomNumber() string {
	// Generate a random number between 100000 and 999999 (6 digits)
	num, err := rand.Int(rand.Reader, big.NewInt(900000))
	if err != nil {
		fmt.Println(err)
		return ""
	}
	return strconv.Itoa(int((num.Int64()) + 100000))
}

// CreateUser creates a User in the database for the given User and password
func (s *StorageManager) CreateUser(up *pb.User, pwd string) (*User, error) {
	//check validity of email
	pattern := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

	if !pattern.MatchString(up.GetEmail()) {
		return nil, status.Error(codes.InvalidArgument, "Incorrect email format")
	}
	//check for user existence
	var user User
	err := s.DB.Where("email = ?", up.GetEmail()).First(&user).Error
	if err != nil && err == gorm.ErrRecordNotFound {
		//db user instance created
		u, err := NewUser(up, pwd)
		if err != nil {
			return nil, err
		}
		//db entry
		result := s.DB.Create(u)
		if result.Error != nil {
			return nil, gormToGrpcError(result.Error)
		}
		//concurrent email verification dispatched
		code, err := s.SetCode(up.GetEmail())
		if err != nil {
			return nil, err
		}
		go auth.EmailVerification(up.GetEmail(), code)
		return u, nil
	}
	return nil, status.Error(codes.AlreadyExists, "user already exists")
}

// UpdateUser updates a User in the database for the given User and password if it exists
func (s *StorageManager) UpdateUser(up *pb.User, pwd string) (*User, error) {
	var old_user User
	result := s.DB.Where("id = ?", up.GetUserId()).First(&old_user)

	if result.Error != nil {
		return nil, gormToGrpcError(result.Error)
	}

	if auth.CheckPassword(pwd, old_user.HashedPassword) != nil {
		return nil, status.Error(codes.PermissionDenied, "password did not match")
	}

	// TODO: add a field mask to update the information
	result = s.DB.Model(old_user).Where("id = ?", old_user.ID).
		Updates(
			map[string]interface{}{
				"email":        up.GetEmail(),
				"username":     up.GetUsername(),
				"is_verified":  up.GetIsEmailVerified(),
				"first_name":   up.GetFirstName(),
				"last_name":    up.GetLastName(),
				"phone_number": up.GetPhoneNumber(),
				"gender":       gender(up.GetGender().String()),
			})

	if result.Error != nil {
		return nil, gormToGrpcError(result.Error)
	}

	var updated_user User
	result = s.DB.Where("id = ?", up.GetUserId()).First(&updated_user)
	if result.Error != nil {
		return nil, gormToGrpcError(result.Error)
	}
	return &updated_user, nil
}
func (s *StorageManager) VerifyUser(email string, code int64) (*pb.VerifyResponse, error) {
	//fetch code from redis
	var storedCode string
	if err := Get(s.Cache, email, &storedCode); err != nil {
		return nil, status.Error(codes.Internal, "Code expired. Re-login.")
	}
	//check if code matches
	strCode := strconv.Itoa(int(code))
	if storedCode != strCode {
		return nil, status.Error(codes.PermissionDenied, "Incorrect code")
	}
	//update in db
	result := s.DB.Model(&User{}).Where("email = ?", email).Updates(map[string]interface{}{
		"is_verified": true,
	})
	if result.Error != nil {
		return nil, gormToGrpcError(result.Error)
	}
	// delete code from cache
	Delete(s.Cache, email)
	return &pb.VerifyResponse{
		Response: "Verified",
	}, nil
}

// GetUserWithEmail returns a User with the given email and password if it exists
func (s *StorageManager) LoginUser(email string, pwd string) (*pb.Tokens, error) {
	var user User
	result := s.DB.Where("email = ?", email).First(&user)
	if result.Error != nil {
		return nil, gormToGrpcError(result.Error)
	}
	// password check
	if auth.CheckPassword(pwd, user.HashedPassword) != nil {
		return nil, status.Error(codes.PermissionDenied, "password did not match")
	}
	// check verification
	if !user.IsVerified {
		//concurrent email verification dispatched
		code, err := s.SetCode(user.Email)
		if err != nil {
			return nil, err
		}
		go auth.EmailVerification(user.Email, code)
		return nil, status.Error(codes.PermissionDenied, "Account unverified")
	}
	// if verified, return tokens
	access, _ := auth.CreateJWT(user.ID.String(), user.Email, auth.TokenType(0))  //access token
	refresh, _ := auth.CreateJWT(user.ID.String(), user.Email, auth.TokenType(1)) //refresh token
	return &pb.Tokens{
		AccessToken:  access,
		RefreshToken: refresh,
	}, nil
}

// GetUserWithEmail returns a User with the given email and password if it exists
func (s *StorageManager) GetUserWithEmail(email string, pwd string) (*User, error) {
	var user User
	result := s.DB.Where("email = ?", email).First(&user)
	if result.Error != nil {
		return nil, gormToGrpcError(result.Error)
	}
	// password check
	if auth.CheckPassword(pwd, user.HashedPassword) != nil {
		return nil, status.Error(codes.PermissionDenied, "password did not match")
	}

	return &user, nil
}

// GetUserWithUsername returns a User with the given username and password if it exists
func (s *StorageManager) GetUserWithUsername(username string, pwd string) (*User, error) {
	var user User
	result := s.DB.Where("username = ?", username).First(&user)
	if result.Error != nil {
		return nil, gormToGrpcError(result.Error)
	}

	if auth.CheckPassword(pwd, user.HashedPassword) != nil {
		return nil, status.Error(codes.PermissionDenied, "password did not match")
	}

	return &user, nil
}

// DeleteUserWithEmail deletes a User with the given email and password if it exists
func (s *StorageManager) DeleteUserWithEmail(email string, pwd string) error {
	user, err := s.GetUserWithEmail(email, pwd)
	if err != nil {
		return err
	}

	result := s.DB.Delete(user, user.ID)
	if result.Error != nil {
		return gormToGrpcError(result.Error)
	}
	return nil
}

// DeleteUserWithUsername deletes a User with the given username and password if it exists
func (s *StorageManager) DeleteUserWithUsername(username string, pwd string) error {
	user, err := s.GetUserWithUsername(username, pwd)
	if err != nil {
		return err
	}

	result := s.DB.Delete(user, user.ID)
	if result.Error != nil {
		return gormToGrpcError(result.Error)
	}
	return nil
}

// CreateMasjid creates a Masjid in the database for the given Masjid proto.
func (s *StorageManager) CreateMasjid(mp *pb.Masjid) (*Masjid, error) {
	masjid := NewMasjid(mp)
	result := s.DB.Create(masjid)

	if result.Error != nil {
		return nil, gormToGrpcError(result.Error)
	}
	return masjid, nil
}

// UpdateMasjid updates a Masjid in the database for the given Masjid proto.
func (s *StorageManager) UpdateMasjid(mp *pb.Masjid) (*Masjid, error) {
	var old_masjid Masjid
	result := s.DB.Where("id = ?", mp.GetId()).First(&old_masjid)
	if result.Error != nil {
		return nil, gormToGrpcError(result.Error)
	}

	new_masjid := NewMasjid(mp)
	result = s.DB.Save(new_masjid)
	if result.Error != nil {
		return nil, gormToGrpcError(result.Error)
	}

	return new_masjid, nil
}

// GetMasjid returns a Masjid with the given id.
func (s *StorageManager) GetMasjid(id string) (*Masjid, error) {
	var masjid Masjid
	result := s.DB.First(&masjid, "id = ?", id)

	if result.Error != nil {
		return nil, gormToGrpcError(result.Error)
	}

	return &masjid, nil
}

// DeleteMasjid deletes a Masjid with the given id.
func (s *StorageManager) DeleteMasjid(id string) error {
	masjid, err := s.GetMasjid(id)
	if err != nil {
		return err
	}

	result := s.DB.Delete(masjid, masjid.ID)
	if result.Error != nil {
		return gormToGrpcError(result.Error)
	}
	return nil
}

// CreateEvent creates a Event in the database for the given Event proto.
func (s *StorageManager) CreateEvent(e *pb.Event) (*Event, error) {
	event, err := NewEvent(e)
	if err != nil {
		return nil, err
	}

	result := s.DB.Create(event)

	if result.Error != nil {
		return nil, gormToGrpcError(result.Error)
	}
	return event, nil
}

// UpdateEvent updates a Event in the database for the given Event proto.
func (s *StorageManager) UpdateEvent(e *pb.Event) (*Event, error) {
	var old_event Event
	result := s.DB.Where("id = ?", e.GetId()).First(&old_event)

	if result.Error != nil {
		return nil, gormToGrpcError(result.Error)
	}

	new_event, err := NewEvent(e)
	if err != nil {
		return nil, err
	}

	result = s.DB.Save(new_event)
	if result.Error != nil {
		return nil, gormToGrpcError(result.Error)
	}

	return new_event, nil
}

// GetEvent returns a Event with the given id.
func (s *StorageManager) GetEvent(id string) (*Event, error) {
	var event Event
	result := s.DB.First(&event, "id = ?", id)

	if result.Error != nil {
		return nil, gormToGrpcError(result.Error)
	}

	return &event, nil
}

// DeleteEvent deletes a Event with the given id.
func (s *StorageManager) DeleteEvent(id string) error {
	event, err := s.GetEvent(id)
	if err != nil {
		return err
	}

	result := s.DB.Delete(event, event.ID)
	if result.Error != nil {
		return gormToGrpcError(result.Error)
	}
	return nil
}

// CreateAdhanFile creates a AdhanFile in the database for the given AdhanFile proto.
func (s *StorageManager) CreateAdhanFile(a *pb.AdhanFile) (*AdhanFile, error) {
	file := NewAdhanFile(a)
	result := s.DB.Create(file)

	if result.Error != nil {
		return nil, gormToGrpcError(result.Error)
	}
	return file, nil
}

// UpdateAdhanFile updates a AdhanFile in the database for the given AdhanFile proto.
func (s *StorageManager) UpdateAdhanFile(a *pb.AdhanFile) (*AdhanFile, error) {
	var old_file AdhanFile
	result := s.DB.Where("id = ?", a.GetId()).First(&old_file)
	if result.Error != nil {
		return nil, gormToGrpcError(result.Error)
	}

	new_file := NewAdhanFile(a)
	result = s.DB.Save(new_file)
	if result.Error != nil {
		return nil, gormToGrpcError(result.Error)
	}

	return new_file, nil
}

// GetAdhanFile returns a AdhanFile with the given id.
func (s *StorageManager) GetAdhanFile(masjid_id string) (*AdhanFile, error) {
	var file AdhanFile
	result := s.DB.First(&file, "masjid_id = ?", masjid_id)

	if result.Error != nil {
		return nil, gormToGrpcError(result.Error)
	}

	return &file, nil
}

// DeleteAdhanFile deletes a AdhanFile with the given id.
func (s *StorageManager) DeleteAdhanFile(id string) error {
	file, err := s.GetAdhanFile(id)
	if err != nil {
		return err
	}

	result := s.DB.Delete(file, file.ID)
	if result.Error != nil {
		return gormToGrpcError(result.Error)
	}
	return nil
}
