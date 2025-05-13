package handler

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	pb "github.com/mnadev/limestone/gen/go"
	"github.com/mnadev/limestone/internal/application/domain/entity"
	services "github.com/mnadev/limestone/internal/application/services"
	auth2 "github.com/mnadev/limestone/internal/infrastructure/auth"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/gorm"
	"time"
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
			return nil, status.Errorf(codes.AlreadyExists, "email or username already exists")
		}
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	return &pb.StandardUserResponse{
		Code:    codes.OK.String(),
		Status:  "success",
		Message: "user created successfully",
		Data: &pb.StandardUserResponse_AddUserResponse{
			AddUserResponse: &pb.User{
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

func (h *UserGrpcHandler) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.StandardUserResponse, error) {
	user, err := h.Svc.GetUser(ctx, req.Id)
	fmt.Println(user)
	if err != nil {
		return nil, status.Errorf(codes.Canceled, err.Error())
	}
	return &pb.StandardUserResponse{
		Code:    codes.OK.String(),
		Status:  "success",
		Message: "user created successfully",
		Data: &pb.StandardUserResponse_GetUserResponse{
			GetUserResponse: &pb.User{
				Id:              user.ID.String(),
				Email:           user.Email,
				Username:        user.Username,
				IsEmailVerified: user.IsVerified,
				FirstName:       user.FirstName,
				LastName:        user.LastName,
				PhoneNumber:     user.PhoneNumber,
				Gender:          pb.User_Gender(pb.User_Gender_value[user.Gender.String()]),
				CreateTime:      timestamppb.New(user.CreatedAt),
				UpdateTime:      timestamppb.New(user.UpdatedAt),
			},
		},
	}, nil
}

func (h *UserGrpcHandler) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.StandardUserResponse, error) {
	fmt.Println(req)
	userIDStr := req.User.Id
	if userIDStr == "" {
		return nil, status.Errorf(codes.InvalidArgument, "id is required")
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid user ID format")
	}

	updateData := &entity.User{
		ID:          userID,
		Email:       req.User.Email,
		Username:    req.User.Username,
		FirstName:   req.User.FirstName,
		LastName:    req.User.LastName,
		PhoneNumber: req.User.PhoneNumber,
		Gender:      entity.Gender(req.User.Gender),
		UpdatedAt:   time.Now(),
	}

	_, err = h.Svc.UpdateUser(ctx, updateData)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, fmt.Sprintf("failed to update user: %v", err))
	}

	updatedUser, err := h.Svc.GetUser(ctx, userIDStr)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "failed to retrieve updated user")
	}

	return &pb.StandardUserResponse{
		Code:    codes.OK.String(),
		Status:  "success",
		Message: "user updated successfully",
		Data: &pb.StandardUserResponse_UpdateUserResponse{
			UpdateUserResponse: &pb.User{
				Id:              updatedUser.ID.String(),
				Email:           updatedUser.Email,
				Username:        updatedUser.Username,
				IsEmailVerified: updatedUser.IsVerified,
				FirstName:       updatedUser.FirstName,
				LastName:        updatedUser.LastName,
				PhoneNumber:     updatedUser.PhoneNumber,
				Gender:          pb.User_Gender(pb.User_Gender_value[updatedUser.Gender.String()]),
				CreateTime:      timestamppb.New(updatedUser.CreatedAt),
				UpdateTime:      timestamppb.New(updatedUser.UpdatedAt),
			},
		},
	}, nil
}

func (h *UserGrpcHandler) DeleteUser(ctx context.Context, req *pb.DeleteUserRequest) (*pb.StandardUserResponse, error) {
	fmt.Println(req)
	userIDStr := req.GetId()
	if userIDStr == "" {
		return nil, status.Errorf(codes.InvalidArgument, "id is required")
	}

	_, err := uuid.Parse(userIDStr)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid user ID format")
	}

	err = h.Svc.DeleteUser(ctx, userIDStr)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "failed to retrieve Delete user")
	}

	return &pb.StandardUserResponse{
		Code:    codes.OK.String(),
		Status:  "success",
		Message: "user deleted successfully",
	}, nil
}
