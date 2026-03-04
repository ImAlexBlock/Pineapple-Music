package scanner

import (
	"pineapple-music/internal/model"

	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// CheckConflict checks if a file already exists in the database.
func CheckConflict(db *gorm.DB, meta *TrackMeta) (*model.Track, string) {
	silent := db.Session(&gorm.Session{Logger: db.Logger.LogMode(logger.Silent)})

	// Check by file path
	var existing model.Track
	if err := silent.Where("file_path = ?", meta.FilePath).First(&existing).Error; err == nil {
		if existing.FileHash == meta.FileHash {
			return &existing, "unchanged"
		}
		return &existing, "updated"
	}

	// Check by hash for duplicates
	if meta.FileHash != "" {
		if err := silent.Where("file_hash = ?", meta.FileHash).First(&existing).Error; err == nil {
			return &existing, "duplicate_hash"
		}
	}

	return nil, "new"
}
