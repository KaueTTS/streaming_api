package repository_mock

import (
	"context"

	models "github.com/KaueTTS/streaming_api/src/models"
	"github.com/stretchr/testify/mock"
)

type ProfileRepositoryMock struct {
	mock.Mock
}

func (m *ProfileRepositoryMock) FindByUserID(ctx context.Context, userID uint, page, perPage int) ([]models.Profile, int64, error) {
	args := m.Called(ctx, userID, page, perPage)
	return args.Get(0).([]models.Profile), args.Get(1).(int64), args.Error(2)
}

func (m *ProfileRepositoryMock) CountByUserID(ctx context.Context, userID uint) (int64, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).(int64), args.Error(1)
}

func (m *ProfileRepositoryMock) Create(ctx context.Context, profile *models.Profile) error {
	args := m.Called(ctx, profile)
	return args.Error(0)
}

func (m *ProfileRepositoryMock) Update(ctx context.Context, profile *models.Profile) error {
	args := m.Called(ctx, profile)
	return args.Error(0)
}

func (m *ProfileRepositoryMock) Delete(ctx context.Context, userID, profileID uint) error {
	args := m.Called(ctx, userID, profileID)
	return args.Error(0)
}

func (m *ProfileRepositoryMock) FindProfileByID(ctx context.Context, profileID uint) (*models.Profile, error) {
	args := m.Called(ctx, profileID)
	profile, _ := args.Get(0).(*models.Profile)
	return profile, args.Error(1)
}

func (m *ProfileRepositoryMock) FindProfileByUserIDAndID(ctx context.Context, userID, profileID uint) (*models.Profile, error) {
	args := m.Called(ctx, userID, profileID)
	profile, _ := args.Get(0).(*models.Profile)
	return profile, args.Error(1)
}
