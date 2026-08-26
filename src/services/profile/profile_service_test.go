package service_profile_test

import (
	"context"
	"errors"
	"testing"
	"time"

	dto_profile "github.com/KaueTTS/streaming_api/src/api/v1/dto/profile"
	repository_mock "github.com/KaueTTS/streaming_api/src/mocks/repositories"
	models "github.com/KaueTTS/streaming_api/src/models"
	service_profile "github.com/KaueTTS/streaming_api/src/services/profile"
	shared_constants_profile "github.com/KaueTTS/streaming_api/src/shared/constants/profile"
	shared_errors "github.com/KaueTTS/streaming_api/src/shared/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func TestListProfiles(t *testing.T) {
	t.Run("should list profiles with pagination", func(t *testing.T) {
		ctx := context.Background()
		repository := new(repository_mock.ProfileRepositoryMock)
		service := service_profile.NewProfileService(repository)
		avatarURL := "https://example.com/avatar.png"
		createdAt := time.Date(2026, 1, 1, 10, 0, 0, 0, time.UTC)
		updatedAt := time.Date(2026, 1, 1, 11, 0, 0, 0, time.UTC)
		profiles := []models.Profile{
			{
				ID:        1,
				UserID:    10,
				Name:      "Main",
				AvatarURL: &avatarURL,
				IsKids:    false,
				CreatedAt: createdAt,
				UpdatedAt: updatedAt,
			},
			{
				ID:        2,
				UserID:    10,
				Name:      "Kids",
				IsKids:    true,
				CreatedAt: createdAt,
				UpdatedAt: updatedAt,
			},
		}

		repository.On("FindByUserID", ctx, uint(10), 2, 2).
			Return(profiles, int64(5), nil).
			Once()

		response, err := service.ListProfiles(ctx, 10, 2, 2)

		require.NoError(t, err)
		require.Len(t, response.Data, 2)
		assert.Equal(t, dto_profile.ProfileDto{
			ID:        1,
			UserID:    10,
			Name:      "Main",
			AvatarURL: &avatarURL,
			IsKids:    false,
			CreatedAt: createdAt,
			UpdatedAt: updatedAt,
		}, response.Data[0])
		assert.Equal(t, dto_profile.ProfileDto{
			ID:        2,
			UserID:    10,
			Name:      "Kids",
			IsKids:    true,
			CreatedAt: createdAt,
			UpdatedAt: updatedAt,
		}, response.Data[1])
		assert.Equal(t, 2, response.Pagination.Page)
		assert.Equal(t, 2, response.Pagination.PerPage)
		assert.Equal(t, 3, response.Pagination.PageCount)
		assert.Equal(t, int64(5), response.Pagination.Total)
		repository.AssertExpectations(t)
	})

	t.Run("should return repository error when listing profiles fails", func(t *testing.T) {
		ctx := context.Background()
		repository := new(repository_mock.ProfileRepositoryMock)
		service := service_profile.NewProfileService(repository)
		expectedErr := errors.New("database unavailable")

		repository.On("FindByUserID", ctx, uint(10), 1, 10).
			Return([]models.Profile(nil), int64(0), expectedErr).
			Once()

		response, err := service.ListProfiles(ctx, 10, 1, 10)

		assert.Empty(t, response)
		assert.ErrorIs(t, err, expectedErr)
		repository.AssertExpectations(t)
	})
}

