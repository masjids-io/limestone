package userservice

import (
	"context"

	"github.com/mnadev/limestone/storage"
	userpb "github.com/mnadev/limestone/user/proto"
	userservicepb "github.com/mnadev/limestone/userservice/proto"
)

type userServiceServer struct {
	sm *storage.StorageManager
}

func (srvr *userServiceServer) CreateUser(ctx context.Context, in *userservicepb.CreateUserRequest) (*userpb.User, error) {
	user, err := srvr.sm.CreateUser(in.GetUser(), in.GetPassword())
	if err != nil {
		return nil, err
	}
	return user.ToProto(), nil
}

func (srvr *userServiceServer) GetUser(ctx context.Context, in *userservicepb.GetUserRequest) (*userpb.User, error) {
	if in.GetEmail() != "" {
		user, err := srvr.sm.GetUserWithEmail(in.GetEmail(), in.GetPassword())
		if err != nil {
			return nil, err
		}
		return user.ToProto(), nil

	}
	user, err := srvr.sm.GetUserWithUsername(in.GetUsername(), in.GetPassword())
	if err != nil {
		return nil, err
	}
	return user.ToProto(), nil
}

func (srvr *userServiceServer) UpdateUser(ctx context.Context, in *userservicepb.UpdateUserRequest) (*userpb.User, error) {
	user, err := srvr.sm.UpdateUser(in.GetUser(), in.GetPassword())
	if err != nil {
		return nil, err
	}
	return user.ToProto(), nil
}

func (srvr *userServiceServer) DeleteUser(ctx context.Context, in *userservicepb.DeleteUserRequest) (*userservicepb.DeleteUserResponse, error) {
	if in.GetEmail() != "" {
		err := srvr.sm.DeleteUserWithEmail(in.GetEmail(), in.GetPassword())
		if err != nil {
			return nil, err
		}
	} else {
		err := srvr.sm.DeleteUserWithUsername(in.GetUsername(), in.GetPassword())
		if err != nil {
			return nil, err
		}
	}
	return &userservicepb.DeleteUserResponse{}, nil
}
