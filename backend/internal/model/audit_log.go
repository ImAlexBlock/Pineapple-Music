package model

import "time"

type AuditLog struct {
	ID        uint   `gorm:"primaryKey" json:"id"`
	Action    string `gorm:"not null;size:64;index" json:"action"`
	Role      string `gorm:"size:16" json:"role"`
	IP        string `gorm:"size:45" json:"ip"`
	Detail    string `gorm:"type:text" json:"detail"`
	CreatedAt time.Time `gorm:"index" json:"created_at"`
}
