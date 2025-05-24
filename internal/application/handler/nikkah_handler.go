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
	"strings"
)

type NikkahIoGrpcHandler struct {
	pb.UnimplementedNikkahIoServiceServer
	NikkahSvc *services.NikkahService
}

func NewNikkahIoGrpcHandler(nikkahSvc *services.NikkahService) *NikkahIoGrpcHandler {
	return &NikkahIoGrpcHandler{
		NikkahSvc: nikkahSvc,
	}
}

func (h *NikkahIoGrpcHandler) CreateNikkahProfile(ctx context.Context, req *pb.CreateNikkahProfileRequest) (*pb.StandardNikkahResponse, error) {
	userID, ok := ctx.Value(auth.UserIDContextKey).(string)
	if !ok || userID == "" {
		return nil, status.Errorf(codes.Unauthenticated, "Unauthenticated: user ID not found in context")
	}

	if req.GetProfile() == nil {
		return nil, status.Errorf(codes.InvalidArgument, "Nikkah profile data is required")
	}
	if req.GetProfile().GetName() == "" {
		return nil, status.Errorf(codes.InvalidArgument, "Name is required")
	}
	if req.GetProfile().GetGender() == pb.NikkahProfile_GENDER_UNSPECIFIED {
		return nil, status.Errorf(codes.InvalidArgument, "Gender must be specified")
	}
	if req.GetProfile().GetBirthDate() == nil {
		return nil, status.Errorf(codes.InvalidArgument, "Birth date is required")
	}

	if req.GetProfile().GetUserId() != "" && req.GetProfile().GetUserId() != userID {
		return nil, status.Errorf(codes.PermissionDenied, "Cannot create profile for a different user ID than the authenticated one")
	}

	profileEntity, err := helper.ToEntityNikkahProfile(req.GetProfile())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "Invalid profile data format: %v", err)
	}
	profileEntity.UserID = userID

	createdProfile, err := h.NikkahSvc.CreateNikkahProfile(ctx, profileEntity)
	if err != nil {
		if errors.Is(err, errors.New("a nikkah profile already exists for this user ID")) {
			return nil, status.Errorf(codes.AlreadyExists, err.Error())
		}
		return nil, status.Errorf(codes.Internal, fmt.Sprintf("Failed to create nikkah profile: %v", err))
	}

	resp, helperErr := helper.StandardNikkahResponse(
		codes.OK,
		"success",
		"Nikkah profile created successfully",
		createdProfile,
	)
	if helperErr != nil {
		return nil, status.Errorf(codes.Internal, fmt.Sprintf("Failed to construct response: %v", helperErr))
	}
	return resp, nil
}

func (h *NikkahIoGrpcHandler) GetSelfNikkahProfile(ctx context.Context, req *pb.GetSelfNikkahProfileRequest) (*pb.StandardNikkahResponse, error) {
	userID, ok := ctx.Value(auth.UserIDContextKey).(string)
	if !ok || userID == "" {
		return nil, status.Errorf(codes.Unauthenticated, "Unauthenticated: user ID not found in context")
	}

	profile, err := h.NikkahSvc.GetNikkahProfileByUserID(ctx, userID)
	if err != nil {
		if errors.Is(err, helper.ErrNotFound) {
			return nil, status.Errorf(codes.NotFound, "Nikkah profile not found for this user")
		}
		return nil, status.Errorf(codes.Internal, fmt.Sprintf("Failed to get self nikkah profile: %v", err))
	}

	resp, helperErr := helper.StandardNikkahResponse(
		codes.OK,
		"success",
		"Self nikkah profile retrieved successfully",
		profile,
	)
	if helperErr != nil {
		return nil, status.Errorf(codes.Internal, fmt.Sprintf("Failed to construct response: %v", helperErr))
	}
	return resp, nil
}

