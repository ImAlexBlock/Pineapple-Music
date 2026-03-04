package api

import (
	"net/http"
	"os"
	"strconv"

	"pineapple-music/internal/model"
	"pineapple-music/internal/util"

	"github.com/dhowden/tag"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

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

		c.Header("Cache-Control", "public, max-age=86400")
		c.Data(http.StatusOK, contentType, pic.Data)
	}
}
