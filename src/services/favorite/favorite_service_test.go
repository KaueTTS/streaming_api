package service_favorite_test

import (
	"context"
	"errors"
	"testing"
	"time"

	dto_content "github.com/KaueTTS/streaming_api/src/api/v1/dto/content"
	dto_favorite "github.com/KaueTTS/streaming_api/src/api/v1/dto/favorite"
	repository_mock "github.com/KaueTTS/streaming_api/src/mocks/repositories"
	models "github.com/KaueTTS/streaming_api/src/models"
	tmdb_dto "github.com/KaueTTS/streaming_api/src/repositories/http/tmdb/dto"
	service_favorite "github.com/KaueTTS/streaming_api/src/services/favorite"
	shared_errors "github.com/KaueTTS/streaming_api/src/shared/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func TestListFavorites(t *testing.T) {
	t.Run("should list favorites with tmdb content and pagination", func(t *testing.T) {
		ctx := context.Background()
		favoriteRepository := new(repository_mock.FavoriteRepositoryMock)
		tmdbRepository := new(repository_mock.TMDBRepositoryMock)
		service := service_favorite.NewFavoriteService(favoriteRepository, tmdbRepository)
		createdAt := time.Date(2026, 1, 1, 10, 0, 0, 0, time.UTC)
		updatedAt := time.Date(2026, 1, 1, 11, 0, 0, 0, time.UTC)
		favorites := []models.Favorite{
			{
				ID:                1,
				UserID:            10,
				ProfileID:         20,
				ContentExternalId: 101,
				Type:              "movie",
				CreatedAt:         createdAt,
				UpdatedAt:         updatedAt,
			},
			{
				ID:                2,
				UserID:            10,
				ProfileID:         20,
				ContentExternalId: 202,
				Type:              "tv",
				CreatedAt:         createdAt.Add(time.Hour),
				UpdatedAt:         updatedAt.Add(time.Hour),
			},
		}
		movieContent := tmdb_dto.ContentDto{
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
		}
		tvContent := tmdb_dto.ContentDto{
			ID:               202,
			Name:             "TV Name",
			OriginalName:     "Original TV",
			Overview:         "TV overview",
			OriginalLanguage: "de",
			FirstAirDate:     "2020-09-01",
			PosterPath:       "/tv-poster.jpg",
			BackdropPath:     "/tv-backdrop.jpg",
			Genres: []tmdb_dto.GenreDto{
				{ID: 18, Name: "Drama"},
				{ID: 80, Name: "Crime"},
			},
			Popularity:  77.2,
			VoteAverage: 8.1,
			VoteCount:   987,
			Adult:       false,
		}

		favoriteRepository.On("FindFavoriteByProfileID", ctx, uint(10), uint(20), 2, 2).
			Return(favorites, int64(5), nil).
			Once()
		tmdbRepository.On("GetContentByID", ctx, "movie", 101, "pt-BR").
			Return(movieContent, nil).
			Once()
		tmdbRepository.On("GetContentByID", ctx, "tv", 202, "pt-BR").
			Return(tvContent, nil).
			Once()

		response, err := service.ListFavorites(ctx, 10, 20, 2, 2, "  pt-BR  ")

		require.NoError(t, err)
		require.Len(t, response.Data, 2)
		assert.Equal(t, dto_favorite.FavoriteDto{
			ID:        1,
			ProfileID: 20,
			Content: &dto_content.ContentDto{
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
			},
			CreatedAt: createdAt,
			UpdatedAt: updatedAt,
		}, response.Data[0])
		assert.Equal(t, dto_favorite.FavoriteDto{
			ID:        2,
			ProfileID: 20,
			Content: &dto_content.ContentDto{
				ExternalID:       202,
				Type:             "tv",
				Title:            "TV Name",
				OriginalTitle:    "Original TV",
				Description:      "TV overview",
				OriginalLanguage: "de",
				ReleaseDate:      "2020-09-01",
				PosterPath:       "/tv-poster.jpg",
				BackdropPath:     "/tv-backdrop.jpg",
				GenreIDs:         []int{18, 80},
				Popularity:       77.2,
				VoteAverage:      8.1,
				VoteCount:        987,
				Adult:            false,
			},
			CreatedAt: createdAt.Add(time.Hour),
			UpdatedAt: updatedAt.Add(time.Hour),
		}, response.Data[1])
		assert.Equal(t, 2, response.Pagination.Page)
		assert.Equal(t, 2, response.Pagination.PerPage)
		assert.Equal(t, 3, response.Pagination.PageCount)
		assert.Equal(t, int64(5), response.Pagination.Total)
		favoriteRepository.AssertExpectations(t)
		tmdbRepository.AssertExpectations(t)
	})

	t.Run("should return profile not found when favorite repository cannot find profile", func(t *testing.T) {
		ctx := context.Background()
		favoriteRepository := new(repository_mock.FavoriteRepositoryMock)
		tmdbRepository := new(repository_mock.TMDBRepositoryMock)
		service := service_favorite.NewFavoriteService(favoriteRepository, tmdbRepository)

		favoriteRepository.On("FindFavoriteByProfileID", ctx, uint(10), uint(999), 1, 10).
			Return([]models.Favorite(nil), int64(0), gorm.ErrRecordNotFound).
			Once()

		response, err := service.ListFavorites(ctx, 10, 999, 1, 10, "pt-BR")

		assert.Empty(t, response)
		assert.ErrorIs(t, err, shared_errors.ErrProfileNotFound)
		tmdbRepository.AssertNotCalled(t, "GetContentByID", mock.Anything, mock.Anything, mock.Anything, mock.Anything)
		favoriteRepository.AssertExpectations(t)
	})

	t.Run("should return repository error when listing favorites fails", func(t *testing.T) {
		ctx := context.Background()
		favoriteRepository := new(repository_mock.FavoriteRepositoryMock)
		tmdbRepository := new(repository_mock.TMDBRepositoryMock)
		service := service_favorite.NewFavoriteService(favoriteRepository, tmdbRepository)
		expectedErr := errors.New("database unavailable")

		favoriteRepository.On("FindFavoriteByProfileID", ctx, uint(10), uint(20), 1, 10).
			Return([]models.Favorite(nil), int64(0), expectedErr).
			Once()

		response, err := service.ListFavorites(ctx, 10, 20, 1, 10, "pt-BR")

		assert.Empty(t, response)
		assert.ErrorIs(t, err, expectedErr)
		tmdbRepository.AssertNotCalled(t, "GetContentByID", mock.Anything, mock.Anything, mock.Anything, mock.Anything)
		favoriteRepository.AssertExpectations(t)
	})

	t.Run("should return tmdb error when loading favorite content fails", func(t *testing.T) {
		ctx := context.Background()
		favoriteRepository := new(repository_mock.FavoriteRepositoryMock)
		tmdbRepository := new(repository_mock.TMDBRepositoryMock)
		service := service_favorite.NewFavoriteService(favoriteRepository, tmdbRepository)
		expectedErr := errors.New("tmdb unavailable")
		favorites := []models.Favorite{
			{
				ID:                1,
				UserID:            10,
				ProfileID:         20,
				ContentExternalId: 101,
				Type:              "movie",
			},
		}

		favoriteRepository.On("FindFavoriteByProfileID", ctx, uint(10), uint(20), 1, 10).
			Return(favorites, int64(1), nil).
			Once()
		tmdbRepository.On("GetContentByID", ctx, "movie", 101, "pt-BR").
			Return(tmdb_dto.ContentDto{}, expectedErr).
			Once()

		response, err := service.ListFavorites(ctx, 10, 20, 1, 10, "pt-BR")

		assert.Empty(t, response)
		assert.ErrorIs(t, err, expectedErr)
		favoriteRepository.AssertExpectations(t)
		tmdbRepository.AssertExpectations(t)
	})
}

