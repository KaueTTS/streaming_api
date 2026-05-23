package service_interface

import (
	"context"

	dto_auth "github.com/KaueTTS/streaming_api/src/api/v1/dto/auth"
)

type AuthServiceInterface interface {
	Register(ctx context.Context, request dto_auth.RegisterRequestDto) (dto_auth.UserResponseDto, error)
	Login(ctx context.Context, request dto_auth.LoginRequestDto) (dto_auth.AuthResponseDto, error)
	Me(ctx context.Context, userID uint) (dto_auth.UserResponseDto, error)
}
