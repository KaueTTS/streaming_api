package v1_controller_profile_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	v1_controller_profile "github.com/KaueTTS/streaming_api/src/api/v1/controllers/profile"
	dto_profile "github.com/KaueTTS/streaming_api/src/api/v1/dto/profile"
	dto_shared "github.com/KaueTTS/streaming_api/src/api/v1/dto/shared"
	service_mock "github.com/KaueTTS/streaming_api/src/mocks/services"
	shared_errors "github.com/KaueTTS/streaming_api/src/shared/errors"
	shared_errors_auth "github.com/KaueTTS/streaming_api/src/shared/errors/auth"
	shared_errors_profile "github.com/KaueTTS/streaming_api/src/shared/errors/profile"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestListProfiles(t *testing.T) {
	t.Run("should list profiles", func(t *testing.T) {
		profileService := new(service_mock.ProfileServiceMock)
		app := setupProfileControllerTestApp(profileService)
		serviceResponse := makeProfileControllerListResponse()

		profileService.On("ListProfiles", mock.Anything, uint(10), 2, 3).
			Return(serviceResponse, nil).
			Once()

		response, body := performProfileControllerRequest(t, app, http.MethodGet, "/profiles?page=2&per_page=3", "")

		assert.Equal(t, fiber.StatusOK, response.StatusCode)

		var responseBody dto_profile.ProfileResponseDto
		decodeProfileControllerJSON(t, body, &responseBody)
		assert.Equal(t, serviceResponse, responseBody)
		profileService.AssertExpectations(t)
	})

	t.Run("should normalize empty pagination", func(t *testing.T) {
		profileService := new(service_mock.ProfileServiceMock)
		app := setupProfileControllerTestApp(profileService)

		profileService.On("ListProfiles", mock.Anything, uint(10), 1, 10).
			Return(dto_profile.ProfileResponseDto{}, nil).
			Once()

		response, _ := performProfileControllerRequest(t, app, http.MethodGet, "/profiles", "")

		assert.Equal(t, fiber.StatusOK, response.StatusCode)
		profileService.AssertExpectations(t)
	})

	t.Run("should return unauthorized when authenticated user id is missing", func(t *testing.T) {
		profileService := new(service_mock.ProfileServiceMock)
		app := setupProfileControllerTestApp(profileService)

		response, body := performProfileControllerRequest(t, app, http.MethodGet, "/profiles/without-user", "")

		assert.Equal(t, fiber.StatusUnauthorized, response.StatusCode)
		assertProfileControllerError(t, body, shared_errors_auth.InvalidToken, dto_shared.Unauthorized)
		profileService.AssertNotCalled(t, "ListProfiles", mock.Anything, mock.Anything, mock.Anything, mock.Anything)
	})

	t.Run("should return bad request when pagination is invalid", func(t *testing.T) {
		profileService := new(service_mock.ProfileServiceMock)
		app := setupProfileControllerTestApp(profileService)

		response, body := performProfileControllerRequest(t, app, http.MethodGet, "/profiles?page=-1&per_page=-1", "")

		assert.Equal(t, fiber.StatusBadRequest, response.StatusCode)

		var errorBody dto_shared.ErrorDto
		decodeProfileControllerJSON(t, body, &errorBody)
		assert.Equal(t, shared_errors.InvalidQueryParameters, errorBody.Message)
		assert.Equal(t, dto_shared.BadRequest, errorBody.CodeMessage)
		assert.Len(t, errorBody.Details, 2)
		profileService.AssertNotCalled(t, "ListProfiles", mock.Anything, mock.Anything, mock.Anything, mock.Anything)
	})

	t.Run("should return internal server error when service fails", func(t *testing.T) {
		profileService := new(service_mock.ProfileServiceMock)
		app := setupProfileControllerTestApp(profileService)

		profileService.On("ListProfiles", mock.Anything, uint(10), 1, 10).
			Return(dto_profile.ProfileResponseDto{}, errors.New("service unavailable")).
			Once()

		response, body := performProfileControllerRequest(t, app, http.MethodGet, "/profiles", "")

		assert.Equal(t, fiber.StatusInternalServerError, response.StatusCode)
		assertProfileControllerError(t, body, shared_errors_profile.FailedToListProfiles, dto_shared.InternalServerError)
		profileService.AssertExpectations(t)
	})
}

