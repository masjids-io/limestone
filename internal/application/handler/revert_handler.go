// Package handler Copyright (c) 2024 Coding-af Limestone Dev
// Licensed under the MIT License.
// file COPYING or http://www.opensource.org/licenses/mit-license.php
package handler

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	pb "github.com/mnadev/limestone/gen/go"
	"github.com/mnadev/limestone/internal/application/domain/entity"
	"github.com/mnadev/limestone/internal/application/helper"
	"github.com/mnadev/limestone/internal/application/services"
	"github.com/mnadev/limestone/internal/infrastructure/auth"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type RevertsIoGrpcHandler struct {
	pb.UnimplementedRevertsIoServiceServer
	RevertSvc *services.RevertService
}

func NewRevertGrpcHandler(revertSvc *services.RevertService) *RevertsIoGrpcHandler {
	return &RevertsIoGrpcHandler{
		RevertSvc: revertSvc,
	}
}

func (h *RevertsIoGrpcHandler) CreateRevertProfile(ctx context.Context, req *pb.CreateRevertProfileRequest) (*pb.StandardRevertResponse, error) {
	// --- Start Authorization (Coarse-Grained) ---
	allowedRolesForAnyUser := []string{
		string(entity.MASJID_ADMIN),
	}
	if err := auth.RequireRole(ctx, allowedRolesForAnyUser, "UpdateUser"); err != nil {
		return nil, err
	}
	// --- End Authorization (Coarse-Grained) ---
	userID, ok := ctx.Value(auth.UserIDContextKey).(string)
	if !ok || userID == "" {
		return nil, status.Errorf(codes.Unauthenticated, "Unauthenticated: user ID not found in context")
	}

	if req.GetProfile() == nil {
		return nil, status.Errorf(codes.InvalidArgument, "Revert profile data is required")
	}

	profileEntity, err := helper.ToEntityRevertProfile(req.GetProfile())
	fmt.Println(req)
	fmt.Println(profileEntity)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "Invalid revert profile data: %v", err)
	}

	profileEntity.UserID = userID

	createdProfile, err := h.RevertSvc.CreateRevertProfile(ctx, profileEntity)
	if err != nil {
		if errors.Is(err, helper.ErrRevertProfileAlreadyExists) {
			return nil, status.Errorf(codes.AlreadyExists, err.Error())
		}
		if errors.Is(err, helper.ErrInvalidRevertProfileData) {
			return nil, status.Errorf(codes.InvalidArgument, err.Error())
		}
		return nil, status.Errorf(codes.Internal, fmt.Sprintf("Failed to create revert profile: %v", err))
	}

	resp, helperErr := helper.StandardRevertResponse(
		codes.OK,
		"success",
		"Revert profile created successfully",
		createdProfile,
	)
	if helperErr != nil {
		return nil, status.Errorf(codes.Internal, fmt.Sprintf("Failed to construct response: %v", helperErr))
	}
	return resp, nil
}

func (h *RevertsIoGrpcHandler) GetSelfRevertProfile(ctx context.Context, req *pb.GetSelfRevertProfileRequest) (*pb.StandardRevertResponse, error) {
	userID, ok := ctx.Value(auth.UserIDContextKey).(string)
	if !ok || userID == "" {
		return nil, status.Errorf(codes.Unauthenticated, "Unauthenticated: user ID not found in context")
	}

	profile, err := h.RevertSvc.GetRevertProfileByUserID(ctx, userID)
	if err != nil {
		if errors.Is(err, helper.ErrNotFound) {
			return nil, status.Errorf(codes.NotFound, "Revert profile not found for user ID: %s", userID)
		}
		return nil, status.Errorf(codes.Internal, fmt.Sprintf("Failed to retrieve self revert profile: %v", err))
	}

	resp, helperErr := helper.StandardRevertResponse(
		codes.OK,
		"success",
		"Self revert profile retrieved successfully",
		profile,
	)
	if helperErr != nil {
		return nil, status.Errorf(codes.Internal, fmt.Sprintf("Failed to construct response: %v", helperErr))
	}
	return resp, nil
}

