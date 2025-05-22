package services

import (
	"context"
	"errors"
	"fmt"
	"github.com/mnadev/limestone/internal/application/helper"
	"math"
	"time"

	"github.com/google/uuid"
	"github.com/mnadev/limestone/internal/application/domain/entity"
	"github.com/mnadev/limestone/internal/application/repository"
)

type NikkahService struct {
	RepoNikkah repository.NikkahRepository
}

func NewNikkahService(nikahRepo repository.NikkahRepository) *NikkahService {
	return &NikkahService{
		RepoNikkah: nikahRepo,
	}
}

func (s *NikkahService) CreateNikkahProfile(ctx context.Context, nikkahProfile *entity.NikkahProfile) (*entity.NikkahProfile, error) {
	if nikkahProfile.ID == uuid.Nil {
		nikkahProfile.ID = uuid.New()
	}
	nikkahProfile.CreatedAt = time.Now()
	nikkahProfile.UpdatedAt = time.Now()

	createdProfile, err := s.RepoNikkah.CreateProfile(ctx, nikkahProfile)
	if err != nil {
		return nil, errors.New("failed to save nikkah profile: " + err.Error())
	}
	return createdProfile, nil
}

func (s *NikkahService) GetNikkahProfileByUserID(ctx context.Context, userID string) (*entity.NikkahProfile, error) {
	profile, err := s.RepoNikkah.GetProfileByUserID(ctx, userID)
	if err != nil {
		return nil, errors.New("failed to retrieve nikkah profile by user ID: " + err.Error())
	}
	return profile, nil
}

func (s *NikkahService) GetNikkahProfileByID(ctx context.Context, profileID uuid.UUID) (*entity.NikkahProfile, error) {
	if profileID == uuid.Nil {
		return nil, errors.New("profile ID cannot be empty")
	}
	profile, err := s.RepoNikkah.GetProfileByID(ctx, profileID)
	if err != nil {
		if errors.Is(err, helper.ErrNotFound) {
			return nil, fmt.Errorf("profile not found: %w", err)
		}
		return nil, fmt.Errorf("service: failed to get profile by ID: %w", err)
	}
	return profile, nil
}

func (s *NikkahService) UpdateNikkahProfile(ctx context.Context, profile *entity.NikkahProfile) (*entity.NikkahProfile, error) {
	if profile.ID == uuid.Nil {
		return nil, errors.New("profile ID cannot be empty for update")
	}

	existingProfile, err := s.RepoNikkah.GetProfileByID(ctx, profile.ID)
	if err != nil {
		if errors.Is(err, helper.ErrNotFound) {
			return nil, fmt.Errorf("profile to update not found: %w", err)
		}
		return nil, fmt.Errorf("service: failed to retrieve existing profile for update: %w", err)
	}

	if profile.UserID != "" && profile.UserID != existingProfile.UserID {
		return nil, errors.New("user ID cannot be changed")
	}
	profile.UserID = existingProfile.UserID
	profile.CreatedAt = existingProfile.CreatedAt
	profile.UpdatedAt = time.Now()

	updatedProfile, err := s.RepoNikkah.UpdateProfile(ctx, profile)
	if err != nil {
		if errors.Is(err, helper.ErrAlreadyExists) {
			return nil, fmt.Errorf("updated profile caused duplicate user ID: %w", err)
		}
		if errors.Is(err, helper.ErrNotFound) {
			return nil, fmt.Errorf("profile to update not found: %w", err)
		}
		return nil, fmt.Errorf("service: failed to update profile: %w", err)
	}
	return updatedProfile, nil
}

func (s *NikkahService) ListNikkahProfiles(ctx context.Context, params *entity.NikkahProfileQueryParams) (*entity.NikkahProfileQueryResult, error) {
	if params.Limit <= 0 {
		params.Limit = 10
	}
	if params.Limit > 100 {
		params.Limit = 100
	}

	if params.Page <= 0 {
		params.Page = 1
	}

	if params.Start < 0 {
		params.Start = 0
	}
	if params.Start == 0 && params.Page > 1 {
		params.Start = (params.Page - 1) * params.Limit
	}
	profiles, totalCount, err := s.RepoNikkah.ListProfiles(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("service: failed to list profiles: %w", err)
	}

	var totalPages int32
	if params.Limit > 0 && totalCount > 0 {
		totalPages = int32(math.Ceil(float64(totalCount) / float64(params.Limit)))
	} else {
		totalPages = 0
	}
	currentPage := params.Page
	if currentPage == 0 && params.Limit > 0 {
		currentPage = (params.Start / params.Limit) + 1
	} else if currentPage == 0 && params.Start == 0 {
		currentPage = 1
	}

	return &entity.NikkahProfileQueryResult{
		Profiles:    profiles,
		TotalCount:  int64(totalCount),
		CurrentPage: currentPage,
		TotalPages:  totalPages,
	}, nil
}

