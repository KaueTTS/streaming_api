package models

import (
	"time"
)

type Favorite struct {
	ID        uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	ProfileID uint      `json:"profile_id" gorm:"not null;index;uniqueIndex:idx_profile_content"`
	ContentID uint      `json:"content_id" gorm:"not null;uniqueIndex:idx_profile_content"`
	Content   *Content  `json:"content,omitempty" gorm:"foreignKey:ContentID;references:ID"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}
