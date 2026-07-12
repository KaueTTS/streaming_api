package models

import (
	"time"
)

type Favorite struct {
	ID          uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	ProfileID   uint      `json:"profile_id" gorm:"not null;index;uniqueIndex:idx_profile_content"`
	Profile     Profile   `json:"-" gorm:"foreignKey:ProfileID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	ContentID   int       `json:"content_id" gorm:"not null;uniqueIndex:idx_profile_content"`
	ContentType string    `json:"content_type" gorm:"type:text;not null;uniqueIndex:idx_profile_content"`
	CreatedAt   time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}
