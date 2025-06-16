// Package storage Copyright (c) 2024 Coding-af Limestone Dev
// Licensed under the MIT License.
// file COPYING or http://www.opensource.org/licenses/mit-license.php
package storage

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/mnadev/limestone/internal/application/domain/entity"
	"gorm.io/gorm"
)

type GormRevertRepository struct {
	db *gorm.DB
}

// NewGormRevertRepository /**
/**
 * NewGormRevertRepository creates a new repository instance with the provided DB connection.
 *
 * @method
 * @name storage#NewGormRevertRepository
 * @param {gorm.DB} db - the GORM database connection
 * @returns {GormRevertRepository} a new repository instance
 */
func NewGormRevertRepository(db *gorm.DB) *GormRevertRepository {
	return &GormRevertRepository{db: db}
}

// CreateRevertProfile /**
/**
 * CreateRevertProfile stores a new RevertProfile entity in the database.
 *
 * @method
 * @name storage#CreateRevertProfile
 * @param {context.Context} ctx
 * @param {RevertProfile} profile
 * @returns {RevertProfile}
 * @throws error on database failure
 */
func (r *GormRevertRepository) CreateRevertProfile(ctx context.Context, profile *entity.RevertProfile) (*entity.RevertProfile, error) {
	if err := r.db.WithContext(ctx).Create(profile).Error; err != nil {
		return nil, fmt.Errorf("failed to create revert profile: %w", err)
	}
	return profile, nil
}

// GetRevertProfileByID /**
/**
 * GetRevertProfileByID fetches a RevertProfile by its unique UUID.
 *
 * @method
 * @name storage#GetRevertProfileByID
 * @param {context.Context} ctx
 * @param {uuid.UUID} revertProfileID
 * @returns {RevertProfile}
 * @throws error if not found or on query failure
 */
func (r *GormRevertRepository) GetRevertProfileByID(ctx context.Context, revertProfileID uuid.UUID) (*entity.RevertProfile, error) {
	var revertProfile entity.RevertProfile
	if err := r.db.WithContext(ctx).First(&revertProfile, "id = ?", revertProfileID.String()).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("Revert profile with ID %s not found: %w", revertProfileID.String(), err)
		}
		return nil, fmt.Errorf("failed to get revert profile by ID %s: %w", revertProfileID.String(), err)
	}
	return &revertProfile, nil
}

// GetRevertProfileByUserID /**
/**
 * GetRevertProfileByUserID fetches a RevertProfile by associated user ID.
 *
 * @method
 * @name storage#GetRevertProfileByUserID
 * @param {context.Context} ctx
 * @param {string} userID
 * @returns {RevertProfile}
 * @throws error if not found or on query failure
 */
func (r *GormRevertRepository) GetRevertProfileByUserID(ctx context.Context, userID string) (*entity.RevertProfile, error) {
	var profile entity.RevertProfile
	if err := r.db.WithContext(ctx).First(&profile, "user_id = ?", userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("revert profile for user ID %s not found: %w", userID, err)
		}
		return nil, fmt.Errorf("failed to get revert profile by user ID %s: %w", userID, err)
	}
	return &profile, nil

}

// UpdateRevertProfile /**
/**
 * UpdateRevertProfile updates an existing RevertProfile in the database.
 *
 * @method
 * @name storage#UpdateRevertProfile
 * @param {context.Context} ctx
 * @param {RevertProfile} profile
 * @returns {RevertProfile}
 * @throws error on update failure
 */
func (r *GormRevertRepository) UpdateRevertProfile(ctx context.Context, profile *entity.RevertProfile) (*entity.RevertProfile, error) {
	if err := r.db.WithContext(ctx).Save(profile).Error; err != nil {
		return nil, fmt.Errorf("failed to update revert profile: %w", err)
	}
	return profile, nil
}

// ListRevertProfiles /**
/**
 * ListRevertProfiles retrieves a paginated list of profiles based on filters.
 *
 * @method
 * @name storage#ListRevertProfiles
 * @param {context.Context} ctx
 * @param {RevertProfileQueryParams} params - includes name filter, pagination data
 * @returns {[RevertProfile[], int64]} list of profiles and total count
 * @throws error on query failure
 */
