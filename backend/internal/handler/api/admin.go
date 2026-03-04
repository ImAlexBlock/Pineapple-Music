package api

import (
	"net/http"
	"strings"

	"pineapple-music/internal/model"
	"pineapple-music/internal/service"
	"pineapple-music/internal/util"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// allowedSettingKeys defines which settings can be modified via the API.
var allowedSettingKeys = map[string]bool{
	"access_mode":      true, // "public" or "private"
	"subsonic_enabled": true, // "true" or "false"
	"site_name":        true,
}

func AdminDashboard(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var trackCount int64
		db.Model(&model.Track{}).Count(&trackCount)

		var totalSize int64
		db.Model(&model.Track{}).Select("COALESCE(SUM(size), 0)").Scan(&totalSize)

		var playCount int64
		db.Model(&model.PlayEvent{}).Count(&playCount)

		var playlistCount int64
		db.Model(&model.Playlist{}).Count(&playlistCount)

		c.JSON(http.StatusOK, gin.H{
			"tracks":     trackCount,
			"total_size": totalSize,
			"plays":      playCount,
			"playlists":  playlistCount,
		})
	}
}

func GetSettings(settingSvc *service.SettingsService) gin.HandlerFunc {
	return func(c *gin.Context) {
		settings, err := settingSvc.GetAll()
		if err != nil {
			util.ErrorResponse(c, http.StatusInternalServerError, "error", err.Error())
			return
		}

		result := make(map[string]string)
		for _, s := range settings {
			// Only expose allowed settings
			if allowedSettingKeys[s.Key] {
				result[s.Key] = s.Value
			}
		}

		c.JSON(http.StatusOK, result)
	}
}

func UpdateSettings(settingSvc *service.SettingsService, auditSvc *service.AuditService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req map[string]string
		if err := c.ShouldBindJSON(&req); err != nil {
			util.ErrorResponse(c, http.StatusBadRequest, "invalid_request", "Invalid JSON")
			return
		}

		var rejected []string
		var errors []string
		for key, value := range req {
			if !allowedSettingKeys[key] {
				rejected = append(rejected, key)
				continue
			}
			if err := settingSvc.Set(key, value); err != nil {
				errors = append(errors, key+": "+err.Error())
			}
		}

		if len(rejected) > 0 {
			util.ErrorResponse(c, http.StatusBadRequest, "invalid_keys",
				"Unknown settings: "+strings.Join(rejected, ", "))
			return
		}
		if len(errors) > 0 {
			util.ErrorResponse(c, http.StatusInternalServerError, "save_error",
				"Failed to save: "+strings.Join(errors, "; "))
			return
		}

		auditSvc.Log("settings_updated", "admin", c.ClientIP(), "")
		c.JSON(http.StatusOK, gin.H{"message": "Settings updated"})
	}
}

func RotateAdminKey(authSvc *service.AuthService, auditSvc *service.AuditService) gin.HandlerFunc {
	return func(c *gin.Context) {
		key, err := authSvc.RotateAdminKey()
		if err != nil {
			util.ErrorResponse(c, http.StatusInternalServerError, "error", err.Error())
			return
		}
		auditSvc.Log("admin_key_rotated", "admin", c.ClientIP(), "")
		c.JSON(http.StatusOK, gin.H{
			"admin_key": key,
			"message":   "Save this key! It will not be shown again.",
		})
	}
}

func RotateGuestKey(authSvc *service.AuthService, auditSvc *service.AuditService) gin.HandlerFunc {
	return func(c *gin.Context) {
		key, err := authSvc.RotateGuestKey()
		if err != nil {
			util.ErrorResponse(c, http.StatusInternalServerError, "error", err.Error())
			return
		}
		auditSvc.Log("guest_key_rotated", "admin", c.ClientIP(), "")
		// Show once — will not be retrievable again
		c.JSON(http.StatusOK, gin.H{
			"guest_key": key,
			"message":   "Save this key! It will not be shown again.",
		})
	}
}
