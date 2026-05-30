package service_profile

import (
	"context"
	"errors"
	"strings"

	dto_profile "github.com/KaueTTS/streaming_api/src/api/v1/dto/profile"
	dto_shared "github.com/KaueTTS/streaming_api/src/api/v1/dto/shared"
	models "github.com/KaueTTS/streaming_api/src/models"
	repository_interface "github.com/KaueTTS/streaming_api/src/repositories/interfaces"
	shared_constants_profile "github.com/KaueTTS/streaming_api/src/shared/constants/profile"
	shared_errors "github.com/KaueTTS/streaming_api/src/shared/errors"
	"gorm.io/gorm"
)

type ProfileService struct {
	ProfileRepositoryInterface repository_interface.ProfileRepositoryInterface
}

func NewProfileService(profileRepositoryInterface repository_interface.ProfileRepositoryInterface) *ProfileService {
	return &ProfileService{
		ProfileRepositoryInterface: profileRepositoryInterface,
	}
}

func (s *ProfileService) ListProfiles(ctx context.Context, userID uint, page, perPage int) (dto_profile.ProfileResponseDto, error) {
	var profiles []models.Profile
	var total int64
	var err error

	profiles, total, err = s.ProfileRepositoryInterface.FindByUserID(ctx, userID, page, perPage)

	if err != nil {
		return dto_profile.ProfileResponseDto{}, err
	}

	data := make([]dto_profile.ProfileDto, 0, len(profiles))
	for _, profile := range profiles {
		data = append(data, dto_profile.ProfileDto{
			ID:        profile.ID,
			UserID:    profile.UserID,
			Name:      profile.Name,
			AvatarURL: profile.AvatarURL,
			IsKids:    profile.IsKids,
			CreatedAt: profile.CreatedAt,
			UpdatedAt: profile.UpdatedAt,
		})
	}

	pageCount := int((total + int64(perPage) - 1) / int64(perPage))

	return dto_profile.ProfileResponseDto{
		Data: data,
		Pagination: dto_shared.PaginationDto{
			Page:      page,
			PerPage:   perPage,
			PageCount: pageCount,
			Total:     total,
		},
	}, nil
}

func (s *ProfileService) CreateProfile(ctx context.Context, userID uint, request dto_profile.ProfileRequestDto) (dto_profile.ProfileDto, error) {
	total, err := s.ProfileRepositoryInterface.CountByUserID(ctx, userID)
	if err != nil {
		return dto_profile.ProfileDto{}, err
	}

	if total >= shared_constants_profile.MaxProfilesPerUser {
		return dto_profile.ProfileDto{}, shared_errors.ErrProfileLimitReached
	}

	profile := models.Profile{
		UserID:    userID,
		Name:      strings.TrimSpace(request.Name),
		AvatarURL: request.AvatarURL,
		IsKids:    request.IsKids,
	}

	if err := s.ProfileRepositoryInterface.Create(ctx, &profile); err != nil {
		return dto_profile.ProfileDto{}, err
	}

	return dto_profile.ProfileDto{
		ID:        profile.ID,
		UserID:    profile.UserID,
		Name:      profile.Name,
		AvatarURL: profile.AvatarURL,
		IsKids:    profile.IsKids,
		CreatedAt: profile.CreatedAt,
		UpdatedAt: profile.UpdatedAt,
	}, nil
}

func (s *ProfileService) UpdateProfile(ctx context.Context, userID, profileID uint, request dto_profile.ProfileRequestDto) (dto_profile.ProfileDto, error) {
	profile := models.Profile{
		ID:        profileID,
		UserID:    userID,
		Name:      strings.TrimSpace(request.Name),
		AvatarURL: request.AvatarURL,
		IsKids:    request.IsKids,
	}

	if err := s.ProfileRepositoryInterface.Update(ctx, &profile); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return dto_profile.ProfileDto{}, shared_errors.ErrProfileNotFound
		}

		return dto_profile.ProfileDto{}, err
	}

	return dto_profile.ProfileDto{
		ID:        profile.ID,
		UserID:    profile.UserID,
		Name:      profile.Name,
		AvatarURL: profile.AvatarURL,
		IsKids:    profile.IsKids,
		CreatedAt: profile.CreatedAt,
		UpdatedAt: profile.UpdatedAt,
	}, nil
}
