package handler

import (
	"context"
	"errors"
	"github.com/google/uuid"
	pb "github.com/mnadev/limestone/gen/go"
	"github.com/mnadev/limestone/internal/application/domain/entity"
	services "github.com/mnadev/limestone/internal/application/services"
	auth2 "github.com/mnadev/limestone/internal/infrastructure/auth"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"
)

type UserGrpcHandler struct {
	pb.UnimplementedUserServiceServer
	Svc *services.UserService
}

func NewUserGrpcHandler(svc *services.UserService) *UserGrpcHandler {
	return &UserGrpcHandler{Svc: svc}
}

func (h *UserGrpcHandler) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.StandardUserResponse, error) {
	if req.GetEmail() == "" {
		return nil, status.Errorf(codes.InvalidArgument, "email is required")
	}
	if req.GetUsername() == "" {
		return nil, status.Errorf(codes.InvalidArgument, "username is required")
	}
	if req.GetPassword() == "" {
		return nil, status.Errorf(codes.InvalidArgument, "password is required")
	}
	if req.GetFirstName() == "" {
		return nil, status.Errorf(codes.InvalidArgument, "first name is required")
	}
	if req.GetLastName() == "" {
		return nil, status.Errorf(codes.InvalidArgument, "last name is required")
	}
	if req.GetPhoneNumber() == "" {
		return nil, status.Errorf(codes.InvalidArgument, "phone number is required")
	}

	password := req.GetPassword()
	if len(password) < 8 {
		return nil, status.Errorf(codes.InvalidArgument, "password must be at least 8 characters")
	}
	hashPassword, err := auth2.HashPassword(password)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
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
		Gender:         entity.Gender(req.GetGender().String()),
	}
	responseCreatedUser, err := h.Svc.CreateUser(ctx, u)
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return &pb.StandardUserResponse{Code: codes.AlreadyExists.String(), Status: "failed", Message: "email or username already exists"}, nil
		}
		return &pb.StandardUserResponse{Code: codes.Internal.String(), Status: "Failed", Message: err.Error()}, nil
	}
	return &pb.StandardUserResponse{
		Code:    codes.OK.String(),
		Status:  "success",
		Message: "user created successfully",
		Data: &pb.StandardUserResponse_User{
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
