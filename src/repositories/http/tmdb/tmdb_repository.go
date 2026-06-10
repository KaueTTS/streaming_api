package repository_http_tmdb

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	env "github.com/KaueTTS/streaming_api/src/configs/env"
	dto "github.com/KaueTTS/streaming_api/src/repositories/http/tmdb/dto"
	shared_constants_content "github.com/KaueTTS/streaming_api/src/shared/constants/content"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

type TMDBRepository struct {
	baseURL    string
	token      string
	httpClient *http.Client
}

func NewTMDBRepository() *TMDBRepository {
	return &TMDBRepository{
		baseURL: env.TMDBBaseURL,
		token:   env.TMDBAccessToken,
		httpClient: &http.Client{
			Timeout:   10 * time.Second,
			Transport: otelhttp.NewTransport(http.DefaultTransport),
		},
	}
}

func (r *TMDBRepository) ListContents(ctx context.Context, filters dto.ContentListFiltersDto) (dto.GetContentResponseDto, error) {
	baseURL := fmt.Sprintf("%s/discover/%s", r.baseURL, filters.Type)

	queryParams := url.Values{}
	addQueryParam(queryParams, "language", filters.Language)
	addQueryParam(queryParams, "sort_by", filters.SortBy)
	addQueryParam(queryParams, "with_genres", filters.WithGenres)

	if filters.Page > 0 {
		queryParams.Set("page", strconv.Itoa(filters.Page))
	}

	if filters.Year > 0 {
		switch filters.Type {
		case shared_constants_content.ContentTypeMovie:
			queryParams.Set("primary_release_year", strconv.Itoa(filters.Year))
		case shared_constants_content.ContentTypeTV:
			queryParams.Set("first_air_date_year", strconv.Itoa(filters.Year))
		}
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, baseURL, nil)
	if err != nil {
		return dto.GetContentResponseDto{}, err
	}
	req.URL.RawQuery = queryParams.Encode()

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", r.token))
	req.Header.Add("Accept", "application/json")

	resp, err := r.httpClient.Do(req)
	if err != nil {
		return dto.GetContentResponseDto{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return dto.GetContentResponseDto{}, fmt.Errorf("código de status inesperado: %d", resp.StatusCode)
	}

	var response dto.GetContentResponseDto
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return dto.GetContentResponseDto{}, err
	}

	return response, nil
}

func (r *TMDBRepository) SearchContents(ctx context.Context, filters dto.ContentSearchFiltersDto) (dto.GetContentResponseDto, error) {
	baseURL := fmt.Sprintf("%s/search/%s", r.baseURL, filters.Type)
	queryParams := url.Values{}
	queryParams.Set("query", strings.TrimSpace(filters.Query))
	addQueryParam(queryParams, "language", filters.Language)

	if filters.Page > 0 {
		queryParams.Set("page", strconv.Itoa(filters.Page))
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, baseURL, nil)
	if err != nil {
		return dto.GetContentResponseDto{}, err
	}
	req.URL.RawQuery = queryParams.Encode()

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", r.token))
	req.Header.Add("Accept", "application/json")

	resp, err := r.httpClient.Do(req)
	if err != nil {
		return dto.GetContentResponseDto{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return dto.GetContentResponseDto{}, fmt.Errorf("código de status inesperado: %d", resp.StatusCode)
	}

	var response dto.GetContentResponseDto
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return dto.GetContentResponseDto{}, err
	}

	return response, nil
}

func addQueryParam(queryParams url.Values, key string, value string) {
	if value != "" {
		queryParams.Set(key, value)
	}
}
