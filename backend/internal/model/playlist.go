package model

import "time"

type Playlist struct {
	ID        uint   `gorm:"primaryKey" json:"id"`
	Name      string `gorm:"not null" json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	Tracks []PlaylistTrack `gorm:"foreignKey:PlaylistID" json:"tracks,omitempty"`
}
