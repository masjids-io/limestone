package user_service

import (
	"context"
	"github.com/mnadev/limestone/internal/infrastructure/storage"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/mnadev/limestone/gen/go"
	"github.com/mnadev/limestone/internal/infrastructure/grpc/auth"
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

func (s *UserServiceServer) AuthenticateUser(ctx context.Context, in *pb.AuthenticateUserRequest) (*pb.AuthenticateUserResponse, error) {
	var user *storage.User
	var err error

	if in.GetUsername() != "" {
		user, err = s.Smgr.GetUserByUsername(in.GetUsername())
	} else if in.GetEmail() != "" {
		user, err = s.Smgr.GetUserByEmail(in.GetEmail())
	} else {
		return nil, status.Errorf(codes.InvalidArgument, "username or email must be provided")
	}

	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "invalid username/email or password")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(in.GetPassword()))
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "invalid username/email or password")
	}

	// Generate JWT tokens
	accessToken, refreshToken, err := auth.GenerateJWT(user.ID.String())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to generate JWT tokens: %v", err)
	}

	return &pb.AuthenticateUserResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *UserServiceServer) RefreshToken(ctx context.Context, in *pb.RefreshTokenRequest) (*pb.RefreshTokenResponse, error) {
	accessToken, refreshToken, err := auth.RefreshToken(in.RefreshToken)
	if err != nil {
		return nil, err
	}
	return &pb.RefreshTokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
