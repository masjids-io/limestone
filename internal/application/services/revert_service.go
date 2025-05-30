// Package services Copyright (c) 2024 Coding-AF Limestone Dev
// Licensed under the MIT License.
// file COPYING or http://www.opensource.org/licenses/mit-license.php
package services

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/mnadev/limestone/internal/application/domain/entity"
	"github.com/mnadev/limestone/internal/application/helper"
	"github.com/mnadev/limestone/internal/application/repository"
	"gorm.io/gorm"
	"math"
	"time"
)

type RevertService struct {
	RepoRevert repository.RevertRepository
}

func NewRevertService(nikahRepo repository.RevertRepository) *RevertService {
	return &RevertService{
		RepoRevert: nikahRepo,
	}
}

// CreateRevertProfile /*
/*
@Method CreateRevertProfile
@name Create a new revert profile
@description This method validates the input, checks for an existing profile by UserID, and creates a new revert profile if none exists.
             It auto-generates an ID if not provided and sets creation and update timestamps.
@param ctx context.Context - the request context used for cancellation and deadlines
@param profile *entity.RevertProfile - the profile data to be created, must include UserID, Name, and Gender
@return *entity.RevertProfile - the newly created profile object
@return error - error if the input is invalid, a profile already exists, or creation fails
*/
func (s *RevertService) CreateRevertProfile(ctx context.Context, profile *entity.RevertProfile) (*entity.RevertProfile, error) {
	if profile.UserID == "" || profile.Name == "" || profile.Gender == "" {
		return nil, helper.ErrInvalidRevertProfileData
	}
	_, err := s.RepoRevert.GetRevertProfileByUserID(ctx, profile.UserID)
	if err == nil {
		return nil, helper.ErrRevertProfileAlreadyExists
	}
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("failed to check for existing profile: %w", err)
	}

	if profile.ID == uuid.Nil {
		profile.ID = uuid.New()
	}
	profile.CreatedAt = time.Now()
	profile.UpdatedAt = time.Now()

	createdProfile, err := s.RepoRevert.CreateRevertProfile(ctx, profile)
	if err != nil {
		return nil, fmt.Errorf("failed to create revert profile: %w", err)
	}
	return createdProfile, nil
}

// GetRevertProfileByID /*
/*
@Method GetRevertProfileByID
@name Retrieve revert profile by ID
@description This method retrieves a revert profile using its unique UUID. It checks for empty ID input and handles not found or query errors.
@param ctx context.Context - the request context used for cancellation and deadlines
@param revertProfileId uuid.UUID - the unique identifier of the revert profile to be retrieved
@return *entity.RevertProfile - the revert profile found by ID
@return error - error if ID is empty, profile not found, or if a query error occurs
*/
func (s *RevertService) GetRevertProfileByID(ctx context.Context, revertProfileId uuid.UUID) (*entity.RevertProfile, error) {
	if revertProfileId == uuid.Nil {
		return nil, errors.New("revert profile ID cannot be empty")
	}
	profile, err := s.RepoRevert.GetRevertProfileByID(ctx, revertProfileId)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, helper.ErrNotFound
		}
		return nil, fmt.Errorf("service: failed to get revert profile by ID: %w", err)
	}
	return profile, nil
}

// GetRevertProfileByUserID /*
/*
@Method GetRevertProfileByUserID
@name Retrieve revert profile by user ID
@description This method retrieves a revert profile based on the given user ID.
             It validates the input and handles "not found" and other database errors accordingly.
@param ctx context.Context - the request context used for cancellation, deadlines, etc.
@param userID string - the unique user identifier to look up the revert profile
@return *entity.RevertProfile - the revert profile associated with the given user ID
@return error - error if user ID is empty, profile not found, or any database query fails
*/
func (s *RevertService) GetRevertProfileByUserID(ctx context.Context, userID string) (*entity.RevertProfile, error) {
	if userID == "" {
		return nil, errors.New("user ID cannot be empty")
	}
	profile, err := s.RepoRevert.GetRevertProfileByUserID(ctx, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, helper.ErrNotFound
		}
		return nil, fmt.Errorf("service: failed to get revert profile by user ID: %w", err)
	}
	return profile, nil
}

