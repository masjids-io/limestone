// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.6
// 	protoc        (unknown)
// source: user_service.proto

package proto

import (
	_ "google.golang.org/genproto/googleapis/api/annotations"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
	reflect "reflect"
	sync "sync"
	unsafe "unsafe"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// Defines the different possible roles a user might have in a masjid.
type MasjidRole_Role int32

const (
	// Default value.
	MasjidRole_ROLE_UNSPECIFIED MasjidRole_Role = 0
	// This role specifies someone who is just a member of the masjid.
	MasjidRole_MASJID_MEMBER MasjidRole_Role = 1
	// This role specifies someone who has some responsibility at the masjid
	// but is not an administrator.
	MasjidRole_MASJID_VOLUNTEER MasjidRole_Role = 2
	// This role specifies someone at the masjid who is involved in a
	// high-level administrative capacity.
	MasjidRole_MASJID_ADMIN MasjidRole_Role = 3
	// This role specifies someone at the masjid who is involved in a
	// religious capacity.
	MasjidRole_MASJID_IMAM MasjidRole_Role = 4
)

// Enum value maps for MasjidRole_Role.
var (
	MasjidRole_Role_name = map[int32]string{
		0: "ROLE_UNSPECIFIED",
		1: "MASJID_MEMBER",
		2: "MASJID_VOLUNTEER",
		3: "MASJID_ADMIN",
		4: "MASJID_IMAM",
	}
	MasjidRole_Role_value = map[string]int32{
		"ROLE_UNSPECIFIED": 0,
		"MASJID_MEMBER":    1,
		"MASJID_VOLUNTEER": 2,
		"MASJID_ADMIN":     3,
		"MASJID_IMAM":      4,
	}
)

func (x MasjidRole_Role) Enum() *MasjidRole_Role {
	p := new(MasjidRole_Role)
	*p = x
	return p
}

func (x MasjidRole_Role) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (MasjidRole_Role) Descriptor() protoreflect.EnumDescriptor {
	return file_user_service_proto_enumTypes[0].Descriptor()
}

func (MasjidRole_Role) Type() protoreflect.EnumType {
	return &file_user_service_proto_enumTypes[0]
}

func (x MasjidRole_Role) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use MasjidRole_Role.Descriptor instead.
func (MasjidRole_Role) EnumDescriptor() ([]byte, []int) {
	return file_user_service_proto_rawDescGZIP(), []int{0, 0}
}

// Defines an enum representing the gender of the user.
type User_Gender int32

const (
	// Default value.
	User_GENDER_UNSPECIFIED User_Gender = 0
	User_MALE               User_Gender = 1
	User_FEMALE             User_Gender = 2
)

// Enum value maps for User_Gender.
var (
	User_Gender_name = map[int32]string{
		0: "GENDER_UNSPECIFIED",
		1: "MALE",
		2: "FEMALE",
	}
	User_Gender_value = map[string]int32{
		"GENDER_UNSPECIFIED": 0,
		"MALE":               1,
		"FEMALE":             2,
	}
)

func (x User_Gender) Enum() *User_Gender {
	p := new(User_Gender)
	*p = x
	return p
}

func (x User_Gender) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (User_Gender) Descriptor() protoreflect.EnumDescriptor {
	return file_user_service_proto_enumTypes[1].Descriptor()
}

func (User_Gender) Type() protoreflect.EnumType {
	return &file_user_service_proto_enumTypes[1]
}

func (x User_Gender) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use User_Gender.Descriptor instead.
func (User_Gender) EnumDescriptor() ([]byte, []int) {
	return file_user_service_proto_rawDescGZIP(), []int{1, 0}
}

// Defines a specific role a user might have for a specific masjid.
type MasjidRole struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// The associated role.
	Role MasjidRole_Role `protobuf:"varint,1,opt,name=role,proto3,enum=limestone.MasjidRole_Role" json:"role,omitempty"`
	// The ID of the masjid associated with this role.
	MasjidId string `protobuf:"bytes,2,opt,name=masjid_id,json=masjidId,proto3" json:"masjid_id,omitempty"`
	// The ID of the associated user.
	UserId        string `protobuf:"bytes,3,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *MasjidRole) Reset() {
	*x = MasjidRole{}
	mi := &file_user_service_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *MasjidRole) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MasjidRole) ProtoMessage() {}

func (x *MasjidRole) ProtoReflect() protoreflect.Message {
	mi := &file_user_service_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MasjidRole.ProtoReflect.Descriptor instead.
func (*MasjidRole) Descriptor() ([]byte, []int) {
	return file_user_service_proto_rawDescGZIP(), []int{0}
}

func (x *MasjidRole) GetRole() MasjidRole_Role {
	if x != nil {
		return x.Role
	}
	return MasjidRole_ROLE_UNSPECIFIED
}

func (x *MasjidRole) GetMasjidId() string {
	if x != nil {
		return x.MasjidId
	}
	return ""
}

func (x *MasjidRole) GetUserId() string {
	if x != nil {
		return x.UserId
	}
	return ""
}

// Defines a User.
type User struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// The ID of the user. This field is output only.
	Id string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	// The email of the user. This field is required and must be unique.
	Email string `protobuf:"bytes,2,opt,name=email,proto3" json:"email,omitempty"`
	// The username of the user. This field is required and must be unique.
	Username string `protobuf:"bytes,3,opt,name=username,proto3" json:"username,omitempty"`
	// Whether or not the email of the user was verified.
	IsEmailVerified bool `protobuf:"varint,4,opt,name=is_email_verified,json=isEmailVerified,proto3" json:"is_email_verified,omitempty"`
	// The first name of the user.
	FirstName string `protobuf:"bytes,5,opt,name=first_name,json=firstName,proto3" json:"first_name,omitempty"`
	// The last name of the user.
	LastName string `protobuf:"bytes,6,opt,name=last_name,json=lastName,proto3" json:"last_name,omitempty"`
	// The phone number associated with the user.
	PhoneNumber string      `protobuf:"bytes,7,opt,name=phone_number,json=phoneNumber,proto3" json:"phone_number,omitempty"`
	Gender      User_Gender `protobuf:"varint,8,opt,name=gender,proto3,enum=limestone.User_Gender" json:"gender,omitempty"`
	// The create time of the user. This field is output only.
	CreateTime *timestamppb.Timestamp `protobuf:"bytes,9,opt,name=create_time,json=createTime,proto3" json:"create_time,omitempty"`
	// The update time of the user. This field is output only.
	UpdateTime    *timestamppb.Timestamp `protobuf:"bytes,10,opt,name=update_time,json=updateTime,proto3" json:"update_time,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *User) Reset() {
	*x = User{}
	mi := &file_user_service_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *User) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*User) ProtoMessage() {}

