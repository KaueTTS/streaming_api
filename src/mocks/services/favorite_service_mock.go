package service_mock

import (
	"context"

	dto_favorite "github.com/KaueTTS/streaming_api/src/api/v1/dto/favorite"
	"github.com/stretchr/testify/mock"
)

type FavoriteServiceMock struct {
	mock.Mock
}

func (m *FavoriteServiceMock) ListFavorites(ctx context.Context, userID, profileID uint, page, perPage int, language string) (dto_favorite.FavoriteResponseDto, error) {
	args := m.Called(ctx, userID, profileID, page, perPage, language)
	return args.Get(0).(dto_favorite.FavoriteResponseDto), args.Error(1)
}

func (m *FavoriteServiceMock) AddFavorite(ctx context.Context, userID uint, request dto_favorite.FavoriteRequestDto) error {
	args := m.Called(ctx, userID, request)
	return args.Error(0)
}

func (m *FavoriteServiceMock) DeleteFavorite(ctx context.Context, userID uint, request dto_favorite.FavoriteRequestDto) error {
	args := m.Called(ctx, userID, request)
	return args.Error(0)
}