func (h *NikkahIoGrpcHandler) UpdateSelfNikkahProfile(ctx context.Context, req *pb.UpdateSelfNikkahProfileRequest) (*pb.StandardNikkahResponse, error) {
	userID, ok := ctx.Value(auth.UserIDContextKey).(string)
	if !ok || userID == "" {
		return nil, status.Errorf(codes.Unauthenticated, "Unauthenticated: user ID not found in context")
	}

	if req.GetProfileId() == "" {
		return nil, status.Errorf(codes.InvalidArgument, "Profile ID is required in the URL path (e.g., /v1/nikkah/profile/YOUR_PROFILE_ID).")
	}
	profileID, err := uuid.Parse(req.GetProfileId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "Invalid profile ID format in URL path: %v", err)
	}

	if req.GetProfile() == nil {
		return nil, status.Errorf(codes.InvalidArgument, "Nikkah profile data is required in the request body.")
	}

	profileEntity, err := helper.ToEntityNikkahProfile(req.GetProfile())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "Invalid profile data format in request body: %v", err)
	}

	existingProfile, err := h.NikkahSvc.GetNikkahProfileByID(ctx, profileID)
	if err != nil {
		if errors.Is(err, helper.ErrNotFound) {
			return nil, status.Errorf(codes.NotFound, "Profile with ID %s to update not found.", profileID.String())
		}
		return nil, status.Errorf(codes.Internal, fmt.Sprintf("Failed to verify existing profile: %v", err))
	}

	if existingProfile.UserID != userID {
		return nil, status.Errorf(codes.PermissionDenied, "You are not authorized to update this profile. This profile belongs to a different user.")
	}

	profileEntity.ID = profileID
	profileEntity.UserID = userID
	profileEntity.CreatedAt = existingProfile.CreatedAt

	updatedProfile, err := h.NikkahSvc.UpdateNikkahProfile(ctx, profileEntity)
	if err != nil {
		if errors.Is(err, helper.ErrAlreadyExists) {
			return nil, status.Errorf(codes.AlreadyExists, "Update caused a data conflict: %v", err.Error())
		}
		if errors.Is(err, helper.ErrNotFound) {
			return nil, status.Errorf(codes.NotFound, "Profile to update not found.")
		}
		if errors.Is(err, errors.New("user ID cannot be changed")) {
			return nil, status.Errorf(codes.InvalidArgument, err.Error())
		}
		return nil, status.Errorf(codes.Internal, fmt.Sprintf("Failed to update self nikkah profile: %v", err))
	}

	resp, helperErr := helper.StandardNikkahResponse(
		codes.OK,
		"success",
		"Self nikkah profile updated successfully",
		updatedProfile,
	)
	if helperErr != nil {
		return nil, status.Errorf(codes.Internal, fmt.Sprintf("Failed to construct response: %v", helperErr))
	}
	return resp, nil
}

func (h *NikkahIoGrpcHandler) GetNikkahProfile(ctx context.Context, req *pb.GetNikkahProfileRequest) (*pb.StandardNikkahResponse, error) {
	if req.GetProfileId() == "" {
		return nil, status.Errorf(codes.InvalidArgument, "Profile ID is required")
	}

	profileID, err := uuid.Parse(req.GetProfileId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "Invalid profile ID format: %v", err)
	}

	profile, err := h.NikkahSvc.GetNikkahProfileByID(ctx, profileID)
	if err != nil {
		if errors.Is(err, errors.New(fmt.Sprintf("profile not found: nikkah profile with ID %s not found", profileID.String()))) {
			return nil, status.Errorf(codes.NotFound, "Nikkah profile not found.")
		}
		return nil, status.Errorf(codes.Internal, fmt.Sprintf("Failed to retrieve nikkah profile: %v", err))
	}

	resp, helperErr := helper.StandardNikkahResponse(
		codes.OK,
		"success",
		"Nikkah profile retrieved successfully",
		profile,
	)
	if helperErr != nil {
		return nil, status.Errorf(codes.Internal, fmt.Sprintf("Failed to construct response: %v", helperErr))
	}
	return resp, nil
}

func (h *NikkahIoGrpcHandler) ListNikkahProfiles(ctx context.Context, req *pb.ListNikkahProfilesRequest) (*pb.StandardNikkahResponse, error) {

	params := &entity.NikkahProfileQueryParams{
		Start: req.GetStart(),
		Limit: req.GetLimit(),
		Page:  req.GetPage(),
		Name:  req.GetName(),
	}

	queryResult, err := h.NikkahSvc.ListNikkahProfiles(ctx, params)
	if err != nil {
		return nil, status.Errorf(codes.Internal, fmt.Sprintf("Failed to list nikkah profiles: %v", err))
	}

	protoProfiles := make([]*pb.NikkahProfile, len(queryResult.Profiles))
	for i, p := range queryResult.Profiles {
		protoProfiles[i] = helper.ToProtoNikkahProfile(p)
	}

	listResponse := &pb.ListNikkahProfilesResponse{
		Profiles:    protoProfiles,
		TotalCount:  int32(queryResult.TotalCount),
		CurrentPage: queryResult.CurrentPage,
		TotalPages:  queryResult.TotalPages,
	}

	resp, helperErr := helper.StandardNikkahResponse(
		codes.OK,
		"success",
		"Nikkah profiles listed successfully",
		listResponse,
	)
	if helperErr != nil {
		return nil, status.Errorf(codes.Internal, fmt.Sprintf("Failed to construct response: %v", helperErr))
	}

	resp.Data = &pb.StandardNikkahResponse_ListProfilesResponse{ListProfilesResponse: listResponse}

	return resp, nil
}

