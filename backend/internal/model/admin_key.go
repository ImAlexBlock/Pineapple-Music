package model

import "time"

type AdminKey struct {
	ID        uint   `gorm:"primaryKey"`
	KeyHash   string `gorm:"not null"`
	MD5Hash   string `gorm:"size:32"` // md5(plaintext) for Subsonic token auth
	CreatedAt time.Time
}
