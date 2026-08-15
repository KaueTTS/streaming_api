package service_content

import (
	"context"
	"errors"
	"strings"

	dto_content "github.com/KaueTTS/streaming_api/src/api/v1/dto/content"
	dto_shared "github.com/KaueTTS/streaming_api/src/api/v1/dto/shared"
	tmdb_dto "github.com/KaueTTS/streaming_api/src/repositories/http/tmdb/dto"
	repository_interface "github.com/KaueTTS/streaming_api/src/repositories/interfaces"
	shared_errors "github.com/KaueTTS/streaming_api/src/shared/errors"
	shared_normalizers "github.com/KaueTTS/streaming_api/src/shared/normalizers"
	"gorm.io/gorm"
)

type ContentService struct {
	TMDBRepositoryInterface    repository_interface.TMDBRepositoryInterface
	ProfileRepositoryInterface repository_interface.ProfileRepositoryInterface
}

func NewContentService(
	tmdbRepositoryInterface repository_interface.TMDBRepositoryInterface,
	profileRepositoryInterface repository_interface.ProfileRepositoryInterface,
) *ContentService {
	return &ContentService{
		TMDBRepositoryInterface:    tmdbRepositoryInterface,
		ProfileRepositoryInterface: profileRepositoryInterface,
	}
}

// ListContents lista os conteúdos da TMDB baseado nos filtros passados pelo usuário
func (s *ContentService) ListContents(ctx context.Context, request dto_content.ContentListRequestDto) (dto_content.ContentResponseDto, error) {
	profile, err := s.ProfileRepositoryInterface.FindProfileByID(ctx, request.ProfileID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return dto_content.ContentResponseDto{}, shared_errors.ErrProfileNotFound
	}

	filters := tmdb_dto.ContentFiltersDto{
		Type:       shared_normalizers.NormalizeString(request.Type),
		Page:       shared_normalizers.NormalizePage(request.Page),
		SortBy:     shared_normalizers.NormalizeString(request.SortBy),
		WithGenres: strings.TrimSpace(request.Genre),
		Language:   strings.TrimSpace(request.Language),
		Year:       request.Year,
		IsKids:     profile.IsKids,
	}

	response, err := s.TMDBRepositoryInterface.ListContents(ctx, filters)
	if err != nil {
		return dto_content.ContentResponseDto{}, err
	}

	return mapTMDBContentResponse(response, filters.Type), nil
}

// SearchContents busca os conteúdos da TMDB baseado nos filtros passados pelo usuário
func (s *ContentService) SearchContents(ctx context.Context, request dto_content.ContentSearchRequestDto) (dto_content.ContentResponseDto, error) {
	filters := tmdb_dto.ContentFiltersDto{
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

// mapTMDBContentResponse mapeia a resposta da TMDB para a resposta do ContentService
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
			Adult:            content.Adult,
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

// firstNonEmpty retorna a primeira string não vazia
func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if value != "" {
			return value
		}
	}

	return ""
}
