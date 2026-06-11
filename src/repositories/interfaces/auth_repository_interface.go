package repository_interface

import (
	"context"

	models "github.com/KaueTTS/streaming_api/src/models"
)

type AuthRepositoryInterface interface {
	Create(ctx context.Context, user *models.User) error
	FindByEmail(ctx context.Context, email string) (*models.User, error)
	FindByID(ctx context.Context, id uint) (*models.User, error)
}
