package dto_auth

import "time"

type UserResponseDto struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type AuthResponseDto struct {
	Token     string          `json:"token"`
	TokenType string          `json:"token_type"`
	ExpiresIn int64           `json:"expires_in"`
	User      UserResponseDto `json:"user"`
}
