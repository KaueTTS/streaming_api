package repository_interface

import (
	"context"

	models "github.com/KaueTTS/streaming_api/src/models"
)

type ProfileRepositoryInterface interface {
	FindByUserID(ctx context.Context, userID uint, page int, perPage int) ([]models.Profile, int64, error)
	CountByUserID(ctx context.Context, userID uint) (int64, error)
	Create(ctx context.Context, profile *models.Profile) error
}
