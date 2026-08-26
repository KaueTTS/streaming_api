package v1_controller_auth_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	v1_controller_auth "github.com/KaueTTS/streaming_api/src/api/v1/controllers/auth"
	dto_auth "github.com/KaueTTS/streaming_api/src/api/v1/dto/auth"
	dto_shared "github.com/KaueTTS/streaming_api/src/api/v1/dto/shared"
	service_mock "github.com/KaueTTS/streaming_api/src/mocks/services"
	shared_constants "github.com/KaueTTS/streaming_api/src/shared/constants"
	shared_errors "github.com/KaueTTS/streaming_api/src/shared/errors"
	shared_errors_auth "github.com/KaueTTS/streaming_api/src/shared/errors/auth"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestRegister(t *testing.T) {
	t.Run("should register user", func(t *testing.T) {
		authService := new(service_mock.AuthServiceMock)
		app := setupAuthControllerTestApp(authService)
		createdAt := time.Date(2026, 1, 1, 10, 0, 0, 0, time.UTC)
		updatedAt := time.Date(2026, 1, 1, 11, 0, 0, 0, time.UTC)
		requestBody := `{"name":"John Doe","email":"john@example.com","password":"abc12345"}`
		expectedRequest := dto_auth.RegisterRequestDto{
			Name:     "John Doe",
			Email:    "john@example.com",
			Password: "abc12345",
		}
		serviceResponse := dto_auth.UserResponseDto{
			ID:        1,
			Name:      "John Doe",
			Email:     "john@example.com",
			Role:      shared_constants.RoleUser,
			CreatedAt: createdAt,
			UpdatedAt: updatedAt,
		}

		authService.On("Register", mock.Anything, expectedRequest).
			Return(serviceResponse, nil).
			Once()

		response, body := performAuthControllerRequest(t, app, http.MethodPost, "/register", requestBody)

		assert.Equal(t, fiber.StatusCreated, response.StatusCode)

		var responseBody dto_auth.UserResponseDto
		decodeAuthControllerJSON(t, body, &responseBody)
		assert.Equal(t, serviceResponse, responseBody)
		authService.AssertExpectations(t)
	})

	t.Run("should return bad request when body is invalid", func(t *testing.T) {
		authService := new(service_mock.AuthServiceMock)
		app := setupAuthControllerTestApp(authService)

		response, body := performAuthControllerRequest(t, app, http.MethodPost, "/register", `{"name":`)

		assert.Equal(t, fiber.StatusBadRequest, response.StatusCode)
		assertAuthControllerError(t, body, shared_errors.InvalidRequestBody, dto_shared.BadRequest)
		authService.AssertNotCalled(t, "Register", mock.Anything, mock.Anything)
	})

	t.Run("should return bad request when request validation fails", func(t *testing.T) {
		authService := new(service_mock.AuthServiceMock)
		app := setupAuthControllerTestApp(authService)
		requestBody := `{"name":"","email":"invalid-email","password":""}`

		response, body := performAuthControllerRequest(t, app, http.MethodPost, "/register", requestBody)

		assert.Equal(t, fiber.StatusBadRequest, response.StatusCode)

		var errorBody dto_shared.ErrorDto
		decodeAuthControllerJSON(t, body, &errorBody)
		assert.Equal(t, shared_errors_auth.InvalidRegisterData, errorBody.Message)
		assert.Equal(t, dto_shared.BadRequest, errorBody.CodeMessage)
		assert.Len(t, errorBody.Details, 3)
		authService.AssertNotCalled(t, "Register", mock.Anything, mock.Anything)
	})

	t.Run("should return conflict when email is already in use", func(t *testing.T) {
		authService := new(service_mock.AuthServiceMock)
		app := setupAuthControllerTestApp(authService)
		requestBody := `{"name":"John Doe","email":"john@example.com","password":"abc12345"}`

		authService.On("Register", mock.Anything, dto_auth.RegisterRequestDto{
			Name:     "John Doe",
			Email:    "john@example.com",
			Password: "abc12345",
		}).
			Return(dto_auth.UserResponseDto{}, shared_errors.ErrEmailAlreadyInUse).
			Once()

		response, body := performAuthControllerRequest(t, app, http.MethodPost, "/register", requestBody)

		assert.Equal(t, fiber.StatusConflict, response.StatusCode)
		assertAuthControllerError(t, body, shared_errors_auth.EmailAlreadyInUse, dto_shared.Conflict)
		authService.AssertExpectations(t)
	})

	t.Run("should return bad request when service returns invalid password", func(t *testing.T) {
		authService := new(service_mock.AuthServiceMock)
		app := setupAuthControllerTestApp(authService)
		requestBody := `{"name":"John Doe","email":"john@example.com","password":"abc12345"}`

		authService.On("Register", mock.Anything, mock.AnythingOfType("dto_auth.RegisterRequestDto")).
			Return(dto_auth.UserResponseDto{}, shared_errors.ErrInvalidPassword).
			Once()

		response, body := performAuthControllerRequest(t, app, http.MethodPost, "/register", requestBody)

		assert.Equal(t, fiber.StatusBadRequest, response.StatusCode)

		var errorBody dto_shared.ErrorDto
		decodeAuthControllerJSON(t, body, &errorBody)
		assert.Equal(t, shared_errors_auth.InvalidPassword, errorBody.Message)
		assert.Equal(t, dto_shared.BadRequest, errorBody.CodeMessage)
		require.Len(t, errorBody.Details, 1)
		assert.Equal(t, shared_constants.Password, errorBody.Details[0].Field)
		assert.Equal(t, shared_constants.Hidden, errorBody.Details[0].Value)
		assert.Equal(t, shared_errors_auth.InvalidPassword, errorBody.Details[0].Message)
		authService.AssertExpectations(t)
	})

	t.Run("should return internal server error when service fails", func(t *testing.T) {
		authService := new(service_mock.AuthServiceMock)
		app := setupAuthControllerTestApp(authService)
		requestBody := `{"name":"John Doe","email":"john@example.com","password":"abc12345"}`

		authService.On("Register", mock.Anything, mock.AnythingOfType("dto_auth.RegisterRequestDto")).
			Return(dto_auth.UserResponseDto{}, errors.New("service unavailable")).
			Once()

		response, body := performAuthControllerRequest(t, app, http.MethodPost, "/register", requestBody)

		assert.Equal(t, fiber.StatusInternalServerError, response.StatusCode)
		assertAuthControllerError(t, body, shared_errors_auth.FailedToRegisterUser, dto_shared.InternalServerError)
		authService.AssertExpectations(t)
	})
}

