package user_service

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/mnadev/limestone/proto"
	"github.com/mnadev/limestone/storage"
)

type UserServiceServer struct {
	SM *storage.StorageManager
	pb.UnimplementedUserServiceServer
}

func (srvr *UserServiceServer) CreateUser(ctx context.Context, in *pb.CreateUserRequest) (*pb.User, error) {
	user, err := srvr.SM.CreateUser(in.GetUser(), in.GetPassword())
	if err != nil {
		return nil, err
	}
	return user.ToProto(), status.Error(codes.OK, codes.OK.String())
}

func (srvr *UserServiceServer) LoginUser(ctx context.Context, in *pb.LoginRequest) (*pb.Tokens, error) {
	tokens, err := srvr.SM.LoginUser(in.GetEmail(), in.GetPassword())
	if err != nil {
		return nil, err
	}
	return tokens, nil
}

func (srvr *UserServiceServer) VerifyUser(ctx context.Context, in *pb.VerifyRequest) (*pb.VerifyResponse, error) {
	res, err := srvr.SM.VerifyUser(in.GetEmail(), in.GetCode())
	if err != nil {
		return nil, err
	}
	return res, nil
}
func (srvr *UserServiceServer) GetUser(ctx context.Context, in *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	if in.GetEmail() != "" {
		user, err := srvr.SM.GetUserWithEmail(in.GetEmail(), in.GetPassword())
		if err != nil {
			return nil, err
		}
		return &pb.GetUserResponse{
			User: user.ToProto(),
		}, status.Error(codes.OK, codes.OK.String())
	}
	user, err := srvr.SM.GetUserWithUsername(in.GetUsername(), in.GetPassword())
	if err != nil {
		return nil, err
	}
	return &pb.GetUserResponse{
		User: user.ToProto(),
	}, status.Error(codes.OK, codes.OK.String())
}

func (srvr *UserServiceServer) UpdateUser(ctx context.Context, in *pb.UpdateUserRequest) (*pb.User, error) {
	user, err := srvr.SM.UpdateUser(in.GetUser(), in.GetPassword())
	if err != nil {
		return nil, err
	}
	return user.ToProto(), status.Error(codes.OK, codes.OK.String())
}

func (srvr *UserServiceServer) DeleteUser(ctx context.Context, in *pb.DeleteUserRequest) (*pb.DeleteUserResponse, error) {
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
	return &pb.DeleteUserResponse{}, status.Error(codes.OK, codes.OK.String())
}