func TestCreateProfile(t *testing.T) {
	t.Run("should create profile", func(t *testing.T) {
		ctx := context.Background()
		repository := new(repository_mock.ProfileRepositoryMock)
		service := service_profile.NewProfileService(repository)
		avatarURL := "https://example.com/avatar.png"
		createdAt := time.Date(2026, 1, 1, 10, 0, 0, 0, time.UTC)
		updatedAt := time.Date(2026, 1, 1, 11, 0, 0, 0, time.UTC)
		request := dto_profile.ProfileRequestDto{
			Name:      "  Main  ",
			AvatarURL: &avatarURL,
			IsKids:    false,
		}

		repository.On("CountByUserID", ctx, uint(10)).
			Return(int64(1), nil).
			Once()
		repository.On("Create", ctx, mock.MatchedBy(func(profile *models.Profile) bool {
			return profile.UserID == 10 &&
				profile.Name == "Main" &&
				profile.AvatarURL == &avatarURL &&
				!profile.IsKids
		})).
			Run(func(args mock.Arguments) {
				profile := args.Get(1).(*models.Profile)
				profile.ID = 1
				profile.CreatedAt = createdAt
				profile.UpdatedAt = updatedAt
			}).
			Return(nil).
			Once()

		response, err := service.CreateProfile(ctx, 10, request)

		require.NoError(t, err)
		assert.Equal(t, dto_profile.ProfileDto{
			ID:        1,
			UserID:    10,
			Name:      "Main",
			AvatarURL: &avatarURL,
			IsKids:    false,
			CreatedAt: createdAt,
			UpdatedAt: updatedAt,
		}, response)
		repository.AssertExpectations(t)
	})

	t.Run("should return profile limit reached when user already has maximum profiles", func(t *testing.T) {
		ctx := context.Background()
		repository := new(repository_mock.ProfileRepositoryMock)
		service := service_profile.NewProfileService(repository)
		request := dto_profile.ProfileRequestDto{
			Name:   "Main",
			IsKids: false,
		}

		repository.On("CountByUserID", ctx, uint(10)).
			Return(int64(shared_constants_profile.MaxProfilesPerUser), nil).
			Once()

		response, err := service.CreateProfile(ctx, 10, request)

		assert.Empty(t, response)
		assert.ErrorIs(t, err, shared_errors.ErrProfileLimitReached)
		repository.AssertNotCalled(t, "Create", mock.Anything, mock.Anything)
		repository.AssertExpectations(t)
	})

	t.Run("should return repository error when counting profiles fails", func(t *testing.T) {
		ctx := context.Background()
		repository := new(repository_mock.ProfileRepositoryMock)
		service := service_profile.NewProfileService(repository)
		expectedErr := errors.New("database unavailable")
		request := dto_profile.ProfileRequestDto{
			Name:   "Main",
			IsKids: false,
		}

		repository.On("CountByUserID", ctx, uint(10)).
			Return(int64(0), expectedErr).
			Once()

		response, err := service.CreateProfile(ctx, 10, request)

		assert.Empty(t, response)
		assert.ErrorIs(t, err, expectedErr)
		repository.AssertNotCalled(t, "Create", mock.Anything, mock.Anything)
		repository.AssertExpectations(t)
	})

	t.Run("should return repository error when creating profile fails", func(t *testing.T) {
		ctx := context.Background()
		repository := new(repository_mock.ProfileRepositoryMock)
		service := service_profile.NewProfileService(repository)
		expectedErr := errors.New("database unavailable")
		request := dto_profile.ProfileRequestDto{
			Name:   "Main",
			IsKids: false,
		}

		repository.On("CountByUserID", ctx, uint(10)).
			Return(int64(1), nil).
			Once()
		repository.On("Create", ctx, mock.AnythingOfType("*models.Profile")).
			Return(expectedErr).
			Once()

		response, err := service.CreateProfile(ctx, 10, request)

		assert.Empty(t, response)
		assert.ErrorIs(t, err, expectedErr)
		repository.AssertExpectations(t)
	})
}