func (s *NikkahService) CreateNikkahLike(ctx context.Context, likerUserID string, likedProfileID uuid.UUID) (*entity.NikkahLike, error) {
	likerProfile, err := s.RepoNikkah.GetProfileByUserID(ctx, likerUserID)
	if err != nil {
		if errors.Is(err, helper.ErrNotFound) {
			return nil, errors.New("service: liker's profile not found")
		}
		return nil, fmt.Errorf("service: failed to get liker's profile: %w", err)
	}
	fmt.Println(likerProfile.ID.String())
	fmt.Println(likedProfileID.String())

	existingLike, err := s.RepoNikkah.GetLikeByLikerAndLikedProfileID(ctx, likerProfile.ID.String(), likedProfileID.String())
	if err == nil && existingLike.Status == entity.LikeStatusInitiated {
		return nil, errors.New("service: you have already sent a like to this profile")
	}
	if err == nil && existingLike.Status == entity.LikeStatusCompleted {
		return nil, errors.New("service: you have already matched with this profile")
	}

	reverseLike, err := s.RepoNikkah.GetLikeByLikerAndLikedProfileID(ctx, likedProfileID.String(), likerProfile.ID.String())

	if err == nil && reverseLike.Status == entity.LikeStatusInitiated {
		like := &entity.NikkahLike{
			ID:             uuid.New(),
			LikerProfileID: likerProfile.ID.String(),
			LikedProfileID: likedProfileID.String(),
			Status:         entity.LikeStatusCompleted,
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
		}
		_, err = s.RepoNikkah.CreateLike(ctx, like)
		if err != nil {
			return nil, fmt.Errorf("service: failed to create liker's completed like: %w", err)
		}

		reverseLike.Status = entity.LikeStatusCompleted
		reverseLike.UpdatedAt = time.Now()
		_, err = s.RepoNikkah.UpdateLike(ctx, reverseLike)
		if err != nil {
			return nil, fmt.Errorf("service: failed to update reverse like to completed: %w", err)
		}

		profileID1Str := likerProfile.ID.String()
		profileID2Str := likedProfileID.String()

		matchProfileID1UUID, err := uuid.Parse(profileID1Str)
		if err != nil {
			return nil, fmt.Errorf("service: failed to parse profileID1 to UUID: %w", err)
		}
		matchProfileID2UUID, err := uuid.Parse(profileID2Str)
		if err != nil {
			return nil, fmt.Errorf("service: failed to parse profileID2 to UUID: %w", err)
		}

		var finalMatchProfileID1, finalMatchProfileID2 uuid.UUID
		if matchProfileID1UUID.String() < matchProfileID2UUID.String() {
			finalMatchProfileID1 = matchProfileID1UUID
			finalMatchProfileID2 = matchProfileID2UUID
		} else {
			finalMatchProfileID1 = matchProfileID2UUID
			finalMatchProfileID2 = matchProfileID1UUID
		}

		match := &entity.NikkahMatch{
			ID:                 uuid.New(),
			InitiatorProfileID: finalMatchProfileID1,
			ReceiverProfileID:  finalMatchProfileID2,
			Status:             entity.MatchStatusInitiated,
			CreatedAt:          time.Now(),
			UpdatedAt:          time.Now(),
		}
		_, err = s.RepoNikkah.CreateMatch(ctx, match)
		if err != nil {
			return nil, fmt.Errorf("service: failed to create match after mutual like: %w", err)
		}

		return like, nil
	}

	like := &entity.NikkahLike{
		ID:             uuid.New(),
		LikerProfileID: likerProfile.ID.String(),
		LikedProfileID: likedProfileID.String(),
		Status:         entity.LikeStatusInitiated,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}
	return s.RepoNikkah.CreateLike(ctx, like)
}

