package models

import (
	"time"

	"gorm.io/gorm"
)

type Episode struct {
	ID              uint           `json:"id" gorm:"primaryKey;autoIncrement"`
	SeasonID        uint           `json:"season_id" gorm:"not null;index"`
	Title           string         `json:"title" gorm:"type:text;not null"`
	EpisodeNumber   int            `json:"episode_number" gorm:"not null"`
	DurationSeconds int            `json:"duration_seconds" gorm:"type:integer;not null;default:0"`
	VideoURL        *string        `json:"video_url"`
	CreatedAt       time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt       time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt       gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}
