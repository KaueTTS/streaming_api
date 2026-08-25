package repository_mock

import (
	"context"

	dto_favorite "github.com/KaueTTS/streaming_api/src/api/v1/dto/favorite"
	models "github.com/KaueTTS/streaming_api/src/models"
	"github.com/stretchr/testify/mock"
)

type FavoriteRepositoryMock struct {
	mock.Mock
}

func (m *FavoriteRepositoryMock) FindFavoriteByProfileID(ctx context.Context, userID, profileID uint, page, perPage int) ([]models.Favorite, int64, error) {
	args := m.Called(ctx, userID, profileID, page, perPage)
	return args.Get(0).([]models.Favorite), args.Get(1).(int64), args.Error(2)
}

func (m *FavoriteRepositoryMock) CreateFavoriteByProfileID(ctx context.Context, userID uint, request dto_favorite.FavoriteRequestDto) error {
	args := m.Called(ctx, userID, request)
	return args.Error(0)
}

func (m *FavoriteRepositoryMock) DeleteFavoriteByProfileID(ctx context.Context, userID uint, request dto_favorite.FavoriteRequestDto) error {
	args := m.Called(ctx, userID, request)
	return args.Error(0)
}