func TestUpdateProfile(t *testing.T) {
	t.Run("should update profile", func(t *testing.T) {
		ctx := context.Background()
		repository := new(repository_mock.ProfileRepositoryMock)
		service := service_profile.NewProfileService(repository)
		avatarURL := "https://example.com/avatar.png"
		createdAt := time.Date(2026, 1, 1, 10, 0, 0, 0, time.UTC)
		updatedAt := time.Date(2026, 1, 1, 11, 0, 0, 0, time.UTC)
		request := dto_profile.ProfileRequestDto{
			Name:      "  Updated Main  ",
			AvatarURL: &avatarURL,
			IsKids:    true,
		}

		repository.On("Update", ctx, mock.MatchedBy(func(profile *models.Profile) bool {
			return profile.ID == 2 &&
				profile.UserID == 10 &&
				profile.Name == "Updated Main" &&
				profile.AvatarURL == &avatarURL &&
				profile.IsKids
		})).
			Run(func(args mock.Arguments) {
				profile := args.Get(1).(*models.Profile)
				profile.CreatedAt = createdAt
				profile.UpdatedAt = updatedAt
			}).
			Return(nil).
			Once()

		response, err := service.UpdateProfile(ctx, 10, 2, request)

		require.NoError(t, err)
		assert.Equal(t, dto_profile.ProfileDto{
			ID:        2,
			UserID:    10,
			Name:      "Updated Main",
			AvatarURL: &avatarURL,
			IsKids:    true,
			CreatedAt: createdAt,
			UpdatedAt: updatedAt,
		}, response)
		repository.AssertExpectations(t)
	})

	t.Run("should return profile not found when profile does not exist", func(t *testing.T) {
		ctx := context.Background()
		repository := new(repository_mock.ProfileRepositoryMock)
		service := service_profile.NewProfileService(repository)
		request := dto_profile.ProfileRequestDto{
			Name:   "Updated Main",
			IsKids: true,
		}

		repository.On("Update", ctx, mock.AnythingOfType("*models.Profile")).
			Return(gorm.ErrRecordNotFound).
			Once()

		response, err := service.UpdateProfile(ctx, 10, 999, request)

		assert.Empty(t, response)
		assert.ErrorIs(t, err, shared_errors.ErrProfileNotFound)
		repository.AssertExpectations(t)
	})

	t.Run("should return repository error when updating profile fails", func(t *testing.T) {
		ctx := context.Background()
		repository := new(repository_mock.ProfileRepositoryMock)
		service := service_profile.NewProfileService(repository)
		expectedErr := errors.New("database unavailable")
		request := dto_profile.ProfileRequestDto{
			Name:   "Updated Main",
			IsKids: true,
		}

		repository.On("Update", ctx, mock.AnythingOfType("*models.Profile")).
			Return(expectedErr).
			Once()

		response, err := service.UpdateProfile(ctx, 10, 2, request)

		assert.Empty(t, response)
		assert.ErrorIs(t, err, expectedErr)
		repository.AssertExpectations(t)
	})
}

func TestDeleteProfile(t *testing.T) {
	t.Run("should delete profile", func(t *testing.T) {
		ctx := context.Background()
		repository := new(repository_mock.ProfileRepositoryMock)
		service := service_profile.NewProfileService(repository)

		repository.On("Delete", ctx, uint(10), uint(2)).
			Return(nil).
			Once()

		err := service.DeleteProfile(ctx, 10, 2)

		require.NoError(t, err)
		repository.AssertExpectations(t)
	})

	t.Run("should return profile not found when profile does not exist", func(t *testing.T) {
		ctx := context.Background()
		repository := new(repository_mock.ProfileRepositoryMock)
		service := service_profile.NewProfileService(repository)

		repository.On("Delete", ctx, uint(10), uint(999)).
			Return(gorm.ErrRecordNotFound).
			Once()

		err := service.DeleteProfile(ctx, 10, 999)

		assert.ErrorIs(t, err, shared_errors.ErrProfileNotFound)
		repository.AssertExpectations(t)
	})

	t.Run("should return repository error when deleting profile fails", func(t *testing.T) {
		ctx := context.Background()
		repository := new(repository_mock.ProfileRepositoryMock)
		service := service_profile.NewProfileService(repository)
		expectedErr := errors.New("database unavailable")

		repository.On("Delete", ctx, uint(10), uint(2)).
			Return(expectedErr).
			Once()

		err := service.DeleteProfile(ctx, 10, 2)

		assert.ErrorIs(t, err, expectedErr)
		repository.AssertExpectations(t)
	})
}
