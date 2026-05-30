package validator

import (
	dto_shared "github.com/KaueTTS/streaming_api/src/api/v1/dto/shared"
)

func NewDetail(field string, value string, message string) dto_shared.DetailErrorDto {
	return dto_shared.DetailErrorDto{
		Field:   field,
		Value:   value,
		Message: message,
	}
}

func ValidatePagination(pagination dto_shared.PaginationDto) (page int, perPage int) {
	if pagination.Page <= 0 {
		pagination.Page = 1
	}

	if pagination.PerPage <= 0 {
		pagination.PerPage = 10
	} else if pagination.PerPage > 100 {
		pagination.PerPage = 100
	}

	return pagination.Page, pagination.PerPage
}
