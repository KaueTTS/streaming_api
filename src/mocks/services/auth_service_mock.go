package service_mock

import (
	"context"

	dto_auth "github.com/KaueTTS/streaming_api/src/api/v1/dto/auth"
	"github.com/stretchr/testify/mock"
)

type AuthServiceMock struct {
	mock.Mock
}

func (m *AuthServiceMock) Register(ctx context.Context, request dto_auth.RegisterRequestDto) (dto_auth.UserResponseDto, error) {
	args := m.Called(ctx, request)
	return args.Get(0).(dto_auth.UserResponseDto), args.Error(1)
}

func (m *AuthServiceMock) Login(ctx context.Context, request dto_auth.LoginRequestDto) (dto_auth.AuthResponseDto, error) {
	args := m.Called(ctx, request)
	return args.Get(0).(dto_auth.AuthResponseDto), args.Error(1)
}

func (m *AuthServiceMock) Me(ctx context.Context, userID uint) (dto_auth.UserResponseDto, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).(dto_auth.UserResponseDto), args.Error(1)
}
