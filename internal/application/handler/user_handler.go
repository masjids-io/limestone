package handler

import (
	"context"
	"github.com/google/uuid"
	pb "github.com/mnadev/limestone/gen/go"
	"github.com/mnadev/limestone/internal/application/domain/entity"
	services "github.com/mnadev/limestone/internal/application/services"
	"github.com/mnadev/limestone/internal/infrastructure/grpc/auth"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type UserGrpcHandler struct {
	pb.UnimplementedUserServiceServer
	Svc *services.UserService
}

func NewUserGrpcHandler(svc *services.UserService) *UserGrpcHandler {
	return &UserGrpcHandler{Svc: svc}
}

func (h *UserGrpcHandler) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.StandardResponse, error) {
	password := req.GetPassword()
	if len(password) < 8 {
		return &pb.StandardResponse{Code: codes.Canceled.String(), Status: "failed", Message: "password must be at least 8 characters"}, nil
	}
	hashPassword, err := auth.HashPassword(password)
	if err != nil {
		return &pb.StandardResponse{Code: codes.Internal.String(), Status: "failed", Message: err.Error()}, nil
	}
	u := &entity.User{
		ID:             uuid.New(),
		Email:          req.GetEmail(),
		Username:       req.GetUsername(),
		HashedPassword: hashPassword,
		IsVerified:     req.GetIsEmailVerified(),
		FirstName:      req.GetFirstName(),
		LastName:       req.GetLastName(),
		PhoneNumber:    req.GetPhoneNumber(),
		Gender:         entity.Gender(req.Gender.String()),
	}
	responseCreatedUser, err := h.Svc.CreateUser(ctx, u)
	if err != nil {
		return &pb.StandardResponse{Code: codes.Internal.String(), Status: "Failed", Message: err.Error()}, nil
	}
	return &pb.StandardResponse{
		Code:    codes.OK.String(),
		Status:  "success",
		Message: "user created successfully",
		Data: &pb.StandardResponse_User{
			User: &pb.User{
				Id:              responseCreatedUser.ID.String(),
				Email:           responseCreatedUser.Email,
				Username:        responseCreatedUser.Username,
				IsEmailVerified: responseCreatedUser.IsVerified,
				FirstName:       responseCreatedUser.FirstName,
				LastName:        responseCreatedUser.LastName,
				PhoneNumber:     responseCreatedUser.PhoneNumber,
				Gender:          pb.User_Gender(pb.User_Gender_value[responseCreatedUser.Gender.String()]),
				CreateTime:      timestamppb.New(responseCreatedUser.CreatedAt),
				UpdateTime:      timestamppb.New(responseCreatedUser.UpdatedAt),
			},
		},
	}, nil
}

func (h *UserGrpcHandler) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	user, err := h.Svc.GetUser(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return &pb.GetUserResponse{User: &pb.User{Id: user.ID.String(), Email: user.Email, Username: user.Username}}, nil
}

func (h *UserGrpcHandler) AuthenticateUser(ctx context.Context, req *pb.AuthenticateUserRequest) (*pb.AuthenticateUserResponse, error) {
	var user *entity.User
	var err error

	if req.GetUsername() != "" {
		user, err = h.Svc.GetUserByUsername(ctx, req.GetUsername())
	} else if req.GetEmail() != "" {
		user, err = h.Svc.GetUserByEmail(ctx, req.GetEmail())
	} else {
		return nil, status.Errorf(codes.InvalidArgument, "username or email must be provided")
	}

	if err != nil || user == nil {
		return nil, status.Errorf(codes.Unauthenticated, "invalid username/email or password")
	}

	// Compare the hashed password
	err = bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(req.GetPassword()))
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