func (x *User) ProtoReflect() protoreflect.Message {
	mi := &file_user_service_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use User.ProtoReflect.Descriptor instead.
func (*User) Descriptor() ([]byte, []int) {
	return file_user_service_proto_rawDescGZIP(), []int{1}
}

func (x *User) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *User) GetEmail() string {
	if x != nil {
		return x.Email
	}
	return ""
}

func (x *User) GetUsername() string {
	if x != nil {
		return x.Username
	}
	return ""
}

func (x *User) GetIsEmailVerified() bool {
	if x != nil {
		return x.IsEmailVerified
	}
	return false
}

func (x *User) GetFirstName() string {
	if x != nil {
		return x.FirstName
	}
	return ""
}

func (x *User) GetLastName() string {
	if x != nil {
		return x.LastName
	}
	return ""
}

func (x *User) GetPhoneNumber() string {
	if x != nil {
		return x.PhoneNumber
	}
	return ""
}

func (x *User) GetGender() User_Gender {
	if x != nil {
		return x.Gender
	}
	return User_GENDER_UNSPECIFIED
}

func (x *User) GetCreateTime() *timestamppb.Timestamp {
	if x != nil {
		return x.CreateTime
	}
	return nil
}

func (x *User) GetUpdateTime() *timestamppb.Timestamp {
	if x != nil {
		return x.UpdateTime
	}
	return nil
}