func TestCreateProfile(t *testing.T) {
	t.Run("should create profile", func(t *testing.T) {
		profileService := new(service_mock.ProfileServiceMock)
		app := setupProfileControllerTestApp(profileService)
		avatarURL := "https://example.com/avatar.png"
		requestBody := `{"name":"Main","avatar_url":"https://example.com/avatar.png","is_kids":false}`
		expectedRequest := dto_profile.ProfileRequestDto{
			Name:      "Main",
			AvatarURL: &avatarURL,
			IsKids:    false,
		}
		serviceResponse := makeProfileControllerProfile(1, "Main", &avatarURL, false)

		profileService.On("CreateProfile", mock.Anything, uint(10), expectedRequest).
			Return(serviceResponse, nil).
			Once()

		response, body := performProfileControllerRequest(t, app, http.MethodPost, "/profiles", requestBody)

		assert.Equal(t, fiber.StatusCreated, response.StatusCode)

		var responseBody dto_profile.ProfileDto
		decodeProfileControllerJSON(t, body, &responseBody)
		assert.Equal(t, serviceResponse, responseBody)
		profileService.AssertExpectations(t)
	})

	t.Run("should return bad request when body is invalid", func(t *testing.T) {
		profileService := new(service_mock.ProfileServiceMock)
		app := setupProfileControllerTestApp(profileService)

		response, body := performProfileControllerRequest(t, app, http.MethodPost, "/profiles", `{"name":`)

		assert.Equal(t, fiber.StatusBadRequest, response.StatusCode)
		assertProfileControllerError(t, body, shared_errors.InvalidRequestBody, dto_shared.BadRequest)
		profileService.AssertNotCalled(t, "CreateProfile", mock.Anything, mock.Anything, mock.Anything)
	})

	t.Run("should return bad request when request validation fails", func(t *testing.T) {
		profileService := new(service_mock.ProfileServiceMock)
		app := setupProfileControllerTestApp(profileService)

		response, body := performProfileControllerRequest(t, app, http.MethodPost, "/profiles", `{"name":"","is_kids":false}`)

		assert.Equal(t, fiber.StatusBadRequest, response.StatusCode)
		assertProfileControllerError(t, body, shared_errors_profile.InvalidCreateProfileData, dto_shared.BadRequest)
		profileService.AssertNotCalled(t, "CreateProfile", mock.Anything, mock.Anything, mock.Anything)
	})

	t.Run("should return conflict when profile limit is reached", func(t *testing.T) {
		profileService := new(service_mock.ProfileServiceMock)
		app := setupProfileControllerTestApp(profileService)

		profileService.On("CreateProfile", mock.Anything, uint(10), mock.AnythingOfType("dto_profile.ProfileRequestDto")).
			Return(dto_profile.ProfileDto{}, shared_errors.ErrProfileLimitReached).
			Once()

		response, body := performProfileControllerRequest(t, app, http.MethodPost, "/profiles", `{"name":"Main","is_kids":false}`)

		assert.Equal(t, fiber.StatusConflict, response.StatusCode)
		assertProfileControllerError(t, body, shared_errors_profile.ProfileLimitReached, dto_shared.Conflict)
		profileService.AssertExpectations(t)
	})

	t.Run("should return internal server error when service fails", func(t *testing.T) {
		profileService := new(service_mock.ProfileServiceMock)
		app := setupProfileControllerTestApp(profileService)

		profileService.On("CreateProfile", mock.Anything, uint(10), mock.AnythingOfType("dto_profile.ProfileRequestDto")).
			Return(dto_profile.ProfileDto{}, errors.New("service unavailable")).
			Once()

		response, body := performProfileControllerRequest(t, app, http.MethodPost, "/profiles", `{"name":"Main","is_kids":false}`)

		assert.Equal(t, fiber.StatusInternalServerError, response.StatusCode)
		assertProfileControllerError(t, body, shared_errors_profile.FailedToCreateProfile, dto_shared.InternalServerError)
		profileService.AssertExpectations(t)
	})
}

