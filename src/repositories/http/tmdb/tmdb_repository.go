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

// ListContents lista os conteúdos da TMBD baseado nos filtros passados pelo usuário
func (r *TMDBRepository) ListContents(ctx context.Context, filters dto.ContentFiltersDto) (dto.GetContentResponseDto, error) {
	baseURL := fmt.Sprintf("%s/discover/%s", r.baseURL, filters.Type)

	queryParams := url.Values{}
	if err := buildQueryParam(queryParams, filters); err != nil {
		return dto.GetContentResponseDto{}, err
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

// SearchContents busca os conteúdos da TMBD baseado nos filtros passados pelo usuário
func (r *TMDBRepository) SearchContents(ctx context.Context, filters dto.ContentFiltersDto) (dto.GetContentResponseDto, error) {
	baseURL := fmt.Sprintf("%s/search/%s", r.baseURL, filters.Type)

	queryParams := url.Values{}
	if err := buildQueryParam(queryParams, filters); err != nil {
		return dto.GetContentResponseDto{}, err
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

// GetContentByID busca um conteúdo da TMBD baseado no tipo, ID externo e idioma passados pelo usuário
func (r *TMDBRepository) GetContentByID(ctx context.Context, contentType string, contentExternalID int, language string) (dto.ContentDto, error) {
	baseURL := fmt.Sprintf("%s/%s/%d", r.baseURL, contentType, contentExternalID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, baseURL, nil)
	if err != nil {
		return dto.ContentDto{}, err
	}

	queryParams := url.Values{}
	if strings.TrimSpace(language) != "" {
		queryParams.Set("language", strings.TrimSpace(language))
	}
	req.URL.RawQuery = queryParams.Encode()

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", r.token))
	req.Header.Add("Accept", "application/json")

	resp, err := r.httpClient.Do(req)
	if err != nil {
		return dto.ContentDto{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return dto.ContentDto{}, fmt.Errorf("código de status inesperado: %d", resp.StatusCode)
	}

	var response dto.ContentDto
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return dto.ContentDto{}, err
	}

	return response, nil
}

// buildQueryParam monta os parâmetros da query string de requisição
func buildQueryParam(queryParams url.Values, filters dto.ContentFiltersDto) error {
	if filters.Query != "" {
		queryParams.Set("query", strings.TrimSpace(filters.Query))
	}

	if filters.Page > 0 {
		queryParams.Set("page", strconv.Itoa(filters.Page))
	}

	if filters.Language != "" {
		queryParams.Set("language", filters.Language)
	}

	if filters.SortBy != "" {
		queryParams.Set("sort_by", filters.SortBy)
	}

	if filters.WithGenres != "" {
		queryParams.Set("with_genres", filters.WithGenres)
	}

	if filters.Year > 0 {
		switch filters.Type {
		case shared_constants_content.ContentTypeMovie:
			queryParams.Set("primary_release_year", strconv.Itoa(filters.Year))
		case shared_constants_content.ContentTypeTV:
			queryParams.Set("first_air_date_year", strconv.Itoa(filters.Year))
		}
	}

	queryParams.Set("include_adult", "false")
	if !filters.IsKids {
		queryParams.Set("include_adult", "true")
	}

	return nil
}
