package services

import (
	"context"
	"errors"
	"fmt"
	"github.com/mnadev/limestone/internal/application/helper"
	"gorm.io/gorm"
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
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, helper.ErrNotFound
		}
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

	existingLike, err := s.RepoNikkah.GetLikeByLikerAndLikedProfileID(ctx, likerProfile.ID.String(), likedProfileID.String())

	if err == nil {
		switch existingLike.Status {
		case entity.LikeStatusInitiated:
			return nil, errors.New("service: you have already sent a like to this profile and it's pending response")
		case entity.LikeStatusCompleted:
			return nil, errors.New("service: you have already matched with this profile")
		default:
			return nil, fmt.Errorf("service: existing like has an unhandled status: %s", existingLike.Status)
		}
	} else if !errors.Is(err, helper.ErrNotFound) {
		return nil, fmt.Errorf("service: failed to check for existing like: %w", err)
	}

	newLike := &entity.NikkahLike{
		ID:             uuid.New(),
		LikerProfileID: likerProfile.ID.String(),
		LikedProfileID: likedProfileID.String(),
		Status:         entity.LikeStatusInitiated,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	createdLike, err := s.RepoNikkah.CreateLike(ctx, newLike)
	if err != nil {
		return nil, fmt.Errorf("service: failed to create new like: %w", err)
	}

	return createdLike, nil
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
			return nil, nil, fmt.Errorf("service: nikkah like with ID %s not found. Please check the ID or if it has been deleted: %w", likeID.String(), err)
		}
		return nil, nil, fmt.Errorf("service: failed to retrieve nikkah like %s: %w", likeID.String(), err)
	}

	requestingUserProfile, err := s.RepoNikkah.GetProfileByUserID(ctx, requestingUserID)
	if err != nil {
		if errors.Is(err, helper.ErrNotFound) {
			return nil, nil, errors.New("service: requesting user's profile not found. Cannot complete like without a valid profile.")
		}
		return nil, nil, fmt.Errorf("service: failed to get requesting user's profile: %w", err)
	}

	if initiatingLike.LikedProfileID != requestingUserProfile.ID.String() {
		return nil, nil, errors.New("service: unauthorized to complete this nikkah like. Only the liked profile (receiver of the like) can complete it.")
	}

	if initiatingLike.Status != entity.LikeStatusInitiated {
		return nil, nil, fmt.Errorf("service: nikkah like with ID %s cannot be completed as its current status is %s. Only likes with 'Initiated' status can be completed.", likeID.String(), initiatingLike.Status)
	}

	reverseLike, err := s.RepoNikkah.GetLikeByLikerAndLikedProfileID(
		ctx,
		initiatingLike.LikedProfileID,
		initiatingLike.LikerProfileID,
	)
	if err != nil {
		if errors.Is(err, helper.ErrNotFound) {
			return nil, nil, errors.New("service: no pending mutual like found from the other profile to complete this match. Ensure the other profile has sent a like that is still pending.")
		}
		return nil, nil, fmt.Errorf("service: failed to check for reverse like: %w", err)
	}

	if reverseLike.Status != entity.LikeStatusInitiated {
		return nil, nil, fmt.Errorf("service: reverse nikkah like from %s to %s cannot be completed as its current status is %s. It must be 'Initiated' to form a match.",
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

	profileA, err := uuid.Parse(initiatingLike.LikerProfileID)
	if err != nil {
		return nil, nil, fmt.Errorf("service: failed to parse LikerProfileID ('%s') to UUID for match creation: %w", initiatingLike.LikerProfileID, err)
	}

	profileB, err := uuid.Parse(initiatingLike.LikedProfileID)
	if err != nil {
		return nil, nil, fmt.Errorf("service: failed to parse LikedProfileID ('%s') to UUID for match creation: %w", initiatingLike.LikedProfileID, err)
	}

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
		return updatedInitiatingLike, nil, fmt.Errorf("service: failed to create nikkah match after successful like completion: %w", err)
	}

	likerProfileData, getLikerProfileErr := s.RepoNikkah.GetProfileByID(ctx, profileA)
	likedProfileData, getLikedProfileErr := s.RepoNikkah.GetProfileByID(ctx, profileB)

	if getLikerProfileErr == nil && getLikedProfileErr == nil {
		fmt.Printf("Notifikasi Match: Selamat! %s dan %s sekarang cocok!\n", likerProfileData.Name, likedProfileData.Name)
	} else {
		fmt.Printf("Notifikasi Match: Selamat! Dua profil telah cocok. (Tidak dapat mengambil nama profil)\n")
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

	if nikkahMatch.ReceiverProfileID != requestingUserProfile.ID {
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

func (s *NikkahService) GetNikkahMatch(ctx context.Context, matchID uuid.UUID, requestingUserID string) (*entity.NikkahMatch, error) {
	if matchID == uuid.Nil {
		return nil, errors.New("service: match ID cannot be empty")
	}

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
			return nil, errors.New("service: requesting user's profile not found. Cannot retrieve match.")
		}
		return nil, fmt.Errorf("service: failed to get requesting user's profile: %w", err)
	}

	requestingProfileID := requestingUserProfile.ID
	if nikkahMatch.InitiatorProfileID != requestingProfileID && nikkahMatch.ReceiverProfileID != requestingProfileID {
		return nil, errors.New("service: unauthorized to view this nikkah match. You are not a participant in this match.")
	}

	return nikkahMatch, nil
}

func (s *NikkahService) RejectNikkahMatchInvite(ctx context.Context, matchID uuid.UUID, requestingUserID string) (*entity.NikkahMatch, error) {
	if matchID == uuid.Nil {
		return nil, errors.New("service: match ID cannot be empty for rejection")
	}

	nikkahMatch, err := s.RepoNikkah.GetMatchByID(ctx, matchID)
	if err != nil {
		if errors.Is(err, helper.ErrNotFound) {
			return nil, fmt.Errorf("service: nikkah match with ID %s not found: %w", matchID.String(), err)
		}
		return nil, fmt.Errorf("service: failed to retrieve nikkah match %s for rejection: %w", matchID.String(), err)
	}

	requestingUserProfile, err := s.RepoNikkah.GetProfileByUserID(ctx, requestingUserID)
	if err != nil {
		if errors.Is(err, helper.ErrNotFound) {
			return nil, errors.New("service: requesting user's profile not found. Cannot reject match.")
		}
		return nil, fmt.Errorf("service: failed to get requesting user's profile for rejection: %w", err)
	}

	requestingProfileID := requestingUserProfile.ID
	if nikkahMatch.ReceiverProfileID != requestingProfileID {
		return nil, errors.New("service: unauthorized to reject this nikkah match. Only the receiver of the match invite can reject it.")
	}

	if nikkahMatch.Status != entity.MatchStatusInitiated {
		return nil, fmt.Errorf("service: nikkah match with ID %s cannot be rejected as its current status is %s. Only 'Initiated' matches can be rejected.", matchID.String(), nikkahMatch.Status)
	}

	nikkahMatch.Status = entity.MatchStatusRejected
	nikkahMatch.UpdatedAt = time.Now()

	updatedMatch, err := s.RepoNikkah.UpdateMatch(ctx, nikkahMatch)
	if err != nil {
		return nil, fmt.Errorf("service: failed to update nikkah match status to ENDED for ID %s: %w", matchID.String(), err)
	}

	return updatedMatch, nil
}

func (s *NikkahService) EndNikkahMatch(ctx context.Context, matchID uuid.UUID, requestingUserID string) (*entity.NikkahMatch, error) {
	if matchID == uuid.Nil {
		return nil, errors.New("service: match ID cannot be empty for ending")
	}
	nikkahMatch, err := s.RepoNikkah.GetMatchByID(ctx, matchID)
	if err != nil {
		if errors.Is(err, helper.ErrNotFound) {
			return nil, fmt.Errorf("service: nikkah match with ID %s not found: %w", matchID.String(), err)
		}
		return nil, fmt.Errorf("service: failed to retrieve nikkah match %s for ending: %w", matchID.String(), err)
	}

	requestingUserProfile, err := s.RepoNikkah.GetProfileByUserID(ctx, requestingUserID)
	if err != nil {
		if errors.Is(err, helper.ErrNotFound) {
			return nil, errors.New("service: requesting user's profile not found. Cannot end match.")
		}
		return nil, fmt.Errorf("service: failed to get requesting user's profile for ending: %w", err)
	}

	initiatorProfileID := nikkahMatch.InitiatorProfileID
	receiverProfileID := nikkahMatch.ReceiverProfileID
	requestingProfileID := requestingUserProfile.ID
	if requestingProfileID != initiatorProfileID && requestingProfileID != receiverProfileID {
		return nil, errors.New("service: unauthorized to end this nikkah match. Only participants can end it.")
	}

	if nikkahMatch.Status == entity.MatchStatusEnded {
		return nil, fmt.Errorf("service: nikkah match with ID %s is already ended.", matchID.String())
	}

	nikkahMatch.Status = entity.MatchStatusEnded
	nikkahMatch.UpdatedAt = time.Now()

	updatedMatch, err := s.RepoNikkah.UpdateMatch(ctx, nikkahMatch)
	if err != nil {
		return nil, fmt.Errorf("service: failed to update nikkah match status to ENDED for ID %s: %w", matchID.String(), err)
	}
	return updatedMatch, nil
}
