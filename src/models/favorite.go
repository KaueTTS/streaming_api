package models

import (
	"time"

	"gorm.io/gorm"
)

type Favorite struct {
	ID                uint           `json:"id" gorm:"primaryKey;autoIncrement"`
	UserID            uint           `json:"user_id" gorm:"not null;index"`
	ProfileID         uint           `json:"profile_id" gorm:"not null;index;uniqueIndex:idx_favorite_profile_content"`
	ContentExternalId int            `json:"content_external_id" gorm:"not null;uniqueIndex:idx_favorite_profile_content"`
	Type              string         `json:"type" gorm:"not null;uniqueIndex:idx_favorite_profile_content"`
	CreatedAt         time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt         time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt         gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}
