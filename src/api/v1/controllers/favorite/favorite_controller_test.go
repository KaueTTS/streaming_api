package v1_controller_favorite_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	v1_controller_favorite "github.com/KaueTTS/streaming_api/src/api/v1/controllers/favorite"
	dto_content "github.com/KaueTTS/streaming_api/src/api/v1/dto/content"
	dto_favorite "github.com/KaueTTS/streaming_api/src/api/v1/dto/favorite"
	dto_shared "github.com/KaueTTS/streaming_api/src/api/v1/dto/shared"
	service_mock "github.com/KaueTTS/streaming_api/src/mocks/services"
	shared_errors "github.com/KaueTTS/streaming_api/src/shared/errors"
	shared_errors_auth "github.com/KaueTTS/streaming_api/src/shared/errors/auth"
	shared_errors_favorite "github.com/KaueTTS/streaming_api/src/shared/errors/favorite"
	shared_errors_profile "github.com/KaueTTS/streaming_api/src/shared/errors/profile"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestListFavorites(t *testing.T) {
	t.Run("should list favorites", func(t *testing.T) {
		favoriteService := new(service_mock.FavoriteServiceMock)
		app := setupFavoriteControllerTestApp(favoriteService)
		serviceResponse := makeFavoriteControllerResponse()

		favoriteService.On("ListFavorites", mock.Anything, uint(10), uint(20), 2, 3, "pt-BR").
			Return(serviceResponse, nil).
			Once()

		response, body := performFavoriteControllerRequest(t, app, http.MethodGet, "/favorites?profile_id=20&page=2&per_page=3&language=pt-BR", "")

		assert.Equal(t, fiber.StatusOK, response.StatusCode)

		var responseBody dto_favorite.FavoriteResponseDto
		decodeFavoriteControllerJSON(t, body, &responseBody)
		assert.Equal(t, serviceResponse, responseBody)
		favoriteService.AssertExpectations(t)
	})

	t.Run("should normalize empty pagination", func(t *testing.T) {
		favoriteService := new(service_mock.FavoriteServiceMock)
		app := setupFavoriteControllerTestApp(favoriteService)

		favoriteService.On("ListFavorites", mock.Anything, uint(10), uint(20), 1, 10, "").
			Return(dto_favorite.FavoriteResponseDto{}, nil).
			Once()

		response, _ := performFavoriteControllerRequest(t, app, http.MethodGet, "/favorites?profile_id=20", "")

		assert.Equal(t, fiber.StatusOK, response.StatusCode)
		favoriteService.AssertExpectations(t)
	})

	t.Run("should return unauthorized when authenticated user id is missing", func(t *testing.T) {
		favoriteService := new(service_mock.FavoriteServiceMock)
		app := setupFavoriteControllerTestApp(favoriteService)

		response, body := performFavoriteControllerRequest(t, app, http.MethodGet, "/favorites/without-user?profile_id=20", "")

		assert.Equal(t, fiber.StatusUnauthorized, response.StatusCode)
		assertFavoriteControllerError(t, body, shared_errors_auth.InvalidToken, dto_shared.Unauthorized)
		favoriteService.AssertNotCalled(t, "ListFavorites", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything)
	})

	t.Run("should return bad request when profile id is missing", func(t *testing.T) {
		favoriteService := new(service_mock.FavoriteServiceMock)
		app := setupFavoriteControllerTestApp(favoriteService)

		response, body := performFavoriteControllerRequest(t, app, http.MethodGet, "/favorites", "")

		assert.Equal(t, fiber.StatusBadRequest, response.StatusCode)
		assertFavoriteControllerError(t, body, shared_errors_favorite.InvalidProfileID, dto_shared.BadRequest)
		favoriteService.AssertNotCalled(t, "ListFavorites", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything)
	})

	t.Run("should return not found when profile does not exist", func(t *testing.T) {
		favoriteService := new(service_mock.FavoriteServiceMock)
		app := setupFavoriteControllerTestApp(favoriteService)

		favoriteService.On("ListFavorites", mock.Anything, uint(10), uint(999), 1, 10, "").
			Return(dto_favorite.FavoriteResponseDto{}, shared_errors.ErrProfileNotFound).
			Once()

		response, body := performFavoriteControllerRequest(t, app, http.MethodGet, "/favorites?profile_id=999", "")

		assert.Equal(t, fiber.StatusNotFound, response.StatusCode)
		assertFavoriteControllerError(t, body, shared_errors_profile.ProfileNotFound, dto_shared.NotFound)
		favoriteService.AssertExpectations(t)
	})

	t.Run("should return internal server error when service fails", func(t *testing.T) {
		favoriteService := new(service_mock.FavoriteServiceMock)
		app := setupFavoriteControllerTestApp(favoriteService)

		favoriteService.On("ListFavorites", mock.Anything, uint(10), uint(20), 1, 10, "").
			Return(dto_favorite.FavoriteResponseDto{}, errors.New("service unavailable")).
			Once()

		response, body := performFavoriteControllerRequest(t, app, http.MethodGet, "/favorites?profile_id=20", "")

		assert.Equal(t, fiber.StatusInternalServerError, response.StatusCode)
		assertFavoriteControllerError(t, body, shared_errors_favorite.FailedToListFavorites, dto_shared.InternalServerError)
		favoriteService.AssertExpectations(t)
	})
}