type CreateUserRequest struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// The user to create.
	User *User `protobuf:"bytes,1,opt,name=user,proto3" json:"user,omitempty"`
	// The password for the user. This password is hashed before being
	// stored in the database.
	// The password must be at least 8 characters.
	Password      string `protobuf:"bytes,2,opt,name=password,proto3" json:"password,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *CreateUserRequest) Reset() {
	*x = CreateUserRequest{}
	mi := &file_user_service_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *CreateUserRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateUserRequest) ProtoMessage() {}

func (x *CreateUserRequest) ProtoReflect() protoreflect.Message {
	mi := &file_user_service_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateUserRequest.ProtoReflect.Descriptor instead.
func (*CreateUserRequest) Descriptor() ([]byte, []int) {
	return file_user_service_proto_rawDescGZIP(), []int{2}
}

func (x *CreateUserRequest) GetUser() *User {
	if x != nil {
		return x.User
	}
	return nil
}

func (x *CreateUserRequest) GetPassword() string {
	if x != nil {
		return x.Password
	}
	return ""
}

type GetUserRequest struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// The id of the user to get.
	Id            string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *GetUserRequest) Reset() {
	*x = GetUserRequest{}
	mi := &file_user_service_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetUserRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetUserRequest) ProtoMessage() {}

func (x *GetUserRequest) ProtoReflect() protoreflect.Message {
	mi := &file_user_service_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetUserRequest.ProtoReflect.Descriptor instead.
func (*GetUserRequest) Descriptor() ([]byte, []int) {
	return file_user_service_proto_rawDescGZIP(), []int{3}
}

func (x *GetUserRequest) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

type GetUserResponse struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// The authenticated user.
	User          *User `protobuf:"bytes,1,opt,name=user,proto3" json:"user,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *GetUserResponse) Reset() {
	*x = GetUserResponse{}
	mi := &file_user_service_proto_msgTypes[4]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *GetUserResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetUserResponse) ProtoMessage() {}

func (x *GetUserResponse) ProtoReflect() protoreflect.Message {
	mi := &file_user_service_proto_msgTypes[4]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetUserResponse.ProtoReflect.Descriptor instead.
func (*GetUserResponse) Descriptor() ([]byte, []int) {
	return file_user_service_proto_rawDescGZIP(), []int{4}
}

func (x *GetUserResponse) GetUser() *User {
	if x != nil {
		return x.User
	}
	return nil
}

type UpdateUserRequest struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// The updated user proto.
	User          *User `protobuf:"bytes,1,opt,name=user,proto3" json:"user,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *UpdateUserRequest) Reset() {
	*x = UpdateUserRequest{}
	mi := &file_user_service_proto_msgTypes[5]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *UpdateUserRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UpdateUserRequest) ProtoMessage() {}

func (x *UpdateUserRequest) ProtoReflect() protoreflect.Message {
	mi := &file_user_service_proto_msgTypes[5]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UpdateUserRequest.ProtoReflect.Descriptor instead.
func (*UpdateUserRequest) Descriptor() ([]byte, []int) {
	return file_user_service_proto_rawDescGZIP(), []int{5}
}

func (x *UpdateUserRequest) GetUser() *User {
	if x != nil {
		return x.User
	}
	return nil
}

type DeleteUserRequest struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// The id of the user to delete.
	Id            string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *DeleteUserRequest) Reset() {
	*x = DeleteUserRequest{}
	mi := &file_user_service_proto_msgTypes[6]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *DeleteUserRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeleteUserRequest) ProtoMessage() {}

func (x *DeleteUserRequest) ProtoReflect() protoreflect.Message {
	mi := &file_user_service_proto_msgTypes[6]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeleteUserRequest.ProtoReflect.Descriptor instead.
func (*DeleteUserRequest) Descriptor() ([]byte, []int) {
	return file_user_service_proto_rawDescGZIP(), []int{6}
}

func (x *DeleteUserRequest) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

// Empty response returned for DeleteUser RPC.
type DeleteUserResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *DeleteUserResponse) Reset() {
	*x = DeleteUserResponse{}
	mi := &file_user_service_proto_msgTypes[7]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *DeleteUserResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeleteUserResponse) ProtoMessage() {}

func (x *DeleteUserResponse) ProtoReflect() protoreflect.Message {
	mi := &file_user_service_proto_msgTypes[7]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeleteUserResponse.ProtoReflect.Descriptor instead.
func (*DeleteUserResponse) Descriptor() ([]byte, []int) {
	return file_user_service_proto_rawDescGZIP(), []int{7}
}

type AuthenticateUserRequest struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// Use either username or email for authentication.
	//
	// Types that are valid to be assigned to Identifier:
	//
	//	*AuthenticateUserRequest_Username
	//	*AuthenticateUserRequest_Email
	Identifier isAuthenticateUserRequest_Identifier `protobuf_oneof:"identifier"`
	// The password for authentication
	Password      string `protobuf:"bytes,3,opt,name=password,proto3" json:"password,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *AuthenticateUserRequest) Reset() {
	*x = AuthenticateUserRequest{}
	mi := &file_user_service_proto_msgTypes[8]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *AuthenticateUserRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AuthenticateUserRequest) ProtoMessage() {}

