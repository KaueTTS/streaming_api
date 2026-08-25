package service_content_test

import (
	"context"
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

func TestListContentsUsesAuthenticatedUserProfile(t *testing.T) {
	ctx := context.Background()
	tmdbRepository := new(repository_mock.TMDBRepositoryMock)
	profileRepository := new(repository_mock.ProfileRepositoryMock)
	service := service_content.NewContentService(tmdbRepository, profileRepository)

	request := dto_content.ContentListRequestDto{
		Type:      " Movie ",
		Page:      0,
		Genre:     "28,12",
		ProfileID: 2,
	}
	profile := &models.Profile{ID: 2, UserID: 1, IsKids: true}
	tmdbResponse := tmdb_dto.GetContentResponseDto{Page: 1, TotalPages: 1}

	profileRepository.On("FindProfileByUserIDAndID", ctx, uint(1), uint(2)).Return(profile, nil).Once()
	tmdbRepository.On("ListContents", ctx, mock.MatchedBy(func(filters tmdb_dto.ContentFiltersDto) bool {
		return filters.Type == "movie" &&
			filters.Page == 1 &&
			filters.WithGenres == "28,12" &&
			filters.IsKids
	})).Return(tmdbResponse, nil).Once()

	response, err := service.ListContents(ctx, 1, request)

	require.NoError(t, err)
	assert.Equal(t, 1, response.Pagination.Page)
	profileRepository.AssertExpectations(t)
	tmdbRepository.AssertExpectations(t)
}

func TestSearchContentsUsesAuthenticatedUserProfile(t *testing.T) {
	ctx := context.Background()
	tmdbRepository := new(repository_mock.TMDBRepositoryMock)
	profileRepository := new(repository_mock.ProfileRepositoryMock)
	service := service_content.NewContentService(tmdbRepository, profileRepository)

	request := dto_content.ContentSearchRequestDto{
		Type:      "TV",
		Page:      2,
		Query:     " avatar ",
		ProfileID: 3,
	}
	profile := &models.Profile{ID: 3, UserID: 1, IsKids: true}
	tmdbResponse := tmdb_dto.GetContentResponseDto{Page: 2, TotalPages: 1}

	profileRepository.On("FindProfileByUserIDAndID", ctx, uint(1), uint(3)).Return(profile, nil).Once()
	tmdbRepository.On("SearchContents", ctx, mock.MatchedBy(func(filters tmdb_dto.ContentFiltersDto) bool {
		return filters.Type == "tv" &&
			filters.Page == 2 &&
			filters.Query == "avatar" &&
			filters.IsKids
	})).Return(tmdbResponse, nil).Once()

	response, err := service.SearchContents(ctx, 1, request)

	require.NoError(t, err)
	assert.Equal(t, 2, response.Pagination.Page)
	profileRepository.AssertExpectations(t)
	tmdbRepository.AssertExpectations(t)
}

func TestSearchContentsReturnsProfileNotFound(t *testing.T) {
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

	assert.ErrorIs(t, err, shared_errors.ErrProfileNotFound)
	assert.Empty(t, response.Data)
	tmdbRepository.AssertNotCalled(t, "SearchContents", mock.Anything, mock.Anything)
	profileRepository.AssertExpectations(t)
}
