package service_interface

import (
	"context"

	dto_favorite "github.com/KaueTTS/streaming_api/src/api/v1/dto/favorite"
)

type FavoriteServiceInterface interface {
	ListFavorites(ctx context.Context, userID, profileID uint, page, perPage int) (dto_favorite.FavoriteResponseDto, error)
}
