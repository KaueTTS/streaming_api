package models

import (
	"time"

	"gorm.io/gorm"
)

type Content struct {
	ID           uint           `json:"id" gorm:"primaryKey;autoIncrement"`
	Title        string         `json:"title" gorm:"type:text;not null;index"`
	Description  *string        `json:"description"`
	Type         string         `json:"type" gorm:"type:text;not null;index"`
	ReleaseYear  *int           `json:"release_year"`
	AgeRating    *int           `json:"age_rating"`
	ThumbnailURL *string        `json:"thumbnail_url"`
	BannerURL    *string        `json:"banner_url"`
	CreatedAt    time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt    time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt    gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}
