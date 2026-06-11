package dto_profile

import (
	"time"

	dto_shared "github.com/KaueTTS/streaming_api/src/api/v1/dto/shared"
)

type ProfileDto struct {
	ID        uint      `json:"id"`
	UserID    uint      `json:"user_id"`
	Name      string    `json:"name"`
	AvatarURL *string   `json:"avatar_url"`
	IsKids    bool      `json:"is_kids"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ProfileResponseDto struct {
	Data       []ProfileDto             `json:"data"`
	Pagination dto_shared.PaginationDto `json:"pagination"`
}
