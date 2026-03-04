package scanner

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"pineapple-music/internal/model"

	"gorm.io/gorm"
)

var (
	scanMu   sync.Mutex
	scanning bool
)

var supportedExts = map[string]bool{
	".mp3":  true,
	".flac": true,
	".ogg":  true,
	".m4a":  true,
	".wav":  true,
	".wma":  true,
	".aac":  true,
}

// Scan walks the music directory and imports tracks into the database.
func Scan(db *gorm.DB, musicDir string, job *model.ScanJob) {
	if !scanMu.TryLock() {
		job.Status = "failed"
		job.ErrorLog = "Another scan is already running"
		db.Save(job)
		return
	}
	defer scanMu.Unlock()
	scanning = true
	defer func() { scanning = false }()

	job.Status = "running"
	job.StartedAt = time.Now()
	db.Save(job)

	// Collect files first
	var files []string
	filepath.Walk(musicDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if info.IsDir() {
			return nil
		}
		ext := strings.ToLower(filepath.Ext(path))
		if supportedExts[ext] {
			files = append(files, path)
		}
		return nil
	})

	job.Total = len(files)
	db.Save(job)

	var errorLines []string

	for _, filePath := range files {
		job.Scanned++

		meta, err := ExtractMetadata(filePath)
		if err != nil {
			job.Errors++
			errorLines = append(errorLines, fmt.Sprintf("metadata error: %s: %v", filePath, err))
			continue
		}

		// Compute hash
		hash, err := HashFile(filePath)
		if err != nil {
			job.Errors++
			errorLines = append(errorLines, fmt.Sprintf("hash error: %s: %v", filePath, err))
			continue
		}
		meta.FileHash = hash

		// Check for external .lrc file
		if !meta.HasLyrics {
			lrcPath := strings.TrimSuffix(filePath, filepath.Ext(filePath)) + ".lrc"
			if lrcContent, err := ReadLRCFile(lrcPath); err == nil && lrcContent != "" {
				meta.Lyrics = lrcContent
				meta.HasLyrics = true
				if _, synced := ParseLRC(lrcContent); synced {
					meta.LyricsType = "synced"
				} else {
					meta.LyricsType = "plain"
				}
			}
		}

		existing, conflict := CheckConflict(db, meta)

		switch conflict {
		case "unchanged":
			// Skip
			continue
		case "updated":
			// Update existing track
			updateTrack(db, existing, meta)
			job.Updated++
		case "duplicate_hash":
			// Record conflict
			db.Create(&model.ScanConflict{
				ScanJobID:  job.ID,
				FilePath:   filePath,
				Reason:     "duplicate_hash",
				ExistingID: existing.ID,
				Resolution: "skip",
			})
		case "new":
			// Insert new track
			insertTrack(db, meta)
			job.Added++
		}

		// Save progress periodically
		if job.Scanned%10 == 0 {
			db.Save(job)
		}
	}

	now := time.Now()
	job.Status = "completed"
	job.FinishedAt = &now
	if len(errorLines) > 0 {
		job.ErrorLog = strings.Join(errorLines, "\n")
	}
	db.Save(job)

	log.Printf("Scan complete: %d total, %d added, %d updated, %d errors",
		job.Total, job.Added, job.Updated, job.Errors)
}

func insertTrack(db *gorm.DB, meta *TrackMeta) {
	track := &model.Track{
		Title:       meta.Title,
		Artist:      meta.Artist,
		Album:       meta.Album,
		AlbumArtist: meta.AlbumArtist,
		Genre:       meta.Genre,
		Year:        meta.Year,
		TrackNumber: meta.TrackNumber,
		DiscNumber:  meta.DiscNumber,
		Duration:    meta.Duration,
		Format:      meta.Format,
		Size:        meta.Size,
		Bitrate:     meta.Bitrate,
		SampleRate:  meta.SampleRate,
		FilePath:    meta.FilePath,
		FileHash:    meta.FileHash,
		HasCover:    meta.HasCover,
		HasLyrics:   meta.HasLyrics,
		MTime:       meta.MTime,
	}

	if err := db.Create(track).Error; err != nil {
		log.Printf("Failed to insert track %s: %v", meta.FilePath, err)
		return
	}

	// Save lyrics if present
	if meta.HasLyrics && meta.Lyrics != "" {
		db.Create(&model.TrackLyric{
			TrackID: track.ID,
			Type:    meta.LyricsType,
			Content: meta.Lyrics,
		})
	}
}

func updateTrack(db *gorm.DB, existing *model.Track, meta *TrackMeta) {
	existing.Title = meta.Title
	existing.Artist = meta.Artist
	existing.Album = meta.Album
	existing.AlbumArtist = meta.AlbumArtist
	existing.Genre = meta.Genre
	existing.Year = meta.Year
	existing.TrackNumber = meta.TrackNumber
	existing.DiscNumber = meta.DiscNumber
	existing.Duration = meta.Duration
	existing.Format = meta.Format
	existing.Size = meta.Size
	existing.Bitrate = meta.Bitrate
	existing.SampleRate = meta.SampleRate
	existing.FileHash = meta.FileHash
	existing.HasCover = meta.HasCover
	existing.HasLyrics = meta.HasLyrics
	existing.MTime = meta.MTime

	db.Save(existing)

	// Update lyrics
	if meta.HasLyrics && meta.Lyrics != "" {
		var lyric model.TrackLyric
		if err := db.Where("track_id = ?", existing.ID).First(&lyric).Error; err == nil {
			lyric.Type = meta.LyricsType
			lyric.Content = meta.Lyrics
			db.Save(&lyric)
		} else {
			db.Create(&model.TrackLyric{
				TrackID: existing.ID,
				Type:    meta.LyricsType,
				Content: meta.Lyrics,
			})
		}
	}
}

// IsScanning returns whether a scan is currently in progress.
func IsScanning() bool {
	return scanning
}
