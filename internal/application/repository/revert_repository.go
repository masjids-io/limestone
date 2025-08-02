// Package repository Package storage Copyright (c) 2024 Coding-AF Limestone Dev
// Licensed under the MIT License.
// file COPYING or http://www.opensource.org/licenses/mit-license.php
package repository

import (
	"context"
	"github.com/google/uuid"
	"github.com/mnadev/limestone/internal/application/domain/entity"
)

type RevertRepository interface {
	CreateRevertProfile(ctx context.Context, profile *entity.RevertProfile) (*entity.RevertProfile, error)
	GetRevertProfileByID(ctx context.Context, id uuid.UUID) (*entity.RevertProfile, error)
	GetRevertProfileByUserID(ctx context.Context, userID string) (*entity.RevertProfile, error)
	UpdateRevertProfile(ctx context.Context, profile *entity.RevertProfile) (*entity.RevertProfile, error)
	ListRevertProfiles(ctx context.Context, params *entity.RevertProfileQueryParams) ([]*entity.RevertProfile, int64, error)
	CreateRevertMatch(ctx context.Context, match *entity.RevertMatch) (*entity.RevertMatch, error)
	GetRevertMatchByID(ctx context.Context, id uuid.UUID) (*entity.RevertMatch, error)
	FindLatestActiveMatch(ctx context.Context, initiatorProfileID, receiverProfileID string) (*entity.RevertMatch, error)
	UpdateRevertMatch(ctx context.Context, match *entity.RevertMatch) (*entity.RevertMatch, error)
}