func (x *AuthenticateUserRequest) ProtoReflect() protoreflect.Message {
	mi := &file_user_service_proto_msgTypes[8]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AuthenticateUserRequest.ProtoReflect.Descriptor instead.
func (*AuthenticateUserRequest) Descriptor() ([]byte, []int) {
	return file_user_service_proto_rawDescGZIP(), []int{8}
}

func (x *AuthenticateUserRequest) GetIdentifier() isAuthenticateUserRequest_Identifier {
	if x != nil {
		return x.Identifier
	}
	return nil
}

func (x *AuthenticateUserRequest) GetUsername() string {
	if x != nil {
		if x, ok := x.Identifier.(*AuthenticateUserRequest_Username); ok {
			return x.Username
		}
	}
	return ""
}

func (x *AuthenticateUserRequest) GetEmail() string {
	if x != nil {
		if x, ok := x.Identifier.(*AuthenticateUserRequest_Email); ok {
			return x.Email
		}
	}
	return ""
}

func (x *AuthenticateUserRequest) GetPassword() string {
	if x != nil {
		return x.Password
	}
	return ""
}

type isAuthenticateUserRequest_Identifier interface {
	isAuthenticateUserRequest_Identifier()
}

type AuthenticateUserRequest_Username struct {
	Username string `protobuf:"bytes,1,opt,name=username,proto3,oneof"`
}

type AuthenticateUserRequest_Email struct {
	Email string `protobuf:"bytes,2,opt,name=email,proto3,oneof"`
}

func (*AuthenticateUserRequest_Username) isAuthenticateUserRequest_Identifier() {}

func (*AuthenticateUserRequest_Email) isAuthenticateUserRequest_Identifier() {}

type AuthenticateUserResponse struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// The access token for subsequent authenticated requests.
	AccessToken string `protobuf:"bytes,1,opt,name=access_token,json=accessToken,proto3" json:"access_token,omitempty"`
	// The refresh token to obtain a new access token when it expires.
	RefreshToken string `protobuf:"bytes,2,opt,name=refresh_token,json=refreshToken,proto3" json:"refresh_token,omitempty"`
	// The ID of the authenticated user (if successful).
	UserId        int32 `protobuf:"varint,3,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *AuthenticateUserResponse) Reset() {
	*x = AuthenticateUserResponse{}
	mi := &file_user_service_proto_msgTypes[9]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *AuthenticateUserResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AuthenticateUserResponse) ProtoMessage() {}

func (x *AuthenticateUserResponse) ProtoReflect() protoreflect.Message {
	mi := &file_user_service_proto_msgTypes[9]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AuthenticateUserResponse.ProtoReflect.Descriptor instead.
func (*AuthenticateUserResponse) Descriptor() ([]byte, []int) {
	return file_user_service_proto_rawDescGZIP(), []int{9}
}

func (x *AuthenticateUserResponse) GetAccessToken() string {
	if x != nil {
		return x.AccessToken
	}
	return ""
}

func (x *AuthenticateUserResponse) GetRefreshToken() string {
	if x != nil {
		return x.RefreshToken
	}
	return ""
}

func (x *AuthenticateUserResponse) GetUserId() int32 {
	if x != nil {
		return x.UserId
	}
	return 0
}

type RefreshTokenRequest struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// The refresh JWT for the user.
	RefreshToken  string `protobuf:"bytes,1,opt,name=refresh_token,json=refreshToken,proto3" json:"refresh_token,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *RefreshTokenRequest) Reset() {
	*x = RefreshTokenRequest{}
	mi := &file_user_service_proto_msgTypes[10]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *RefreshTokenRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RefreshTokenRequest) ProtoMessage() {}

