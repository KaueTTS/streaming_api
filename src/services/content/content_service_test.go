package service_content_test

import (
	"context"
	"errors"
	"testing"

	dto_content "github.com/KaueTTS/streaming_api/src/api/v1/dto/content"
	repository_mock "github.com/KaueTTS/streaming_api/src/mocks/repositories"
	models "github.com/KaueTTS/streaming_api/src/models"
	tmdb_dto "github.com/KaueTTS/streaming_api/src/repositories/http/tmdb/dto"
	service_content "github.com/KaueTTS/streaming_api/src/services/content"
	shared_errors "github.com/KaueTTS/streaming_api/src/shared/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func TestListContents(t *testing.T) {
	t.Run("should list contents using authenticated user profile", func(t *testing.T) {
		ctx := context.Background()
		tmdbRepository := new(repository_mock.TMDBRepositoryMock)
		profileRepository := new(repository_mock.ProfileRepositoryMock)
		service := service_content.NewContentService(tmdbRepository, profileRepository)
		request := dto_content.ContentListRequestDto{
			Type:      " Movie ",
			Page:      0,
			SortBy:    " Popularity.Desc ",
			Genre:     " 28,12 ",
			Language:  " pt-BR ",
			Year:      2026,
			ProfileID: 2,
		}
		profile := &models.Profile{ID: 2, UserID: 1, IsKids: true}
		tmdbResponse := tmdb_dto.GetContentResponseDto{
			Page:         1,
			TotalPages:   3,
			TotalResults: 21,
			Results: []tmdb_dto.ContentDto{
				{
					ID:               101,
					Title:            "Movie Title",
					OriginalTitle:    "Original Movie",
					Overview:         "Movie overview",
					OriginalLanguage: "en",
					ReleaseDate:      "2026-01-01",
					PosterPath:       "/movie-poster.jpg",
					BackdropPath:     "/movie-backdrop.jpg",
					GenreIDs:         []int{28, 12},
					Popularity:       91.5,
					VoteAverage:      8.7,
					VoteCount:        1234,
					Adult:            false,
				},
			},
		}

		profileRepository.On("FindProfileByUserIDAndID", ctx, uint(1), uint(2)).
			Return(profile, nil).
			Once()
		tmdbRepository.On("ListContents", ctx, mock.MatchedBy(func(filters tmdb_dto.ContentFiltersDto) bool {
			return filters.Type == "movie" &&
				filters.Page == 1 &&
				filters.SortBy == "popularity.desc" &&
				filters.WithGenres == "28,12" &&
				filters.Language == "pt-BR" &&
				filters.Year == 2026 &&
				filters.IsKids
		})).
			Return(tmdbResponse, nil).
			Once()

		response, err := service.ListContents(ctx, 1, request)

		require.NoError(t, err)
		require.Len(t, response.Data, 1)
		assert.Equal(t, dto_content.ContentDto{
			ExternalID:       101,
			Type:             "movie",
			Title:            "Movie Title",
			OriginalTitle:    "Original Movie",
			Description:      "Movie overview",
			OriginalLanguage: "en",
			ReleaseDate:      "2026-01-01",
			PosterPath:       "/movie-poster.jpg",
			BackdropPath:     "/movie-backdrop.jpg",
			GenreIDs:         []int{28, 12},
			Popularity:       91.5,
			VoteAverage:      8.7,
			VoteCount:        1234,
			Adult:            false,
		}, response.Data[0])
		assert.Equal(t, 1, response.Pagination.Page)
		assert.Equal(t, 1, response.Pagination.PerPage)
		assert.Equal(t, 3, response.Pagination.PageCount)
		assert.Equal(t, int64(21), response.Pagination.Total)
		profileRepository.AssertExpectations(t)
		tmdbRepository.AssertExpectations(t)
	})

	t.Run("should return profile not found when profile does not belong to user", func(t *testing.T) {
		ctx := context.Background()
		tmdbRepository := new(repository_mock.TMDBRepositoryMock)
		profileRepository := new(repository_mock.ProfileRepositoryMock)
		service := service_content.NewContentService(tmdbRepository, profileRepository)
		request := dto_content.ContentListRequestDto{
			Type:      "movie",
			ProfileID: 99,
		}

		profileRepository.On("FindProfileByUserIDAndID", ctx, uint(1), uint(99)).
			Return((*models.Profile)(nil), gorm.ErrRecordNotFound).
			Once()

		response, err := service.ListContents(ctx, 1, request)

		assert.Empty(t, response)
		assert.ErrorIs(t, err, shared_errors.ErrProfileNotFound)
		tmdbRepository.AssertNotCalled(t, "ListContents", mock.Anything, mock.Anything)
		profileRepository.AssertExpectations(t)
	})

	t.Run("should return repository error when finding profile fails", func(t *testing.T) {
		ctx := context.Background()
		tmdbRepository := new(repository_mock.TMDBRepositoryMock)
		profileRepository := new(repository_mock.ProfileRepositoryMock)
		service := service_content.NewContentService(tmdbRepository, profileRepository)
		expectedErr := errors.New("database unavailable")
		request := dto_content.ContentListRequestDto{
			Type:      "movie",
			ProfileID: 2,
		}

		profileRepository.On("FindProfileByUserIDAndID", ctx, uint(1), uint(2)).
			Return((*models.Profile)(nil), expectedErr).
			Once()

		response, err := service.ListContents(ctx, 1, request)

		assert.Empty(t, response)
		assert.ErrorIs(t, err, expectedErr)
		tmdbRepository.AssertNotCalled(t, "ListContents", mock.Anything, mock.Anything)
		profileRepository.AssertExpectations(t)
	})

	t.Run("should return tmdb error when listing contents fails", func(t *testing.T) {
		ctx := context.Background()
		tmdbRepository := new(repository_mock.TMDBRepositoryMock)
		profileRepository := new(repository_mock.ProfileRepositoryMock)
		service := service_content.NewContentService(tmdbRepository, profileRepository)
		expectedErr := errors.New("tmdb unavailable")
		request := dto_content.ContentListRequestDto{
			Type:      "movie",
			Page:      1,
			ProfileID: 2,
		}
		profile := &models.Profile{ID: 2, UserID: 1, IsKids: false}

		profileRepository.On("FindProfileByUserIDAndID", ctx, uint(1), uint(2)).
			Return(profile, nil).
			Once()
		tmdbRepository.On("ListContents", ctx, mock.MatchedBy(func(filters tmdb_dto.ContentFiltersDto) bool {
			return filters.Type == "movie" &&
				filters.Page == 1 &&
				!filters.IsKids
		})).
			Return(tmdb_dto.GetContentResponseDto{}, expectedErr).
			Once()

		response, err := service.ListContents(ctx, 1, request)

		assert.Empty(t, response)
		assert.ErrorIs(t, err, expectedErr)
		profileRepository.AssertExpectations(t)
		tmdbRepository.AssertExpectations(t)
	})
}

