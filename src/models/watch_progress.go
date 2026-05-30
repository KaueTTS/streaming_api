package models

import (
	"time"

	"gorm.io/gorm"
)

type WatchProgress struct {
	ID              uint           `json:"id" gorm:"primaryKey;autoIncrement"`
	ProfileID       uint           `json:"profile_id" gorm:"not null;index;uniqueIndex:idx_profile_progress"`
	ContentID       uint           `json:"content_id" gorm:"not null;index;uniqueIndex:idx_profile_progress"`
	EpisodeID       *uint          `json:"episode_id,omitempty" gorm:"index;uniqueIndex:idx_profile_progress"`
	ProgressSeconds int            `json:"progress_seconds" gorm:"type:integer;not null;default:0"`
	Completed       bool           `json:"completed" gorm:"not null;default:false"`
	Content         *Content       `json:"content,omitempty" gorm:"foreignKey:ContentID;references:ID"`
	CreatedAt       time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt       time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt       gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}
