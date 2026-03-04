package model

import "time"

type ScanConflict struct {
	ID          uint   `gorm:"primaryKey" json:"id"`
	ScanJobID   uint   `gorm:"index" json:"scan_job_id"`
	FilePath    string `gorm:"not null" json:"file_path"`
	Reason      string `gorm:"not null;size:64" json:"reason"` // "duplicate_hash", "duplicate_path"
	ExistingID  uint   `json:"existing_id"`
	Resolution  string `gorm:"size:16" json:"resolution"` // "skip", "replace", "pending"
	CreatedAt   time.Time `json:"created_at"`
}
