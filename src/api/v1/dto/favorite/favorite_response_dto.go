package dto_favorite

import (
	"time"

	dto_content "github.com/KaueTTS/streaming_api/src/api/v1/dto/content"
	dto_shared "github.com/KaueTTS/streaming_api/src/api/v1/dto/shared"
)

type FavoriteDto struct {
	ID        uint                    `json:"id"`
	ProfileID uint                    `json:"profile_id"`
	Content   *dto_content.ContentDto `json:"content,omitempty"`
	CreatedAt time.Time               `json:"created_at"`
	UpdatedAt time.Time               `json:"updated_at"`
}

type FavoriteResponseDto struct {
	Data       []FavoriteDto            `json:"data"`
	Pagination dto_shared.PaginationDto `json:"pagination"`
}