func (h *NikkahIoGrpcHandler) InitiateNikkahLike(ctx context.Context, req *pb.InitiateNikkahLikeRequest) (*pb.StandardNikkahResponse, error) {
	likerUserID, ok := ctx.Value(auth.UserIDContextKey).(string)
	if !ok || likerUserID == "" {
		return nil, status.Errorf(codes.Unauthenticated, "Unauthenticated: user ID not found in context")
	}

	protoLike := req.GetLike()
	if protoLike == nil {
		return nil, status.Errorf(codes.InvalidArgument, "Like object is required in the request body.")
	}

	if protoLike.GetLikedProfileId() == "" {
		return nil, status.Errorf(codes.InvalidArgument, "Liked Profile ID is required.")
	}
	likedProfileID, err := uuid.Parse(protoLike.GetLikedProfileId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "Invalid Liked Profile ID format: %v", err)
	}

	createdLike, err := h.NikkahSvc.CreateNikkahLike(ctx, likerUserID, likedProfileID)
	if err != nil {
		if errors.Is(err, errors.New("service: cannot like your own profile")) {
			return nil, status.Errorf(codes.InvalidArgument, err.Error())
		}
		if errors.Is(err, errors.New("service: liker's profile not found. Cannot create like")) {
			return nil, status.Errorf(codes.NotFound, err.Error())
		}

		if errors.Is(err, fmt.Errorf("service: like already exists from profile %s to profile %s", likerUserID, likedProfileID.String())) {
			return nil, status.Errorf(codes.AlreadyExists, err.Error())
		}
		return nil, status.Errorf(codes.Internal, fmt.Sprintf("Failed to initiate nikkah like: %v", err))
	}

	protoResponseLike := helper.ToProtoNikkahLike(createdLike)
	resp, helperErr := helper.StandardNikkahResponse(
		codes.OK,
		"success",
		"Nikkah like initiated successfully",
		nil,
	)
	if helperErr != nil {
		return nil, status.Errorf(codes.Internal, fmt.Sprintf("Failed to construct response: %v", helperErr))
	}

	resp.Data = &pb.StandardNikkahResponse_Like{Like: protoResponseLike}

	return resp, nil
}

func (h *NikkahIoGrpcHandler) GetNikkahLike(ctx context.Context, req *pb.GetNikkahLikeRequest) (*pb.StandardNikkahResponse, error) {
	if req.GetLikeId() == "" {
		return nil, status.Errorf(codes.InvalidArgument, "Like ID is required.")
	}
	likeID, err := uuid.Parse(req.GetLikeId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "Invalid Like ID format: %v", err)
	}

	retrievedLike, err := h.NikkahSvc.GetNikkahLikeByID(ctx, likeID)
	if err != nil {
		if errors.Is(err, errors.New("service: like ID cannot be empty")) || strings.Contains(err.Error(), "invalid Like ID format") {
			return nil, status.Errorf(codes.InvalidArgument, err.Error())
		}
		if errors.Is(err, errors.New("service: nikkah like not found")) {
			return nil, status.Errorf(codes.NotFound, err.Error())
		}
		return nil, status.Errorf(codes.Internal, fmt.Sprintf("Failed to retrieve nikkah like: %v", err))
	}

	protoLike := helper.ToProtoNikkahLike(retrievedLike)

	resp, helperErr := helper.StandardNikkahResponse(
		codes.OK,
		"success",
		"Nikkah like retrieved successfully",
		nil,
	)
	if helperErr != nil {
		return nil, status.Errorf(codes.Internal, fmt.Sprintf("Failed to construct response: %v", helperErr))
	}
	resp.Data = &pb.StandardNikkahResponse_Like{Like: protoLike}

	return resp, nil
}

