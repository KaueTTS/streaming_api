package v1_controller_content_test

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	v1_controller_content "github.com/KaueTTS/streaming_api/src/api/v1/controllers/content"
	dto_content "github.com/KaueTTS/streaming_api/src/api/v1/dto/content"
	dto_shared "github.com/KaueTTS/streaming_api/src/api/v1/dto/shared"
	service_mock "github.com/KaueTTS/streaming_api/src/mocks/services"
	shared_errors "github.com/KaueTTS/streaming_api/src/shared/errors"
	shared_errors_auth "github.com/KaueTTS/streaming_api/src/shared/errors/auth"
	shared_errors_content "github.com/KaueTTS/streaming_api/src/shared/errors/content"
	shared_errors_profile "github.com/KaueTTS/streaming_api/src/shared/errors/profile"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestListContents(t *testing.T) {
	t.Run("should list contents", func(t *testing.T) {
		contentService := new(service_mock.ContentServiceMock)
		app := setupContentControllerTestApp(contentService)
		expectedRequest := dto_content.ContentListRequestDto{
			Type:      "movie",
			Page:      2,
			SortBy:    "popularity.desc",
			Genre:     "28,12",
			Language:  "pt-BR",
			Year:      2026,
			ProfileID: 20,
		}
		serviceResponse := makeContentControllerResponse("movie")

		contentService.On("ListContents", mock.Anything, uint(10), expectedRequest).
			Return(serviceResponse, nil).
			Once()

		response, body := performContentControllerRequest(t, app, "/contents?type=movie&page=2&sort_by=popularity.desc&with_genres=28,12&language=pt-BR&year=2026&profile_id=20")

		assert.Equal(t, fiber.StatusOK, response.StatusCode)

		var responseBody dto_content.ContentResponseDto
		decodeContentControllerJSON(t, body, &responseBody)
		assert.Equal(t, serviceResponse, responseBody)
		contentService.AssertExpectations(t)
	})

	t.Run("should return unauthorized when authenticated user id is missing", func(t *testing.T) {
		contentService := new(service_mock.ContentServiceMock)
		app := setupContentControllerTestApp(contentService)

		response, body := performContentControllerRequest(t, app, "/contents/without-user?type=movie&profile_id=20")

		assert.Equal(t, fiber.StatusUnauthorized, response.StatusCode)
		assertContentControllerError(t, body, shared_errors_auth.InvalidToken, dto_shared.Unauthorized)
		contentService.AssertNotCalled(t, "ListContents", mock.Anything, mock.Anything, mock.Anything)
	})

	t.Run("should return bad request when query validation fails", func(t *testing.T) {
		contentService := new(service_mock.ContentServiceMock)
		app := setupContentControllerTestApp(contentService)

		response, body := performContentControllerRequest(t, app, "/contents?type=book&page=-1&with_genres=abc&language=portuguese&year=-1&profile_id=0")

		assert.Equal(t, fiber.StatusBadRequest, response.StatusCode)

		var errorBody dto_shared.ErrorDto
		decodeContentControllerJSON(t, body, &errorBody)
		assert.Equal(t, shared_errors_content.InvalidContentQueryParameters, errorBody.Message)
		assert.Equal(t, dto_shared.BadRequest, errorBody.CodeMessage)
		assert.Len(t, errorBody.Details, 6)
		contentService.AssertNotCalled(t, "ListContents", mock.Anything, mock.Anything, mock.Anything)
	})

	t.Run("should return not found when profile does not exist", func(t *testing.T) {
		contentService := new(service_mock.ContentServiceMock)
		app := setupContentControllerTestApp(contentService)
		expectedRequest := dto_content.ContentListRequestDto{
			Type:      "movie",
			ProfileID: 999,
		}

		contentService.On("ListContents", mock.Anything, uint(10), expectedRequest).
			Return(dto_content.ContentResponseDto{}, shared_errors.ErrProfileNotFound).
			Once()

		response, body := performContentControllerRequest(t, app, "/contents?type=movie&profile_id=999")

		assert.Equal(t, fiber.StatusNotFound, response.StatusCode)
		assertContentControllerError(t, body, shared_errors_profile.ProfileNotFound, dto_shared.NotFound)
		contentService.AssertExpectations(t)
	})

	t.Run("should return bad gateway when service fails", func(t *testing.T) {
		contentService := new(service_mock.ContentServiceMock)
		app := setupContentControllerTestApp(contentService)

		contentService.On("ListContents", mock.Anything, uint(10), mock.AnythingOfType("dto_content.ContentListRequestDto")).
			Return(dto_content.ContentResponseDto{}, errors.New("tmdb unavailable")).
			Once()

		response, body := performContentControllerRequest(t, app, "/contents?type=movie&profile_id=20")

		assert.Equal(t, fiber.StatusBadGateway, response.StatusCode)
		assertContentControllerError(t, body, shared_errors_content.FailedToListContents, dto_shared.BadGateway)
		contentService.AssertExpectations(t)
	})
}

