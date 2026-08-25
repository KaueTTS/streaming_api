package service_mock

import (
	"context"

	dto_content "github.com/KaueTTS/streaming_api/src/api/v1/dto/content"
	"github.com/stretchr/testify/mock"
)

type ContentServiceMock struct {
	mock.Mock
}

func (m *ContentServiceMock) ListContents(ctx context.Context, userID uint, request dto_content.ContentListRequestDto) (dto_content.ContentResponseDto, error) {
	args := m.Called(ctx, userID, request)
	return args.Get(0).(dto_content.ContentResponseDto), args.Error(1)
}

func (m *ContentServiceMock) SearchContents(ctx context.Context, userID uint, request dto_content.ContentSearchRequestDto) (dto_content.ContentResponseDto, error) {
	args := m.Called(ctx, userID, request)
	return args.Get(0).(dto_content.ContentResponseDto), args.Error(1)
}
