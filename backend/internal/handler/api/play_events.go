package api

import (
	"net/http"
	"time"

	"pineapple-music/internal/middleware"
	"pineapple-music/internal/model"
	"pineapple-music/internal/util"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RecordPlayEvent(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			TrackID uint `json:"track_id" binding:"required"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			util.ErrorResponse(c, http.StatusBadRequest, "invalid_request", "track_id is required")
			return
		}

		// Verify track exists
		var track model.Track
		if err := db.First(&track, req.TrackID).Error; err != nil {
			util.ErrorResponse(c, http.StatusNotFound, "not_found", "Track not found")
			return
		}

		sessionID, _ := c.Get(middleware.ContextSession)
		sid, _ := sessionID.(string)

		event := model.PlayEvent{
			TrackID:   req.TrackID,
			SessionID: sid,
			IP:        c.ClientIP(),
			PlayedAt:  time.Now(),
		}
		db.Create(&event)

		c.JSON(http.StatusOK, gin.H{"message": "recorded"})
	}
}
