package service_interface

import (
	"context"

	dto_profile "github.com/KaueTTS/streaming_api/src/api/v1/dto/profile"
)

type ProfileServiceInterface interface {
	GetProfiles(ctx context.Context, userID uint, role string) ([]dto_profile.ProfileResponseDto, error)
}
