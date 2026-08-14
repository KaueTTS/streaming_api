package repository_sqlite_favorite

import (
	"context"

	models "github.com/KaueTTS/streaming_api/src/models"
	"gorm.io/gorm"
)

type FavoriteRepository struct {
	db *gorm.DB
}

func NewFavoriteRepository(db *gorm.DB) *FavoriteRepository {
	return &FavoriteRepository{
		db: db,
	}
}

// FindByProfileID busca os favoritos de um perfil, verificando se o perfil pertence ao usuário
func (r *FavoriteRepository) FindByProfileID(ctx context.Context, userID, profileID uint, page, perPage int) ([]models.Favorite, int64, error) {
	var profileCount int64

	if err := r.db.WithContext(ctx).
		Model(&models.Profile{}).
		Where("id = ? AND user_id = ?", profileID, userID).
		Count(&profileCount).Error; err != nil {
		return nil, 0, err
	}

	if profileCount == 0 {
		return nil, 0, gorm.ErrRecordNotFound
	}

	var favorites []models.Favorite
	var total int64

	query := r.db.WithContext(ctx).
		Model(&models.Favorite{}).
		Where("profile_id = ?", profileID)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * perPage

	if err := query.
		Order("created_at desc").
		Limit(perPage).
		Offset(offset).
		Find(&favorites).Error; err != nil {
		return nil, 0, err
	}

	return favorites, total, nil
}
