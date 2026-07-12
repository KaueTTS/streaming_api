package service_favorite

import (
	"context"
	"errors"

	dto_favorite "github.com/KaueTTS/streaming_api/src/api/v1/dto/favorite"
	dto_shared "github.com/KaueTTS/streaming_api/src/api/v1/dto/shared"
	models "github.com/KaueTTS/streaming_api/src/models"
	repository_interface "github.com/KaueTTS/streaming_api/src/repositories/interfaces"
	shared_errors "github.com/KaueTTS/streaming_api/src/shared/errors"
	"gorm.io/gorm"
)

type FavoriteService struct {
	FavoriteRepositoryInterface repository_interface.FavoriteRepositoryInterface
}

func NewFavoriteService(favoriteRepositoryInterface repository_interface.FavoriteRepositoryInterface) *FavoriteService {
	return &FavoriteService{
		FavoriteRepositoryInterface: favoriteRepositoryInterface,
	}
}

func (s *FavoriteService) ListFavorites(ctx context.Context, userID, profileID uint, page, perPage int) (dto_favorite.FavoriteResponseDto, error) {
	var favorites []models.Favorite
	var total int64
	var err error

	favorites, total, err = s.FavoriteRepositoryInterface.FindByProfileID(ctx, userID, profileID, page, perPage)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return dto_favorite.FavoriteResponseDto{}, shared_errors.ErrProfileNotFound
		}

		return dto_favorite.FavoriteResponseDto{}, err
	}

	data := make([]dto_favorite.FavoriteDto, 0, len(favorites))
	for _, favorite := range favorites {
		data = append(data, dto_favorite.FavoriteDto{
			ID:          favorite.ID,
			ProfileID:   favorite.ProfileID,
			ContentID:   favorite.ContentID,
			ContentType: favorite.ContentType,
			CreatedAt:   favorite.CreatedAt,
			UpdatedAt:   favorite.UpdatedAt,
		})
	}

	pageCount := int((total + int64(perPage) - 1) / int64(perPage))

	return dto_favorite.FavoriteResponseDto{
		Data: data,
		Pagination: dto_shared.PaginationDto{
			Page:      page,
			PerPage:   perPage,
			PageCount: pageCount,
			Total:     total,
		},
	}, nil
}
