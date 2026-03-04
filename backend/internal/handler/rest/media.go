package rest

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

func Stream(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Query("id")
		trackID, err := strconv.ParseUint(id, 10, 32)
		if err != nil {
			util.SubsonicErrorResp(c, 70, "Song not found")
			return
		}

		var track model.Track
		if err := db.First(&track, trackID).Error; err != nil {
			util.SubsonicErrorResp(c, 70, "Song not found")
			return
		}

		if _, err := os.Stat(track.FilePath); os.IsNotExist(err) {
			util.SubsonicErrorResp(c, 70, "Song file not found")
			return
		}

		http.ServeFile(c.Writer, c.Request, track.FilePath)
	}
}

func GetCoverArt(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Query("id")
		trackID, err := strconv.ParseUint(id, 10, 32)
		if err != nil {
			util.SubsonicErrorResp(c, 70, "Cover not found")
			return
		}

		var track model.Track
		if err := db.First(&track, trackID).Error; err != nil {
			util.SubsonicErrorResp(c, 70, "Cover not found")
			return
		}

		f, err := os.Open(track.FilePath)
		if err != nil {
			util.SubsonicErrorResp(c, 70, "File not found")
			return
		}
		defer f.Close()

		m, err := tag.ReadFrom(f)
		if err != nil {
			util.SubsonicErrorResp(c, 70, "Cannot read tags")
			return
		}

		pic := m.Picture()
		if pic == nil {
			util.SubsonicErrorResp(c, 70, "No cover art")
			return
		}

		contentType := pic.MIMEType
		if contentType == "" {
			contentType = "image/jpeg"
		}

		c.Data(http.StatusOK, contentType, pic.Data)
	}
}
