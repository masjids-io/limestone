package handler

import (
	"context"
	pb "github.com/mnadev/limestone/gen/go"
	services "github.com/mnadev/limestone/internal/application/services"
	"github.com/mnadev/limestone/internal/infrastructure/auth"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
		return nil, status.Errorf(codes.Canceled, "username or email must be provided")
	}

	user, err := h.Svc.AuthenticateUser(ctx, identifier, req.GetPassword())
	if err != nil {
		return nil, status.Errorf(codes.Canceled, "invalid username/email or password")
	}

	accessToken, refreshToken, err := auth.GenerateJWT(user.ID.String(), user.Role.String())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to generate JWT tokens: %v", err)
	}

	return &pb.StandardAuthResponse{
		Code:    codes.OK.String(),
		Status:  "success",
		Message: "Authentication successful",
		Datas: &pb.StandardAuthResponse_AuthenticateUserData{
			AuthenticateUserData: &pb.DataAuthenticateUserResponse{
				AccessToken:  accessToken,
				RefreshToken: refreshToken,
				UserId:       user.ID.String(),
			},
		},
	}, nil
}

func (h *AuthGrpcHandler) RefreshToken(ctx context.Context, req *pb.RefreshTokenRequest) (*pb.StandardAuthResponse, error) {
	newAccessToken, newRefreshToken, err := h.Svc.RefreshToken(ctx, req.GetRefreshToken())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to refresh access token: %v", err)
	}
	return &pb.StandardAuthResponse{
		Code:    codes.OK.String(),
		Status:  "success",
		Message: "Token refreshed",
		Datas: &pb.StandardAuthResponse_RefreshTokenData{
			RefreshTokenData: &pb.DataRefreshTokenResponse{
				AccessToken:  newAccessToken,
				RefreshToken: newRefreshToken,
			},
		},
	}, nil
}
