package storage

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/mnadev/limestone/internal/application/domain/entity"
	"github.com/mnadev/limestone/internal/application/helper"
	"gorm.io/gorm"
)

type GormNikkahRepository struct {
	db *gorm.DB
}

func NewGormNikkahRepository(db *gorm.DB) *GormNikkahRepository {
	return &GormNikkahRepository{db: db}
}

func (r *GormNikkahRepository) CreateProfile(ctx context.Context, profile *entity.NikkahProfile) (*entity.NikkahProfile, error) {
	if err := r.db.WithContext(ctx).Create(profile).Error; err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code.Name() == "unique_violation" {
			return nil, fmt.Errorf("%w: user ID already exists", helper.ErrAlreadyExists)
		}
		return nil, fmt.Errorf("failed to create nikkah profile: %w", err)
	}
	return profile, nil
}

func (r *GormNikkahRepository) GetProfileByID(ctx context.Context, profileID uuid.UUID) (*entity.NikkahProfile, error) {
	var profile entity.NikkahProfile
	if err := r.db.WithContext(ctx).First(&profile, "id = ?", profileID.String()).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("nikkah profile with ID %s not found: %w", profileID.String(), err)
		}
		return nil, fmt.Errorf("failed to get nikkah profile by ID %s: %w", profileID.String(), err)
	}
	return &profile, nil
}

func (r *GormNikkahRepository) GetProfileByUserID(ctx context.Context, userID string) (*entity.NikkahProfile, error) {
	var profile entity.NikkahProfile
	if err := r.db.WithContext(ctx).First(&profile, "user_id = ?", userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("nikkah profile for user ID %s not found: %w", userID, err)
		}
		return nil, fmt.Errorf("failed to get nikkah profile by user ID %s: %w", userID, err)
	}
	return &profile, nil
}

func (r *GormNikkahRepository) UpdateProfile(ctx context.Context, profile *entity.NikkahProfile) (*entity.NikkahProfile, error) {
	_, err := r.GetProfileByID(ctx, profile.ID)
	if err != nil {
		return nil, fmt.Errorf("cannot update profile, target profile not found: %w", err)
	}

	if err := r.db.WithContext(ctx).Save(profile).Error; err != nil {
		return nil, fmt.Errorf("failed to update nikkah profile with ID %s: %w", profile.ID.String(), err)
	}

	return profile, nil
}

func (r *GormNikkahRepository) ListProfiles(ctx context.Context, params *entity.NikkahProfileQueryParams) ([]*entity.NikkahProfile, int64, error) {
	db := r.db.WithContext(ctx).Model(&entity.NikkahProfile{})

	fmt.Println(params)

	if params.Name != "" {
		db = db.Where("name ILIKE ?", "%"+params.Name+"%")
	}

	if params.Gender != "" {
		db = db.Where("gender = ?", params.Gender)
	}

	// Location filters
	if params.Location != nil {
		if params.Location.Country != "" {
			db = db.Where("location_country ILIKE ?", "%"+params.Location.Country+"%")
		}
		if params.Location.City != "" {
			db = db.Where("location_city ILIKE ?", "%"+params.Location.City+"%")
		}
		if params.Location.State != "" {
			db = db.Where("location_state ILIKE ?", "%"+params.Location.State+"%")
		}
		if params.Location.ZipCode != "" {
			db = db.Where("location_zip_code = ?", params.Location.ZipCode)
		}
		// Note: Latitude and Longitude filtering would require more complex geospatial queries
		// For now, we'll skip them as they would need proper geospatial indexing
	}

	// Education filter
	if params.Education != nil {
		db = db.Where("education = ?", *params.Education)
	}

	// Occupation filter
	if params.Occupation != "" {
		db = db.Where("occupation ILIKE ?", "%"+params.Occupation+"%")
	}

	// Height filter
	if params.Height != nil {
		db = db.Where("height_cm = ?", params.Height.Cm)
	}

	// Sect filter
	if params.Sect != nil {
		db = db.Where("sect = ?", *params.Sect)
	}

	// Hobbies filter - this is more complex as it's a JSONB array
	if len(params.Hobbies) > 0 {
		// For JSONB array contains, we need to check if the hobbies array contains any of the specified hobbies
		// This is a simplified approach - in production you might want more sophisticated matching
		for _, hobby := range params.Hobbies {
			db = db.Where("hobbies @> ?", fmt.Sprintf(`["%s"]`, hobby.String()))
		}
	}

	var totalCount int64
	if err := db.Count(&totalCount).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count nikkah profiles: %w", err)
	}

	var profiles []*entity.NikkahProfile
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

	result := db.Find(&profiles)
	if result.Error != nil {
		return nil, 0, fmt.Errorf("failed to retrieve nikkah profiles: %w", result.Error)
	}

	return profiles, totalCount, nil
}

