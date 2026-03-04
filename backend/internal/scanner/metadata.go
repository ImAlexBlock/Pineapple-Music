package scanner

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/dhowden/tag"
)

type TrackMeta struct {
	Title       string
	Artist      string
	Album       string
	AlbumArtist string
	Genre       string
	Year        int
	TrackNumber int
	DiscNumber  int
	Duration    float64
	Format      string
	Size        int64
	Bitrate     int
	SampleRate  int
	FilePath    string
	FileHash    string
	HasCover    bool
	HasLyrics   bool
	MTime       time.Time
	Lyrics      string
	LyricsType  string // "plain" or "synced"
}

// ExtractMetadata reads audio metadata from a file.
func ExtractMetadata(filePath string) (*TrackMeta, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("open file: %w", err)
	}
	defer f.Close()

	info, err := f.Stat()
	if err != nil {
		return nil, fmt.Errorf("stat file: %w", err)
	}

	meta := &TrackMeta{
		FilePath: filePath,
		Size:     info.Size(),
		MTime:    info.ModTime(),
		Format:   strings.TrimPrefix(strings.ToLower(filepath.Ext(filePath)), "."),
	}

	// Parse tags
	m, err := tag.ReadFrom(f)
	if err == nil {
		meta.Title = m.Title()
		meta.Artist = m.Artist()
		meta.Album = m.Album()
		meta.AlbumArtist = m.AlbumArtist()
		meta.Genre = m.Genre()
		meta.Year = m.Year()
		tn, _ := m.Track()
		meta.TrackNumber = tn
		dn, _ := m.Disc()
		meta.DiscNumber = dn
		meta.HasCover = m.Picture() != nil

		// Check for embedded lyrics
		if raw := m.Raw(); raw != nil {
			if lyrics, ok := raw["lyrics"]; ok {
				if s, ok := lyrics.(string); ok && s != "" {
					meta.Lyrics = s
					meta.HasLyrics = true
					if isLRCSynced(s) {
						meta.LyricsType = "synced"
					} else {
						meta.LyricsType = "plain"
					}
				}
			}
		}
	}

	// Default title to filename
	if meta.Title == "" {
		meta.Title = strings.TrimSuffix(filepath.Base(filePath), filepath.Ext(filePath))
	}

	return meta, nil
}

// HashFile computes the SHA-256 hash of a file.
func HashFile(filePath string) (string, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer f.Close()

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}

// isLRCSynced checks if the lyrics text contains LRC time tags.
func isLRCSynced(s string) bool {
	// Simple check: look for [mm:ss pattern
	return strings.Contains(s, "[0") || strings.Contains(s, "[1") ||
		strings.Contains(s, "[2") || strings.Contains(s, "[3")
}
