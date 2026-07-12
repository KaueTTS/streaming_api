package repository_interface

import (
	"context"

	models "github.com/KaueTTS/streaming_api/src/models"
)

type FavoriteRepositoryInterface interface {
	FindByProfileID(ctx context.Context, userID, profileID uint, page, perPage int) ([]models.Favorite, int64, error)
}
