package repository_http_tmdb_test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	env "github.com/KaueTTS/streaming_api/src/configs/env"
	repository_http_tmdb "github.com/KaueTTS/streaming_api/src/repositories/http/tmdb"
	tmdb_dto "github.com/KaueTTS/streaming_api/src/repositories/http/tmdb/dto"
	shared_constants_content "github.com/KaueTTS/streaming_api/src/shared/constants/content"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestListContents(t *testing.T) {
	t.Run("should list movie contents using filters", func(t *testing.T) {
		server := setupTMDBTestServer(t, func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, http.MethodGet, r.Method)
			assert.Equal(t, "/discover/movie", r.URL.Path)
			assertTMDBTestHeaders(t, r)

			query := r.URL.Query()
			assert.Equal(t, "2", query.Get("page"))
			assert.Equal(t, "pt-BR", query.Get("language"))
			assert.Equal(t, "popularity.desc", query.Get("sort_by"))
			assert.Equal(t, "28,12", query.Get("with_genres"))
			assert.Equal(t, "2026", query.Get("primary_release_year"))
			assert.Empty(t, query.Get("first_air_date_year"))
			assert.Equal(t, "false", query.Get("include_adult"))

			fmt.Fprint(w, `{
				"page": 2,
				"results": [
					{
						"adult": false,
						"backdrop_path": "/movie-backdrop.jpg",
						"genre_ids": [28, 12],
						"id": 101,
						"original_language": "en",
						"original_title": "Original Movie",
						"overview": "Movie overview",
						"popularity": 91.5,
						"poster_path": "/movie-poster.jpg",
						"release_date": "2026-01-01",
						"title": "Movie Title",
						"video": false,
						"vote_average": 8.7,
						"vote_count": 1234
					}
				],
				"total_pages": 5,
				"total_results": 50
			}`)
		})
		repository := newTMDBTestRepository(t, server.URL)
		filters := tmdb_dto.ContentFiltersDto{
			Type:       shared_constants_content.ContentTypeMovie,
			Page:       2,
			Language:   "pt-BR",
			SortBy:     "popularity.desc",
			WithGenres: "28,12",
			Year:       2026,
			IsKids:     true,
		}

		response, err := repository.ListContents(context.Background(), filters)

		require.NoError(t, err)
		assert.Equal(t, 2, response.Page)
		assert.Equal(t, 5, response.TotalPages)
		assert.Equal(t, 50, response.TotalResults)
		require.Len(t, response.Results, 1)
		assert.Equal(t, 101, response.Results[0].ID)
		assert.Equal(t, "Movie Title", response.Results[0].Title)
		assert.Equal(t, []int{28, 12}, response.Results[0].GenreIDs)
		assert.Equal(t, 8.7, response.Results[0].VoteAverage)
	})

	t.Run("should return error when response status is not ok", func(t *testing.T) {
		server := setupTMDBTestServer(t, func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/discover/movie", r.URL.Path)
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprint(w, `{"status_message":"Invalid API key"}`)
		})
		repository := newTMDBTestRepository(t, server.URL)
		filters := tmdb_dto.ContentFiltersDto{
			Type: shared_constants_content.ContentTypeMovie,
		}

		response, err := repository.ListContents(context.Background(), filters)

		assert.Empty(t, response)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "código de status inesperado: 401")
	})
}

func TestSearchContents(t *testing.T) {
	t.Run("should search tv contents using filters", func(t *testing.T) {
		server := setupTMDBTestServer(t, func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, http.MethodGet, r.Method)
			assert.Equal(t, "/search/tv", r.URL.Path)
			assertTMDBTestHeaders(t, r)

			query := r.URL.Query()
			assert.Equal(t, "dark", query.Get("query"))
			assert.Equal(t, "3", query.Get("page"))
			assert.Equal(t, "en-US", query.Get("language"))
			assert.Equal(t, "2020", query.Get("first_air_date_year"))
			assert.Empty(t, query.Get("primary_release_year"))
			assert.Equal(t, "true", query.Get("include_adult"))

			fmt.Fprint(w, `{
				"page": 3,
				"results": [
					{
						"adult": false,
						"backdrop_path": "/tv-backdrop.jpg",
						"genre_ids": [18],
						"id": 202,
						"original_language": "de",
						"original_name": "Original TV",
						"overview": "TV overview",
						"popularity": 77.2,
						"poster_path": "/tv-poster.jpg",
						"first_air_date": "2020-09-01",
						"name": "TV Name",
						"vote_average": 8.1,
						"vote_count": 987
					}
				],
				"total_pages": 4,
				"total_results": 40
			}`)
		})
		repository := newTMDBTestRepository(t, server.URL)
		filters := tmdb_dto.ContentFiltersDto{
			Type:     shared_constants_content.ContentTypeTV,
			Query:    "  dark  ",
			Page:     3,
			Language: "en-US",
			Year:     2020,
			IsKids:   false,
		}

		response, err := repository.SearchContents(context.Background(), filters)

		require.NoError(t, err)
		assert.Equal(t, 3, response.Page)
		assert.Equal(t, 4, response.TotalPages)
		assert.Equal(t, 40, response.TotalResults)
		require.Len(t, response.Results, 1)
		assert.Equal(t, 202, response.Results[0].ID)
		assert.Equal(t, "TV Name", response.Results[0].Name)
		assert.Equal(t, "2020-09-01", response.Results[0].FirstAirDate)
		assert.Equal(t, []int{18}, response.Results[0].GenreIDs)
	})

	t.Run("should return error when response body is invalid json", func(t *testing.T) {
		server := setupTMDBTestServer(t, func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/search/movie", r.URL.Path)
			fmt.Fprint(w, `{"page":`)
		})
		repository := newTMDBTestRepository(t, server.URL)
		filters := tmdb_dto.ContentFiltersDto{
			Type: shared_constants_content.ContentTypeMovie,
		}

		response, err := repository.SearchContents(context.Background(), filters)

		assert.Empty(t, response)
		assert.Error(t, err)
	})
}

