package repository

import (
	"context"
	"github.com/google/uuid"
	"github.com/mnadev/limestone/internal/application/domain/entity"
)

type NikkahRepository interface {
	CreateProfile(ctx context.Context, profile *entity.NikkahProfile) (*entity.NikkahProfile, error)
	GetProfileByID(ctx context.Context, profileID uuid.UUID) (*entity.NikkahProfile, error)
	GetProfileByUserID(ctx context.Context, userID string) (*entity.NikkahProfile, error)
	UpdateProfile(ctx context.Context, profile *entity.NikkahProfile) (*entity.NikkahProfile, error)
	ListProfiles(ctx context.Context, params *entity.NikkahProfileQueryParams) ([]*entity.NikkahProfile, int64, error)

	CreateLike(ctx context.Context, like *entity.NikkahLike) (*entity.NikkahLike, error)
	GetLikeByID(ctx context.Context, likeID uuid.UUID) (*entity.NikkahLike, error)
	GetLikeByLikerAndLikedProfileID(ctx context.Context, likerUserID string, likedProfileID string) (*entity.NikkahLike, error)
	UpdateLike(ctx context.Context, like *entity.NikkahLike) (*entity.NikkahLike, error)
	DeleteLike(ctx context.Context, likeID uuid.UUID) error
	ListLikesByLikerUserID(ctx context.Context, likerUserID string, params *entity.NikkahLikeQueryParams) ([]*entity.NikkahLike, int32, error)
	ListLikesForLikedProfileID(ctx context.Context, likedProfileID uuid.UUID, params *entity.NikkahLikeQueryParams) ([]*entity.NikkahLike, int32, error)

	CreateMatch(ctx context.Context, match *entity.NikkahMatch) (*entity.NikkahMatch, error)
	GetMatchByID(ctx context.Context, matchID uuid.UUID) (*entity.NikkahMatch, error)
	UpdateMatch(ctx context.Context, match *entity.NikkahMatch) (*entity.NikkahMatch, error)
}
