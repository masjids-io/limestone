package user_service

import (
	"context"

	"github.com/mnadev/limestone/storage"
	userservicepb "github.com/mnadev/limestone/user_service/proto"
)

type UserServiceServer struct {
	SM *storage.StorageManager
	userservicepb.UnimplementedUserServiceServer
}

func (srvr *UserServiceServer) CreateUser(ctx context.Context, in *userservicepb.CreateUserRequest) (*userservicepb.User, error) {
	user, err := srvr.SM.CreateUser(in.GetUser(), in.GetPassword())
	if err != nil {
		return nil, err
	}
	return user.ToProto(), nil
}

func (srvr *UserServiceServer) GetUser(ctx context.Context, in *userservicepb.GetUserRequest) (*userservicepb.GetUserResponse, error) {
	if in.GetEmail() != "" {
		user, err := srvr.SM.GetUserWithEmail(in.GetEmail(), in.GetPassword())
		if err != nil {
			return nil, err
		}
		return &userservicepb.GetUserResponse{
			User: user.ToProto(),
		}, nil
	}
	user, err := srvr.SM.GetUserWithUsername(in.GetUsername(), in.GetPassword())
	if err != nil {
		return nil, err
	}
	return &userservicepb.GetUserResponse{
		User: user.ToProto(),
	}, nil
}

func (srvr *UserServiceServer) UpdateUser(ctx context.Context, in *userservicepb.UpdateUserRequest) (*userservicepb.User, error) {
	user, err := srvr.SM.UpdateUser(in.GetUser(), in.GetPassword())
	if err != nil {
		return nil, err
	}
	return user.ToProto(), nil
}

func (srvr *UserServiceServer) DeleteUser(ctx context.Context, in *userservicepb.DeleteUserRequest) (*userservicepb.DeleteUserResponse, error) {
	if in.GetEmail() != "" {
		err := srvr.SM.DeleteUserWithEmail(in.GetEmail(), in.GetPassword())
		if err != nil {
			return nil, err
		}
	} else {
		err := srvr.SM.DeleteUserWithUsername(in.GetUsername(), in.GetPassword())
		if err != nil {
			return nil, err
		}
	}
	return &userservicepb.DeleteUserResponse{}, nil
}
