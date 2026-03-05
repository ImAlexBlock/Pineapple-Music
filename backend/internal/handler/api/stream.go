package api

import (
	"crypto/sha256"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"sync"

	"pineapple-music/internal/model"
	"pineapple-music/internal/util"

	"github.com/dhowden/tag"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// coverCache provides an in-memory LRU-ish cache for extracted cover art.
// This avoids re-reading and parsing audio files for every cover request.
var coverCache = struct {
	sync.RWMutex
	items map[string]*coverEntry
}{items: make(map[string]*coverEntry)}

type coverEntry struct {
	data        []byte
	contentType string
}

const maxCoverCacheSize = 200

func getCachedCover(key string) *coverEntry {
	coverCache.RLock()
	defer coverCache.RUnlock()
	return coverCache.items[key]
}

func putCachedCover(key string, entry *coverEntry) {
	coverCache.Lock()
	defer coverCache.Unlock()
	// Simple eviction: if cache is full, clear half
	if len(coverCache.items) >= maxCoverCacheSize {
		i := 0
		for k := range coverCache.items {
			delete(coverCache.items, k)
			i++
			if i >= maxCoverCacheSize/2 {
				break
			}
		}
	}
	coverCache.items[key] = entry
}

func coverCacheKey(trackID uint, filePath string) string {
	h := sha256.Sum256([]byte(fmt.Sprintf("%d:%s", trackID, filepath.Clean(filePath))))
	return fmt.Sprintf("%x", h[:8])
}

func StreamTrack(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.ParseUint(c.Param("id"), 10, 32)
		if err != nil {
			util.ErrorResponse(c, http.StatusBadRequest, "invalid_id", "Invalid track ID")
			return
		}

		var track model.Track
		if err := db.First(&track, id).Error; err != nil {
			util.ErrorResponse(c, http.StatusNotFound, "not_found", "Track not found")
			return
		}

		if _, err := os.Stat(track.FilePath); os.IsNotExist(err) {
			util.ErrorResponse(c, http.StatusNotFound, "file_missing", "Audio file not found")
			return
		}

		// http.ServeFile handles Range requests automatically
		c.Header("Accept-Ranges", "bytes")
		http.ServeFile(c.Writer, c.Request, track.FilePath)
	}
}

func TrackCover(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.ParseUint(c.Param("id"), 10, 32)
		if err != nil {
			util.ErrorResponse(c, http.StatusBadRequest, "invalid_id", "Invalid track ID")
			return
		}

		var track model.Track
		if err := db.First(&track, id).Error; err != nil {
			util.ErrorResponse(c, http.StatusNotFound, "not_found", "Track not found")
			return
		}

		// Check cache first
		cacheKey := coverCacheKey(track.ID, track.FilePath)
		if entry := getCachedCover(cacheKey); entry != nil {
			c.Header("Cache-Control", "public, max-age=86400")
			c.Data(http.StatusOK, entry.contentType, entry.data)
			return
		}

		f, err := os.Open(track.FilePath)
		if err != nil {
			util.ErrorResponse(c, http.StatusNotFound, "file_missing", "Audio file not found")
			return
		}
		defer f.Close()

		m, err := tag.ReadFrom(f)
		if err != nil {
			util.ErrorResponse(c, http.StatusNotFound, "no_tags", "Cannot read tags")
			return
		}

		pic := m.Picture()
		if pic == nil {
			util.ErrorResponse(c, http.StatusNotFound, "no_cover", "No cover art")
			return
		}

		contentType := pic.MIMEType
		if contentType == "" {
			contentType = "image/jpeg"
		}

		// Cache the cover
		putCachedCover(cacheKey, &coverEntry{data: pic.Data, contentType: contentType})

		c.Header("Cache-Control", "public, max-age=86400")
		c.Data(http.StatusOK, contentType, pic.Data)
	}
}
