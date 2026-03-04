package model

import "time"

type Session struct {
	ID        string `gorm:"primaryKey;size:64"`
	Role      string `gorm:"not null;size:16"` // "admin" or "guest"
	IP        string `gorm:"size:45"`
	UserAgent string `gorm:"size:512"`
	ExpiresAt time.Time `gorm:"not null;index"`
	CreatedAt time.Time
}
