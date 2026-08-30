package repository_mock

import (
	"context"

	dto "github.com/KaueTTS/streaming_api/src/repositories/http/tmdb/dto"
	"github.com/stretchr/testify/mock"
)

type TMDBRepositoryMock struct {
	mock.Mock
}

func (m *TMDBRepositoryMock) ListContents(ctx context.Context, filters dto.ContentFiltersDto) (dto.GetContentResponseDto, error) {
	args := m.Called(ctx, filters)
	return args.Get(0).(dto.GetContentResponseDto), args.Error(1)
}

func (m *TMDBRepositoryMock) SearchContents(ctx context.Context, filters dto.ContentFiltersDto) (dto.GetContentResponseDto, error) {
	args := m.Called(ctx, filters)
	return args.Get(0).(dto.GetContentResponseDto), args.Error(1)
}

func (m *TMDBRepositoryMock) GetContentByID(ctx context.Context, contentType string, contentExternalID int, language string) (dto.ContentDto, error) {
	args := m.Called(ctx, contentType, contentExternalID, language)
	return args.Get(0).(dto.ContentDto), args.Error(1)
}