func (r *GormRevertRepository) ListRevertProfiles(ctx context.Context, params *entity.RevertProfileQueryParams) ([]*entity.RevertProfile, int64, error) {
	db := r.db.WithContext(ctx).Model(&entity.RevertProfile{})

	if params.Name != "" {
		db = db.Where("name ILIKE ?", "%"+params.Name+"%")
	}

	if params.Gender != "" {
		db = db.Where("gender = ?", params.Gender)
	}

	var totalCount int64
	if err := db.Count(&totalCount).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count revert profiles: %w", err)
	}

	var revertProfiles []*entity.RevertProfile
	pageSize := params.Limit
	if pageSize <= 0 {
		pageSize = 10
	}

	offset := int((params.Page - 1) * pageSize)
	if offset < 0 {
		offset = 0
	}
	db = db.Offset(offset).Limit(int(pageSize))

	db = db.Order("created_at DESC, id ASC")

	result := db.Find(&revertProfiles)
	if result.Error != nil {
		return nil, 0, fmt.Errorf("failed to retrieve nikkah profiles: %w", result.Error)
	}

	return revertProfiles, totalCount, nil
}

// CreateRevertMatch /**
/**
 * CreateRevertMatch stores a new RevertMatch entity in the database.
 *
 * @method
 * @name storage#CreateRevertMatch
 * @param {context.Context} ctx
 * @param {RevertMatch} match
 * @returns {RevertMatch}
 * @throws error on insert failure
 */
func (r *GormRevertRepository) CreateRevertMatch(ctx context.Context, match *entity.RevertMatch) (*entity.RevertMatch, error) {
	if err := r.db.WithContext(ctx).Create(match).Error; err != nil {
		return nil, fmt.Errorf("failed to create revert match: %w", err)
	}
	return match, nil
}

// GetRevertMatchByID /**
/**
 * GetRevertMatchByID retrieves a RevertMatch by its UUID.
 *
 * @method
 * @name storage#GetRevertMatchByID
 * @param {context.Context} ctx
 * @param {uuid.UUID} id
 * @returns {RevertMatch|null}
 * @throws error on query failure
 */
func (r *GormRevertRepository) GetRevertMatchByID(ctx context.Context, id uuid.UUID) (*entity.RevertMatch, error) {
	var match entity.RevertMatch
	if err := r.db.WithContext(ctx).First(&match, "id = ?", id.String()).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get revert match by ID: %w", err)
	}
	return &match, nil
}

// FindLatestActiveMatch /**
/**
 * FindLatestActiveMatch finds the most recent active match between two profiles.
 *
 * @method
 * @name storage#FindLatestActiveMatch
 * @param {context.Context} ctx
 * @param {string} initiatorProfileID
 * @param {string} receiverProfileID
 * @returns {RevertMatch|null}
 * @throws error on query failure
 */
func (r *GormRevertRepository) FindLatestActiveMatch(ctx context.Context, initiatorProfileID, receiverProfileID string) (*entity.RevertMatch, error) {
	var match entity.RevertMatch
	err := r.db.WithContext(ctx).
		Where("(initiator_profile_id = ? AND receiver_profile_id = ?) OR (initiator_profile_id = ? AND receiver_profile_id = ?)",
			initiatorProfileID, receiverProfileID, receiverProfileID, initiatorProfileID).
		Where("status IN (?)", []entity.MatchStatus{entity.MatchStatusInitiated, entity.MatchStatusAccepted}).
		Order("created_at DESC").
		First(&match).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find latest active match: %w", err)
	}
	return &match, nil
}

// UpdateRevertMatch /**
/**
 * UpdateRevertMatch updates an existing RevertMatch in the database.
 *
 * @method
 * @name storage#UpdateRevertMatch
 * @param {context.Context} ctx
 * @param {RevertMatch} match
 * @returns {RevertMatch}
 * @throws error on update failure
 */

func (r *GormRevertRepository) UpdateRevertMatch(ctx context.Context, match *entity.RevertMatch) (*entity.RevertMatch, error) {
	if err := r.db.WithContext(ctx).Save(match).Error; err != nil {
		return nil, fmt.Errorf("failed to update revert match: %w", err)
	}
	return match, nil
}
