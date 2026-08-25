package repository_mock

import (
	"context"

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