func TestUpdateProfile(t *testing.T) {
	t.Run("should update profile", func(t *testing.T) {
		profileService := new(service_mock.ProfileServiceMock)
		app := setupProfileControllerTestApp(profileService)
		avatarURL := "https://example.com/avatar.png"
		requestBody := `{"name":"Updated Main","avatar_url":"https://example.com/avatar.png","is_kids":true}`
		expectedRequest := dto_profile.ProfileRequestDto{
			Name:      "Updated Main",
			AvatarURL: &avatarURL,
			IsKids:    true,
		}
		serviceResponse := makeProfileControllerProfile(2, "Updated Main", &avatarURL, true)

		profileService.On("UpdateProfile", mock.Anything, uint(10), uint(2), expectedRequest).
			Return(serviceResponse, nil).
			Once()

		response, body := performProfileControllerRequest(t, app, http.MethodPut, "/profiles/2", requestBody)

		assert.Equal(t, fiber.StatusOK, response.StatusCode)

		var responseBody dto_profile.ProfileDto
		decodeProfileControllerJSON(t, body, &responseBody)
		assert.Equal(t, serviceResponse, responseBody)
		profileService.AssertExpectations(t)
	})

	t.Run("should return bad request when profile id is invalid", func(t *testing.T) {
		profileService := new(service_mock.ProfileServiceMock)
		app := setupProfileControllerTestApp(profileService)

		response, body := performProfileControllerRequest(t, app, http.MethodPut, "/profiles/invalid", `{"name":"Main","is_kids":false}`)

		assert.Equal(t, fiber.StatusBadRequest, response.StatusCode)
		assertProfileControllerError(t, body, shared_errors_profile.InvalidProfileID, dto_shared.BadRequest)
		profileService.AssertNotCalled(t, "UpdateProfile", mock.Anything, mock.Anything, mock.Anything, mock.Anything)
	})

	t.Run("should return bad request when request validation fails", func(t *testing.T) {
		profileService := new(service_mock.ProfileServiceMock)
		app := setupProfileControllerTestApp(profileService)

		response, body := performProfileControllerRequest(t, app, http.MethodPut, "/profiles/2", `{"name":"","is_kids":false}`)

		assert.Equal(t, fiber.StatusBadRequest, response.StatusCode)
		assertProfileControllerError(t, body, shared_errors_profile.InvalidUpdateProfileData, dto_shared.BadRequest)
		profileService.AssertNotCalled(t, "UpdateProfile", mock.Anything, mock.Anything, mock.Anything, mock.Anything)
	})

	t.Run("should return not found when profile does not exist", func(t *testing.T) {
		profileService := new(service_mock.ProfileServiceMock)
		app := setupProfileControllerTestApp(profileService)

		profileService.On("UpdateProfile", mock.Anything, uint(10), uint(999), mock.AnythingOfType("dto_profile.ProfileRequestDto")).
			Return(dto_profile.ProfileDto{}, shared_errors.ErrProfileNotFound).
			Once()

		response, body := performProfileControllerRequest(t, app, http.MethodPut, "/profiles/999", `{"name":"Main","is_kids":false}`)

		assert.Equal(t, fiber.StatusNotFound, response.StatusCode)
		assertProfileControllerError(t, body, shared_errors_profile.ProfileNotFound, dto_shared.NotFound)
		profileService.AssertExpectations(t)
	})

	t.Run("should return internal server error when service fails", func(t *testing.T) {
		profileService := new(service_mock.ProfileServiceMock)
		app := setupProfileControllerTestApp(profileService)

		profileService.On("UpdateProfile", mock.Anything, uint(10), uint(2), mock.AnythingOfType("dto_profile.ProfileRequestDto")).
			Return(dto_profile.ProfileDto{}, errors.New("service unavailable")).
			Once()

		response, body := performProfileControllerRequest(t, app, http.MethodPut, "/profiles/2", `{"name":"Main","is_kids":false}`)

		assert.Equal(t, fiber.StatusInternalServerError, response.StatusCode)
		assertProfileControllerError(t, body, shared_errors_profile.FailedToUpdateProfile, dto_shared.InternalServerError)
		profileService.AssertExpectations(t)
	})
}

