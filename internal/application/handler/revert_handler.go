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

// CreateRevertProfile /*
/*
@Method CreateRevertProfile
@name Handle creation of a revert profile via gRPC
@description Validates the user authentication and request payload,
             converts the gRPC profile message to the domain entity,
             sets the authenticated user ID,
             calls the RevertService to create the profile,
             and returns a standardized gRPC response.
@param ctx context.Context - request context, must include authenticated user ID
@param req *pb.CreateRevertProfileRequest - incoming gRPC request containing profile data
@return *pb.StandardRevertResponse - standardized response with created profile or error details
@return error - gRPC error with appropriate status codes on failure
*/
func (h *RevertsIoGrpcHandler) CreateRevertProfile(ctx context.Context, req *pb.CreateRevertProfileRequest) (*pb.StandardRevertResponse, error) {
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

// GetSelfRevertProfile /*
/*
@Method GetSelfRevertProfile
@name Retrieve the authenticated user's own revert profile
@description Extracts the user ID from context, fetches the revert profile associated with that user,
             handles possible errors like unauthenticated user or profile not found,
             and returns a standardized gRPC response containing the profile data.
@param ctx context.Context - request context, must include authenticated user ID
@param req *pb.GetSelfRevertProfileRequest - incoming gRPC request (empty payload)
@return *pb.StandardRevertResponse - standardized response with the user's revert profile or error details
@return error - gRPC error with appropriate status codes on failure
*/
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

// UpdateSelfRevertProfile /*
/*
@Method UpdateSelfRevertProfile
@name Update the authenticated user's revert profile
@description Retrieves the user ID from the context to ensure authentication,
             validates the provided profile ID and update data,
             checks ownership to prevent unauthorized updates,
             converts the incoming protobuf profile data to entity format,
             and updates the profile in the service layer.
             Returns a standardized gRPC response with the updated profile or an appropriate error.
@param ctx context.Context - the request context containing authentication info
@param req *pb.UpdateSelfRevertProfileRequest - the gRPC request containing profile ID and updated profile data
@return *pb.StandardRevertResponse - the standardized response containing updated profile data on success
@return error - gRPC error with relevant status code if authentication, validation, permission, or update fails
*/
func (h *RevertsIoGrpcHandler) UpdateSelfRevertProfile(ctx context.Context, req *pb.UpdateSelfRevertProfileRequest) (*pb.StandardRevertResponse, error) {
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

// ListRevertProfiles /*
/*
@Method ListRevertProfiles
@name List revert profiles with optional pagination and filtering
@description Handles listing of revert profiles based on query parameters such as start index, limit, page number, and optional name filter.
             Calls the service layer to retrieve paginated results, converts domain entities to protobuf representations,
             constructs a detailed list response including pagination metadata,
             and returns a standardized gRPC response.
@param ctx context.Context - the request context
@param req *pb.ListRevertProfilesRequest - the gRPC request containing pagination and filter parameters
@return *pb.StandardRevertResponse - the standardized response containing a list of revert profiles and pagination info on success
@return error - gRPC error with relevant status code if the list operation or response construction fails
*/
func (h *RevertsIoGrpcHandler) ListRevertProfiles(ctx context.Context, req *pb.ListRevertProfilesRequest) (*pb.StandardRevertResponse, error) {
	params := &entity.RevertProfileQueryParams{
		Start: req.GetStart(),
		Limit: req.GetLimit(),
		Page:  req.GetPage(),
		Name:  req.GetName(),
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

// CreateRevertMatchInvite /*
/*
@Method CreateRevertMatchInvite
@name Create a new revert match invitation
@description Handles the creation of a revert match invitation from an initiator to a receiver.
             Validates the initiator's identity from the context, ensures required request fields are provided,
             delegates creation to the RevertService, and handles various domain-specific errors
             such as invalid data, self-invitation attempts, duplicates, or missing profiles.
@param ctx context.Context - the request context, which must contain the initiator's user ID in context
@param req *pb.CreateRevertMatchInviteRequest - the gRPC request containing the receiver profile ID
@return *pb.StandardRevertResponse - a standardized response including the created revert match entity on success
@return error - gRPC error with appropriate status code for authentication, validation, or internal failures
*/
func (h *RevertsIoGrpcHandler) CreateRevertMatchInvite(ctx context.Context, req *pb.CreateRevertMatchInviteRequest) (*pb.StandardRevertResponse, error) {
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

// GetRevertMatch /*
/*
@Method GetRevertMatch
@name Retrieve a revert match by its ID
@description Fetches a revert match entity based on the provided match ID.
             Validates the match ID presence, delegates retrieval to the RevertService,
             and handles errors including invalid ID format and not found cases.
@param ctx context.Context - the request context
@param req *pb.GetRevertMatchRequest - the gRPC request containing the match ID
@return *pb.StandardRevertResponse - a standardized response containing the revert match on success
@return error - gRPC error indicating validation, not found, or internal errors
*/
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

// AcceptRevertMatchInvite /*
/*
@Method AcceptRevertMatchInvite
@name Accept a revert match invitation
@description Accepts a revert match invite identified by the provided match ID.
             Requires authentication and verifies the presence of the match ID.
             Handles errors such as invalid match ID, match not found, and invalid match status.
@param ctx context.Context - the request context containing authentication info
@param req *pb.AcceptRevertMatchInviteRequest - the gRPC request containing the match ID to accept
@return *pb.StandardRevertResponse - a standardized response containing the accepted revert match on success
@return error - gRPC error indicating authentication failure, invalid arguments, not found, failed precondition, or internal errors
*/
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

// RejectRevertMatchInvite /*
/*
@Method RejectRevertMatchInvite
@name Reject a revert match invitation
@description Rejects a revert match invite specified by the given match ID.
             Requires authentication and validates the presence of the match ID.
             Handles errors such as invalid match ID, match not found, and invalid match status.
@param ctx context.Context - the request context containing authentication information
@param req *pb.RejectRevertMatchInviteRequest - the gRPC request containing the match ID to reject
@return *pb.StandardRevertResponse - a standardized response containing the rejected revert match on success
@return error - gRPC error indicating authentication failure, invalid arguments, not found, failed precondition, or internal errors
*/
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

// EndRevertMatch /*
/*
@Method EndRevertMatch
@name End a revert match
@description Ends an active revert match specified by the given match ID.
             Requires authentication and validates the presence of the match ID.
             Returns appropriate gRPC errors for invalid match ID, not found matches, or invalid match status.
@param ctx context.Context - the request context containing authentication information
@param req *pb.EndRevertMatchRequest - the gRPC request containing the match ID to end
@return *pb.StandardRevertResponse - a standardized response containing the ended revert match on success
@return error - gRPC error indicating authentication failure, invalid arguments, not found, failed precondition, or internal errors
*/
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
