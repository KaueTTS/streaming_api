package dto_profile

import "time"

type ProfileResponseDto struct {
	ID        uint      `json:"id"`
	UserID    uint      `json:"user_id"`
	Name      string    `json:"name"`
	AvatarURL *string   `json:"avatar_url"`
	IsKids    bool      `json:"is_kids"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
