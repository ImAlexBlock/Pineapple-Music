package rest

import (
	"fmt"
	"strconv"

	"pineapple-music/internal/model"
	"pineapple-music/internal/util"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetPlaylists(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var playlists []model.Playlist
		db.Find(&playlists)

		result := make([]gin.H, 0, len(playlists))
		for _, pl := range playlists {
			var count int64
			db.Model(&model.PlaylistTrack{}).Where("playlist_id = ?", pl.ID).Count(&count)
			result = append(result, gin.H{
				"id":        fmt.Sprintf("%d", pl.ID),
				"name":      pl.Name,
				"songCount": count,
				"public":    true,
				"created":   pl.CreatedAt.Format("2006-01-02T15:04:05"),
			})
		}

		util.SubsonicOK(c, gin.H{
			"playlists": gin.H{
				"playlist": result,
			},
		})
	}
}

func GetPlaylistSub(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Query("id")
		plID, err := strconv.ParseUint(id, 10, 32)
		if err != nil {
			util.SubsonicErrorResp(c, 70, "Playlist not found")
			return
		}

		var pl model.Playlist
		if err := db.First(&pl, plID).Error; err != nil {
			util.SubsonicErrorResp(c, 70, "Playlist not found")
			return
		}

		var pts []model.PlaylistTrack
		db.Where("playlist_id = ?", pl.ID).Order("position").Preload("Track").Find(&pts)

		songs := make([]gin.H, 0, len(pts))
		for _, pt := range pts {
			songs = append(songs, trackToSubsonic(pt.Track))
		}

		util.SubsonicOK(c, gin.H{
			"playlist": gin.H{
				"id":        fmt.Sprintf("%d", pl.ID),
				"name":      pl.Name,
				"songCount": len(songs),
				"entry":     songs,
			},
		})
	}
}