func (h *NikkahIoGrpcHandler) CancelNikkahLike(ctx context.Context, req *pb.CancelNikkahLikeRequest) (*pb.StandardNikkahResponse, error) {
	requestingUserID, ok := ctx.Value(auth.UserIDContextKey).(string)
	if !ok || requestingUserID == "" {
		return nil, status.Errorf(codes.Unauthenticated, "Unauthenticated: user ID not found in context")
	}

	if req.GetLikeId() == "" {
		return nil, status.Errorf(codes.InvalidArgument, "Like ID is required for cancellation.")
	}
	likeID, err := uuid.Parse(req.GetLikeId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "Invalid Like ID format for cancellation: %v", err)
	}

	updatedLike, err := h.NikkahSvc.CancelNikkahLike(ctx, likeID, requestingUserID)
	if err != nil {
		if errors.Is(err, errors.New("service: nikkah like with ID "+likeID.String()+" not found")) {
			return nil, status.Errorf(codes.NotFound, err.Error())
		}
		if errors.Is(err, errors.New("service: unauthorized to cancel this nikkah like. Only the liker can cancel it")) {
			return nil, status.Errorf(codes.PermissionDenied, err.Error())
		}
		if strings.Contains(err.Error(), "cannot be cancelled as its current status is") {
			return nil, status.Errorf(codes.FailedPrecondition, err.Error())
		}
		if errors.Is(err, errors.New("service: requesting user's profile not found. Cannot cancel like")) {
			return nil, status.Errorf(codes.NotFound, err.Error())
		}
		return nil, status.Errorf(codes.Internal, fmt.Sprintf("Failed to cancel nikkah like: %v", err))
	}

	protoLike := helper.ToProtoNikkahLike(updatedLike)

	resp, helperErr := helper.StandardNikkahResponse(
		codes.OK,
		"success",
		"Nikkah like cancelled successfully",
		nil,
	)
	if helperErr != nil {
		return nil, status.Errorf(codes.Internal, fmt.Sprintf("Failed to construct response: %v", helperErr))
	}

	resp.Data = &pb.StandardNikkahResponse_Like{Like: protoLike}

	return resp, nil
}

func (h *NikkahIoGrpcHandler) CompleteNikkahLike(ctx context.Context, req *pb.CompleteNikkahLikeRequest) (*pb.StandardNikkahResponse, error) {
	requestingUserID, ok := ctx.Value(auth.UserIDContextKey).(string)
	if !ok || requestingUserID == "" {
		return nil, status.Errorf(codes.Unauthenticated, "Unauthenticated: user ID not found in context")
	}
	if req.GetLikeId() == "" {
		return nil, status.Errorf(codes.InvalidArgument, "Like ID is required for completion.")
	}
	likeID, err := uuid.Parse(req.GetLikeId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "Invalid Like ID format for completion: %v", err)
	}

	updatedLike, createdMatch, err := h.NikkahSvc.CompleteNikkahLike(ctx, likeID, requestingUserID)
	if err != nil {
		if errors.Is(err, errors.New("service: nikkah like with ID "+likeID.String()+" not found")) {
			return nil, status.Errorf(codes.NotFound, err.Error())
		}
		if errors.Is(err, errors.New("service: requesting user's profile not found. Cannot complete like")) {
			return nil, status.Errorf(codes.NotFound, err.Error())
		}
		if errors.Is(err, errors.New("service: unauthorized to complete this nikkah like. Only the liked profile can complete it")) {
			return nil, status.Errorf(codes.PermissionDenied, err.Error())
		}
		if strings.Contains(err.Error(), "cannot be completed as its current status is") {
			return nil, status.Errorf(codes.FailedPrecondition, err.Error())
		}
		if errors.Is(err, errors.New("service: no mutual like found from the other profile to complete this match")) {
			return nil, status.Errorf(codes.FailedPrecondition, err.Error())
		}
		if strings.Contains(err.Error(), "reverse nikkah like from") && strings.Contains(err.Error(), "cannot be completed as its current status is") {
			return nil, status.Errorf(codes.FailedPrecondition, err.Error())
		}
		return nil, status.Errorf(codes.Internal, fmt.Sprintf("Failed to complete nikkah like: %v", err))
	}

	protoLike := helper.ToProtoNikkahLike(updatedLike)
	protoMatch := helper.ToProtoNikkahMatch(createdMatch)

	completeResp := &pb.CompleteNikkahLikeResponse{
		Like:  protoLike,
		Match: protoMatch,
	}

	resp, helperErr := helper.StandardNikkahResponse(
		codes.OK,
		"success",
		"Nikkah like completed and match created successfully",
		nil,
	)
	if helperErr != nil {
		return nil, status.Errorf(codes.Internal, fmt.Sprintf("Failed to construct response: %v", helperErr))
	}

	resp.Data = &pb.StandardNikkahResponse_CompleteNikahLike{CompleteNikahLike: completeResp}

	return resp, nil
}