func (h *RevertsIoGrpcHandler) UpdateSelfRevertProfile(ctx context.Context, req *pb.UpdateSelfRevertProfileRequest) (*pb.StandardRevertResponse, error) {
	// --- Start Authorization (Coarse-Grained) ---
	allowedRolesForAnyUser := []string{
		string(entity.MASJID_ADMIN),
		string(entity.MASJID_VOLUNTEER),
	}
	if err := auth.RequireRole(ctx, allowedRolesForAnyUser, "UpdateUser"); err != nil {
		return nil, err
	}
	// --- End Authorization (Coarse-Grained) ---
	userID, ok := ctx.Value(auth.UserIDContextKey).(string)
	if !ok || userID == "" {
		return nil, status.Errorf(codes.Unauthenticated, "Unauthenticated: user ID not found in context")
	}

	if req.GetProfile() == nil {
		return nil, status.Errorf(codes.InvalidArgument, "Update profile data is required")
	}

	profileID, err := uuid.Parse(req.GetProfileId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "Invalid profile ID format: %v", err)
	}

	existingProfile, err := h.RevertSvc.GetRevertProfileByID(ctx, profileID)
	if err != nil {
		if errors.Is(err, helper.ErrNotFound) {
			return nil, status.Errorf(codes.NotFound, "Revert profile not found for ID: %s", profileID.String())
		}
		return nil, status.Errorf(codes.Internal, fmt.Sprintf("Failed to get existing profile: %v", err))
	}

	if existingProfile.UserID != userID {
		return nil, status.Errorf(codes.PermissionDenied, "Permission denied: You can only update your own profile")
	}

	updatedProfileEntity, err := helper.ToEntityRevertProfile(req.GetProfile())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "Invalid update profile data: %v", err)
	}
	updatedProfileEntity.ID = existingProfile.ID
	updatedProfileEntity.UserID = existingProfile.UserID

	updatedProfile, err := h.RevertSvc.UpdateRevertProfile(ctx, updatedProfileEntity)
	if err != nil {
		if errors.Is(err, helper.ErrNotFound) {
			return nil, status.Errorf(codes.NotFound, err.Error())
		}
		if errors.Is(err, helper.ErrInvalidRevertProfileData) {
			return nil, status.Errorf(codes.InvalidArgument, err.Error())
		}
		return nil, status.Errorf(codes.Internal, fmt.Sprintf("Failed to update revert profile: %v", err))
	}

	resp, helperErr := helper.StandardRevertResponse(
		codes.OK,
		"success",
		"Revert profile updated successfully",
		updatedProfile,
	)
	if helperErr != nil {
		return nil, status.Errorf(codes.Internal, fmt.Sprintf("Failed to construct response: %v", helperErr))
	}
	return resp, nil
}

func (h *RevertsIoGrpcHandler) ListRevertProfiles(ctx context.Context, req *pb.ListRevertProfilesRequest) (*pb.StandardRevertResponse, error) {
	params := &entity.RevertProfileQueryParams{
		Start:  req.GetStart(),
		Limit:  req.GetLimit(),
		Page:   req.GetPage(),
		Name:   req.GetName(),
		Gender: req.GetGender(),
	}

	queryResult, err := h.RevertSvc.ListRevertProfiles(ctx, params)
	if err != nil {
		return nil, status.Errorf(codes.Internal, fmt.Sprintf("Failed to list revert profiles: %v", err))
	}

	protoRevertProfiles := make([]*pb.RevertProfile, len(queryResult.RevertProfile))
	for i, p := range queryResult.RevertProfile {
		protoRevertProfiles[i] = helper.ToProtoRevertProfile(p)
	}

	listResponse := &pb.ListRevertProfilesResponse{
		Profiles:    protoRevertProfiles,
		TotalCount:  int32(queryResult.TotalCount),
		CurrentPage: queryResult.CurrentPage,
		TotalPages:  queryResult.TotalPages,
	}

	resp, helperErr := helper.StandardRevertResponse(
		codes.OK,
		"success",
		"Nikkah profiles listed successfully",
		listResponse,
	)
	if helperErr != nil {
		return nil, status.Errorf(codes.Internal, fmt.Sprintf("Failed to construct response: %v", helperErr))
	}

	resp.Data = &pb.StandardRevertResponse_ListRevertProfilesResponse{ListRevertProfilesResponse: listResponse}

	return resp, nil
}

