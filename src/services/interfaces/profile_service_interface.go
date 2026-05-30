package service_interface

import (
	"context"

	dto_profile "github.com/KaueTTS/streaming_api/src/api/v1/dto/profile"
)

type ProfileServiceInterface interface {
	ListProfiles(ctx context.Context, userID uint, page, perPage int) (dto_profile.ProfileResponseDto, error)
	CreateProfile(ctx context.Context, userID uint, request dto_profile.ProfileRequestDto) (dto_profile.ProfileDto, error)
	UpdateProfile(ctx context.Context, userID, profileID uint, request dto_profile.ProfileRequestDto) (dto_profile.ProfileDto, error)
}
