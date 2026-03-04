package model

import "time"

type PlayEvent struct {
	ID        uint   `gorm:"primaryKey" json:"id"`
	TrackID   uint   `gorm:"not null;index" json:"track_id"`
	SessionID string `gorm:"size:64" json:"session_id"`
	IP        string `gorm:"size:45" json:"ip"`
	PlayedAt  time.Time `gorm:"not null;index" json:"played_at"`
}
