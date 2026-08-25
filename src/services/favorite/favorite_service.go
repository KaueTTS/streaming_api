package service_favorite

import (
	"context"
	"errors"
	"strings"

	dto_content "github.com/KaueTTS/streaming_api/src/api/v1/dto/content"
	dto_favorite "github.com/KaueTTS/streaming_api/src/api/v1/dto/favorite"
	dto_shared "github.com/KaueTTS/streaming_api/src/api/v1/dto/shared"
	models "github.com/KaueTTS/streaming_api/src/models"
	tmdb_dto "github.com/KaueTTS/streaming_api/src/repositories/http/tmdb/dto"
	repository_interface "github.com/KaueTTS/streaming_api/src/repositories/interfaces"
	shared_errors "github.com/KaueTTS/streaming_api/src/shared/errors"
	"gorm.io/gorm"
)

type FavoriteService struct {
	FavoriteRepositoryInterface repository_interface.FavoriteRepositoryInterface
	TMDBRepositoryInterface     repository_interface.TMDBRepositoryInterface
}

func NewFavoriteService(
	favoriteRepositoryInterface repository_interface.FavoriteRepositoryInterface,
	tmdbRepositoryInterface repository_interface.TMDBRepositoryInterface,
) *FavoriteService {
	return &FavoriteService{
		FavoriteRepositoryInterface: favoriteRepositoryInterface,
		TMDBRepositoryInterface:     tmdbRepositoryInterface,
	}
}

// ListFavorites lista os favoritos de um perfil com os dados do TMDB.
func (s *FavoriteService) ListFavorites(ctx context.Context, userID, profileID uint, page, perPage int, language string) (dto_favorite.FavoriteResponseDto, error) {
	var favorites []models.Favorite
	var total int64
	var err error

	favorites, total, err = s.FavoriteRepositoryInterface.FindFavoriteByProfileID(ctx, userID, profileID, page, perPage)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return dto_favorite.FavoriteResponseDto{}, shared_errors.ErrProfileNotFound
		}

		return dto_favorite.FavoriteResponseDto{}, err
	}

	data := make([]dto_favorite.FavoriteDto, 0, len(favorites))
	language = strings.TrimSpace(language)
	for _, favorite := range favorites {
		content, err := s.TMDBRepositoryInterface.GetContentByID(ctx, favorite.Type, favorite.ContentExternalId, language)
		if err != nil {
			return dto_favorite.FavoriteResponseDto{}, err
		}

		contentDto := mapTMDBFavoriteContent(content, favorite.Type)

		data = append(data, dto_favorite.FavoriteDto{
			ID:        favorite.ID,
			ProfileID: favorite.ProfileID,
			Content:   &contentDto,
			CreatedAt: favorite.CreatedAt,
			UpdatedAt: favorite.UpdatedAt,
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

// AddFavorite adiciona um conteudo aos favoritos de um perfil.
func (s *FavoriteService) AddFavorite(ctx context.Context, userID uint, request dto_favorite.FavoriteRequestDto) error {
	err := s.FavoriteRepositoryInterface.CreateFavoriteByProfileID(ctx, userID, request)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return shared_errors.ErrProfileNotFound
		}

		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return shared_errors.ErrFavoriteAlreadyExists
		}

		return err
	}

	return nil
}

// DeleteFavorite remove um conteudo dos favoritos de um perfil.
func (s *FavoriteService) DeleteFavorite(ctx context.Context, userID uint, request dto_favorite.FavoriteRequestDto) error {
	err := s.FavoriteRepositoryInterface.DeleteFavoriteByProfileID(ctx, userID, request)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return shared_errors.ErrProfileNotFound
		}

		return err
	}

	return nil
}

// mapTMDBFavoriteContent converte o conteudo do TMDB para o dto da API.
func mapTMDBFavoriteContent(content tmdb_dto.ContentDto, contentType string) dto_content.ContentDto {
	return dto_content.ContentDto{
		ExternalID:       content.ID,
		Type:             contentType,
		Title:            firstNonEmpty(content.Title, content.Name),
		OriginalTitle:    firstNonEmpty(content.OriginalTitle, content.OriginalName),
		Description:      content.Overview,
		OriginalLanguage: content.OriginalLanguage,
		ReleaseDate:      firstNonEmpty(content.ReleaseDate, content.FirstAirDate),
		PosterPath:       content.PosterPath,
		BackdropPath:     content.BackdropPath,
		GenreIDs:         getGenreIDs(content),
		Popularity:       content.Popularity,
		VoteAverage:      content.VoteAverage,
		VoteCount:        content.VoteCount,
		Adult:            content.Adult,
	}
}

// getGenreIDs extrai os ids dos generos do conteudo.
func getGenreIDs(content tmdb_dto.ContentDto) []int {
	if len(content.GenreIDs) > 0 {
		return content.GenreIDs
	}

	genreIDs := make([]int, 0, len(content.Genres))
	for _, genre := range content.Genres {
		genreIDs = append(genreIDs, genre.ID)
	}

	return genreIDs
}

// firstNonEmpty retorna o primeiro texto nao vazio.
func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if value != "" {
			return value
		}
	}

	return ""
}
