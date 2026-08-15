package dto_content

import dto_shared "github.com/KaueTTS/streaming_api/src/api/v1/dto/shared"

type ContentDto struct {
	ExternalID       int     `json:"external_id"`
	Type             string  `json:"type"`
	Title            string  `json:"title"`
	OriginalTitle    string  `json:"original_title"`
	Description      string  `json:"description"`
	OriginalLanguage string  `json:"original_language"`
	ReleaseDate      string  `json:"release_date"`
	PosterPath       string  `json:"poster_path"`
	BackdropPath     string  `json:"backdrop_path"`
	GenreIDs         []int   `json:"genre_ids"`
	Popularity       float64 `json:"popularity"`
	VoteAverage      float64 `json:"vote_average"`
	VoteCount        int     `json:"vote_count"`
	Adult            bool    `json:"adult"`
}

type ContentResponseDto struct {
	Data       []ContentDto             `json:"data"`
	Pagination dto_shared.PaginationDto `json:"pagination"`
}