func (x *RefreshTokenRequest) ProtoReflect() protoreflect.Message {
	mi := &file_user_service_proto_msgTypes[10]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RefreshTokenRequest.ProtoReflect.Descriptor instead.
func (*RefreshTokenRequest) Descriptor() ([]byte, []int) {
	return file_user_service_proto_rawDescGZIP(), []int{10}
}

func (x *RefreshTokenRequest) GetRefreshToken() string {
	if x != nil {
		return x.RefreshToken
	}
	return ""
}

type RefreshTokenResponse struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// The access token for subsequent authenticated requests.
	AccessToken   string `protobuf:"bytes,1,opt,name=access_token,json=accessToken,proto3" json:"access_token,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *RefreshTokenResponse) Reset() {
	*x = RefreshTokenResponse{}
	mi := &file_user_service_proto_msgTypes[11]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *RefreshTokenResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RefreshTokenResponse) ProtoMessage() {}

func (x *RefreshTokenResponse) ProtoReflect() protoreflect.Message {
	mi := &file_user_service_proto_msgTypes[11]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RefreshTokenResponse.ProtoReflect.Descriptor instead.
func (*RefreshTokenResponse) Descriptor() ([]byte, []int) {
	return file_user_service_proto_rawDescGZIP(), []int{11}
}

func (x *RefreshTokenResponse) GetAccessToken() string {
	if x != nil {
		return x.AccessToken
	}
	return ""
}

var File_user_service_proto protoreflect.FileDescriptor

const file_user_service_proto_rawDesc = "" +
	"\n" +
	"\x12user_service.proto\x12\tlimestone\x1a\x1cgoogle/api/annotations.proto\x1a\x17google/api/client.proto\x1a\x1fgoogle/api/field_behavior.proto\x1a\x1fgoogle/protobuf/timestamp.proto\"\xdc\x01\n" +
	"\n" +
	"MasjidRole\x12.\n" +
	"\x04role\x18\x01 \x01(\x0e2\x1a.limestone.MasjidRole.RoleR\x04role\x12\x1b\n" +
	"\tmasjid_id\x18\x02 \x01(\tR\bmasjidId\x12\x17\n" +
	"\auser_id\x18\x03 \x01(\tR\x06userId\"h\n" +
	"\x04Role\x12\x14\n" +
	"\x10ROLE_UNSPECIFIED\x10\x00\x12\x11\n" +
	"\rMASJID_MEMBER\x10\x01\x12\x14\n" +
	"\x10MASJID_VOLUNTEER\x10\x02\x12\x10\n" +
	"\fMASJID_ADMIN\x10\x03\x12\x0f\n" +
	"\vMASJID_IMAM\x10\x04\"\xb5\x03\n" +
	"\x04User\x12\x0e\n" +
	"\x02id\x18\x01 \x01(\tR\x02id\x12\x14\n" +
	"\x05email\x18\x02 \x01(\tR\x05email\x12\x1a\n" +
	"\busername\x18\x03 \x01(\tR\busername\x12*\n" +
	"\x11is_email_verified\x18\x04 \x01(\bR\x0fisEmailVerified\x12\x1d\n" +
	"\n" +
	"first_name\x18\x05 \x01(\tR\tfirstName\x12\x1b\n" +
	"\tlast_name\x18\x06 \x01(\tR\blastName\x12!\n" +
	"\fphone_number\x18\a \x01(\tR\vphoneNumber\x12.\n" +
	"\x06gender\x18\b \x01(\x0e2\x16.limestone.User.GenderR\x06gender\x12;\n" +
	"\vcreate_time\x18\t \x01(\v2\x1a.google.protobuf.TimestampR\n" +
	"createTime\x12;\n" +
	"\vupdate_time\x18\n" +
	" \x01(\v2\x1a.google.protobuf.TimestampR\n" +
	"updateTime\"6\n" +
	"\x06Gender\x12\x16\n" +
	"\x12GENDER_UNSPECIFIED\x10\x00\x12\b\n" +
	"\x04MALE\x10\x01\x12\n" +
	"\n" +
	"\x06FEMALE\x10\x02\"`\n" +
	"\x11CreateUserRequest\x12)\n" +
	"\x04user\x18\x01 \x01(\v2\x0f.limestone.UserB\x04\xe2A\x01\x02R\x04user\x12 \n" +
	"\bpassword\x18\x02 \x01(\tB\x04\xe2A\x01\x02R\bpassword\"&\n" +
	"\x0eGetUserRequest\x12\x14\n" +
	"\x02id\x18\x01 \x01(\tB\x04\xe2A\x01\x02R\x02id\"6\n" +
	"\x0fGetUserResponse\x12#\n" +
	"\x04user\x18\x01 \x01(\v2\x0f.limestone.UserR\x04user\"8\n" +
	"\x11UpdateUserRequest\x12#\n" +
	"\x04user\x18\x01 \x01(\v2\x0f.limestone.UserR\x04user\")\n" +
	"\x11DeleteUserRequest\x12\x14\n" +
	"\x02id\x18\x01 \x01(\tB\x04\xe2A\x01\x02R\x02id\"\x14\n" +
	"\x12DeleteUserResponse\"y\n" +
	"\x17AuthenticateUserRequest\x12\x1c\n" +
	"\busername\x18\x01 \x01(\tH\x00R\busername\x12\x16\n" +
	"\x05email\x18\x02 \x01(\tH\x00R\x05email\x12\x1a\n" +
	"\bpassword\x18\x03 \x01(\tR\bpasswordB\f\n" +
	"\n" +
	"identifier\"{\n" +
	"\x18AuthenticateUserResponse\x12!\n" +
	"\faccess_token\x18\x01 \x01(\tR\vaccessToken\x12#\n" +
	"\rrefresh_token\x18\x02 \x01(\tR\frefreshToken\x12\x17\n" +
	"\auser_id\x18\x03 \x01(\x05R\x06userId\":\n" +
	"\x13RefreshTokenRequest\x12#\n" +
	"\rrefresh_token\x18\x01 \x01(\tR\frefreshToken\"9\n" +
	"\x14RefreshTokenResponse\x12!\n" +
	"\faccess_token\x18\x01 \x01(\tR\vaccessToken2\x96\x05\n" +
	"\vUserService\x12d\n" +
	"\n" +
	"CreateUser\x12\x1c.limestone.CreateUserRequest\x1a\x0f.limestone.User\"'\xdaA\ruser,password\x82\xd3\xe4\x93\x02\x11:\x04user\"\t/v1/users\x12_\n" +
	"\aGetUser\x12\x19.limestone.GetUserRequest\x1a\x1a.limestone.GetUserResponse\"\x1d\xdaA\x02id\x82\xd3\xe4\x93\x02\x12\x12\x10/v1/{id=users/*}\x12g\n" +
	"\n" +
	"UpdateUser\x12\x1c.limestone.UpdateUserRequest\x1a\x0f.limestone.User\"*\xdaA\x04user\x82\xd3\xe4\x93\x02\x1d:\x04user2\x15/v1/{user.id=users/*}\x12h\n" +
	"\n" +
	"DeleteUser\x12\x1c.limestone.DeleteUserRequest\x1a\x1d.limestone.DeleteUserResponse\"\x1d\xdaA\x02id\x82\xd3\xe4\x93\x02\x12*\x10/v1/{id=users/*}\x12{\n" +
	"\x10AuthenticateUser\x12\".limestone.AuthenticateUserRequest\x1a#.limestone.AuthenticateUserResponse\"\x1e\x82\xd3\xe4\x93\x02\x18:\x01*\"\x13/users/authenticate\x12p\n" +
	"\fRefreshToken\x12\x1e.limestone.RefreshTokenRequest\x1a\x1f.limestone.RefreshTokenResponse\"\x1f\x82\xd3\xe4\x93\x02\x19:\x01*\"\x14/users/refresh_tokenB\x88\x01\n" +
	"\rcom.limestoneB\x10UserServiceProtoP\x01Z!github.com/mnadev/limestone/proto\xa2\x02\x03LXX\xaa\x02\tLimestone\xca\x02\tLimestone\xe2\x02\x15Limestone\\GPBMetadata\xea\x02\tLimestoneb\x06proto3"

var (
	file_user_service_proto_rawDescOnce sync.Once
	file_user_service_proto_rawDescData []byte
)

func file_user_service_proto_rawDescGZIP() []byte {
	file_user_service_proto_rawDescOnce.Do(func() {
		file_user_service_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_user_service_proto_rawDesc), len(file_user_service_proto_rawDesc)))
	})
	return file_user_service_proto_rawDescData
}