func (r *GormNikkahRepository) CreateLike(ctx context.Context, like *entity.NikkahLike) (*entity.NikkahLike, error) {
	if err := r.db.WithContext(ctx).Create(like).Error; err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code.Name() == "unique_violation" {
			return nil, fmt.Errorf("%w: user ID already exists", helper.ErrAlreadyExists)
		}
		return nil, fmt.Errorf("failed to create nikkah like: %w", err)
	}
	return like, nil
}

func (r *GormNikkahRepository) GetLikeByID(ctx context.Context, likeID uuid.UUID) (*entity.NikkahLike, error) {
	var like entity.NikkahLike
	result := r.db.WithContext(ctx).First(&like, "id = ?", likeID.String())

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("repository: nikkah like not found: %w", helper.ErrNotFound)
		}
		return nil, fmt.Errorf("repository: failed to get nikkah like by ID: %w", result.Error)
	}
	return &like, nil
}

func (r *GormNikkahRepository) UpdateLike(ctx context.Context, like *entity.NikkahLike) (*entity.NikkahLike, error) {
	if err := r.db.WithContext(ctx).Save(like).Error; err != nil {
		return nil, fmt.Errorf("repository: failed to update like: %w", err)
	}
	return like, nil
}

func (r *GormNikkahRepository) CreateMatch(ctx context.Context, match *entity.NikkahMatch) (*entity.NikkahMatch, error) {
	if match.ID == uuid.Nil {
		match.ID = uuid.New()
	}

	result := r.db.WithContext(ctx).Create(match)
	if result.Error != nil {
		if errors.Is(result.Error, &pq.Error{}) {
			pqErr := result.Error.(*pq.Error)
			if pqErr.Code.Name() == "unique_violation" {
				return nil, fmt.Errorf("repository: match already exists between these profiles: %w", helper.ErrAlreadyExists)
			}
		}
		return nil, fmt.Errorf("repository: failed to create match: %w", result.Error)
	}

	return match, nil
}

func (r *GormNikkahRepository) GetLikeByLikerAndLikedProfileID(ctx context.Context, likerProfileID string, likedProfileID string) (*entity.NikkahLike, error) {
	var like entity.NikkahLike
	result := r.db.WithContext(ctx).
		Where("liker_profile_id = ?", likerProfileID).
		Where("liked_profile_id = ?", likedProfileID).
		First(&like)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("repository: nikkah like not found for liker %s and liked %s: %w", likerProfileID, likedProfileID, helper.ErrNotFound)
		}
		return nil, fmt.Errorf("repository: failed to get nikkah like by liker and liked profile IDs: %w", result.Error)
	}

	return &like, nil
}

func (r *GormNikkahRepository) DeleteLike(ctx context.Context, likeID uuid.UUID) error {
	//TODO implement me
	panic("implement me")
}

func (r *GormNikkahRepository) ListLikesByLikerUserID(ctx context.Context, likerUserID string, params *entity.NikkahLikeQueryParams) ([]*entity.NikkahLike, int32, error) {
	//TODO implement me
	panic("implement me")
}

func (r *GormNikkahRepository) ListLikesForLikedProfileID(ctx context.Context, likedProfileID uuid.UUID, params *entity.NikkahLikeQueryParams) ([]*entity.NikkahLike, int32, error) {
	//TODO implement me
	panic("implement me")
}

func (r *GormNikkahRepository) GetMatchByID(ctx context.Context, matchID uuid.UUID) (*entity.NikkahMatch, error) {
	var match entity.NikkahMatch
	result := r.db.WithContext(ctx).First(&match, "id = ?", matchID.String())

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("repository: nikkah match not found: %w", helper.ErrNotFound)
		}
		return nil, fmt.Errorf("repository: failed to get nikkah match by ID: %w", result.Error)
	}

	return &match, nil
}

func (r *GormNikkahRepository) UpdateMatch(ctx context.Context, match *entity.NikkahMatch) (*entity.NikkahMatch, error) {
	result := r.db.WithContext(ctx).Save(match)
	if result.Error != nil {
		return nil, fmt.Errorf("repository: failed to update match: %w", result.Error)
	}
	return match, nil
}
