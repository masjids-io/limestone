package storage

import (
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"

	pb "github.com/mnadev/limestone/proto"
)

type StorageManager struct {
	DB *gorm.DB
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

// CreateUser creates a User in the database for the given User and password.
func (s *StorageManager) CreateUser(up *pb.User, pwd string) (*User, error) {
	user, err := NewUser(up, pwd)
	if err != nil {
		return nil, err
	}

	result := s.DB.Create(user)
	if result.Error != nil {
		return nil, gormToGrpcError(result.Error)
	}
	return user, nil
}

// UpdateUser updates a User in the database for the given User if it exists.
// We expect the user to be authenticated via a token before this step.
func (s *StorageManager) UpdateUser(up *pb.User) (*User, error) {
	var old_user User
	result := s.DB.Where("id = ?", up.GetId()).First(&old_user)

	if result.Error != nil {
		return nil, gormToGrpcError(result.Error)
	}

	// TODO: add a field mask to update the information.
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
	result = s.DB.Where("id = ?", up.GetId()).First(&updated_user)
	if result.Error != nil {
		return nil, gormToGrpcError(result.Error)
	}
	return &updated_user, nil
}

// GetUser returns a User with the given id if it exists.
// We expect the user to be authenticated via a token before this step.
func (s *StorageManager) GetUser(id string) (*User, error) {
	var user User
	result := s.DB.Where("id = ?", id).First(&user)
	if result.Error != nil {
		return nil, gormToGrpcError(result.Error)
	}

	return &user, nil
}

// DeleteUser deletes a User with the given id if it exists.
// We expect the user to be authenticated via a token before this step.
func (s *StorageManager) DeleteUser(id string) error {
	user, err := s.GetUser(id)
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
func (s *StorageManager) GetAdhanFile(id string) (*AdhanFile, error) {
	var file AdhanFile
	result := s.DB.First(&file, "id = ?", id)

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

// CreateNikkahProfile creates a NikkahProfile in the database for the given NikkahProfile proto.
func (s *StorageManager) CreateNikkahProfile(np *pb.NikkahProfile) (*NikkahProfile, error) {
	if np == nil {
		return nil, status.Error(codes.InvalidArgument, "profile cannot be nil")
	}

	if np.GetUserId() == "" {
		return nil, status.Error(codes.InvalidArgument, "user cannot be nil")
	}

	// Check if the user exists.
	_, err := s.GetUser(np.GetUserId())
	if err != nil {
		return nil, err
	}

	profile, err := NewNikkahProfile(np)
	if err != nil {
		return nil, err
	}
	result := s.DB.Create(profile)

	if result.Error != nil {
		return nil, gormToGrpcError(result.Error)
	}
	return profile, nil
}

// UpdateNikkahProfile updates a NikkahProfile in the database for the given NikkahProfile proto.
func (s *StorageManager) UpdateNikkahProfile(np *pb.NikkahProfile) (*NikkahProfile, error) {
	if np == nil {
		return nil, status.Error(codes.InvalidArgument, "profile cannot be nil")
	}
	profile, err := s.GetNikkahProfile(np.GetId())
	if err != nil {
		return nil, err
	}

	if np.GetUserId() != profile.UserID {
		return nil, status.Error(codes.InvalidArgument, "user does not match profile")
	}

	new_profile, err := NewNikkahProfile(np)
	if err != nil {
		return nil, err
	}

	result := s.DB.Save(new_profile)
	if result.Error != nil {
		return nil, gormToGrpcError(result.Error)
	}
	return new_profile, nil
}

// GetNikkahProfile returns a NikkahProfile with the given id.
func (s *StorageManager) GetNikkahProfile(id string) (*NikkahProfile, error) {
	if id == "" {
		return nil, status.Error(codes.InvalidArgument, "id cannot be null")
	}

	var profile NikkahProfile
	result := s.DB.First(&profile, "id = ?", id)

	if result.Error != nil {
		return nil, gormToGrpcError(result.Error)
	}

	return &profile, nil
}

// DeleteNikkahProfile deletes a NikkahProfile with the given id.
func (s *StorageManager) DeleteNikkahProfile(id string) error {
	profile, err := s.GetNikkahProfile(id)
	if err != nil {
		return err
	}

	result := s.DB.Delete(profile, profile.ID)
	if result.Error != nil {
		return gormToGrpcError(result.Error)
	}
	return nil
}

// CreateNikkahLike creates a NikkahLike in the database for the given NikkahLike proto.
func (s *StorageManager) CreateNikkahLike(l *pb.NikkahLike) (*NikkahLike, error) {
	if l == nil {
		return nil, status.Error(codes.InvalidArgument, "like cannot be nil")
	}

	_, err := s.GetNikkahProfile(l.GetLikerProfileId())
	if err != nil {
		return nil, err
	}

	_, err = s.GetNikkahProfile(l.GetLikedProfileId())
	if err != nil {
		return nil, err
	}

	like, err := NewNikkahLike(l)
	if err != nil {
		return nil, err
	}
	result := s.DB.Create(like)
	if result.Error != nil {
		return nil, gormToGrpcError(result.Error)
	}
	return like, nil
}

// UpdateNikkahLike updates a NikkahLike in the database for the given NikkahLike proto.
func (s *StorageManager) UpdateNikkahLike(l *pb.NikkahLike) (*NikkahLike, error) {
	if l == nil {
		return nil, status.Error(codes.InvalidArgument, "like cannot be nil")
	}

	like, err := s.GetNikkahLike(l.GetLikeId())
	if err != nil {
		return nil, err
	}

	if like.LikerProfileID != l.GetLikerProfileId() {
		return nil, status.Error(codes.InvalidArgument, "liker profile does not match")
	}

	if like.LikedProfileID != l.GetLikedProfileId() {
		return nil, status.Error(codes.InvalidArgument, "liked profile does not match")
	}

	new_like, err := NewNikkahLike(l)
	if err != nil {
		return nil, err
	}

	result := s.DB.Save(new_like)
	if result.Error != nil {
		return nil, gormToGrpcError(result.Error)
	}
	return new_like, nil
}

// GetNikkahLike returns a NikkahLike with the given id.
func (s *StorageManager) GetNikkahLike(id string) (*NikkahLike, error) {
	if id == "" {
		return nil, status.Error(codes.InvalidArgument, "id cannot be empty")
	}
	var like NikkahLike
	result := s.DB.First(&like, "id = ?", id)

	if result.Error != nil {
		return nil, gormToGrpcError(result.Error)
	}
	return &like, nil
}

// CreateNikkahMatch creates a NikkahMatch in the database for the given NikkahMatch proto.
func (s *StorageManager) CreateNikkahMatch(m *pb.NikkahMatch) (*NikkahMatch, error) {
	if m == nil {
		return nil, status.Error(codes.InvalidArgument, "match cannot be nil")
	}

	_, err := s.GetNikkahProfile(m.GetInitiatorProfileId())
	if err != nil {
		return nil, err
	}

	_, err = s.GetNikkahProfile(m.GetReceiverProfileId())
	if err != nil {
		return nil, err
	}

	match, err := NewNikkahMatch(m)
	if err != nil {
		return nil, err
	}

	result := s.DB.Create(match)
	if result.Error != nil {
		return nil, gormToGrpcError(result.Error)
	}
	return match, nil
}

// UpdateNikkahMatch updates a NikkahMatch in the database for the given NikkahMatch proto.
func (s *StorageManager) UpdateNikkahMatch(m *pb.NikkahMatch) (*NikkahMatch, error) {
	if m == nil {
		return nil, status.Error(codes.InvalidArgument, "match cannot be nil")
	}

	match, err := s.GetNikkahMatch(m.GetMatchId())
	if err != nil {
		return nil, err
	}

	if match.InitiatorProfileID.String() != m.GetInitiatorProfileId() {
		return nil, status.Error(codes.InvalidArgument, "initiator profile does not match")
	}

	if match.ReceiverProfileID.String() != m.GetReceiverProfileId() {
		return nil, status.Error(codes.InvalidArgument, "receiver profile does not match")
	}

	new_match, err := NewNikkahMatch(m)
	if err != nil {
		return nil, err
	}

	result := s.DB.Save(new_match)
	if result.Error != nil {
		return nil, gormToGrpcError(result.Error)
	}
	return new_match, nil
}

// GetNikkahMatch returns a NikkahMatch with the given id.
func (s *StorageManager) GetNikkahMatch(id string) (*NikkahMatch, error) {
	if id == "" {
		return nil, status.Error(codes.InvalidArgument, "id cannot be empty")
	}
	var match NikkahMatch
	result := s.DB.First(&match, "id = ?", id)

	if result.Error != nil {
		return nil, gormToGrpcError(result.Error)
	}
	return &match, nil
}

func (s *StorageManager) GetUserByUsername(username string) (*User, error) {
	var user User
	result := s.DB.Where("username = ?", username).First(&user)
	if result.Error != nil {
		return nil, gormToGrpcError(result.Error)
	}

	return &user, nil
}

func (s *StorageManager) GetUserByEmail(email string) (*User, error) {
	var user User
	result := s.DB.Where("email = ?", email).First(&user)
	if result.Error != nil {
		return nil, gormToGrpcError(result.Error)
	}

	return &user, nil
}