func TestAddFavorite(t *testing.T) {
	t.Run("should add favorite", func(t *testing.T) {
		favoriteService := new(service_mock.FavoriteServiceMock)
		app := setupFavoriteControllerTestApp(favoriteService)
		requestBody := `{"profile_id":20,"content_external_id":101,"type":"movie"}`
		expectedRequest := dto_favorite.FavoriteRequestDto{
			ProfileID:         20,
			ContentExternalID: 101,
			Type:              "movie",
		}

		favoriteService.On("AddFavorite", mock.Anything, uint(10), expectedRequest).
			Return(nil).
			Once()

		response, _ := performFavoriteControllerRequest(t, app, http.MethodPost, "/favorites", requestBody)

		assert.Equal(t, fiber.StatusCreated, response.StatusCode)
		favoriteService.AssertExpectations(t)
	})

	t.Run("should return bad request when body is invalid", func(t *testing.T) {
		favoriteService := new(service_mock.FavoriteServiceMock)
		app := setupFavoriteControllerTestApp(favoriteService)

		response, body := performFavoriteControllerRequest(t, app, http.MethodPost, "/favorites", `{"profile_id":`)

		assert.Equal(t, fiber.StatusBadRequest, response.StatusCode)
		assertFavoriteControllerError(t, body, shared_errors.InvalidRequestBody, dto_shared.BadRequest)
		favoriteService.AssertNotCalled(t, "AddFavorite", mock.Anything, mock.Anything, mock.Anything)
	})

	t.Run("should return bad request when request validation fails", func(t *testing.T) {
		favoriteService := new(service_mock.FavoriteServiceMock)
		app := setupFavoriteControllerTestApp(favoriteService)

		response, body := performFavoriteControllerRequest(t, app, http.MethodPost, "/favorites", `{"profile_id":0,"content_external_id":0,"type":"book"}`)

		assert.Equal(t, fiber.StatusBadRequest, response.StatusCode)

		var errorBody dto_shared.ErrorDto
		decodeFavoriteControllerJSON(t, body, &errorBody)
		assert.Equal(t, shared_errors_favorite.InvalidCreateFavoriteData, errorBody.Message)
		assert.Equal(t, dto_shared.BadRequest, errorBody.CodeMessage)
		assert.Len(t, errorBody.Details, 3)
		favoriteService.AssertNotCalled(t, "AddFavorite", mock.Anything, mock.Anything, mock.Anything)
	})

	t.Run("should return not found when profile does not exist", func(t *testing.T) {
		favoriteService := new(service_mock.FavoriteServiceMock)
		app := setupFavoriteControllerTestApp(favoriteService)

		favoriteService.On("AddFavorite", mock.Anything, uint(10), mock.AnythingOfType("dto_favorite.FavoriteRequestDto")).
			Return(shared_errors.ErrProfileNotFound).
			Once()

		response, body := performFavoriteControllerRequest(t, app, http.MethodPost, "/favorites", `{"profile_id":999,"content_external_id":101,"type":"movie"}`)

		assert.Equal(t, fiber.StatusNotFound, response.StatusCode)
		assertFavoriteControllerError(t, body, shared_errors_profile.ProfileNotFound, dto_shared.NotFound)
		favoriteService.AssertExpectations(t)
	})

	t.Run("should return conflict when favorite already exists", func(t *testing.T) {
		favoriteService := new(service_mock.FavoriteServiceMock)
		app := setupFavoriteControllerTestApp(favoriteService)

		favoriteService.On("AddFavorite", mock.Anything, uint(10), mock.AnythingOfType("dto_favorite.FavoriteRequestDto")).
			Return(shared_errors.ErrFavoriteAlreadyExists).
			Once()

		response, body := performFavoriteControllerRequest(t, app, http.MethodPost, "/favorites", `{"profile_id":20,"content_external_id":101,"type":"movie"}`)

		assert.Equal(t, fiber.StatusConflict, response.StatusCode)
		assertFavoriteControllerError(t, body, shared_errors_favorite.FavoriteAlreadyExists, dto_shared.Conflict)
		favoriteService.AssertExpectations(t)
	})

	t.Run("should return internal server error when service fails", func(t *testing.T) {
		favoriteService := new(service_mock.FavoriteServiceMock)
		app := setupFavoriteControllerTestApp(favoriteService)

		favoriteService.On("AddFavorite", mock.Anything, uint(10), mock.AnythingOfType("dto_favorite.FavoriteRequestDto")).
			Return(errors.New("service unavailable")).
			Once()

		response, body := performFavoriteControllerRequest(t, app, http.MethodPost, "/favorites", `{"profile_id":20,"content_external_id":101,"type":"movie"}`)

		assert.Equal(t, fiber.StatusInternalServerError, response.StatusCode)
		assertFavoriteControllerError(t, body, shared_errors_favorite.FailedToCreateFavorite, dto_shared.InternalServerError)
		favoriteService.AssertExpectations(t)
	})
}

