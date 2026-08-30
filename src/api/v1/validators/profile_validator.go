package validator

import (
	"strconv"
	"strings"
	"unicode/utf8"

	dto_profile "github.com/KaueTTS/streaming_api/src/api/v1/dto/profile"
	dto_shared "github.com/KaueTTS/streaming_api/src/api/v1/dto/shared"
	shared_constants "github.com/KaueTTS/streaming_api/src/shared/constants"
	shared_errors "github.com/KaueTTS/streaming_api/src/shared/errors"
)

// ValidateProfileRequest valida o dto de perfil
func ValidateProfileRequest(request dto_profile.ProfileRequestDto) []dto_shared.DetailErrorDto {
	var details []dto_shared.DetailErrorDto

	name := strings.TrimSpace(request.Name)
	if name == "" {
		details = append(
			details,
			NewDetail(shared_constants.Name, name, shared_errors.NameRequired),
		)
	} else if utf8.RuneCountInString(name) < shared_constants.MinNameLength {
		details = append(
			details,
			NewDetail(shared_constants.Name, name, shared_errors.NameMustLeast2Character),
		)
	} else if utf8.RuneCountInString(name) > shared_constants.MaxNameLength {
		details = append(
			details,
			NewDetail(shared_constants.Name, name, shared_errors.NameMustMaximum120Character),
		)
	}

	return details
}

// ValidateProfilePagination valida os parâmetros de paginação da listagem de perfis.
func ValidateProfilePagination(pagination dto_shared.PaginationDto) []dto_shared.DetailErrorDto {
	var details []dto_shared.DetailErrorDto

	if pagination.Page < 0 {
		details = append(
			details,
			NewDetail(shared_constants.Page, strconv.Itoa(pagination.Page), shared_errors.PageMustBePositive),
		)
	}

	if pagination.PerPage < 0 {
		details = append(
			details,
			NewDetail(shared_constants.PerPage, strconv.Itoa(pagination.PerPage), shared_errors.PerPageMustBePositive),
		)
	}

	return details
}