func TestDeleteProfile(t *testing.T) {
	t.Run("should delete profile", func(t *testing.T) {
		profileService := new(service_mock.ProfileServiceMock)
		app := setupProfileControllerTestApp(profileService)

		profileService.On("DeleteProfile", mock.Anything, uint(10), uint(2)).
			Return(nil).
			Once()

		response, body := performProfileControllerRequest(t, app, http.MethodDelete, "/profiles/2", "")

		assert.Equal(t, fiber.StatusNoContent, response.StatusCode)
		assert.Empty(t, body)
		profileService.AssertExpectations(t)
	})

	t.Run("should return bad request when profile id is invalid", func(t *testing.T) {
		profileService := new(service_mock.ProfileServiceMock)
		app := setupProfileControllerTestApp(profileService)

		response, body := performProfileControllerRequest(t, app, http.MethodDelete, "/profiles/invalid", "")

		assert.Equal(t, fiber.StatusBadRequest, response.StatusCode)
		assertProfileControllerError(t, body, shared_errors_profile.InvalidProfileID, dto_shared.BadRequest)
		profileService.AssertNotCalled(t, "DeleteProfile", mock.Anything, mock.Anything, mock.Anything)
	})

	t.Run("should return not found when profile does not exist", func(t *testing.T) {
		profileService := new(service_mock.ProfileServiceMock)
		app := setupProfileControllerTestApp(profileService)

		profileService.On("DeleteProfile", mock.Anything, uint(10), uint(999)).
			Return(shared_errors.ErrProfileNotFound).
			Once()

		response, body := performProfileControllerRequest(t, app, http.MethodDelete, "/profiles/999", "")

		assert.Equal(t, fiber.StatusNotFound, response.StatusCode)
		assertProfileControllerError(t, body, shared_errors_profile.ProfileNotFound, dto_shared.NotFound)
		profileService.AssertExpectations(t)
	})

	t.Run("should return internal server error when service fails", func(t *testing.T) {
		profileService := new(service_mock.ProfileServiceMock)
		app := setupProfileControllerTestApp(profileService)

		profileService.On("DeleteProfile", mock.Anything, uint(10), uint(2)).
			Return(errors.New("service unavailable")).
			Once()

		response, body := performProfileControllerRequest(t, app, http.MethodDelete, "/profiles/2", "")

		assert.Equal(t, fiber.StatusInternalServerError, response.StatusCode)
		assertProfileControllerError(t, body, shared_errors_profile.FailedToDeleteProfile, dto_shared.InternalServerError)
		profileService.AssertExpectations(t)
	})
}

func setupProfileControllerTestApp(profileService *service_mock.ProfileServiceMock) *fiber.App {
	app := fiber.New()
	controller := v1_controller_profile.NewProfileController(profileService)

	app.Get("/profiles", setProfileControllerUserID(10), controller.ListProfiles)
	app.Get("/profiles/without-user", controller.ListProfiles)
	app.Post("/profiles", setProfileControllerUserID(10), controller.CreateProfile)
	app.Put("/profiles/:id", setProfileControllerUserID(10), controller.UpdateProfile)
	app.Delete("/profiles/:id", setProfileControllerUserID(10), controller.DeleteProfile)

	return app
}

func setProfileControllerUserID(userID uint) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		ctx.Locals("user_id", userID)
		return ctx.Next()
	}
}

func makeProfileControllerListResponse() dto_profile.ProfileResponseDto {
	avatarURL := "https://example.com/avatar.png"

	return dto_profile.ProfileResponseDto{
		Data: []dto_profile.ProfileDto{
			makeProfileControllerProfile(1, "Main", &avatarURL, false),
			makeProfileControllerProfile(2, "Kids", nil, true),
		},
		Pagination: dto_shared.PaginationDto{
			Page:      2,
			PerPage:   3,
			PageCount: 4,
			Total:     10,
		},
	}
}

func makeProfileControllerProfile(id uint, name string, avatarURL *string, isKids bool) dto_profile.ProfileDto {
	return dto_profile.ProfileDto{
		ID:        id,
		UserID:    10,
		Name:      name,
		AvatarURL: avatarURL,
		IsKids:    isKids,
		CreatedAt: time.Date(2026, 1, 1, 10, 0, 0, 0, time.UTC),
		UpdatedAt: time.Date(2026, 1, 1, 11, 0, 0, 0, time.UTC),
	}
}

func performProfileControllerRequest(t *testing.T, app *fiber.App, method string, path string, body string) (*http.Response, string) {
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

func decodeProfileControllerJSON(t *testing.T, body string, target any) {
	t.Helper()

	require.NoError(t, json.Unmarshal([]byte(body), target))
}

func assertProfileControllerError(t *testing.T, body string, message string, codeMessage string) {
	t.Helper()

	var errorBody dto_shared.ErrorDto
	decodeProfileControllerJSON(t, body, &errorBody)
	assert.Equal(t, message, errorBody.Message)
	assert.Equal(t, codeMessage, errorBody.CodeMessage)
}
