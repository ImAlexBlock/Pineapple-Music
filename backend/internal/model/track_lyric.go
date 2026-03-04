package model

import "time"

type TrackLyric struct {
	ID        uint   `gorm:"primaryKey" json:"id"`
	TrackID   uint   `gorm:"uniqueIndex;not null" json:"track_id"`
	Type      string `gorm:"size:16;not null" json:"type"` // "plain" or "synced"
	Content   string `gorm:"type:text;not null" json:"content"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