// UpdateRevertProfile /*
/*
@Method UpdateRevertProfile
@name Update an existing revert profile
@description This method updates the details of an existing revert profile.
             It first verifies the existence of the profile by its ID,
             then selectively updates the fields if new values are provided,
             and finally persists the changes to the database.
@param ctx context.Context - the context for request-scoped values, cancellation, and deadlines
@param profile *entity.RevertProfile - the revert profile containing updated fields and a valid ID
@return *entity.RevertProfile - the updated revert profile after persistence
@return error - error if the profile ID is missing, not found, or update operation fails
*/
func (s *RevertService) UpdateRevertProfile(ctx context.Context, profile *entity.RevertProfile) (*entity.RevertProfile, error) {
	fmt.Println(profile)
	if profile.ID == uuid.Nil {
		return nil, errors.New("profile ID is required for update")
	}

	existingProfile, err := s.RepoRevert.GetRevertProfileByID(ctx, profile.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, helper.ErrNotFound
		}
		return nil, fmt.Errorf("failed to retrieve existing profile for update: %w", err)
	}

	if profile.Name != "" {
		existingProfile.Name = profile.Name
	}
	if profile.Gender != "" && profile.Gender != entity.RevertProfileGender_UNSPECIFIED {
		existingProfile.Gender = profile.Gender
	}
	if profile.BirthDate.Year != 0 || profile.BirthDate.Month != 0 || profile.BirthDate.Day != 0 {
		existingProfile.BirthDate = profile.BirthDate
	}

	existingProfile.UpdatedAt = time.Now()

	updatedProfile, err := s.RepoRevert.UpdateRevertProfile(ctx, existingProfile)
	if err != nil {
		return nil, fmt.Errorf("failed to update revert profile: %w", err)
	}
	return updatedProfile, nil
}

// ListRevertProfiles /*
/*
@Method ListRevertProfiles
@name List revert profiles with pagination and filtering
@description This method retrieves a list of revert profiles based on the provided query parameters.
             It supports pagination, limit control, and calculates total pages based on the result count.
             The method adjusts invalid or missing pagination parameters to ensure safe query execution.
@param ctx context.Context - the context for request-scoped values, cancellation, and deadlines
@param params *entity.RevertProfileQueryParams - the parameters for pagination and filtering (limit, page, start, etc.)
@return *entity.RevertProfileQueryResult - a result structure containing the list of revert profiles,
        total count, current page, and total pages
@return error - error if the profile retrieval operation fails
*/
func (s *RevertService) ListRevertProfiles(ctx context.Context, params *entity.RevertProfileQueryParams) (*entity.RevertProfileQueryResult, error) {
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
	revertProfiles, totalCount, err := s.RepoRevert.ListRevertProfiles(ctx, params)
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

	return &entity.RevertProfileQueryResult{
		RevertProfile: revertProfiles,
		TotalCount:    int64(totalCount),
		CurrentPage:   currentPage,
		TotalPages:    totalPages,
	}, nil
}

// CreateRevertMatchInvite /*
/*
@Method CreateRevertMatchInvite
@name Create revert match invite
@description This method creates a match invitation between two revert profiles.
             It checks for valid initiator and receiver IDs, and prevents self-invitation.
             Upon validation, it initializes a new RevertMatch entity with `INITIATED` status
             and persists it through the repository.
@param ctx context.Context - the context for request-scoped values, cancellation, and deadlines
@param initiatorProfileID string - the ID of the profile initiating the match
@param receiverProfileID string - the ID of the profile receiving the match invitation
@return *entity.RevertMatch - the created revert match entity
@return error - error if validation fails or the repository fails to persist the match
*/
func (s *RevertService) CreateRevertMatchInvite(ctx context.Context, initiatorProfileID, receiverProfileID string) (*entity.RevertMatch, error) {
	if initiatorProfileID == "" || receiverProfileID == "" {
		return nil, helper.ErrInvalidMatchData
	}

	if initiatorProfileID == receiverProfileID {
		return nil, helper.ErrSelfInvitation
	}

	match := &entity.RevertMatch{
		ID:                 uuid.New(),
		InitiatorProfileID: initiatorProfileID,
		ReceiverProfileID:  receiverProfileID,
		Status:             entity.RevertMatchStatus_INITIATED,
		CreatedAt:          time.Now(),
		UpdatedAt:          time.Now(),
	}

	createdMatch, err := s.RepoRevert.CreateRevertMatch(ctx, match)
	if err != nil {
		return nil, fmt.Errorf("failed to create revert match invite: %w", err)
	}
	return createdMatch, nil
}

// GetRevertMatch /*
/*
@Method GetRevertMatch
@name Get revert match by ID
@description Retrieves a specific revert match entity using its UUID string.
             It validates the match ID, parses it to UUID format, and queries the repository.
             If the match is not found or an error occurs during retrieval, an appropriate error is returned.
@param ctx context.Context - the context for controlling request lifetime
@param matchID string - the string representation of the match UUID
@return *entity.RevertMatch - the retrieved revert match entity if found
@return error - error if validation fails, ID is malformed, or the match is not found
*/
func (s *RevertService) GetRevertMatch(ctx context.Context, matchID string) (*entity.RevertMatch, error) {
	if matchID == "" {
		return nil, helper.ErrInvalidMatchID
	}

	parsedID, err := uuid.Parse(matchID)
	if err != nil {
		return nil, fmt.Errorf("invalid match ID format: %w", err)
	}

	match, err := s.RepoRevert.GetRevertMatchByID(ctx, parsedID)
	if err != nil {
		return nil, fmt.Errorf("failed to get revert match: %w", err)
	}
	if match == nil {
		return nil, helper.ErrMatchNotFound
	}
	return match, nil
}