func (s *NikkahService) GetNikkahLikeByID(ctx context.Context, likeID uuid.UUID) (*entity.NikkahLike, error) {
	if likeID == uuid.Nil {
		return nil, errors.New("service: like ID cannot be empty")
	}

	nikkahLike, err := s.RepoNikkah.GetLikeByID(ctx, likeID)
	if err != nil {
		if errors.Is(err, helper.ErrNotFound) {
			return nil, fmt.Errorf("service: nikkah like not found: %w", err)
		}
		return nil, fmt.Errorf("service: failed to retrieve nikkah like: %w", err)
	}

	return nikkahLike, nil
}

func (s *NikkahService) CancelNikkahLike(ctx context.Context, likeID uuid.UUID, requestingUserID string) (*entity.NikkahLike, error) {
	nikkahLike, err := s.RepoNikkah.GetLikeByID(ctx, likeID)
	if err != nil {
		if errors.Is(err, helper.ErrNotFound) {
			return nil, fmt.Errorf("service: nikkah like with ID %s not found: %w", likeID.String(), err)
		}
		return nil, fmt.Errorf("service: failed to retrieve nikkah like %s: %w", likeID.String(), err)
	}

	requestingUserProfile, err := s.RepoNikkah.GetProfileByUserID(ctx, requestingUserID)
	fmt.Printf("xxx %v", requestingUserProfile)
	if err != nil {
		if errors.Is(err, helper.ErrNotFound) {
			return nil, errors.New("service: requesting user's profile not found. Cannot cancel like.")
		}
		return nil, fmt.Errorf("service: failed to get requesting user's profile: %w", err)
	}

	if nikkahLike.LikerProfileID != requestingUserProfile.UserID {
		return nil, errors.New("service: unauthorized to cancel this nikkah like. Only the liker can cancel it.")
	}

	if nikkahLike.Status != entity.LikeStatusInitiated {
		return nil, fmt.Errorf("service: nikkah like with ID %s cannot be cancelled as its current status is %s", likeID.String(), nikkahLike.Status)
	}

	nikkahLike.Status = entity.LikeStatusCancelled
	nikkahLike.UpdatedAt = time.Now()

	updatedLike, err := s.RepoNikkah.UpdateLike(ctx, nikkahLike)
	if err != nil {
		return nil, fmt.Errorf("service: failed to update nikkah like status to CANCELLED for ID %s: %w", likeID.String(), err)
	}

	return updatedLike, nil
}

