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

func (repository *ProfileRepository) FindAll(ctx context.Context) ([]models.Profile, error) {
	var profiles []models.Profile

	if err := repository.db.WithContext(ctx).Order("name asc").Find(&profiles).Error; err != nil {
		return nil, err
	}

	return profiles, nil
}

func (repository *ProfileRepository) FindByUserID(ctx context.Context, userID uint) ([]models.Profile, error) {
	var profiles []models.Profile

	if err := repository.db.WithContext(ctx).Where("user_id = ?", userID).Order("name asc").Find(&profiles).Error; err != nil {
		return nil, err
	}

	return profiles, nil
}
