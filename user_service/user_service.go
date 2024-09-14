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

func (s *UserServiceServer) GetUser(ctx context.Context, in *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	user, err := s.Smgr.GetUser(in.GetId())
	if err != nil {
		return nil, err
	}
	return &pb.GetUserResponse{
		User: user.ToProto(),
	}, status.Error(codes.OK, codes.OK.String())
}

func (s *UserServiceServer) UpdateUser(ctx context.Context, in *pb.UpdateUserRequest) (*pb.User, error) {
	user, err := s.Smgr.UpdateUser(in.GetUser())
	if err != nil {
		return nil, err
	}
	return user.ToProto(), status.Error(codes.OK, codes.OK.String())
}

func (s *UserServiceServer) DeleteUser(ctx context.Context, in *pb.DeleteUserRequest) (*pb.DeleteUserResponse, error) {
	err := s.Smgr.DeleteUser(in.GetId())
	if err != nil {
		return nil, err
	}
	return &pb.DeleteUserResponse{}, status.Error(codes.OK, codes.OK.String())
}