func (s *NikkahService) CompleteNikkahLike(ctx context.Context, likeID uuid.UUID, requestingUserID string) (*entity.NikkahLike, *entity.NikkahMatch, error) {
	initiatingLike, err := s.RepoNikkah.GetLikeByID(ctx, likeID)
	if err != nil {
		if errors.Is(err, helper.ErrNotFound) {
			return nil, nil, fmt.Errorf("service: nikkah like with ID %s not found: %w", likeID.String(), err)
		}
		return nil, nil, fmt.Errorf("service: failed to retrieve nikkah like %s: %w", likeID.String(), err)
	}

	requestingUserProfile, err := s.RepoNikkah.GetProfileByUserID(ctx, requestingUserID)
	if err != nil {
		if errors.Is(err, helper.ErrNotFound) {
			return nil, nil, errors.New("service: requesting user's profile not found. Cannot complete like.")
		}
		return nil, nil, fmt.Errorf("service: failed to get requesting user's profile: %w", err)
	}

	if initiatingLike.LikedProfileID != requestingUserProfile.UserID {
		return nil, nil, errors.New("service: unauthorized to complete this nikkah like. Only the liked profile can complete it.")
	}

	if initiatingLike.Status != entity.LikeStatusInitiated {
		return nil, nil, fmt.Errorf("service: nikkah like with ID %s cannot be completed as its current status is %s", likeID.String(), initiatingLike.Status)
	}

	reverseLike, err := s.RepoNikkah.GetLikeByLikerAndLikedProfileID(
		ctx,
		initiatingLike.LikedProfileID,
		initiatingLike.LikerProfileID,
	)
	if err != nil {
		if errors.Is(err, helper.ErrNotFound) {
			return nil, nil, errors.New("service: no mutual like found from the other profile to complete this match")
		}
		return nil, nil, fmt.Errorf("service: failed to check for reverse like: %w", err)
	}

	if reverseLike.Status != entity.LikeStatusInitiated {
		return nil, nil, fmt.Errorf("service: reverse nikkah like from %s to %s cannot be completed as its current status is %s",
			reverseLike.LikerProfileID, reverseLike.LikedProfileID, reverseLike.Status)
	}

	initiatingLike.Status = entity.LikeStatusCompleted
	initiatingLike.UpdatedAt = time.Now()

	reverseLike.Status = entity.LikeStatusCompleted
	reverseLike.UpdatedAt = time.Now()

	updatedInitiatingLike, err := s.RepoNikkah.UpdateLike(ctx, initiatingLike)
	if err != nil {
		return nil, nil, fmt.Errorf("service: failed to update initiating nikkah like to COMPLETED for ID %s: %w", likeID.String(), err)
	}

	_, err = s.RepoNikkah.UpdateLike(ctx, reverseLike)
	if err != nil {
		return nil, nil, fmt.Errorf("service: failed to update reverse nikkah like to COMPLETED for ID %s: %w", reverseLike.ID.String(), err)
	}

	profileA := uuid.MustParse(initiatingLike.LikerProfileID)
	profileB, _ := uuid.Parse(initiatingLike.LikedProfileID)

	var matchProfileID1, matchProfileID2 uuid.UUID
	if profileA.String() < profileB.String() {
		matchProfileID1 = profileA
		matchProfileID2 = profileB
	} else {
		matchProfileID1 = profileB
		matchProfileID2 = profileA
	}

	match := &entity.NikkahMatch{
		ID:                 uuid.New(),
		InitiatorProfileID: matchProfileID1,
		ReceiverProfileID:  matchProfileID2,
		Status:             entity.MatchStatusInitiated,
		CreatedAt:          time.Now(),
		UpdatedAt:          time.Now(),
	}

	createdMatch, err := s.RepoNikkah.CreateMatch(ctx, match)
	if err != nil {
		return updatedInitiatingLike, nil, fmt.Errorf("service: failed to create nikkah match: %w", err)
	}

	return updatedInitiatingLike, createdMatch, nil
}

func (s *NikkahService) AcceptNikkahMatchInvite(ctx context.Context, matchID uuid.UUID, requestingUserID string) (*entity.NikkahMatch, error) {
	nikkahMatch, err := s.RepoNikkah.GetMatchByID(ctx, matchID)
	if err != nil {
		if errors.Is(err, helper.ErrNotFound) {
			return nil, fmt.Errorf("service: nikkah match with ID %s not found: %w", matchID.String(), err)
		}
		return nil, fmt.Errorf("service: failed to retrieve nikkah match %s: %w", matchID.String(), err)
	}

	requestingUserProfile, err := s.RepoNikkah.GetProfileByUserID(ctx, requestingUserID)
	if err != nil {
		if errors.Is(err, helper.ErrNotFound) {
			return nil, errors.New("service: requesting user's profile not found. Cannot accept match.")
		}
		return nil, fmt.Errorf("service: failed to get requesting user's profile: %w", err)
	}

	requestingUserProfileTemp, _ := uuid.Parse(requestingUserProfile.UserID)
	if nikkahMatch.ReceiverProfileID != requestingUserProfileTemp {
		return nil, errors.New("service: unauthorized to accept this nikkah match. Only the receiver of the match invite can accept it.")
	}

	if nikkahMatch.Status != entity.MatchStatusInitiated {
		return nil, fmt.Errorf("service: nikkah match with ID %s cannot be accepted as its current status is %s", matchID.String(), nikkahMatch.Status)
	}

	nikkahMatch.Status = entity.MatchStatusAccepted
	nikkahMatch.UpdatedAt = time.Now()

	updatedMatch, err := s.RepoNikkah.UpdateMatch(ctx, nikkahMatch)
	if err != nil {
		return nil, fmt.Errorf("service: failed to update nikkah match status to ACCEPTED for ID %s: %w", matchID.String(), err)
	}

	return updatedMatch, nil
}
