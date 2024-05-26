package user_service

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/mnadev/limestone/proto"
	"github.com/mnadev/limestone/storage"
)

type UserServiceServer struct {
	Smgr *storage.StorageManager
	pb.UnimplementedUserServiceServer
}

func (s *UserServiceServer) CreateUser(ctx context.Context, in *pb.CreateUserRequest) (*pb.User, error) {
	user, err := s.Smgr.CreateUser(in.GetUser(), in.GetPassword())
	if err != nil {
		return nil, err
	}
	return user.ToProto(), status.Error(codes.OK, codes.OK.String())
}

func (s *UserServiceServer) LoginUser(ctx context.Context, in *pb.LoginRequest) (*pb.Tokens, error) {
	tokens, err := s.Smgr.LoginUser(in.GetEmail(), in.GetPassword())
	if err != nil {
		return nil, err
	}
	return tokens, nil
}

func (s *UserServiceServer) VerifyUser(ctx context.Context, in *pb.VerifyRequest) (*pb.VerifyResponse, error) {
	res, err := s.Smgr.VerifyUser(in.GetEmail(), in.GetCode())
	if err != nil {
		return nil, err
	}
	return res, nil
}
func (s *UserServiceServer) GetUser(ctx context.Context, in *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	if in.GetEmail() != "" {
		user, err := s.Smgr.GetUserWithEmail(in.GetEmail(), in.GetPassword())
		if err != nil {
			return nil, err
		}
		return &pb.GetUserResponse{
			User: user.ToProto(),
		}, status.Error(codes.OK, codes.OK.String())
	}
	user, err := s.Smgr.GetUserWithUsername(in.GetUsername(), in.GetPassword())
	if err != nil {
		return nil, err
	}
	return &pb.GetUserResponse{
		User: user.ToProto(),
	}, status.Error(codes.OK, codes.OK.String())
}

func (s *UserServiceServer) UpdateUser(ctx context.Context, in *pb.UpdateUserRequest) (*pb.User, error) {
	user, err := s.Smgr.UpdateUser(in.GetUser(), in.GetPassword())
	if err != nil {
		return nil, err
	}
	return user.ToProto(), status.Error(codes.OK, codes.OK.String())
}

func (s *UserServiceServer) DeleteUser(ctx context.Context, in *pb.DeleteUserRequest) (*pb.DeleteUserResponse, error) {
	if in.GetEmail() != "" {
		err := s.Smgr.DeleteUserWithEmail(in.GetEmail(), in.GetPassword())
		if err != nil {
			return nil, err
		}
	} else {
		err := s.Smgr.DeleteUserWithUsername(in.GetUsername(), in.GetPassword())
		if err != nil {
			return nil, err
		}
	}
	return &pb.DeleteUserResponse{}, status.Error(codes.OK, codes.OK.String())
}