var file_user_service_proto_enumTypes = make([]protoimpl.EnumInfo, 2)
var file_user_service_proto_msgTypes = make([]protoimpl.MessageInfo, 12)
var file_user_service_proto_goTypes = []any{
	(MasjidRole_Role)(0),             // 0: limestone.MasjidRole.Role
	(User_Gender)(0),                 // 1: limestone.User.Gender
	(*MasjidRole)(nil),               // 2: limestone.MasjidRole
	(*User)(nil),                     // 3: limestone.User
	(*CreateUserRequest)(nil),        // 4: limestone.CreateUserRequest
	(*GetUserRequest)(nil),           // 5: limestone.GetUserRequest
	(*GetUserResponse)(nil),          // 6: limestone.GetUserResponse
	(*UpdateUserRequest)(nil),        // 7: limestone.UpdateUserRequest
	(*DeleteUserRequest)(nil),        // 8: limestone.DeleteUserRequest
	(*DeleteUserResponse)(nil),       // 9: limestone.DeleteUserResponse
	(*AuthenticateUserRequest)(nil),  // 10: limestone.AuthenticateUserRequest
	(*AuthenticateUserResponse)(nil), // 11: limestone.AuthenticateUserResponse
	(*RefreshTokenRequest)(nil),      // 12: limestone.RefreshTokenRequest
	(*RefreshTokenResponse)(nil),     // 13: limestone.RefreshTokenResponse
	(*timestamppb.Timestamp)(nil),    // 14: google.protobuf.Timestamp
}
var file_user_service_proto_depIdxs = []int32{
	0,  // 0: limestone.MasjidRole.role:type_name -> limestone.MasjidRole.Role
	1,  // 1: limestone.User.gender:type_name -> limestone.User.Gender
	14, // 2: limestone.User.create_time:type_name -> google.protobuf.Timestamp
	14, // 3: limestone.User.update_time:type_name -> google.protobuf.Timestamp
	3,  // 4: limestone.CreateUserRequest.user:type_name -> limestone.User
	3,  // 5: limestone.GetUserResponse.user:type_name -> limestone.User
	3,  // 6: limestone.UpdateUserRequest.user:type_name -> limestone.User
	4,  // 7: limestone.UserService.CreateUser:input_type -> limestone.CreateUserRequest
	5,  // 8: limestone.UserService.GetUser:input_type -> limestone.GetUserRequest
	7,  // 9: limestone.UserService.UpdateUser:input_type -> limestone.UpdateUserRequest
	8,  // 10: limestone.UserService.DeleteUser:input_type -> limestone.DeleteUserRequest
	10, // 11: limestone.UserService.AuthenticateUser:input_type -> limestone.AuthenticateUserRequest
	12, // 12: limestone.UserService.RefreshToken:input_type -> limestone.RefreshTokenRequest
	3,  // 13: limestone.UserService.CreateUser:output_type -> limestone.User
	6,  // 14: limestone.UserService.GetUser:output_type -> limestone.GetUserResponse
	3,  // 15: limestone.UserService.UpdateUser:output_type -> limestone.User
	9,  // 16: limestone.UserService.DeleteUser:output_type -> limestone.DeleteUserResponse
	11, // 17: limestone.UserService.AuthenticateUser:output_type -> limestone.AuthenticateUserResponse
	13, // 18: limestone.UserService.RefreshToken:output_type -> limestone.RefreshTokenResponse
	13, // [13:19] is the sub-list for method output_type
	7,  // [7:13] is the sub-list for method input_type
	7,  // [7:7] is the sub-list for extension type_name
	7,  // [7:7] is the sub-list for extension extendee
	0,  // [0:7] is the sub-list for field type_name
}

func init() { file_user_service_proto_init() }
func file_user_service_proto_init() {
	if File_user_service_proto != nil {
		return
	}
	file_user_service_proto_msgTypes[8].OneofWrappers = []any{
		(*AuthenticateUserRequest_Username)(nil),
		(*AuthenticateUserRequest_Email)(nil),
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_user_service_proto_rawDesc), len(file_user_service_proto_rawDesc)),
			NumEnums:      2,
			NumMessages:   12,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_user_service_proto_goTypes,
		DependencyIndexes: file_user_service_proto_depIdxs,
		EnumInfos:         file_user_service_proto_enumTypes,
		MessageInfos:      file_user_service_proto_msgTypes,
	}.Build()
	File_user_service_proto = out.File
	file_user_service_proto_goTypes = nil
	file_user_service_proto_depIdxs = nil
}