func (h *RevertsIoGrpcHandler) CreateRevertMatchInvite(ctx context.Context, req *pb.CreateRevertMatchInviteRequest) (*pb.StandardRevertResponse, error) {
	// --- Start Authorization (Coarse-Grained) ---
	allowedRolesForAnyUser := []string{
		string(entity.MASJID_ADMIN),
	}
	if err := auth.RequireRole(ctx, allowedRolesForAnyUser, "UpdateUser"); err != nil {
		return nil, err
	}
	// --- End Authorization (Coarse-Grained) ---
	initiatorProfileID, ok := ctx.Value(auth.UserIDContextKey).(string)
	if !ok || initiatorProfileID == "" {
		return nil, status.Errorf(codes.Unauthenticated, "initiator profile ID not found in context")
	}

	if req.GetReceiverProfileId() == "" {
		return nil, status.Errorf(codes.InvalidArgument, "receiver_profile_id is required")
	}

	revertMatch, err := h.RevertSvc.CreateRevertMatchInvite(ctx, initiatorProfileID, req.GetReceiverProfileId())
	if err != nil {
		if errors.Is(err, helper.ErrInvalidMatchData) {
			return nil, status.Errorf(codes.InvalidArgument, err.Error())
		}
		if errors.Is(err, helper.ErrSelfInvitation) {
			return nil, status.Errorf(codes.InvalidArgument, err.Error())
		}
		if errors.Is(err, helper.ErrMatchAlreadyExists) {
			return nil, status.Errorf(codes.AlreadyExists, err.Error())
		}
		if errors.Is(err, helper.ErrProfileNotFound) {
			return nil, status.Errorf(codes.NotFound, err.Error())
		}
		return nil, status.Errorf(codes.Internal, fmt.Sprintf("Failed to create revert match invite: %v", err))
	}

	resp, helperErr := helper.StandardRevertResponse(
		codes.OK,
		"success",
		"Revert match invite created successfully",
		revertMatch,
	)
	if helperErr != nil {
		return nil, status.Errorf(codes.Internal, fmt.Sprintf("Failed to construct response: %v", helperErr))
	}
	return resp, nil
}

func (h *RevertsIoGrpcHandler) GetRevertMatch(ctx context.Context, req *pb.GetRevertMatchRequest) (*pb.StandardRevertResponse, error) {
	if req.GetMatchId() == "" {
		return nil, status.Errorf(codes.InvalidArgument, "match_id is required")
	}

	revertMatch, err := h.RevertSvc.GetRevertMatch(ctx, req.GetMatchId())
	if err != nil {
		if errors.Is(err, helper.ErrInvalidMatchID) {
			return nil, status.Errorf(codes.InvalidArgument, err.Error())
		}
		if errors.Is(err, helper.ErrMatchNotFound) {
			return nil, status.Errorf(codes.NotFound, err.Error())
		}
		return nil, status.Errorf(codes.Internal, fmt.Sprintf("Failed to get revert match: %v", err))
	}

	if revertMatch == nil {
		return nil, status.Errorf(codes.NotFound, "Revert match not found")
	}

	resp, helperErr := helper.StandardRevertResponse(
		codes.OK,
		"success",
		"Revert match retrieved successfully",
		revertMatch,
	)
	if helperErr != nil {
		return nil, status.Errorf(codes.Internal, fmt.Sprintf("Failed to construct response: %v", helperErr))
	}

	return resp, nil
}

func (h *RevertsIoGrpcHandler) AcceptRevertMatchInvite(ctx context.Context, req *pb.AcceptRevertMatchInviteRequest) (*pb.StandardRevertResponse, error) {
	currentProfileID, ok := ctx.Value(auth.UserIDContextKey).(string)
	if !ok || currentProfileID == "" {
		return nil, status.Errorf(codes.Unauthenticated, "authenticated profile ID not found in context")
	}

	if req.GetMatchId() == "" {
		return nil, status.Errorf(codes.InvalidArgument, "match_id is required")
	}

	revertMatch, err := h.RevertSvc.AcceptRevertMatchInvite(ctx, req.GetMatchId())
	if err != nil {
		if errors.Is(err, helper.ErrInvalidMatchID) {
			return nil, status.Errorf(codes.InvalidArgument, err.Error())
		}
		if errors.Is(err, helper.ErrMatchNotFound) {
			return nil, status.Errorf(codes.NotFound, err.Error())
		}
		if errors.Is(err, helper.ErrMatchStatusInvalid) {
			return nil, status.Errorf(codes.FailedPrecondition, err.Error())
		}
		return nil, status.Errorf(codes.Internal, fmt.Sprintf("Failed to accept revert match invite: %v", err))
	}

	if revertMatch == nil {
		return nil, status.Errorf(codes.Internal, "Unexpected: Accepted match returned nil")
	}

	resp, helperErr := helper.StandardRevertResponse(
		codes.OK,
		"success",
		"Revert match invite accepted successfully",
		revertMatch,
	)
	if helperErr != nil {
		return nil, status.Errorf(codes.Internal, fmt.Sprintf("Failed to construct response: %v", helperErr))
	}

	return resp, nil
}

