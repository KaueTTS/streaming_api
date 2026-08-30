package repository_interface

import (
	"context"

	tmdb_dto "github.com/KaueTTS/streaming_api/src/repositories/http/tmdb/dto"
)

type TMDBRepositoryInterface interface {
	ListContents(ctx context.Context, filters tmdb_dto.ContentFiltersDto) (tmdb_dto.GetContentResponseDto, error)
	SearchContents(ctx context.Context, filters tmdb_dto.ContentFiltersDto) (tmdb_dto.GetContentResponseDto, error)
	GetContentByID(ctx context.Context, contentType string, contentExternalID int, language string) (tmdb_dto.ContentDto, error)
}
