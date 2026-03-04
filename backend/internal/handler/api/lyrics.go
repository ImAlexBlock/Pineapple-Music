package api

import (
	"net/http"
	"strconv"

	"pineapple-music/internal/model"
	"pineapple-music/internal/util"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetLyrics(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.ParseUint(c.Param("id"), 10, 32)
		if err != nil {
			util.ErrorResponse(c, http.StatusBadRequest, "invalid_id", "Invalid track ID")
			return
		}

		var lyric model.TrackLyric
		if err := db.Where("track_id = ?", id).First(&lyric).Error; err != nil {
			util.ErrorResponse(c, http.StatusNotFound, "no_lyrics", "No lyrics found")
			return
		}

		c.JSON(http.StatusOK, lyric)
	}
}