func TestLogin(t *testing.T) {
	t.Run("should login user", func(t *testing.T) {
		authService := new(service_mock.AuthServiceMock)
		app := setupAuthControllerTestApp(authService)
		requestBody := `{"email":"john@example.com","password":"abc12345"}`
		expectedRequest := dto_auth.LoginRequestDto{
			Email:    "john@example.com",
			Password: "abc12345",
		}
		serviceResponse := dto_auth.AuthResponseDto{
			Token:     "token",
			TokenType: "Bearer",
			ExpiresIn: 7200,
			User: dto_auth.UserResponseDto{
				ID:    1,
				Name:  "John Doe",
				Email: "john@example.com",
				Role:  shared_constants.RoleUser,
			},
		}

		authService.On("Login", mock.Anything, expectedRequest).
			Return(serviceResponse, nil).
			Once()

		response, body := performAuthControllerRequest(t, app, http.MethodPost, "/login", requestBody)

		assert.Equal(t, fiber.StatusOK, response.StatusCode)

		var responseBody dto_auth.AuthResponseDto
		decodeAuthControllerJSON(t, body, &responseBody)
		assert.Equal(t, serviceResponse, responseBody)
		authService.AssertExpectations(t)
	})

	t.Run("should return bad request when body is invalid", func(t *testing.T) {
		authService := new(service_mock.AuthServiceMock)
		app := setupAuthControllerTestApp(authService)

		response, body := performAuthControllerRequest(t, app, http.MethodPost, "/login", `{"email":`)

		assert.Equal(t, fiber.StatusBadRequest, response.StatusCode)
		assertAuthControllerError(t, body, shared_errors.InvalidRequestBody, dto_shared.BadRequest)
		authService.AssertNotCalled(t, "Login", mock.Anything, mock.Anything)
	})

	t.Run("should return bad request when request validation fails", func(t *testing.T) {
		authService := new(service_mock.AuthServiceMock)
		app := setupAuthControllerTestApp(authService)
		requestBody := `{"email":"invalid-email","password":""}`

		response, body := performAuthControllerRequest(t, app, http.MethodPost, "/login", requestBody)

		assert.Equal(t, fiber.StatusBadRequest, response.StatusCode)

		var errorBody dto_shared.ErrorDto
		decodeAuthControllerJSON(t, body, &errorBody)
		assert.Equal(t, shared_errors_auth.InvalidLoginData, errorBody.Message)
		assert.Equal(t, dto_shared.BadRequest, errorBody.CodeMessage)
		assert.Len(t, errorBody.Details, 2)
		authService.AssertNotCalled(t, "Login", mock.Anything, mock.Anything)
	})

	t.Run("should return unauthorized when credentials are invalid", func(t *testing.T) {
		authService := new(service_mock.AuthServiceMock)
		app := setupAuthControllerTestApp(authService)
		requestBody := `{"email":"john@example.com","password":"wrong12345"}`

		authService.On("Login", mock.Anything, dto_auth.LoginRequestDto{
			Email:    "john@example.com",
			Password: "wrong12345",
		}).
			Return(dto_auth.AuthResponseDto{}, shared_errors.ErrInvalidCredentials).
			Once()

		response, body := performAuthControllerRequest(t, app, http.MethodPost, "/login", requestBody)

		assert.Equal(t, fiber.StatusUnauthorized, response.StatusCode)
		assertAuthControllerError(t, body, shared_errors_auth.InvalidCredentials, dto_shared.Unauthorized)
		authService.AssertExpectations(t)
	})

	t.Run("should return internal server error when service fails", func(t *testing.T) {
		authService := new(service_mock.AuthServiceMock)
		app := setupAuthControllerTestApp(authService)
		requestBody := `{"email":"john@example.com","password":"abc12345"}`

		authService.On("Login", mock.Anything, mock.AnythingOfType("dto_auth.LoginRequestDto")).
			Return(dto_auth.AuthResponseDto{}, errors.New("service unavailable")).
			Once()

		response, body := performAuthControllerRequest(t, app, http.MethodPost, "/login", requestBody)

		assert.Equal(t, fiber.StatusInternalServerError, response.StatusCode)
		assertAuthControllerError(t, body, shared_errors_auth.FailedToLogin, dto_shared.InternalServerError)
		authService.AssertExpectations(t)
	})
}

