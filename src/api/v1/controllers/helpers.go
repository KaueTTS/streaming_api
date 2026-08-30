package controllers_helpers

import (
	"strconv"

	dto_shared "github.com/KaueTTS/streaming_api/src/api/v1/dto/shared"
	shared_errors "github.com/KaueTTS/streaming_api/src/shared/errors"
	"github.com/gofiber/fiber/v2"
)

// GetAuthenticatedUserID pega o user_id do token e verifica se ele é válido
func GetAuthenticatedUserID(ctx *fiber.Ctx) (uint, bool) {
	userID, ok := ctx.Locals("user_id").(uint)
	return userID, ok && userID > 0
}

// ParseBody pega o body da requisição e verifica se ele é válido
func ParseBody[T any](ctx *fiber.Ctx) (T, []dto_shared.DetailErrorDto, bool) {
	var request T

	if err := ctx.BodyParser(&request); err != nil {
		return request, []dto_shared.DetailErrorDto{
			{
				Field:   "",
				Value:   "",
				Message: shared_errors.InvalidRequestBody,
			},
		}, false
	}

	return request, nil, true
}

// ParseQuery pega os parâmetros da query e verifica se eles são válidos
func ParseQuery[T any](ctx *fiber.Ctx) (T, []dto_shared.DetailErrorDto, bool) {
	var request T

	if err := ctx.QueryParser(&request); err != nil {
		return request, []dto_shared.DetailErrorDto{
			{
				Field:   "",
				Value:   "",
				Message: shared_errors.InvalidQueryParameters,
			},
		}, false
	}

	return request, nil, true
}

// ParseUintParam pega um parâmetro de rota do tipo uint e verifica se ele é válido
func ParseUintParam(ctx *fiber.Ctx, param string) (uint, string, bool) {
	value := ctx.Params(param)

	id, err := strconv.ParseUint(value, 10, 64)
	if err != nil || id == 0 {
		return 0, value, false
	}

	return uint(id), value, true
}
