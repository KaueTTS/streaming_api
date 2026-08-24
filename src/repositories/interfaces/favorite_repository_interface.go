package repository_interface

import (
	"context"

	dto_favorite "github.com/KaueTTS/streaming_api/src/api/v1/dto/favorite"
	models "github.com/KaueTTS/streaming_api/src/models"
)

type FavoriteRepositoryInterface interface {
	FindFavoriteByProfileID(ctx context.Context, userID, profileID uint, page, perPage int) ([]models.Favorite, int64, error)
	CreateFavoriteByProfileID(ctx context.Context, userID uint, request dto_favorite.FavoriteRequestDto) error
}