func TestGetContentByID(t *testing.T) {
	t.Run("should return content by id", func(t *testing.T) {
		server := setupTMDBTestServer(t, func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, http.MethodGet, r.Method)
			assert.Equal(t, "/movie/101", r.URL.Path)
			assertTMDBTestHeaders(t, r)

			query := r.URL.Query()
			assert.Equal(t, "pt-BR", query.Get("language"))

			fmt.Fprint(w, `{
				"adult": false,
				"backdrop_path": "/movie-backdrop.jpg",
				"genres": [
					{"id": 28, "name": "Action"},
					{"id": 12, "name": "Adventure"}
				],
				"id": 101,
				"original_language": "en",
				"original_title": "Original Movie",
				"overview": "Movie overview",
				"popularity": 91.5,
				"poster_path": "/movie-poster.jpg",
				"release_date": "2026-01-01",
				"title": "Movie Title",
				"video": false,
				"vote_average": 8.7,
				"vote_count": 1234
			}`)
		})
		repository := newTMDBTestRepository(t, server.URL)

		content, err := repository.GetContentByID(context.Background(), shared_constants_content.ContentTypeMovie, 101, "  pt-BR  ")

		require.NoError(t, err)
		assert.Equal(t, 101, content.ID)
		assert.Equal(t, "Movie Title", content.Title)
		assert.Equal(t, "Original Movie", content.OriginalTitle)
		assert.Equal(t, "2026-01-01", content.ReleaseDate)
		require.Len(t, content.Genres, 2)
		assert.Equal(t, 28, content.Genres[0].ID)
		assert.Equal(t, "Action", content.Genres[0].Name)
		assert.Equal(t, 12, content.Genres[1].ID)
		assert.Equal(t, "Adventure", content.Genres[1].Name)
	})

	t.Run("should omit language query when language is blank", func(t *testing.T) {
		server := setupTMDBTestServer(t, func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/tv/202", r.URL.Path)
			assert.Empty(t, r.URL.Query().Get("language"))

			fmt.Fprint(w, `{
				"id": 202,
				"name": "TV Name",
				"original_name": "Original TV",
				"first_air_date": "2020-09-01"
			}`)
		})
		repository := newTMDBTestRepository(t, server.URL)

		content, err := repository.GetContentByID(context.Background(), shared_constants_content.ContentTypeTV, 202, "   ")

		require.NoError(t, err)
		assert.Equal(t, 202, content.ID)
		assert.Equal(t, "TV Name", content.Name)
		assert.Equal(t, "Original TV", content.OriginalName)
		assert.Equal(t, "2020-09-01", content.FirstAirDate)
	})

	t.Run("should return error when response status is not ok", func(t *testing.T) {
		server := setupTMDBTestServer(t, func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/movie/999", r.URL.Path)
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprint(w, `{"status_message":"Not found"}`)
		})
		repository := newTMDBTestRepository(t, server.URL)

		content, err := repository.GetContentByID(context.Background(), shared_constants_content.ContentTypeMovie, 999, "pt-BR")

		assert.Empty(t, content)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "código de status inesperado: 404")
	})
}

func setupTMDBTestServer(t *testing.T, handler http.HandlerFunc) *httptest.Server {
	t.Helper()

	server := httptest.NewServer(handler)
	t.Cleanup(server.Close)

	return server
}

func newTMDBTestRepository(t *testing.T, baseURL string) *repository_http_tmdb.TMDBRepository {
	t.Helper()

	originalBaseURL := env.TMDBBaseURL
	originalAccessToken := env.TMDBAccessToken
	env.TMDBBaseURL = baseURL
	env.TMDBAccessToken = "test-token"
	t.Cleanup(func() {
		env.TMDBBaseURL = originalBaseURL
		env.TMDBAccessToken = originalAccessToken
	})

	return repository_http_tmdb.NewTMDBRepository()
}

func assertTMDBTestHeaders(t *testing.T, r *http.Request) {
	t.Helper()

	assert.Equal(t, "Bearer test-token", r.Header.Get("Authorization"))
	assert.Equal(t, "application/json", r.Header.Get("Accept"))
}
