package auth_middleware_test

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	responses_dto "github.com/KaueTTS/streaming_api/src/api/v1/dto/shared"
	env "github.com/KaueTTS/streaming_api/src/configs/env"
	auth_middleware "github.com/KaueTTS/streaming_api/src/middlewares"
	security "github.com/KaueTTS/streaming_api/src/security"
	shared_errors "github.com/KaueTTS/streaming_api/src/shared/errors"
	shared_errors_auth "github.com/KaueTTS/streaming_api/src/shared/errors/auth"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAuthRequired(t *testing.T) {
	t.Run("should return unauthorized when authorization header is missing", func(t *testing.T) {
		setupAuthEnv(t)

		app := fiber.New()

		app.Get("/private", auth_middleware.AuthRequired(), func(ctx *fiber.Ctx) error {
			return ctx.SendStatus(fiber.StatusOK)
		})

		response, body := performRequest(t, app, http.MethodGet, "/private", "")

		assert.Equal(t, fiber.StatusUnauthorized, response.StatusCode)

		var errorResponse responses_dto.ErrorDto
		require.NoError(t, json.Unmarshal([]byte(body), &errorResponse))

		assert.Equal(t, shared_errors_auth.TokenMissingOrInvalid, errorResponse.Message)
		assert.Equal(t, responses_dto.Unauthorized, errorResponse.CodeMessage)
	})

	t.Run("should return unauthorized when authorization header is invalid", func(t *testing.T) {
		setupAuthEnv(t)

		app := fiber.New()

		app.Get("/private", auth_middleware.AuthRequired(), func(ctx *fiber.Ctx) error {
			return ctx.SendStatus(fiber.StatusOK)
		})

		response, body := performRequest(t, app, http.MethodGet, "/private", "invalid-token")

		assert.Equal(t, fiber.StatusUnauthorized, response.StatusCode)

		var errorResponse responses_dto.ErrorDto
		require.NoError(t, json.Unmarshal([]byte(body), &errorResponse))

		assert.Equal(t, shared_errors_auth.TokenMissingOrInvalid, errorResponse.Message)
		assert.Equal(t, responses_dto.Unauthorized, errorResponse.CodeMessage)
	})

	t.Run("should return unauthorized when token is invalid", func(t *testing.T) {
		setupAuthEnv(t)

		app := fiber.New()

		app.Get("/private", auth_middleware.AuthRequired(), func(ctx *fiber.Ctx) error {
			return ctx.SendStatus(fiber.StatusOK)
		})

		response, body := performRequest(t, app, http.MethodGet, "/private", "Bearer invalid-token")

		assert.Equal(t, fiber.StatusUnauthorized, response.StatusCode)

		var errorResponse responses_dto.ErrorDto
		require.NoError(t, json.Unmarshal([]byte(body), &errorResponse))

		assert.Equal(t, shared_errors_auth.TokenInvalidOrExpired, errorResponse.Message)
		assert.Equal(t, responses_dto.Unauthorized, errorResponse.CodeMessage)
	})

	t.Run("should call next when token is valid", func(t *testing.T) {
		setupAuthEnv(t)

		token, err := security.GenerateToken(1, "admin@test.com", "admin")
		require.NoError(t, err)

		app := fiber.New()

		app.Get("/private", auth_middleware.AuthRequired(), func(ctx *fiber.Ctx) error {
			return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
				"user_id":    ctx.Locals("user_id"),
				"user_email": ctx.Locals("user_email"),
				"user_role":  ctx.Locals("user_role"),
			})
		})

		response, body := performRequest(t, app, http.MethodGet, "/private", "Bearer "+token)

		assert.Equal(t, fiber.StatusOK, response.StatusCode)

		assert.JSONEq(t, `{
			"user_id": 1,
			"user_email": "admin@test.com",
			"user_role": "admin"
		}`, body)
	})
}

func TestAdminRequired(t *testing.T) {
	t.Run("should return forbidden when user role is missing", func(t *testing.T) {
		app := fiber.New()

		app.Get("/admin", auth_middleware.AdminRequired(), func(ctx *fiber.Ctx) error {
			return ctx.SendStatus(fiber.StatusOK)
		})

		response, body := performRequest(t, app, http.MethodGet, "/admin", "")

		assert.Equal(t, fiber.StatusForbidden, response.StatusCode)

		var errorResponse responses_dto.ErrorDto
		require.NoError(t, json.Unmarshal([]byte(body), &errorResponse))

		assert.Equal(t, shared_errors.AccessAdminOnly, errorResponse.Message)
		assert.Equal(t, responses_dto.Forbidden, errorResponse.CodeMessage)
	})

	t.Run("should return forbidden when user role is not admin", func(t *testing.T) {
		app := fiber.New()

		app.Get("/admin",
			func(ctx *fiber.Ctx) error {
				ctx.Locals("user_role", "user")
				return ctx.Next()
			},
			auth_middleware.AdminRequired(),
			func(ctx *fiber.Ctx) error {
				return ctx.SendStatus(fiber.StatusOK)
			},
		)

		response, body := performRequest(t, app, http.MethodGet, "/admin", "")

		assert.Equal(t, fiber.StatusForbidden, response.StatusCode)

		var errorResponse responses_dto.ErrorDto
		require.NoError(t, json.Unmarshal([]byte(body), &errorResponse))

		assert.Equal(t, shared_errors.AccessAdminOnly, errorResponse.Message)
		assert.Equal(t, responses_dto.Forbidden, errorResponse.CodeMessage)
	})

	t.Run("should call next when user role is admin", func(t *testing.T) {
		app := fiber.New()

		app.Get("/admin",
			func(ctx *fiber.Ctx) error {
				ctx.Locals("user_role", "admin")
				return ctx.Next()
			},
			auth_middleware.AdminRequired(),
			func(ctx *fiber.Ctx) error {
				return ctx.SendStatus(fiber.StatusOK)
			},
		)

		response, _ := performRequest(t, app, http.MethodGet, "/admin", "")

		assert.Equal(t, fiber.StatusOK, response.StatusCode)
	})
}

func setupAuthEnv(t *testing.T) {
	oldJWTSecret := env.JWTSecret
	oldAppName := env.AppName
	oldAuthTokenExpirationTimeInHours := env.AuthTokenExpirationTimeInHours

	env.JWTSecret = "test-secret"
	env.AppName = "streaming_api_test"
	env.AuthTokenExpirationTimeInHours = 8

	t.Cleanup(func() {
		env.JWTSecret = oldJWTSecret
		env.AppName = oldAppName
		env.AuthTokenExpirationTimeInHours = oldAuthTokenExpirationTimeInHours
	})
}

func performRequest(t *testing.T, app *fiber.App, method string, path string, authorizationHeader string) (*http.Response, string) {
	request := httptest.NewRequest(method, path, nil)

	if authorizationHeader != "" {
		request.Header.Set("Authorization", authorizationHeader)
	}

	response, err := app.Test(request, -1)
	require.NoError(t, err)
	require.NotNil(t, response)

	t.Cleanup(func() {
		require.NoError(t, response.Body.Close())
	})

	body, err := io.ReadAll(response.Body)
	require.NoError(t, err)

	return response, string(body)
}
