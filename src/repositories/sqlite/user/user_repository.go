package repository_sqlite_user

import (
	"context"

	models "github.com/KaueTTS/streaming_api/src/models"
	shared_normalizers "github.com/KaueTTS/streaming_api/src/shared/normalizers"
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

// Create cria um novo usuário
func (r *UserRepository) Create(ctx context.Context, user *models.User) error {
	return r.db.WithContext(ctx).Create(user).Error
}

// FindByEmail busca um usuário pelo email
func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	normalizedEmail := shared_normalizers.NormalizeString(email)

	if err := r.db.WithContext(ctx).Where("email = ?", normalizedEmail).First(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

// FindByID busca um usuário pelo ID
func (r *UserRepository) FindByID(ctx context.Context, id uint) (*models.User, error) {
	var user models.User

	if err := r.db.WithContext(ctx).First(&user, id).Error; err != nil {
		return nil, err
	}

	return &user, nil
}
