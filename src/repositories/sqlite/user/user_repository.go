package repository_sqlite_user

import (
	"context"
	"strings"

	models "github.com/KaueTTS/streaming_api/src/models"
	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (repository *UserRepository) Create(ctx context.Context, user *models.User) error {
	return repository.db.WithContext(ctx).Create(user).Error
}

func (repository *UserRepository) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	normalizedEmail := strings.ToLower(strings.TrimSpace(email))

	if err := repository.db.WithContext(ctx).Where("email = ?", normalizedEmail).First(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func (repository *UserRepository) FindByID(ctx context.Context, id uint) (*models.User, error) {
	var user models.User

	if err := repository.db.WithContext(ctx).First(&user, id).Error; err != nil {
		return nil, err
	}

	return &user, nil
}
