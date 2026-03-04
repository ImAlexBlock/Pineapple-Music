package api

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"pineapple-music/internal/config"
	"pineapple-music/internal/model"
	"pineapple-music/internal/scanner"
	"pineapple-music/internal/service"
	"pineapple-music/internal/util"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// safeFilename strips directory components and non-safe characters.
var unsafeChars = regexp.MustCompile(`[^\w\-. ]`)

func sanitizeFilename(name string) string {
	// Use only the base name (no directory traversal)
	name = filepath.Base(name)
	// Remove any remaining unsafe characters
	ext := filepath.Ext(name)
	base := strings.TrimSuffix(name, ext)
	base = unsafeChars.ReplaceAllString(base, "_")
	if base == "" || base == "." || base == "_" {
		base = "upload"
	}
	return base + ext
}

func Upload(cfg *config.Config, scanSvc *service.ScanService, auditSvc *service.AuditService, db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		file, err := c.FormFile("file")
		if err != nil {
			util.ErrorResponse(c, http.StatusBadRequest, "no_file", "No file uploaded")
			return
		}

		if file.Size > cfg.MaxUploadSize {
			util.ErrorResponse(c, http.StatusRequestEntityTooLarge, "file_too_large",
				fmt.Sprintf("File exceeds %dMB limit", cfg.MaxUploadSize/1024/1024))
			return
		}

		// Validate extension
		ext := strings.ToLower(filepath.Ext(file.Filename))
		validExts := map[string]bool{".mp3": true, ".flac": true, ".ogg": true, ".m4a": true, ".wav": true, ".aac": true}
		if !validExts[ext] {
			util.ErrorResponse(c, http.StatusBadRequest, "invalid_format", "Unsupported audio format")
			return
		}

		// Sanitize filename to prevent directory traversal
		safeName := sanitizeFilename(file.Filename)
		destPath := filepath.Join(cfg.MusicDir(), safeName)

		// Verify the resolved path is still inside MusicDir
		absMusic, _ := filepath.Abs(cfg.MusicDir())
		absDest, _ := filepath.Abs(destPath)
		if !strings.HasPrefix(absDest, absMusic+string(os.PathSeparator)) && absDest != absMusic {
			util.ErrorResponse(c, http.StatusBadRequest, "invalid_path", "Invalid file path")
			return
		}

		// Auto-rename if exists
		if _, err := os.Stat(destPath); err == nil {
			base := strings.TrimSuffix(safeName, ext)
			for i := 1; ; i++ {
				destPath = filepath.Join(cfg.MusicDir(), fmt.Sprintf("%s_%d%s", base, i, ext))
				if _, err := os.Stat(destPath); os.IsNotExist(err) {
					break
				}
				if i > 1000 {
					util.ErrorResponse(c, http.StatusConflict, "too_many_copies", "Too many copies of this file")
					return
				}
			}
		}

		if err := c.SaveUploadedFile(file, destPath); err != nil {
			util.ErrorResponse(c, http.StatusInternalServerError, "save_failed", "Failed to save file")
			return
		}

		// Extract metadata and insert into database
		meta, err := scanner.ExtractMetadata(destPath)
		if err != nil {
			// File saved but metadata failed — still insert with basic info
			track := &model.Track{
				Title:    strings.TrimSuffix(safeName, ext),
				Format:   strings.TrimPrefix(ext, "."),
				Size:     file.Size,
				FilePath: destPath,
			}
			db.Create(track)
			auditSvc.Log("upload", "admin", c.ClientIP(), safeName)
			c.JSON(http.StatusOK, gin.H{
				"message": "File uploaded (metadata extraction failed)",
				"track": gin.H{
					"id":    track.ID,
					"title": track.Title,
				},
			})
			return
		}

		hash, _ := scanner.HashFile(destPath)
		meta.FileHash = hash

		// Check for external .lrc
		if !meta.HasLyrics {
			lrcPath := strings.TrimSuffix(destPath, ext) + ".lrc"
			if lrcContent, err := scanner.ReadLRCFile(lrcPath); err == nil && lrcContent != "" {
				meta.Lyrics = lrcContent
				meta.HasLyrics = true
				if _, synced := scanner.ParseLRC(lrcContent); synced {
					meta.LyricsType = "synced"
				} else {
					meta.LyricsType = "plain"
				}
			}
		}

		// Insert track into DB
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
			FilePath:    destPath,
			FileHash:    meta.FileHash,
			HasCover:    meta.HasCover,
			HasLyrics:   meta.HasLyrics,
			MTime:       meta.MTime,
		}
		if err := db.Create(track).Error; err != nil {
			util.ErrorResponse(c, http.StatusInternalServerError, "db_error", "Failed to save track to database")
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

		auditSvc.Log("upload", "admin", c.ClientIP(), safeName)

		c.JSON(http.StatusOK, gin.H{
			"message": "File uploaded and added to library",
			"track": gin.H{
				"id":     track.ID,
				"title":  track.Title,
				"artist": track.Artist,
				"album":  track.Album,
			},
		})
	}
}

func StartScan(scanSvc *service.ScanService, auditSvc *service.AuditService) gin.HandlerFunc {
	return func(c *gin.Context) {
		job, err := scanSvc.StartScan()
		if err != nil {
			util.ErrorResponse(c, http.StatusConflict, "scan_running", err.Error())
			return
		}

		auditSvc.Log("scan_started", "admin", c.ClientIP(), "")

		c.JSON(http.StatusOK, gin.H{
			"message": "Scan started",
			"job_id":  job.ID,
		})
	}
}

func ScanStatus(scanSvc *service.ScanService) gin.HandlerFunc {
	return func(c *gin.Context) {
		job, err := scanSvc.GetLatestJob()
		if err != nil {
			util.ErrorResponse(c, http.StatusNotFound, "no_scan", "No scan jobs found")
			return
		}
		c.JSON(http.StatusOK, job)
	}
}
