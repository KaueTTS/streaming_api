package repository_sqlite_profile

import (
	"context"

	models "github.com/KaueTTS/streaming_api/src/models"
	"gorm.io/gorm"
)

type ProfileRepository struct {
	db *gorm.DB
}

func NewProfileRepository(db *gorm.DB) *ProfileRepository {
	return &ProfileRepository{
		db: db,
	}
}

func (r *ProfileRepository) FindAll(ctx context.Context) ([]models.Profile, error) {
	var profiles []models.Profile

	if err := r.db.WithContext(ctx).Order("name asc").Find(&profiles).Error; err != nil {
		return nil, err
	}

	return profiles, nil
}

func (r *ProfileRepository) FindByUserID(ctx context.Context, userID uint, page int, perPage int) ([]models.Profile, int64, error) {
	var profiles []models.Profile
	var total int64

	query := r.db.WithContext(ctx).Model(&models.Profile{}).Where("user_id = ?", userID)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * perPage

	if err := query.
		Order("name asc").
		Limit(perPage).
		Offset(offset).
		Find(&profiles).Error; err != nil {
		return nil, 0, err
	}

	return profiles, total, nil
}

func (r *ProfileRepository) Create(ctx context.Context, profile *models.Profile) error {
	if err := r.db.WithContext(ctx).Create(profile).Error; err != nil {
		return err
	}

	return nil
}