func TestMe(t *testing.T) {
	t.Run("should return authenticated user", func(t *testing.T) {
		authService := new(service_mock.AuthServiceMock)
		app := setupAuthControllerTestApp(authService)
		serviceResponse := dto_auth.UserResponseDto{
			ID:    1,
			Name:  "John Doe",
			Email: "john@example.com",
			Role:  shared_constants.RoleUser,
		}

		authService.On("Me", mock.Anything, uint(1)).
			Return(serviceResponse, nil).
			Once()

		response, body := performAuthControllerRequest(t, app, http.MethodGet, "/me", "")

		assert.Equal(t, fiber.StatusOK, response.StatusCode)

		var responseBody dto_auth.UserResponseDto
		decodeAuthControllerJSON(t, body, &responseBody)
		assert.Equal(t, serviceResponse, responseBody)
		authService.AssertExpectations(t)
	})

	t.Run("should return unauthorized when authenticated user id is missing", func(t *testing.T) {
		authService := new(service_mock.AuthServiceMock)
		app := setupAuthControllerTestApp(authService)

		response, body := performAuthControllerRequest(t, app, http.MethodGet, "/me/without-user", "")

		assert.Equal(t, fiber.StatusUnauthorized, response.StatusCode)
		assertAuthControllerError(t, body, shared_errors_auth.InvalidToken, dto_shared.Unauthorized)
		authService.AssertNotCalled(t, "Me", mock.Anything, mock.Anything)
	})

	t.Run("should return not found when user does not exist", func(t *testing.T) {
		authService := new(service_mock.AuthServiceMock)
		app := setupAuthControllerTestApp(authService)

		authService.On("Me", mock.Anything, uint(999)).
			Return(dto_auth.UserResponseDto{}, shared_errors.ErrUserNotFound).
			Once()

		response, body := performAuthControllerRequest(t, app, http.MethodGet, "/me/missing-user", "")

		assert.Equal(t, fiber.StatusNotFound, response.StatusCode)
		assertAuthControllerError(t, body, shared_errors_auth.UserNotFound, dto_shared.NotFound)
		authService.AssertExpectations(t)
	})

	t.Run("should return internal server error when service fails", func(t *testing.T) {
		authService := new(service_mock.AuthServiceMock)
		app := setupAuthControllerTestApp(authService)

		authService.On("Me", mock.Anything, uint(1)).
			Return(dto_auth.UserResponseDto{}, errors.New("service unavailable")).
			Once()

		response, body := performAuthControllerRequest(t, app, http.MethodGet, "/me", "")

		assert.Equal(t, fiber.StatusInternalServerError, response.StatusCode)
		assertAuthControllerError(t, body, shared_errors_auth.FailedToGetAuthenticatedUser, dto_shared.InternalServerError)
		authService.AssertExpectations(t)
	})
}

func setupAuthControllerTestApp(authService *service_mock.AuthServiceMock) *fiber.App {
	app := fiber.New()
	controller := v1_controller_auth.NewAuthController(authService)

	app.Post("/register", controller.Register)
	app.Post("/login", controller.Login)
	app.Get("/me", setAuthControllerUserID(1), controller.Me)
	app.Get("/me/missing-user", setAuthControllerUserID(999), controller.Me)
	app.Get("/me/without-user", controller.Me)

	return app
}

func setAuthControllerUserID(userID uint) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		ctx.Locals("user_id", userID)
		return ctx.Next()
	}
}

func performAuthControllerRequest(t *testing.T, app *fiber.App, method string, path string, body string) (*http.Response, string) {
	t.Helper()

	request := httptest.NewRequest(method, path, bytes.NewBufferString(body))
	request.Header.Set("Content-Type", "application/json")

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

func decodeAuthControllerJSON(t *testing.T, body string, target any) {
	t.Helper()

	require.NoError(t, json.Unmarshal([]byte(body), target))
}

func assertAuthControllerError(t *testing.T, body string, message string, codeMessage string) {
	t.Helper()

	var errorBody dto_shared.ErrorDto
	decodeAuthControllerJSON(t, body, &errorBody)
	assert.Equal(t, message, errorBody.Message)
	assert.Equal(t, codeMessage, errorBody.CodeMessage)
}
