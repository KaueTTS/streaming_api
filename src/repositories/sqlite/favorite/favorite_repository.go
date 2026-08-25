package repository_sqlite_favorite

import (
	"context"

	dto_favorite "github.com/KaueTTS/streaming_api/src/api/v1/dto/favorite"
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

// FindFavoriteByProfileID busca os favoritos de um perfil do usuario.
func (r *FavoriteRepository) FindFavoriteByProfileID(ctx context.Context, userID, profileID uint, page, perPage int) ([]models.Favorite, int64, error) {
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

// CreateFavoriteByProfileID cria um favorito em um perfil do usuario.
func (r *FavoriteRepository) CreateFavoriteByProfileID(ctx context.Context, userID uint, request dto_favorite.FavoriteRequestDto) error {
	if err := r.db.WithContext(ctx).
		Select("id").
		Where("id = ? AND user_id = ?", request.ProfileID, userID).
		Take(&models.Profile{}).Error; err != nil {
		return err
	}

	favorite := models.Favorite{
		UserID:            userID,
		ProfileID:         request.ProfileID,
		ContentExternalId: request.ContentExternalID,
		Type:              request.Type,
	}

	if err := r.db.WithContext(ctx).Create(&favorite).Error; err != nil {
		return err
	}

	return nil
}

// DeleteFavoriteByProfileID remove um favorito de um perfil do usuario.
func (r *FavoriteRepository) DeleteFavoriteByProfileID(ctx context.Context, userID uint, request dto_favorite.FavoriteRequestDto) error {
	if err := r.db.WithContext(ctx).
		Select("id").
		Where("id = ? AND user_id = ?", request.ProfileID, userID).
		Take(&models.Profile{}).Error; err != nil {
		return err
	}

	if err := r.db.WithContext(ctx).
		Where(
			"user_id = ? AND profile_id = ? AND content_external_id = ? AND type = ?",
			userID,
			request.ProfileID,
			request.ContentExternalID,
			request.Type,
		).
		Delete(&models.Favorite{}).Error; err != nil {
		return err
	}

	return nil
}
