package service_mock

import (
	"context"

	dto_profile "github.com/KaueTTS/streaming_api/src/api/v1/dto/profile"
	"github.com/stretchr/testify/mock"
)

type ProfileServiceMock struct {
	mock.Mock
}

func (m *ProfileServiceMock) ListProfiles(ctx context.Context, userID uint, page, perPage int) (dto_profile.ProfileResponseDto, error) {
	args := m.Called(ctx, userID, page, perPage)
	return args.Get(0).(dto_profile.ProfileResponseDto), args.Error(1)
}

func (m *ProfileServiceMock) CreateProfile(ctx context.Context, userID uint, request dto_profile.ProfileRequestDto) (dto_profile.ProfileDto, error) {
	args := m.Called(ctx, userID, request)
	return args.Get(0).(dto_profile.ProfileDto), args.Error(1)
}

func (m *ProfileServiceMock) UpdateProfile(ctx context.Context, userID, profileID uint, request dto_profile.ProfileRequestDto) (dto_profile.ProfileDto, error) {
	args := m.Called(ctx, userID, profileID, request)
	return args.Get(0).(dto_profile.ProfileDto), args.Error(1)
}

func (m *ProfileServiceMock) DeleteProfile(ctx context.Context, userID, profileID uint) error {
	args := m.Called(ctx, userID, profileID)
	return args.Error(0)
}
