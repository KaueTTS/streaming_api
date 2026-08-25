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

// FindByUserID lista os perfis de um usuário, verificando se o perfil pertence ao usuário
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

// CountByUserID conta quantos perfis um usuário tem, verificando se o perfil pertence ao usuário
func (r *ProfileRepository) CountByUserID(ctx context.Context, userID uint) (int64, error) {
	var total int64

	if err := r.db.WithContext(ctx).
		Model(&models.Profile{}).
		Where("user_id = ?", userID).
		Count(&total).Error; err != nil {
		return 0, err
	}

	return total, nil
}

// Create cria um novo perfil para um usuário, verificando se o perfil pertence ao usuário
func (r *ProfileRepository) Create(ctx context.Context, profile *models.Profile) error {
	if err := r.db.WithContext(ctx).Create(profile).Error; err != nil {
		return err
	}

	return nil
}

// Update atualiza um perfil, verificando se o perfil pertence ao usuário
func (r *ProfileRepository) Update(ctx context.Context, profile *models.Profile) error {
	var existingProfile models.Profile

	if err := r.db.WithContext(ctx).
		Where("id = ? AND user_id = ?", profile.ID, profile.UserID).
		First(&existingProfile).Error; err != nil {
		return err
	}

	existingProfile.Name = profile.Name
	existingProfile.AvatarURL = profile.AvatarURL
	existingProfile.IsKids = profile.IsKids

	if err := r.db.WithContext(ctx).Save(&existingProfile).Error; err != nil {
		return err
	}

	*profile = existingProfile
	return nil
}

// Delete deleta um perfil, verificando se o perfil pertence ao usuário
func (r *ProfileRepository) Delete(ctx context.Context, userID, profileID uint) error {
	result := r.db.WithContext(ctx).
		Where("id = ? AND user_id = ?", profileID, userID).
		Delete(&models.Profile{})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

// FindProfileByID busca um perfil pelo id, sem verificar se o perfil pertence ao usuário
func (r *ProfileRepository) FindProfileByID(ctx context.Context, profileID uint) (*models.Profile, error) {
	var profile models.Profile

	if err := r.db.WithContext(ctx).
		Where("id = ?", profileID).
		First(&profile).Error; err != nil {
		return nil, err
	}

	return &profile, nil
}

// FindProfileByUserIDAndID busca um perfil pelo id, verificando se pertence ao usuário.
func (r *ProfileRepository) FindProfileByUserIDAndID(ctx context.Context, userID, profileID uint) (*models.Profile, error) {
	var profile models.Profile

	if err := r.db.WithContext(ctx).
		Where("id = ? AND user_id = ?", profileID, userID).
		First(&profile).Error; err != nil {
		return nil, err
	}

	return &profile, nil
}