// AcceptRevertMatchInvite /*
/*
@Method AcceptRevertMatchInvite
@name Accept a revert match invitation
@description Validates and accepts a revert match invitation by its ID.
             It ensures the match ID is valid, the match exists, and is currently in the INITIATED state.
             Upon acceptance, the match status is updated to ACCEPTED and saved.
@param ctx context.Context - context for managing request lifecycle
@param matchID string - string representation of the revert match UUID to accept
@return *entity.RevertMatch - the updated revert match entity after acceptance
@return error - error if validation fails, match is not found, or update fails
*/
func (s *RevertService) AcceptRevertMatchInvite(ctx context.Context, matchID string) (*entity.RevertMatch, error) {
	if matchID == "" {
		return nil, helper.ErrInvalidMatchID
	}

	parsedID, err := uuid.Parse(matchID)
	if err != nil {
		return nil, fmt.Errorf("invalid match ID format: %w", err)
	}

	match, err := s.RepoRevert.GetRevertMatchByID(ctx, parsedID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve match for acceptance: %w", err)
	}
	if match == nil {
		return nil, helper.ErrMatchNotFound
	}

	if match.Status != entity.RevertMatchStatus_INITIATED {
		return nil, helper.ErrMatchStatusInvalid
	}

	match.Status = entity.RevertMatchStatus_ACCEPTED
	match.UpdatedAt = time.Now()

	updatedMatch, err := s.RepoRevert.UpdateRevertMatch(ctx, match)
	if err != nil {
		return nil, fmt.Errorf("failed to accept revert match invite: %w", err)
	}
	return updatedMatch, nil
}

// RejectRevertMatchInvite /*
/*
@Method RejectRevertMatchInvite
@name Reject a revert match invitation
@description Validates and rejects a revert match invitation by its ID.
             It checks if the match ID is valid, the match exists, and is in the INITIATED state.
             If valid, the match status is updated to REJECTED and persisted.
@param ctx context.Context - context for managing request lifecycle
@param matchID string - string representation of the revert match UUID to reject
@return *entity.RevertMatch - the updated revert match entity after rejection
@return error - error if validation fails, match is not found, or update fails
*/
func (s *RevertService) RejectRevertMatchInvite(ctx context.Context, matchID string) (*entity.RevertMatch, error) {
	if matchID == "" {
		return nil, helper.ErrInvalidMatchID
	}

	parsedID, err := uuid.Parse(matchID)
	if err != nil {
		return nil, fmt.Errorf("invalid match ID format: %w", err)
	}

	match, err := s.RepoRevert.GetRevertMatchByID(ctx, parsedID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve match for rejection: %w", err)
	}
	if match == nil {
		return nil, helper.ErrMatchNotFound
	}

	if match.Status != entity.RevertMatchStatus_INITIATED {
		return nil, helper.ErrMatchStatusInvalid
	}

	match.Status = entity.RevertMatchStatus_REJECTED
	match.UpdatedAt = time.Now()

	updatedMatch, err := s.RepoRevert.UpdateRevertMatch(ctx, match)
	if err != nil {
		return nil, fmt.Errorf("failed to reject revert match invite: %w", err)
	}
	return updatedMatch, nil
}

// EndRevertMatch /*
/*
@Method EndRevertMatch
@name End a revert match
@description Validates and ends a revert match by its ID.
             Ensures the match ID is valid, the match exists, and the match status is ACCEPTED.
             Updates the match status to ENDED and saves the changes.
@param ctx context.Context - context for managing request lifecycle
@param matchID string - string representation of the revert match UUID to end
@return *entity.RevertMatch - the updated revert match entity after ending
@return error - error if validation fails, match is not found, or update fails
*/
func (s *RevertService) EndRevertMatch(ctx context.Context, matchID string) (*entity.RevertMatch, error) {
	if matchID == "" {
		return nil, helper.ErrInvalidMatchID
	}

	parsedID, err := uuid.Parse(matchID)
	if err != nil {
		return nil, fmt.Errorf("invalid match ID format: %w", err)
	}

	match, err := s.RepoRevert.GetRevertMatchByID(ctx, parsedID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve match for ending: %w", err)
	}
	if match == nil {
		return nil, helper.ErrMatchNotFound
	}

	if match.Status != entity.RevertMatchStatus_ACCEPTED {
		return nil, helper.ErrMatchStatusInvalid
	}

	match.Status = entity.RevertMatchStatus_ENDED
	match.UpdatedAt = time.Now()

	updatedMatch, err := s.RepoRevert.UpdateRevertMatch(ctx, match)
	if err != nil {
		return nil, fmt.Errorf("failed to end revert match: %w", err)
	}
	return updatedMatch, nil
}
