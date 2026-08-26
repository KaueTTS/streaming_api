package responses_test

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	dto_shared "github.com/KaueTTS/streaming_api/src/api/v1/dto/shared"
	responses "github.com/KaueTTS/streaming_api/src/api/v1/responses"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestErrorResponses(t *testing.T) {
	details := []dto_shared.DetailErrorDto{
		{
			Field:   "email",
			Value:   "invalid-email",
			Message: "e-mail invalido",
		},
	}

	tests := []struct {
		name           string
		path           string
		expectedStatus int
		expectedBody   dto_shared.ErrorDto
		handler        func(ctx *fiber.Ctx) error
	}{
		{
			name:           "should return bad request",
			path:           "/bad-request",
			expectedStatus: fiber.StatusBadRequest,
			expectedBody: dto_shared.ErrorDto{
				Message:     "invalid request",
				CodeMessage: dto_shared.BadRequest,
				Details:     details,
			},
			handler: func(ctx *fiber.Ctx) error {
				return responses.BadRequest(ctx, "invalid request", details)
			},
		},
		{
			name:           "should return unauthorized",
			path:           "/unauthorized",
			expectedStatus: fiber.StatusUnauthorized,
			expectedBody: dto_shared.ErrorDto{
				Message:     "unauthorized",
				CodeMessage: dto_shared.Unauthorized,
			},
			handler: func(ctx *fiber.Ctx) error {
				return responses.Unauthorized(ctx, "unauthorized")
			},
		},
		{
			name:           "should return forbidden",
			path:           "/forbidden",
			expectedStatus: fiber.StatusForbidden,
			expectedBody: dto_shared.ErrorDto{
				Message:     "forbidden",
				CodeMessage: dto_shared.Forbidden,
			},
			handler: func(ctx *fiber.Ctx) error {
				return responses.Forbidden(ctx, "forbidden")
			},
		},
		{
			name:           "should return not found",
			path:           "/not-found",
			expectedStatus: fiber.StatusNotFound,
			expectedBody: dto_shared.ErrorDto{
				Message:     "not found",
				CodeMessage: dto_shared.NotFound,
			},
			handler: func(ctx *fiber.Ctx) error {
				return responses.NotFound(ctx, "not found")
			},
		},
		{
			name:           "should return conflict",
			path:           "/conflict",
			expectedStatus: fiber.StatusConflict,
			expectedBody: dto_shared.ErrorDto{
				Message:     "conflict",
				CodeMessage: dto_shared.Conflict,
			},
			handler: func(ctx *fiber.Ctx) error {
				return responses.Conflict(ctx, "conflict")
			},
		},
		{
			name:           "should return too many requests",
			path:           "/too-many-requests",
			expectedStatus: fiber.StatusTooManyRequests,
			expectedBody: dto_shared.ErrorDto{
				Message:     "too many requests",
				CodeMessage: dto_shared.TooManyRequests,
			},
			handler: func(ctx *fiber.Ctx) error {
				return responses.TooManyRequests(ctx, "too many requests")
			},
		},
		{
			name:           "should return internal server error",
			path:           "/internal-server-error",
			expectedStatus: fiber.StatusInternalServerError,
			expectedBody: dto_shared.ErrorDto{
				Message:     "internal server error",
				CodeMessage: dto_shared.InternalServerError,
			},
			handler: func(ctx *fiber.Ctx) error {
				return responses.InternalServerError(ctx, "internal server error")
			},
		},
		{
			name:           "should return bad gateway",
			path:           "/bad-gateway",
			expectedStatus: fiber.StatusBadGateway,
			expectedBody: dto_shared.ErrorDto{
				Message:     "bad gateway",
				CodeMessage: dto_shared.BadGateway,
			},
			handler: func(ctx *fiber.Ctx) error {
				return responses.BadGateway(ctx, "bad gateway")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := fiber.New()
			app.Get(tt.path, tt.handler)

			response, body := performErrorResponseRequest(t, app, tt.path)

			assert.Equal(t, tt.expectedStatus, response.StatusCode)

			var responseBody dto_shared.ErrorDto
			require.NoError(t, json.Unmarshal([]byte(body), &responseBody))
			assert.Equal(t, tt.expectedBody, responseBody)
		})
	}
}

func performErrorResponseRequest(t *testing.T, app *fiber.App, path string) (*http.Response, string) {
	t.Helper()

	request := httptest.NewRequest(http.MethodGet, path, nil)
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