func TestAddFavorite(t *testing.T) {
	t.Run("should add favorite", func(t *testing.T) {
		ctx := context.Background()
		favoriteRepository := new(repository_mock.FavoriteRepositoryMock)
		tmdbRepository := new(repository_mock.TMDBRepositoryMock)
		service := service_favorite.NewFavoriteService(favoriteRepository, tmdbRepository)
		request := dto_favorite.FavoriteRequestDto{
			ProfileID:         20,
			ContentExternalID: 101,
			Type:              "movie",
		}

		favoriteRepository.On("CreateFavoriteByProfileID", ctx, uint(10), request).
			Return(nil).
			Once()

		err := service.AddFavorite(ctx, 10, request)

		require.NoError(t, err)
		tmdbRepository.AssertNotCalled(t, "GetContentByID", mock.Anything, mock.Anything, mock.Anything, mock.Anything)
		favoriteRepository.AssertExpectations(t)
	})

	t.Run("should return profile not found when profile does not exist", func(t *testing.T) {
		ctx := context.Background()
		favoriteRepository := new(repository_mock.FavoriteRepositoryMock)
		tmdbRepository := new(repository_mock.TMDBRepositoryMock)
		service := service_favorite.NewFavoriteService(favoriteRepository, tmdbRepository)
		request := dto_favorite.FavoriteRequestDto{
			ProfileID:         999,
			ContentExternalID: 101,
			Type:              "movie",
		}

		favoriteRepository.On("CreateFavoriteByProfileID", ctx, uint(10), request).
			Return(gorm.ErrRecordNotFound).
			Once()

		err := service.AddFavorite(ctx, 10, request)

		assert.ErrorIs(t, err, shared_errors.ErrProfileNotFound)
		favoriteRepository.AssertExpectations(t)
	})

	t.Run("should return favorite already exists when favorite is duplicated", func(t *testing.T) {
		ctx := context.Background()
		favoriteRepository := new(repository_mock.FavoriteRepositoryMock)
		tmdbRepository := new(repository_mock.TMDBRepositoryMock)
		service := service_favorite.NewFavoriteService(favoriteRepository, tmdbRepository)
		request := dto_favorite.FavoriteRequestDto{
			ProfileID:         20,
			ContentExternalID: 101,
			Type:              "movie",
		}

		favoriteRepository.On("CreateFavoriteByProfileID", ctx, uint(10), request).
			Return(gorm.ErrDuplicatedKey).
			Once()

		err := service.AddFavorite(ctx, 10, request)

		assert.ErrorIs(t, err, shared_errors.ErrFavoriteAlreadyExists)
		favoriteRepository.AssertExpectations(t)
	})

	t.Run("should return repository error when adding favorite fails", func(t *testing.T) {
		ctx := context.Background()
		favoriteRepository := new(repository_mock.FavoriteRepositoryMock)
		tmdbRepository := new(repository_mock.TMDBRepositoryMock)
		service := service_favorite.NewFavoriteService(favoriteRepository, tmdbRepository)
		expectedErr := errors.New("database unavailable")
		request := dto_favorite.FavoriteRequestDto{
			ProfileID:         20,
			ContentExternalID: 101,
			Type:              "movie",
		}

		favoriteRepository.On("CreateFavoriteByProfileID", ctx, uint(10), request).
			Return(expectedErr).
			Once()

		err := service.AddFavorite(ctx, 10, request)

		assert.ErrorIs(t, err, expectedErr)
		favoriteRepository.AssertExpectations(t)
	})
}

