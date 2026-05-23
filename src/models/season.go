package models

import (
	"time"

	"gorm.io/gorm"
)

type Season struct {
	ID           uint           `json:"id" gorm:"primaryKey;autoIncrement"`
	ContentID    uint           `json:"content_id" gorm:"not null;index"`
	SeasonNumber int            `json:"season_number" gorm:"not null"`
	Episodes     []Episode      `json:"episodes,omitempty" gorm:"foreignKey:SeasonID"`
	CreatedAt    time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt    time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt    gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}
