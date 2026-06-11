package validator

import (
	"regexp"
	"strconv"
	"strings"

	dto_content "github.com/KaueTTS/streaming_api/src/api/v1/dto/content"
	dto_shared "github.com/KaueTTS/streaming_api/src/api/v1/dto/shared"
	shared_constants "github.com/KaueTTS/streaming_api/src/shared/constants"
	shared_constants_content "github.com/KaueTTS/streaming_api/src/shared/constants/content"
	shared_errors "github.com/KaueTTS/streaming_api/src/shared/errors"
	shared_errors_content "github.com/KaueTTS/streaming_api/src/shared/errors/content"
	shared_normalizers "github.com/KaueTTS/streaming_api/src/shared/normalizers"
)

var (
	contentLanguageRegex = regexp.MustCompile(`^[a-zA-Z]{2}([-_][a-zA-Z]{2})?$`)
	contentGenresRegex   = regexp.MustCompile(`^\d+([,|]\d+)*$`)
)

var contentSortByOptions = map[string]map[string]struct{}{
	shared_constants_content.ContentTypeMovie: {
		"original_title.asc":        {},
		"original_title.desc":       {},
		"popularity.asc":            {},
		"popularity.desc":           {},
		"primary_release_date.asc":  {},
		"primary_release_date.desc": {},
		"revenue.asc":               {},
		"revenue.desc":              {},
		"title.asc":                 {},
		"title.desc":                {},
		"vote_average.asc":          {},
		"vote_average.desc":         {},
		"vote_count.asc":            {},
		"vote_count.desc":           {},
	},
	shared_constants_content.ContentTypeTV: {
		"first_air_date.asc":  {},
		"first_air_date.desc": {},
		"name.asc":            {},
		"name.desc":           {},
		"original_name.asc":   {},
		"original_name.desc":  {},
		"popularity.asc":      {},
		"popularity.desc":     {},
		"vote_average.asc":    {},
		"vote_average.desc":   {},
		"vote_count.asc":      {},
		"vote_count.desc":     {},
	},
}

func ValidateContentListRequest(request dto_content.ContentListRequestDto) []dto_shared.DetailErrorDto {
	var details []dto_shared.DetailErrorDto

	details = append(details, validateContentType(request.Type)...)
	details = append(details, validateContentPage(request.Page)...)
	details = append(details, validateContentSortBy(request.Type, request.SortBy)...)
	details = append(details, validateContentGenres(request.WithGenres)...)
	details = append(details, validateContentLanguage(request.Language)...)
	details = append(details, validateContentYear(request.Year)...)

	return details
}

func ValidateContentSearchRequest(request dto_content.ContentSearchRequestDto) []dto_shared.DetailErrorDto {
	var details []dto_shared.DetailErrorDto

	query := strings.TrimSpace(request.Query)
	if query == "" {
		details = append(details, NewDetail(shared_constants.Query, request.Query, shared_errors_content.SearchQueryRequired))
	}

	details = append(details, validateContentType(request.Type)...)
	details = append(details, validateContentPage(request.Page)...)
	details = append(details, validateContentLanguage(request.Language)...)

	return details
}

func validateContentType(contentType string) []dto_shared.DetailErrorDto {
	normalizedType := shared_normalizers.NormalizeString(contentType)
	if normalizedType == "" {
		return []dto_shared.DetailErrorDto{
			NewDetail(shared_constants.Type, contentType, shared_errors_content.ContentTypeRequired),
		}
	}

	switch normalizedType {
	case shared_constants_content.ContentTypeMovie,
		shared_constants_content.ContentTypeTV:
		return nil
	default:
		return []dto_shared.DetailErrorDto{
			NewDetail(shared_constants.Type, contentType, shared_errors_content.InvalidContentType),
		}
	}
}

func validateContentPage(page int) []dto_shared.DetailErrorDto {
	if page < 0 {
		return []dto_shared.DetailErrorDto{
			NewDetail(shared_constants.Page, "", shared_errors.PageMustBePositive),
		}
	}

	return nil
}

func validateContentSortBy(contentType string, sortBy string) []dto_shared.DetailErrorDto {
	normalizedSortBy := shared_normalizers.NormalizeString(sortBy)
	if normalizedSortBy == "" {
		return nil
	}

	allowedSortBy, ok := contentSortByOptions[shared_normalizers.NormalizeString(contentType)]
	if !ok {
		return nil
	}

	if _, ok := allowedSortBy[normalizedSortBy]; !ok {
		return []dto_shared.DetailErrorDto{
			NewDetail(shared_constants.SortBy, sortBy, shared_errors_content.InvalidContentSortBy),
		}
	}

	return nil
}

func validateContentGenres(withGenres string) []dto_shared.DetailErrorDto {
	normalizedGenres := strings.TrimSpace(withGenres)
	if normalizedGenres == "" {
		return nil
	}

	if !contentGenresRegex.MatchString(normalizedGenres) {
		return []dto_shared.DetailErrorDto{
			NewDetail(shared_constants.WithGenres, withGenres, shared_errors_content.InvalidContentGenres),
		}
	}

	return nil
}

func validateContentLanguage(language string) []dto_shared.DetailErrorDto {
	normalizedLanguage := strings.TrimSpace(language)
	if normalizedLanguage == "" {
		return nil
	}

	if !contentLanguageRegex.MatchString(normalizedLanguage) {
		return []dto_shared.DetailErrorDto{
			NewDetail(shared_constants.Language, language, shared_errors_content.InvalidContentLanguage),
		}
	}

	return nil
}

func validateContentYear(year int) []dto_shared.DetailErrorDto {
	if year < 0 {
		return []dto_shared.DetailErrorDto{
			NewDetail(shared_constants.Year, strconv.Itoa(year), shared_errors_content.InvalidContentYear),
		}
	}

	return nil
}
