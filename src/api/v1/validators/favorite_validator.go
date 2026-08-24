package validator

import (
	"strconv"

	dto_favorite "github.com/KaueTTS/streaming_api/src/api/v1/dto/favorite"
	dto_shared "github.com/KaueTTS/streaming_api/src/api/v1/dto/shared"
	shared_constants "github.com/KaueTTS/streaming_api/src/shared/constants"
	shared_errors_favorite "github.com/KaueTTS/streaming_api/src/shared/errors/favorite"
)

// ValidateFavoriteRequest valida o dto de criação de favorito.
func ValidateFavoriteRequest(request dto_favorite.FavoriteRequestDto) []dto_shared.DetailErrorDto {
	var details []dto_shared.DetailErrorDto

	if request.ProfileID == 0 {
		details = append(
			details,
			NewDetail(shared_constants.ProfileID, strconv.Itoa(int(request.ProfileID)), shared_errors_favorite.InvalidProfileID),
		)
	}

	if request.ContentExternalID <= 0 {
		details = append(
			details,
			NewDetail(shared_constants.ContentExternalID, strconv.Itoa(request.ContentExternalID), shared_errors_favorite.InvalidContentExternalID),
		)
	}

	if request.Type != "tv" && request.Type != "movie" {
		details = append(
			details,
			NewDetail(shared_constants.Type, request.Type, shared_errors_favorite.InvalidContentType),
		)
	}

	return details
}