func (h *NikkahIoGrpcHandler) AcceptNikkahMatchInvite(ctx context.Context, req *pb.AcceptNikkahMatchInviteRequest) (*pb.StandardNikkahResponse, error) {
	requestingUserID, ok := ctx.Value(auth.UserIDContextKey).(string)
	if !ok || requestingUserID == "" {
		return nil, status.Errorf(codes.Unauthenticated, "Authentication required: user ID not found in context.")
	}

	if req.GetMatchId() == "" {
		return nil, status.Errorf(codes.InvalidArgument, "Match ID is required to accept the invite.")
	}
	matchID, err := uuid.Parse(req.GetMatchId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "Invalid Match ID format: %v", err)
	}

	updatedMatch, err := h.NikkahSvc.AcceptNikkahMatchInvite(ctx, matchID, requestingUserID)
	if err != nil {
		if errors.Is(err, helper.ErrNotFound) {
			return nil, status.Errorf(codes.NotFound, "Match or requesting user's profile not found: %v", err)
		}
		if errors.Is(err, errors.New("service: unauthorized to accept this nikkah match. Only the receiver of the match invite can accept it.")) {
			return nil, status.Errorf(codes.PermissionDenied, "Unauthorized: You do not have permission to accept this match.")
		}
		if strings.Contains(err.Error(), "cannot be accepted as its current status is") {
			return nil, status.Errorf(codes.FailedPrecondition, "Cannot accept match: %v", err)
		}
		return nil, status.Errorf(codes.Internal, "Failed to accept nikkah match invite due to internal service error: %v", err)
	}

	protoMatch := helper.ToProtoNikkahMatch(updatedMatch)

	resp, helperErr := helper.StandardNikkahResponse(
		codes.OK,
		"success",
		"Nikkah match invite accepted successfully",
		nil,
	)
	if helperErr != nil {
		return nil, status.Errorf(codes.Internal, fmt.Sprintf("Failed to construct response: %v", helperErr))
	}
	resp.Data = &pb.StandardNikkahResponse_Match{Match: protoMatch}

	return resp, nil
}

func (h *NikkahIoGrpcHandler) GetNikkahMatch(ctx context.Context, req *pb.GetNikkahMatchRequest) (*pb.StandardNikkahResponse, error) {
	requestingUserID, ok := ctx.Value(auth.UserIDContextKey).(string)
	if !ok || requestingUserID == "" {
		return nil, status.Errorf(codes.Unauthenticated, "Authentication required: user ID not found in context.")
	}

	if req.GetMatchId() == "" {
		return nil, status.Errorf(codes.InvalidArgument, "Match ID is required to retrieve match details.")
	}
	matchID, err := uuid.Parse(req.GetMatchId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "Invalid Match ID format: %v", err)
	}

	nikkahMatch, err := h.NikkahSvc.GetNikkahMatch(ctx, matchID, requestingUserID)
	if err != nil {
		if errors.Is(err, helper.ErrNotFound) {
			return nil, status.Errorf(codes.NotFound, "Nikkah match or requesting user's profile not found: %v", err)
		}
		if errors.Is(err, errors.New("service: unauthorized to view this nikkah match. You are not a participant in this match.")) {
			return nil, status.Errorf(codes.PermissionDenied, "Unauthorized: You do not have permission to view this match.")
		}
		return nil, status.Errorf(codes.Internal, "Failed to retrieve nikkah match due to internal service error: %v", err)
	}

	protoMatch := helper.ToProtoNikkahMatch(nikkahMatch)

	resp, helperErr := helper.StandardNikkahResponse(
		codes.OK,
		"success",
		"Nikkah match retrieved successfully.",
		nil,
	)
	if helperErr != nil {
		return nil, status.Errorf(codes.Internal, fmt.Sprintf("Failed to construct response: %v", helperErr))
	}
	resp.Data = &pb.StandardNikkahResponse_Match{Match: protoMatch}

	return resp, nil
}