func TestSearchContents(t *testing.T) {
	t.Run("should search contents", func(t *testing.T) {
		contentService := new(service_mock.ContentServiceMock)
		app := setupContentControllerTestApp(contentService)
		expectedRequest := dto_content.ContentSearchRequestDto{
			Type:      "tv",
			Page:      2,
			Language:  "en-US",
			Query:     "dark",
			ProfileID: 20,
		}
		serviceResponse := makeContentControllerResponse("tv")

		contentService.On("SearchContents", mock.Anything, uint(10), expectedRequest).
			Return(serviceResponse, nil).
			Once()

		response, body := performContentControllerRequest(t, app, "/contents/search?type=tv&page=2&language=en-US&query=dark&profile_id=20")

		assert.Equal(t, fiber.StatusOK, response.StatusCode)

		var responseBody dto_content.ContentResponseDto
		decodeContentControllerJSON(t, body, &responseBody)
		assert.Equal(t, serviceResponse, responseBody)
		contentService.AssertExpectations(t)
	})

	t.Run("should return unauthorized when authenticated user id is missing", func(t *testing.T) {
		contentService := new(service_mock.ContentServiceMock)
		app := setupContentControllerTestApp(contentService)

		response, body := performContentControllerRequest(t, app, "/contents/search/without-user?type=movie&query=matrix&profile_id=20")

		assert.Equal(t, fiber.StatusUnauthorized, response.StatusCode)
		assertContentControllerError(t, body, shared_errors_auth.InvalidToken, dto_shared.Unauthorized)
		contentService.AssertNotCalled(t, "SearchContents", mock.Anything, mock.Anything, mock.Anything)
	})

	t.Run("should return bad request when query validation fails", func(t *testing.T) {
		contentService := new(service_mock.ContentServiceMock)
		app := setupContentControllerTestApp(contentService)

		response, body := performContentControllerRequest(t, app, "/contents/search?type=book&page=-1&language=portuguese&query=&profile_id=0")

		assert.Equal(t, fiber.StatusBadRequest, response.StatusCode)

		var errorBody dto_shared.ErrorDto
		decodeContentControllerJSON(t, body, &errorBody)
		assert.Equal(t, shared_errors_content.InvalidContentQueryParameters, errorBody.Message)
		assert.Equal(t, dto_shared.BadRequest, errorBody.CodeMessage)
		assert.Len(t, errorBody.Details, 5)
		contentService.AssertNotCalled(t, "SearchContents", mock.Anything, mock.Anything, mock.Anything)
	})

	t.Run("should return not found when profile does not exist", func(t *testing.T) {
		contentService := new(service_mock.ContentServiceMock)
		app := setupContentControllerTestApp(contentService)
		expectedRequest := dto_content.ContentSearchRequestDto{
			Type:      "movie",
			Query:     "matrix",
			ProfileID: 999,
		}

		contentService.On("SearchContents", mock.Anything, uint(10), expectedRequest).
			Return(dto_content.ContentResponseDto{}, shared_errors.ErrProfileNotFound).
			Once()

		response, body := performContentControllerRequest(t, app, "/contents/search?type=movie&query=matrix&profile_id=999")

		assert.Equal(t, fiber.StatusNotFound, response.StatusCode)
		assertContentControllerError(t, body, shared_errors_profile.ProfileNotFound, dto_shared.NotFound)
		contentService.AssertExpectations(t)
	})

	t.Run("should return bad gateway when service fails", func(t *testing.T) {
		contentService := new(service_mock.ContentServiceMock)
		app := setupContentControllerTestApp(contentService)

		contentService.On("SearchContents", mock.Anything, uint(10), mock.AnythingOfType("dto_content.ContentSearchRequestDto")).
			Return(dto_content.ContentResponseDto{}, errors.New("tmdb unavailable")).
			Once()

		response, body := performContentControllerRequest(t, app, "/contents/search?type=movie&query=matrix&profile_id=20")

		assert.Equal(t, fiber.StatusBadGateway, response.StatusCode)
		assertContentControllerError(t, body, shared_errors_content.FailedToSearchContents, dto_shared.BadGateway)
		contentService.AssertExpectations(t)
	})
}

func setupContentControllerTestApp(contentService *service_mock.ContentServiceMock) *fiber.App {
	app := fiber.New()
	controller := v1_controller_content.NewContentController(contentService)

	app.Get("/contents", setContentControllerUserID(10), controller.ListContents)
	app.Get("/contents/without-user", controller.ListContents)
	app.Get("/contents/search", setContentControllerUserID(10), controller.SearchContents)
	app.Get("/contents/search/without-user", controller.SearchContents)

	return app
}

func setContentControllerUserID(userID uint) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		ctx.Locals("user_id", userID)
		return ctx.Next()
	}
}

func makeContentControllerResponse(contentType string) dto_content.ContentResponseDto {
	return dto_content.ContentResponseDto{
		Data: []dto_content.ContentDto{
			{
				ExternalID:       101,
				Type:             contentType,
				Title:            "Content Title",
				OriginalTitle:    "Original Content",
				Description:      "Content overview",
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
		},
		Pagination: dto_shared.PaginationDto{
			Page:      2,
			PerPage:   1,
			PageCount: 5,
			Total:     50,
		},
	}
}

func performContentControllerRequest(t *testing.T, app *fiber.App, path string) (*http.Response, string) {
	t.Helper()

	request := httptest.NewRequest(http.MethodGet, path, nil)

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

func decodeContentControllerJSON(t *testing.T, body string, target any) {
	t.Helper()

	require.NoError(t, json.Unmarshal([]byte(body), target))
}

func assertContentControllerError(t *testing.T, body string, message string, codeMessage string) {
	t.Helper()

	var errorBody dto_shared.ErrorDto
	decodeContentControllerJSON(t, body, &errorBody)
	assert.Equal(t, message, errorBody.Message)
	assert.Equal(t, codeMessage, errorBody.CodeMessage)
}
