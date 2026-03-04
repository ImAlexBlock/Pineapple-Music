package model

import "time"

type Track struct {
	ID          uint    `gorm:"primaryKey" json:"id"`
	Title       string  `gorm:"not null;index" json:"title"`
	Artist      string  `gorm:"index" json:"artist"`
	Album       string  `gorm:"index" json:"album"`
	AlbumArtist string  `json:"album_artist"`
	Genre       string  `json:"genre"`
	Year        int     `json:"year"`
	TrackNumber int     `json:"track_number"`
	DiscNumber  int     `json:"disc_number"`
	Duration    float64 `json:"duration"` // seconds
	Format      string  `gorm:"size:16" json:"format"`
	Size        int64   `json:"size"` // bytes
	Bitrate     int     `json:"bitrate"`
	SampleRate  int     `json:"sample_rate"`
	FilePath    string  `gorm:"uniqueIndex;not null" json:"-"` // internal only
	FileHash    string  `gorm:"index;size:64" json:"-"`        // internal only
	HasCover    bool    `json:"has_cover"`
	HasLyrics   bool    `json:"has_lyrics"`
	MTime       time.Time `json:"-"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
