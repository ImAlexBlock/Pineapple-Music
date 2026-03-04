package model

import "time"

type ScanJob struct {
	ID         uint   `gorm:"primaryKey" json:"id"`
	Status     string `gorm:"not null;size:16" json:"status"` // "running", "completed", "failed"
	Total      int    `json:"total"`
	Scanned    int    `json:"scanned"`
	Added      int    `json:"added"`
	Updated    int    `json:"updated"`
	Errors     int    `json:"errors"`
	ErrorLog   string `gorm:"type:text" json:"error_log,omitempty"`
	StartedAt  time.Time `json:"started_at"`
	FinishedAt *time.Time `json:"finished_at,omitempty"`
}