func TestSearchContents(t *testing.T) {
	t.Run("should search contents using authenticated user profile", func(t *testing.T) {
		ctx := context.Background()
		tmdbRepository := new(repository_mock.TMDBRepositoryMock)
		profileRepository := new(repository_mock.ProfileRepositoryMock)
		service := service_content.NewContentService(tmdbRepository, profileRepository)
		request := dto_content.ContentSearchRequestDto{
			Type:      "TV",
			Page:      2,
			Language:  " en-US ",
			Query:     " avatar ",
			ProfileID: 3,
		}
		profile := &models.Profile{ID: 3, UserID: 1, IsKids: true}
		tmdbResponse := tmdb_dto.GetContentResponseDto{
			Page:         2,
			TotalPages:   4,
			TotalResults: 31,
			Results: []tmdb_dto.ContentDto{
				{
					ID:               202,
					Name:             "TV Name",
					OriginalName:     "Original TV",
					Overview:         "TV overview",
					OriginalLanguage: "de",
					FirstAirDate:     "2020-09-01",
					PosterPath:       "/tv-poster.jpg",
					BackdropPath:     "/tv-backdrop.jpg",
					GenreIDs:         []int{18},
					Popularity:       77.2,
					VoteAverage:      8.1,
					VoteCount:        987,
					Adult:            false,
				},
			},
		}

		profileRepository.On("FindProfileByUserIDAndID", ctx, uint(1), uint(3)).
			Return(profile, nil).
			Once()
		tmdbRepository.On("SearchContents", ctx, mock.MatchedBy(func(filters tmdb_dto.ContentFiltersDto) bool {
			return filters.Type == "tv" &&
				filters.Page == 2 &&
				filters.Language == "en-US" &&
				filters.Query == "avatar" &&
				filters.IsKids
		})).
			Return(tmdbResponse, nil).
			Once()

		response, err := service.SearchContents(ctx, 1, request)

		require.NoError(t, err)
		require.Len(t, response.Data, 1)
		assert.Equal(t, dto_content.ContentDto{
			ExternalID:       202,
			Type:             "tv",
			Title:            "TV Name",
			OriginalTitle:    "Original TV",
			Description:      "TV overview",
			OriginalLanguage: "de",
			ReleaseDate:      "2020-09-01",
			PosterPath:       "/tv-poster.jpg",
			BackdropPath:     "/tv-backdrop.jpg",
			GenreIDs:         []int{18},
			Popularity:       77.2,
			VoteAverage:      8.1,
			VoteCount:        987,
			Adult:            false,
		}, response.Data[0])
		assert.Equal(t, 2, response.Pagination.Page)
		assert.Equal(t, 1, response.Pagination.PerPage)
		assert.Equal(t, 4, response.Pagination.PageCount)
		assert.Equal(t, int64(31), response.Pagination.Total)
		profileRepository.AssertExpectations(t)
		tmdbRepository.AssertExpectations(t)
	})

	t.Run("should return profile not found when profile does not belong to user", func(t *testing.T) {
		ctx := context.Background()
		tmdbRepository := new(repository_mock.TMDBRepositoryMock)
		profileRepository := new(repository_mock.ProfileRepositoryMock)
		service := service_content.NewContentService(tmdbRepository, profileRepository)
		request := dto_content.ContentSearchRequestDto{
			Type:      "movie",
			Query:     "matrix",
			ProfileID: 99,
		}

		profileRepository.On("FindProfileByUserIDAndID", ctx, uint(1), uint(99)).
			Return((*models.Profile)(nil), gorm.ErrRecordNotFound).
			Once()

		response, err := service.SearchContents(ctx, 1, request)

		assert.Empty(t, response)
		assert.ErrorIs(t, err, shared_errors.ErrProfileNotFound)
		tmdbRepository.AssertNotCalled(t, "SearchContents", mock.Anything, mock.Anything)
		profileRepository.AssertExpectations(t)
	})

	t.Run("should return repository error when finding profile fails", func(t *testing.T) {
		ctx := context.Background()
		tmdbRepository := new(repository_mock.TMDBRepositoryMock)
		profileRepository := new(repository_mock.ProfileRepositoryMock)
		service := service_content.NewContentService(tmdbRepository, profileRepository)
		expectedErr := errors.New("database unavailable")
		request := dto_content.ContentSearchRequestDto{
			Type:      "movie",
			Query:     "matrix",
			ProfileID: 2,
		}

		profileRepository.On("FindProfileByUserIDAndID", ctx, uint(1), uint(2)).
			Return((*models.Profile)(nil), expectedErr).
			Once()

		response, err := service.SearchContents(ctx, 1, request)

		assert.Empty(t, response)
		assert.ErrorIs(t, err, expectedErr)
		tmdbRepository.AssertNotCalled(t, "SearchContents", mock.Anything, mock.Anything)
		profileRepository.AssertExpectations(t)
	})

	t.Run("should return tmdb error when searching contents fails", func(t *testing.T) {
		ctx := context.Background()
		tmdbRepository := new(repository_mock.TMDBRepositoryMock)
		profileRepository := new(repository_mock.ProfileRepositoryMock)
		service := service_content.NewContentService(tmdbRepository, profileRepository)
		expectedErr := errors.New("tmdb unavailable")
		request := dto_content.ContentSearchRequestDto{
			Type:      "movie",
			Query:     "matrix",
			ProfileID: 2,
		}
		profile := &models.Profile{ID: 2, UserID: 1, IsKids: false}

		profileRepository.On("FindProfileByUserIDAndID", ctx, uint(1), uint(2)).
			Return(profile, nil).
			Once()
		tmdbRepository.On("SearchContents", ctx, mock.MatchedBy(func(filters tmdb_dto.ContentFiltersDto) bool {
			return filters.Type == "movie" &&
				filters.Page == 1 &&
				filters.Query == "matrix" &&
				!filters.IsKids
		})).
			Return(tmdb_dto.GetContentResponseDto{}, expectedErr).
			Once()

		response, err := service.SearchContents(ctx, 1, request)

		assert.Empty(t, response)
		assert.ErrorIs(t, err, expectedErr)
		profileRepository.AssertExpectations(t)
		tmdbRepository.AssertExpectations(t)
	})
}