func TestRemoveFavorite(t *testing.T) {
	t.Run("should remove favorite", func(t *testing.T) {
		favoriteService := new(service_mock.FavoriteServiceMock)
		app := setupFavoriteControllerTestApp(favoriteService)
		requestBody := `{"profile_id":20,"content_external_id":101,"type":"movie"}`
		expectedRequest := dto_favorite.FavoriteRequestDto{
			ProfileID:         20,
			ContentExternalID: 101,
			Type:              "movie",
		}

		favoriteService.On("DeleteFavorite", mock.Anything, uint(10), expectedRequest).
			Return(nil).
			Once()

		response, _ := performFavoriteControllerRequest(t, app, http.MethodDelete, "/favorites", requestBody)

		assert.Equal(t, fiber.StatusNoContent, response.StatusCode)
		favoriteService.AssertExpectations(t)
	})

	t.Run("should return bad request when body is invalid", func(t *testing.T) {
		favoriteService := new(service_mock.FavoriteServiceMock)
		app := setupFavoriteControllerTestApp(favoriteService)

		response, body := performFavoriteControllerRequest(t, app, http.MethodDelete, "/favorites", `{"profile_id":`)

		assert.Equal(t, fiber.StatusBadRequest, response.StatusCode)
		assertFavoriteControllerError(t, body, shared_errors.InvalidRequestBody, dto_shared.BadRequest)
		favoriteService.AssertNotCalled(t, "DeleteFavorite", mock.Anything, mock.Anything, mock.Anything)
	})

	t.Run("should return bad request when request validation fails", func(t *testing.T) {
		favoriteService := new(service_mock.FavoriteServiceMock)
		app := setupFavoriteControllerTestApp(favoriteService)

		response, body := performFavoriteControllerRequest(t, app, http.MethodDelete, "/favorites", `{"profile_id":0,"content_external_id":0,"type":"book"}`)

		assert.Equal(t, fiber.StatusBadRequest, response.StatusCode)
		assertFavoriteControllerError(t, body, shared_errors_favorite.InvalidCreateFavoriteData, dto_shared.BadRequest)
		favoriteService.AssertNotCalled(t, "DeleteFavorite", mock.Anything, mock.Anything, mock.Anything)
	})

	t.Run("should return not found when profile does not exist", func(t *testing.T) {
		favoriteService := new(service_mock.FavoriteServiceMock)
		app := setupFavoriteControllerTestApp(favoriteService)

		favoriteService.On("DeleteFavorite", mock.Anything, uint(10), mock.AnythingOfType("dto_favorite.FavoriteRequestDto")).
			Return(shared_errors.ErrProfileNotFound).
			Once()

		response, body := performFavoriteControllerRequest(t, app, http.MethodDelete, "/favorites", `{"profile_id":999,"content_external_id":101,"type":"movie"}`)

		assert.Equal(t, fiber.StatusNotFound, response.StatusCode)
		assertFavoriteControllerError(t, body, shared_errors_profile.ProfileNotFound, dto_shared.NotFound)
		favoriteService.AssertExpectations(t)
	})

	t.Run("should return not found when favorite does not exist", func(t *testing.T) {
		favoriteService := new(service_mock.FavoriteServiceMock)
		app := setupFavoriteControllerTestApp(favoriteService)

		favoriteService.On("DeleteFavorite", mock.Anything, uint(10), mock.AnythingOfType("dto_favorite.FavoriteRequestDto")).
			Return(shared_errors.ErrFavoriteNotFound).
			Once()

		response, body := performFavoriteControllerRequest(t, app, http.MethodDelete, "/favorites", `{"profile_id":20,"content_external_id":999,"type":"movie"}`)

		assert.Equal(t, fiber.StatusNotFound, response.StatusCode)
		assertFavoriteControllerError(t, body, shared_errors_favorite.FavoriteNotFound, dto_shared.NotFound)
		favoriteService.AssertExpectations(t)
	})

	t.Run("should return internal server error when service fails", func(t *testing.T) {
		favoriteService := new(service_mock.FavoriteServiceMock)
		app := setupFavoriteControllerTestApp(favoriteService)

		favoriteService.On("DeleteFavorite", mock.Anything, uint(10), mock.AnythingOfType("dto_favorite.FavoriteRequestDto")).
			Return(errors.New("service unavailable")).
			Once()

		response, body := performFavoriteControllerRequest(t, app, http.MethodDelete, "/favorites", `{"profile_id":20,"content_external_id":101,"type":"movie"}`)

		assert.Equal(t, fiber.StatusInternalServerError, response.StatusCode)
		assertFavoriteControllerError(t, body, shared_errors_favorite.FailedToDeleteFavorite, dto_shared.InternalServerError)
		favoriteService.AssertExpectations(t)
	})
}

