package service_mock

import (
	"context"

	"github.com/stretchr/testify/mock"
)

type SummaryServiceMock struct {
	mock.Mock
}

func (m *ProfileServiceMock) ListFavorites(ctx context.Context, userID, profileID uint, page, perPage int) error {
	args := m.Called(ctx, userID, profileID, page, perPage)
	return args.Error(0)
}
