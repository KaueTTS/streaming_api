package validator

import (
	dto_shared "github.com/KaueTTS/streaming_api/src/api/v1/dto/shared"
	shared_normalizers "github.com/KaueTTS/streaming_api/src/shared/normalizers"
)

func NewDetail(field string, value string, message string) dto_shared.DetailErrorDto {
	return dto_shared.DetailErrorDto{
		Field:   field,
		Value:   value,
		Message: message,
	}
}

func ValidatePagination(pagination dto_shared.PaginationDto) (page int, perPage int) {
	page = shared_normalizers.NormalizePage(pagination.Page)
	perPage = shared_normalizers.NormalizePerPage(pagination.PerPage)

	return page, perPage
}