func setupFavoriteControllerTestApp(favoriteService *service_mock.FavoriteServiceMock) *fiber.App {
	app := fiber.New()
	controller := v1_controller_favorite.NewFavoriteController(favoriteService)

	app.Get("/favorites", setFavoriteControllerUserID(10), controller.ListFavorites)
	app.Get("/favorites/without-user", controller.ListFavorites)
	app.Post("/favorites", setFavoriteControllerUserID(10), controller.AddFavorite)
	app.Delete("/favorites", setFavoriteControllerUserID(10), controller.RemoveFavorite)

	return app
}

func setFavoriteControllerUserID(userID uint) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		ctx.Locals("user_id", userID)
		return ctx.Next()
	}
}

func makeFavoriteControllerResponse() dto_favorite.FavoriteResponseDto {
	createdAt := time.Date(2026, 1, 1, 10, 0, 0, 0, time.UTC)
	updatedAt := time.Date(2026, 1, 1, 11, 0, 0, 0, time.UTC)

	return dto_favorite.FavoriteResponseDto{
		Data: []dto_favorite.FavoriteDto{
			{
				ID:        1,
				ProfileID: 20,
				Content: &dto_content.ContentDto{
					ExternalID:       101,
					Type:             "movie",
					Title:            "Movie Title",
					OriginalTitle:    "Original Movie",
					Description:      "Movie overview",
					OriginalLanguage: "en",
					ReleaseDate:      "2026-01-01",
					PosterPath:       "/poster.jpg",
					BackdropPath:     "/backdrop.jpg",
					GenreIDs:         []int{28, 12},
					Popularity:       91.5,
					VoteAverage:      8.7,
					VoteCount:        1234,
					Adult:            false,
				},
				CreatedAt: createdAt,
				UpdatedAt: updatedAt,
			},
		},
		Pagination: dto_shared.PaginationDto{
			Page:      2,
			PerPage:   3,
			PageCount: 4,
			Total:     10,
		},
	}
}

func performFavoriteControllerRequest(t *testing.T, app *fiber.App, method string, path string, body string) (*http.Response, string) {
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

func decodeFavoriteControllerJSON(t *testing.T, body string, target any) {
	t.Helper()

	require.NoError(t, json.Unmarshal([]byte(body), target))
}

func assertFavoriteControllerError(t *testing.T, body string, message string, codeMessage string) {
	t.Helper()

	var errorBody dto_shared.ErrorDto
	decodeFavoriteControllerJSON(t, body, &errorBody)
	assert.Equal(t, message, errorBody.Message)
	assert.Equal(t, codeMessage, errorBody.CodeMessage)
}