func (h *RevertsIoGrpcHandler) RejectRevertMatchInvite(ctx context.Context, req *pb.RejectRevertMatchInviteRequest) (*pb.StandardRevertResponse, error) {
	currentProfileID, ok := ctx.Value(auth.UserIDContextKey).(string)
	if !ok || currentProfileID == "" {
		return nil, status.Errorf(codes.Unauthenticated, "authenticated profile ID not found in context")
	}

	if req.GetMatchId() == "" {
		return nil, status.Errorf(codes.InvalidArgument, "match_id is required")
	}

	revertMatch, err := h.RevertSvc.RejectRevertMatchInvite(ctx, req.GetMatchId())
	if err != nil {
		if errors.Is(err, helper.ErrInvalidMatchID) {
			return nil, status.Errorf(codes.InvalidArgument, err.Error())
		}
		if errors.Is(err, helper.ErrMatchNotFound) {
			return nil, status.Errorf(codes.NotFound, err.Error())
		}
		if errors.Is(err, helper.ErrMatchStatusInvalid) {
			return nil, status.Errorf(codes.FailedPrecondition, err.Error())
		}
		return nil, status.Errorf(codes.Internal, fmt.Sprintf("Failed to reject revert match invite: %v", err))
	}

	if revertMatch == nil {
		return nil, status.Errorf(codes.Internal, "Unexpected: Rejected match returned nil")
	}

	resp, helperErr := helper.StandardRevertResponse(
		codes.OK,
		"success",
		"Revert match invite rejected successfully",
		revertMatch,
	)
	if helperErr != nil {
		return nil, status.Errorf(codes.Internal, fmt.Sprintf("Failed to construct response: %v", helperErr))
	}

	return resp, nil
}

func (h *RevertsIoGrpcHandler) EndRevertMatch(ctx context.Context, req *pb.EndRevertMatchRequest) (*pb.StandardRevertResponse, error) {
	currentProfileID, ok := ctx.Value(auth.UserIDContextKey).(string)
	if !ok || currentProfileID == "" {
		return nil, status.Errorf(codes.Unauthenticated, "authenticated profile ID not found in context")
	}

	if req.GetMatchId() == "" {
		return nil, status.Errorf(codes.InvalidArgument, "match_id is required")
	}

	revertMatch, err := h.RevertSvc.EndRevertMatch(ctx, req.GetMatchId())
	if err != nil {
		if errors.Is(err, helper.ErrInvalidMatchID) {
			return nil, status.Errorf(codes.InvalidArgument, err.Error())
		}
		if errors.Is(err, helper.ErrMatchNotFound) {
			return nil, status.Errorf(codes.NotFound, err.Error())
		}
		if errors.Is(err, helper.ErrMatchStatusInvalid) {
			return nil, status.Errorf(codes.FailedPrecondition, err.Error())
		}
		return nil, status.Errorf(codes.Internal, fmt.Sprintf("Failed to end revert match: %v", err))
	}

	if revertMatch == nil {
		return nil, status.Errorf(codes.Internal, "Unexpected: Ended match returned nil")
	}

	resp, helperErr := helper.StandardRevertResponse(
		codes.OK,
		"success",
		"Revert match ended successfully",
		revertMatch,
	)
	if helperErr != nil {
		return nil, status.Errorf(codes.Internal, fmt.Sprintf("Failed to construct response: %v", helperErr))
	}

	return resp, nil
}
