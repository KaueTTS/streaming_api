package repository_http_tmdb

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	httptrace "github.com/DataDog/dd-trace-go/contrib/net/http/v2"
	env "github.com/KaueTTS/streaming_api/src/configs/env"
	dto "github.com/KaueTTS/streaming_api/src/repositories/http/tmdb/dto"
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
			Timeout: 10 * time.Second,
		},
	}
}

func (r *TMDBRepository) GetMovies(ctx context.Context) (dto.GetMovieResponseDto, error) {
	baseURL := fmt.Sprintf("%s/discover/movie", r.baseURL)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, baseURL, nil)
	if err != nil {
		return dto.GetMovieResponseDto{}, err
	}

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", r.token))
	req.Header.Add("Accept", "application/json")

	client := httptrace.WrapClient(r.httpClient)
	resp, err := client.Do(req)
	if err != nil {
		return dto.GetMovieResponseDto{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return dto.GetMovieResponseDto{}, fmt.Errorf("código de status inesperado: %d", resp.StatusCode)
	}

	var response dto.GetMovieResponseDto
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return dto.GetMovieResponseDto{}, err
	}

	return response, nil
}
