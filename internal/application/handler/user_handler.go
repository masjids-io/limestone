package handler

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	pb "github.com/mnadev/limestone/gen/go"
	"github.com/mnadev/limestone/internal/application/domain/entity"
	"github.com/mnadev/limestone/internal/application/helper"
	"github.com/mnadev/limestone/internal/application/services"
	"github.com/mnadev/limestone/internal/infrastructure/auth"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"strings"
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
	if req.GetRole() == pb.CreateUserRequest_ROLE_UNSPECIFIED {
		return nil, status.Errorf(codes.InvalidArgument, "role is required and cannot be unspecified.")
	}

	password := req.GetPassword()
	if len(password) < 8 {
		return nil, status.Errorf(codes.InvalidArgument, "password must be at least 8 characters")
	}
	hashPassword, err := auth.HashPassword(password)
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
		Role:           entity.Role(req.GetRole().String()),
	}

	responseCreatedUser, err := h.Svc.CreateUser(ctx, u)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key") || strings.Contains(err.Error(), "already exists") {
			return nil, status.Errorf(codes.AlreadyExists, "email or username already exists")
		}
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	return helper.StandardUserResponse(codes.OK, "success", "user created successfully", responseCreatedUser, nil)
}

func (h *UserGrpcHandler) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.StandardUserResponse, error) {
	// --- Start Authorization (Coarse-Grained) ---
	allowedRolesForAnyUser := []string{
		string(entity.MASJID_ADMIN),
		string(entity.MASJID_MEMBER),
		string(entity.MASJID_VOLUNTEER),
		string(entity.MASJID_IMAM),
	}
	if err := auth.RequireRole(ctx, allowedRolesForAnyUser, "GetUser"); err != nil {
		return nil, err
	}
	// --- End Authorization (Coarse-Grained) ---

	user, err := h.Svc.GetUser(ctx, req.Id)
	if err != nil {
		return nil, status.Errorf(codes.Canceled, err.Error())
	}
	return helper.StandardUserResponse(codes.OK, "success", "user retrieved successfully", user, nil)
}

func (h *UserGrpcHandler) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.StandardUserResponse, error) {
	// --- Start Authorization (Coarse-Grained) ---
	allowedRolesForAnyUser := []string{
		string(entity.MASJID_ADMIN),
		string(entity.MASJID_VOLUNTEER),
	}
	if err := auth.RequireRole(ctx, allowedRolesForAnyUser, "GetUser"); err != nil {
		return nil, err
	}
	// --- End Authorization (Coarse-Grained) ---
	userIDStr := req.User.GetId()
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

	return helper.StandardUserResponse(codes.OK, "success", "user updated successfully", updatedUser, nil)
}

func (h *UserGrpcHandler) DeleteUser(ctx context.Context, req *pb.DeleteUserRequest) (*pb.StandardUserResponse, error) {
	// --- Start Authorization (Coarse-Grained) ---
	allowedRolesForAnyUser := []string{
		string(entity.MASJID_ADMIN),
	}
	if err := auth.RequireRole(ctx, allowedRolesForAnyUser, "GetUser"); err != nil {
		return nil, err
	}
	// --- End Authorization (Coarse-Grained) ---

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

	return helper.StandardUserResponse(codes.OK, "success", "user deleted successfully", nil, &pb.DeleteUserResponse{})

}
