package service_profile

import (
	"context"

	dto_profile "github.com/KaueTTS/streaming_api/src/api/v1/dto/profile"
	models "github.com/KaueTTS/streaming_api/src/models"
	repository_interface "github.com/KaueTTS/streaming_api/src/repositories/interfaces"
	shared_constants "github.com/KaueTTS/streaming_api/src/shared/constants"
)

type ProfileService struct {
	ProfileRepositoryInterface repository_interface.ProfileRepositoryInterface
}

func NewProfileService(profileRepositoryInterface repository_interface.ProfileRepositoryInterface) *ProfileService {
	return &ProfileService{
		ProfileRepositoryInterface: profileRepositoryInterface,
	}
}

func (s *ProfileService) GetProfiles(ctx context.Context, userID uint, role string) ([]dto_profile.ProfileResponseDto, error) {
	var profiles []models.Profile
	var err error

	if role == shared_constants.RoleAdmin {
		profiles, err = s.ProfileRepositoryInterface.FindAll(ctx)
	} else {
		profiles, err = s.ProfileRepositoryInterface.FindByUserID(ctx, userID)
	}

	if err != nil {
		return nil, err
	}

	response := make([]dto_profile.ProfileResponseDto, 0, len(profiles))

	for _, profile := range profiles {
		response = append(response, dto_profile.ProfileResponseDto{
			ID:        profile.ID,
			UserID:    profile.UserID,
			Name:      profile.Name,
			AvatarURL: profile.AvatarURL,
			IsKids:    profile.IsKids,
			CreatedAt: profile.CreatedAt,
			UpdatedAt: profile.UpdatedAt,
		})
	}

	return response, nil
}
