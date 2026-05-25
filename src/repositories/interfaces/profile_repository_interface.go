package repository_interface

import (
	"context"

	models "github.com/KaueTTS/streaming_api/src/models"
)

type ProfileRepositoryInterface interface {
	FindAll(ctx context.Context) ([]models.Profile, error)
	FindByUserID(ctx context.Context, userID uint) ([]models.Profile, error)
}
