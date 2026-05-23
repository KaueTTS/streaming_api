package models

import (
	"time"

	"gorm.io/gorm"
)

type Profile struct {
	ID        uint           `json:"id" gorm:"primaryKey;autoIncrement"`
	UserID    uint           `json:"user_id" gorm:"not null;index"`
	Name      string         `json:"name" gorm:"type:text;not null"`
	AvatarURL *string        `json:"avatar_url"`
	IsKids    bool           `json:"is_kids"`
	CreatedAt time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}
