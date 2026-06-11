package service_content

import (
	"context"
	"strings"

	dto_content "github.com/KaueTTS/streaming_api/src/api/v1/dto/content"
	dto_shared "github.com/KaueTTS/streaming_api/src/api/v1/dto/shared"
	tmdb_dto "github.com/KaueTTS/streaming_api/src/repositories/http/tmdb/dto"
	repository_interface "github.com/KaueTTS/streaming_api/src/repositories/interfaces"
	shared_normalizers "github.com/KaueTTS/streaming_api/src/shared/normalizers"
)

type ContentService struct {
	TMDBRepositoryInterface repository_interface.TMDBRepositoryInterface
}

func NewContentService(tmdbRepositoryInterface repository_interface.TMDBRepositoryInterface) *ContentService {
	return &ContentService{
		TMDBRepositoryInterface: tmdbRepositoryInterface,
	}
}

func (s *ContentService) ListContents(ctx context.Context, request dto_content.ContentListRequestDto) (dto_content.ContentResponseDto, error) {
	filters := tmdb_dto.ContentListFiltersDto{
		Type:       shared_normalizers.NormalizeString(request.Type),
		Page:       shared_normalizers.NormalizePage(request.Page),
		SortBy:     shared_normalizers.NormalizeString(request.SortBy),
		WithGenres: strings.TrimSpace(request.WithGenres),
		Language:   strings.TrimSpace(request.Language),
		Year:       request.Year,
	}

	response, err := s.TMDBRepositoryInterface.ListContents(ctx, filters)
	if err != nil {
		return dto_content.ContentResponseDto{}, err
	}

	return mapTMDBContentResponse(response, filters.Type), nil
}

func (s *ContentService) SearchContents(ctx context.Context, request dto_content.ContentSearchRequestDto) (dto_content.ContentResponseDto, error) {
	filters := tmdb_dto.ContentSearchFiltersDto{
		Type:     shared_normalizers.NormalizeString(request.Type),
		Page:     shared_normalizers.NormalizePage(request.Page),
		Language: strings.TrimSpace(request.Language),
		Query:    strings.TrimSpace(request.Query),
	}

	response, err := s.TMDBRepositoryInterface.SearchContents(ctx, filters)
	if err != nil {
		return dto_content.ContentResponseDto{}, err
	}

	return mapTMDBContentResponse(response, filters.Type), nil
}

func mapTMDBContentResponse(response tmdb_dto.GetContentResponseDto, contentType string) dto_content.ContentResponseDto {
	data := make([]dto_content.ContentDto, 0, len(response.Results))
	for _, content := range response.Results {
		data = append(data, dto_content.ContentDto{
			ExternalID:       content.ID,
			Type:             contentType,
			Title:            firstNonEmpty(content.Title, content.Name),
			OriginalTitle:    firstNonEmpty(content.OriginalTitle, content.OriginalName),
			Description:      content.Overview,
			OriginalLanguage: content.OriginalLanguage,
			ReleaseDate:      firstNonEmpty(content.ReleaseDate, content.FirstAirDate),
			PosterPath:       content.PosterPath,
			BackdropPath:     content.BackdropPath,
			GenreIDs:         content.GenreIDs,
			Popularity:       content.Popularity,
			VoteAverage:      content.VoteAverage,
			VoteCount:        content.VoteCount,
		})
	}

	return dto_content.ContentResponseDto{
		Data: data,
		Pagination: dto_shared.PaginationDto{
			Page:      response.Page,
			PerPage:   len(response.Results),
			PageCount: response.TotalPages,
			Total:     int64(response.TotalResults),
		},
	}
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if value != "" {
			return value
		}
	}

	return ""
}
