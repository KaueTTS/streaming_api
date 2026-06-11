package repository_interface

import (
	"context"

	tmdb_dto "github.com/KaueTTS/streaming_api/src/repositories/http/tmdb/dto"
)

type TMDBRepositoryInterface interface {
	ListContents(ctx context.Context, filters tmdb_dto.ContentListFiltersDto) (tmdb_dto.GetContentResponseDto, error)
	SearchContents(ctx context.Context, filters tmdb_dto.ContentSearchFiltersDto) (tmdb_dto.GetContentResponseDto, error)
}
