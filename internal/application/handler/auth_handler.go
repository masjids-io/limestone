package handler

import (
	"context"
	"fmt"
	pb "github.com/mnadev/limestone/gen/go"
	services "github.com/mnadev/limestone/internal/application/services"
	auth2 "github.com/mnadev/limestone/internal/infrastructure/auth"
	"google.golang.org/grpc/codes"
)

type AuthGrpcHandler struct {
	pb.UnimplementedAuthServiceServer
	Svc *services.AuthService
}

func NewAuthGrpcHandler(svc *services.AuthService) *AuthGrpcHandler {
	return &AuthGrpcHandler{Svc: svc}
}

func (h *AuthGrpcHandler) AuthenticateUser(ctx context.Context, req *pb.AuthenticateUserRequest) (*pb.StandardAuthResponse, error) {
	identifier := ""
	if req.GetUsername() != "" {
		identifier = req.GetUsername()
	} else if req.GetEmail() != "" {
		identifier = req.GetEmail()
	} else {
		return &pb.StandardAuthResponse{Code: codes.Canceled.String(), Status: "error", Message: "username or email must be provided"}, nil
	}

	user, err := h.Svc.AuthenticateUser(ctx, identifier, req.GetPassword())
	if err != nil {
		return &pb.StandardAuthResponse{Code: codes.Canceled.String(), Status: "error", Message: "invalid username/email or password"}, nil
	}

	// Generate JWT tokens
	accessToken, refreshToken, err := auth2.GenerateJWT(user.ID.String())
	if err != nil {
		return &pb.StandardAuthResponse{Code: codes.Internal.String(), Status: "error", Message: fmt.Sprintf("failed to generate JWT tokens: %v", err)}, nil
	}

	return &pb.StandardAuthResponse{
		Code:    codes.OK.String(),
		Status:  "success",
		Message: "authentication successful",
		Datas: &pb.StandardAuthResponse_Data{
			Data: &pb.DataAuthenticateUserResponse{
				AccessToken:  accessToken,
				RefreshToken: refreshToken,
				UserId:       user.ID.String(),
			},
		},
	}, nil
}