func (h *NikkahIoGrpcHandler) RejectNikkahMatchInvite(ctx context.Context, req *pb.RejectNikkahMatchInviteRequest) (*pb.StandardNikkahResponse, error) {
	requestingUserID, ok := ctx.Value(auth.UserIDContextKey).(string)
	if !ok || requestingUserID == "" {
		return nil, status.Errorf(codes.Unauthenticated, "Authentication required: user ID not found in context.")
	}

	if req.GetMatchId() == "" {
		return nil, status.Errorf(codes.InvalidArgument, "Match ID is required to reject the invite.")
	}
	matchID, err := uuid.Parse(req.GetMatchId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "Invalid Match ID format for rejection: %v", err)
	}

	updatedMatch, err := h.NikkahSvc.RejectNikkahMatchInvite(ctx, matchID, requestingUserID)
	if err != nil {
		if errors.Is(err, helper.ErrNotFound) {
			return nil, status.Errorf(codes.NotFound, "Match or requesting user's profile not found: %v", err)
		}
		if errors.Is(err, errors.New("service: unauthorized to reject this nikkah match. Only the receiver of the match invite can reject it.")) {
			return nil, status.Errorf(codes.PermissionDenied, "Unauthorized: You do not have permission to reject this match.")
		}
		if strings.Contains(err.Error(), "cannot be rejected as its current status is") {
			return nil, status.Errorf(codes.FailedPrecondition, "Cannot reject match: %v", err)
		}
		return nil, status.Errorf(codes.Internal, "Failed to reject nikkah match invite due to internal service error: %v", err)
	}

	protoMatch := helper.ToProtoNikkahMatch(updatedMatch)

	resp, helperErr := helper.StandardNikkahResponse(
		codes.OK,
		"success",
		"Nikkah match invite rejected successfully.",
		nil,
	)
	if helperErr != nil {
		return nil, status.Errorf(codes.Internal, fmt.Sprintf("Failed to construct response: %v", helperErr))
	}
	resp.Data = &pb.StandardNikkahResponse_Match{Match: protoMatch}

	return resp, nil

}

func (h *NikkahIoGrpcHandler) EndNikkahMatch(ctx context.Context, req *pb.EndNikkahMatchRequest) (*pb.StandardNikkahResponse, error) {
	requestingUserID, ok := ctx.Value(auth.UserIDContextKey).(string)
	if !ok || requestingUserID == "" {
		return nil, status.Errorf(codes.Unauthenticated, "Authentication required: user ID not found in context.")
	}

	if req.GetMatchId() == "" {
		return nil, status.Errorf(codes.InvalidArgument, "Match ID is required to end the match.")
	}
	matchID, err := uuid.Parse(req.GetMatchId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "Invalid Match ID format for ending: %v", err)
	}

	updatedMatch, err := h.NikkahSvc.EndNikkahMatch(ctx, matchID, requestingUserID)
	if err != nil {
		if errors.Is(err, helper.ErrNotFound) {
			return nil, status.Errorf(codes.NotFound, "Match or requesting user's profile not found: %v", err)
		}
		if strings.Contains(err.Error(), "unauthorized to end this nikkah match") {
			return nil, status.Errorf(codes.PermissionDenied, "Unauthorized: You do not have permission to end this match.")
		}
		if strings.Contains(err.Error(), "is already ended") {
			return nil, status.Errorf(codes.FailedPrecondition, "Cannot end match: %v", err)
		}
		return nil, status.Errorf(codes.Internal, "Failed to end nikkah match due to internal service error: %v", err)
	}

	protoMatch := helper.ToProtoNikkahMatch(updatedMatch)

	resp, helperErr := helper.StandardNikkahResponse(
		codes.OK,
		"success",
		"Nikkah match ended successfully.",
		nil,
	)
	if helperErr != nil {
		return nil, status.Errorf(codes.Internal, fmt.Sprintf("Failed to construct response: %v", helperErr))
	}
	resp.Data = &pb.StandardNikkahResponse_Match{Match: protoMatch}

	return resp, nil
}