func TestDeleteFavorite(t *testing.T) {
	t.Run("should delete favorite", func(t *testing.T) {
		ctx := context.Background()
		favoriteRepository := new(repository_mock.FavoriteRepositoryMock)
		tmdbRepository := new(repository_mock.TMDBRepositoryMock)
		service := service_favorite.NewFavoriteService(favoriteRepository, tmdbRepository)
		request := dto_favorite.FavoriteRequestDto{
			ProfileID:         20,
			ContentExternalID: 101,
			Type:              "movie",
		}

		favoriteRepository.On("DeleteFavoriteByProfileID", ctx, uint(10), request).
			Return(nil).
			Once()

		err := service.DeleteFavorite(ctx, 10, request)

		require.NoError(t, err)
		tmdbRepository.AssertNotCalled(t, "GetContentByID", mock.Anything, mock.Anything, mock.Anything, mock.Anything)
		favoriteRepository.AssertExpectations(t)
	})

	t.Run("should return profile not found when profile does not exist", func(t *testing.T) {
		ctx := context.Background()
		favoriteRepository := new(repository_mock.FavoriteRepositoryMock)
		tmdbRepository := new(repository_mock.TMDBRepositoryMock)
		service := service_favorite.NewFavoriteService(favoriteRepository, tmdbRepository)
		request := dto_favorite.FavoriteRequestDto{
			ProfileID:         999,
			ContentExternalID: 101,
			Type:              "movie",
		}

		favoriteRepository.On("DeleteFavoriteByProfileID", ctx, uint(10), request).
			Return(gorm.ErrRecordNotFound).
			Once()

		err := service.DeleteFavorite(ctx, 10, request)

		assert.ErrorIs(t, err, shared_errors.ErrProfileNotFound)
		favoriteRepository.AssertExpectations(t)
	})

	t.Run("should return favorite not found when favorite does not exist", func(t *testing.T) {
		ctx := context.Background()
		favoriteRepository := new(repository_mock.FavoriteRepositoryMock)
		tmdbRepository := new(repository_mock.TMDBRepositoryMock)
		service := service_favorite.NewFavoriteService(favoriteRepository, tmdbRepository)
		request := dto_favorite.FavoriteRequestDto{
			ProfileID:         20,
			ContentExternalID: 999,
			Type:              "movie",
		}

		favoriteRepository.On("DeleteFavoriteByProfileID", ctx, uint(10), request).
			Return(shared_errors.ErrFavoriteNotFound).
			Once()

		err := service.DeleteFavorite(ctx, 10, request)

		assert.ErrorIs(t, err, shared_errors.ErrFavoriteNotFound)
		favoriteRepository.AssertExpectations(t)
	})

	t.Run("should return repository error when deleting favorite fails", func(t *testing.T) {
		ctx := context.Background()
		favoriteRepository := new(repository_mock.FavoriteRepositoryMock)
		tmdbRepository := new(repository_mock.TMDBRepositoryMock)
		service := service_favorite.NewFavoriteService(favoriteRepository, tmdbRepository)
		expectedErr := errors.New("database unavailable")
		request := dto_favorite.FavoriteRequestDto{
			ProfileID:         20,
			ContentExternalID: 101,
			Type:              "movie",
		}

		favoriteRepository.On("DeleteFavoriteByProfileID", ctx, uint(10), request).
			Return(expectedErr).
			Once()

		err := service.DeleteFavorite(ctx, 10, request)

		assert.ErrorIs(t, err, expectedErr)
		favoriteRepository.AssertExpectations(t)
	})
}
