package service_interface

import (
	"context"

	dto_profile "github.com/KaueTTS/streaming_api/src/api/v1/dto/profile"
)

type ProfileServiceInterface interface {
	ListProfiles(ctx context.Context, userID uint, page int, perPage int) (dto_profile.ProfileResponseDto, error)
	CreateProfile(ctx context.Context, userID uint, request dto_profile.CreateProfileRequestDto) (dto_profile.ProfileDto, error)
}
