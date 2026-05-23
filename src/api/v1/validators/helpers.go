package validator

import dto_shared "github.com/KaueTTS/streaming_api/src/api/v1/dto/shared"

func NewDetail(field string, value string, message string) dto_shared.DetailErrorDto {
	return dto_shared.DetailErrorDto{
		Field:   field,
		Value:   value,
		Message: message,
	}
}
