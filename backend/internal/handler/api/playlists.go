package api

import (
	"net/http"
	"strconv"

	"pineapple-music/internal/service"
	"pineapple-music/internal/util"

	"github.com/gin-gonic/gin"
)

func ListPlaylists(plSvc *service.PlaylistService) gin.HandlerFunc {
	return func(c *gin.Context) {
		playlists, err := plSvc.List()
		if err != nil {
			util.ErrorResponse(c, http.StatusInternalServerError, "error", err.Error())
			return
		}
		c.JSON(http.StatusOK, playlists)
	}
}

func GetPlaylist(plSvc *service.PlaylistService) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.ParseUint(c.Param("id"), 10, 32)
		if err != nil {
			util.ErrorResponse(c, http.StatusBadRequest, "invalid_id", "Invalid playlist ID")
			return
		}

		pl, err := plSvc.Get(uint(id))
		if err != nil {
			util.ErrorResponse(c, http.StatusNotFound, "not_found", "Playlist not found")
			return
		}

		c.JSON(http.StatusOK, pl)
	}
}

func CreatePlaylist(plSvc *service.PlaylistService, auditSvc *service.AuditService) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req struct {
			Name string `json:"name" binding:"required"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			util.ErrorResponse(c, http.StatusBadRequest, "invalid_request", "Name is required")
			return
		}

		pl, err := plSvc.Create(req.Name)
		if err != nil {
			util.ErrorResponse(c, http.StatusInternalServerError, "error", err.Error())
			return
		}

		auditSvc.Log("playlist_created", "admin", c.ClientIP(), req.Name)
		c.JSON(http.StatusCreated, pl)
	}
}

func DeletePlaylist(plSvc *service.PlaylistService, auditSvc *service.AuditService) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.ParseUint(c.Param("id"), 10, 32)
		if err != nil {
			util.ErrorResponse(c, http.StatusBadRequest, "invalid_id", "Invalid playlist ID")
			return
		}

		if err := plSvc.Delete(uint(id)); err != nil {
			util.ErrorResponse(c, http.StatusInternalServerError, "error", err.Error())
			return
		}

		auditSvc.Log("playlist_deleted", "admin", c.ClientIP(), c.Param("id"))
		c.JSON(http.StatusOK, gin.H{"message": "Playlist deleted"})
	}
}

func AddTrackToPlaylist(plSvc *service.PlaylistService) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
		var req struct {
			TrackID uint `json:"track_id" binding:"required"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			util.ErrorResponse(c, http.StatusBadRequest, "invalid_request", "track_id is required")
			return
		}

		if err := plSvc.AddTrack(uint(id), req.TrackID); err != nil {
			util.ErrorResponse(c, http.StatusInternalServerError, "error", err.Error())
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Track added"})
	}
}

func RemoveTrackFromPlaylist(plSvc *service.PlaylistService) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
		trackID, _ := strconv.ParseUint(c.Param("trackId"), 10, 32)

		if err := plSvc.RemoveTrack(uint(id), uint(trackID)); err != nil {
			util.ErrorResponse(c, http.StatusInternalServerError, "error", err.Error())
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Track removed"})
	}
}

func ReorderPlaylistTracks(plSvc *service.PlaylistService) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
		var req struct {
			TrackIDs []uint `json:"track_ids" binding:"required"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			util.ErrorResponse(c, http.StatusBadRequest, "invalid_request", "track_ids is required")
			return
		}

		if err := plSvc.ReorderTracks(uint(id), req.TrackIDs); err != nil {
			util.ErrorResponse(c, http.StatusInternalServerError, "error", err.Error())
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Tracks reordered"})
	}
}
