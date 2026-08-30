package responses

import (
	dto_shared "github.com/KaueTTS/streaming_api/src/api/v1/dto/shared"
	"github.com/gofiber/fiber/v2"
)

// Error 400
func BadRequest(ctx *fiber.Ctx, message string, details []dto_shared.DetailErrorDto) error {
	return ctx.Status(fiber.StatusBadRequest).JSON(dto_shared.ErrorDto{
		Message:     message,
		CodeMessage: dto_shared.BadRequest,
		Details:     details,
	})
}

// Error 401
func Unauthorized(ctx *fiber.Ctx, message string) error {
	return ctx.Status(fiber.StatusUnauthorized).JSON(dto_shared.ErrorDto{
		Message:     message,
		CodeMessage: dto_shared.Unauthorized,
	})
}

// Error 403
func Forbidden(ctx *fiber.Ctx, message string) error {
	return ctx.Status(fiber.StatusForbidden).JSON(dto_shared.ErrorDto{
		Message:     message,
		CodeMessage: dto_shared.Forbidden,
	})
}

// Error 404
func NotFound(ctx *fiber.Ctx, message string) error {
	return ctx.Status(fiber.StatusNotFound).JSON(dto_shared.ErrorDto{
		Message:     message,
		CodeMessage: dto_shared.NotFound,
	})
}

// Error 409
func Conflict(ctx *fiber.Ctx, message string) error {
	return ctx.Status(fiber.StatusConflict).JSON(dto_shared.ErrorDto{
		Message:     message,
		CodeMessage: dto_shared.Conflict,
	})
}

// Error 429
func TooManyRequests(ctx *fiber.Ctx, message string) error {
	return ctx.Status(fiber.StatusTooManyRequests).JSON(dto_shared.ErrorDto{
		Message:     message,
		CodeMessage: dto_shared.TooManyRequests,
	})
}

// Error 500
func InternalServerError(ctx *fiber.Ctx, message string) error {
	return ctx.Status(fiber.StatusInternalServerError).JSON(dto_shared.ErrorDto{
		Message:     message,
		CodeMessage: dto_shared.InternalServerError,
	})
}

// Error 502
func BadGateway(ctx *fiber.Ctx, message string) error {
	return ctx.Status(fiber.StatusBadGateway).JSON(dto_shared.ErrorDto{
		Message:     message,
		CodeMessage: dto_shared.BadGateway,
	})
}
