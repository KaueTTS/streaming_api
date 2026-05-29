package auth_middleware

import (
	responses "github.com/KaueTTS/streaming_api/src/api/v1/responses"
	security "github.com/KaueTTS/streaming_api/src/security"
	shared_errors "github.com/KaueTTS/streaming_api/src/shared/errors"
	shared_errors_auth "github.com/KaueTTS/streaming_api/src/shared/errors/auth"
	"github.com/gofiber/fiber/v2"
)

func AuthRequired() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		token, err := security.ExtractBearerToken(ctx.Get("Authorization"))
		if err != nil {
			return responses.Unauthorized(ctx, shared_errors_auth.TokenMissingOrInvalid)
		}

		claims, err := security.ValidateToken(token)
		if err != nil {
			return responses.Unauthorized(ctx, shared_errors_auth.TokenInvalidOrExpired)
		}

		ctx.Locals("user_id", claims.UserID)
		ctx.Locals("user_email", claims.Email)
		ctx.Locals("user_role", claims.Role)

		return ctx.Next()
	}
}

func AdminRequired() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		role, ok := ctx.Locals("user_role").(string)
		if !ok || role != "admin" {
			return responses.Forbidden(ctx, shared_errors.AccessAdminOnly)
		}

		return ctx.Next()
	}
}
