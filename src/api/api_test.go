package api_test

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	api "github.com/KaueTTS/streaming_api/src/api"
	dto_shared "github.com/KaueTTS/streaming_api/src/api/v1/dto/shared"
	env "github.com/KaueTTS/streaming_api/src/configs/env"
	shared_errors_auth "github.com/KaueTTS/streaming_api/src/shared/errors/auth"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInit(t *testing.T) {
	t.Run("should create app with configured name", func(t *testing.T) {
		setupAPIEnv(t)

		app := api.Init(nil, nil)

		assert.Equal(t, "streaming_api_test", app.Config().AppName)
	})

	t.Run("should register health route", func(t *testing.T) {
		setupAPIEnv(t)

		app := api.Init(nil, nil)

		response, body := performAPIRequest(t, app, http.MethodGet, "/health", "")

		assert.Equal(t, fiber.StatusOK, response.StatusCode)
		assert.JSONEq(t, `{"status":"ok"}`, body)
	})

	t.Run("should protect private routes", func(t *testing.T) {
		setupAPIEnv(t)

		app := api.Init(nil, nil)

		response, body := performAPIRequest(t, app, http.MethodGet, "/v1/profiles", "")

		assert.Equal(t, fiber.StatusUnauthorized, response.StatusCode)

		var errorResponse dto_shared.ErrorDto
		require.NoError(t, json.Unmarshal([]byte(body), &errorResponse))

		assert.Equal(t, shared_errors_auth.TokenMissingOrInvalid, errorResponse.Message)
		assert.Equal(t, dto_shared.Unauthorized, errorResponse.CodeMessage)
	})

	t.Run("should register public auth routes", func(t *testing.T) {
		setupAPIEnv(t)

		app := api.Init(nil, nil)

		response, body := performAPIRequest(
			t,
			app,
			http.MethodPost,
			"/v1/auth/login",
			`{"email":"invalid-email","password":""}`,
		)

		assert.Equal(t, fiber.StatusBadRequest, response.StatusCode)

		var errorResponse dto_shared.ErrorDto
		require.NoError(t, json.Unmarshal([]byte(body), &errorResponse))

		assert.Equal(t, shared_errors_auth.InvalidLoginData, errorResponse.Message)
		assert.Equal(t, dto_shared.BadRequest, errorResponse.CodeMessage)
	})

	t.Run("should configure cors preflight", func(t *testing.T) {
		setupAPIEnv(t)

		app := api.Init(nil, nil)

		request := httptest.NewRequest(http.MethodOptions, "/v1/auth/login", nil)
		request.Header.Set("Origin", "https://example.com")
		request.Header.Set("Access-Control-Request-Method", http.MethodPost)

		response, err := app.Test(request, -1)
		require.NoError(t, err)
		require.NotNil(t, response)

		t.Cleanup(func() {
			require.NoError(t, response.Body.Close())
		})

		assert.Equal(t, fiber.StatusNoContent, response.StatusCode)
		assert.Equal(t, "*", response.Header.Get("Access-Control-Allow-Origin"))
		assert.Contains(t, response.Header.Get("Access-Control-Allow-Methods"), http.MethodPost)
		assert.Contains(t, response.Header.Get("Access-Control-Allow-Headers"), "Authorization")
		assert.Contains(t, response.Header.Get("Access-Control-Allow-Headers"), "Content-Type")
	})
}

func setupAPIEnv(t *testing.T) {
	t.Helper()

	oldAppName := env.AppName
	oldAppEnv := env.AppEnv
	oldJWTSecret := env.JWTSecret
	oldAuthTokenExpirationTimeInHours := env.AuthTokenExpirationTimeInHours
	oldTMDBBaseURL := env.TMDBBaseURL
	oldTMDBAccessToken := env.TMDBAccessToken

	env.AppName = "streaming_api_test"
	env.AppEnv = "test"
	env.JWTSecret = "test-secret"
	env.AuthTokenExpirationTimeInHours = 8
	env.TMDBBaseURL = "https://tmdb.test"
	env.TMDBAccessToken = "test-token"

	t.Cleanup(func() {
		env.AppName = oldAppName
		env.AppEnv = oldAppEnv
		env.JWTSecret = oldJWTSecret
		env.AuthTokenExpirationTimeInHours = oldAuthTokenExpirationTimeInHours
		env.TMDBBaseURL = oldTMDBBaseURL
		env.TMDBAccessToken = oldTMDBAccessToken
	})
}

func performAPIRequest(t *testing.T, app *fiber.App, method string, path string, body string) (*http.Response, string) {
	t.Helper()

	var requestBody io.Reader
	if body != "" {
		requestBody = strings.NewReader(body)
	}

	request := httptest.NewRequest(method, path, requestBody)
	if body != "" {
		request.Header.Set("Content-Type", "application/json")
	}

	response, err := app.Test(request, -1)
	require.NoError(t, err)
	require.NotNil(t, response)

	t.Cleanup(func() {
		require.NoError(t, response.Body.Close())
	})

	responseBody, err := io.ReadAll(response.Body)
	require.NoError(t, err)

	return response, string(responseBody)
}
