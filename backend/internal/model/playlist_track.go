package model

import "time"

type PlaylistTrack struct {
	ID         uint `gorm:"primaryKey" json:"id"`
	PlaylistID uint `gorm:"not null;index:idx_playlist_position,unique" json:"playlist_id"`
	TrackID    uint `gorm:"not null" json:"track_id"`
	Position   int  `gorm:"not null;index:idx_playlist_position,unique" json:"position"`
	CreatedAt  time.Time `json:"created_at"`

	Track Track `gorm:"foreignKey:TrackID" json:"track,omitempty"`
}
