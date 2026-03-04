package rest

import (
	"pineapple-music/internal/model"
	"pineapple-music/internal/util"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Search3(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		query := c.Query("query")
		if query == "" {
			util.SubsonicErrorResp(c, 10, "Required parameter is missing: query")
			return
		}

		like := "%" + query + "%"

		// Search artists
		var artistRows []struct{ Artist string }
		db.Model(&model.Track{}).
			Select("DISTINCT artist").
			Where("artist LIKE ?", like).
			Limit(20).
			Find(&artistRows)

		artists := make([]gin.H, 0)
		for i, a := range artistRows {
			artists = append(artists, gin.H{
				"id":   i,
				"name": a.Artist,
			})
		}

		// Search albums
		var albumRows []struct {
			Album  string
			Artist string
		}
		db.Model(&model.Track{}).
			Select("DISTINCT album, MIN(artist) as artist").
			Where("album LIKE ?", like).
			Group("album").
			Limit(20).
			Find(&albumRows)

		albums := make([]gin.H, 0)
		for i, a := range albumRows {
			albums = append(albums, gin.H{
				"id":     i,
				"name":   a.Album,
				"artist": a.Artist,
			})
		}

		// Search songs
		var tracks []model.Track
		db.Where("title LIKE ? OR artist LIKE ?", like, like).
			Limit(50).
			Find(&tracks)

		songs := make([]gin.H, 0)
		for _, t := range tracks {
			songs = append(songs, trackToSubsonic(t))
		}

		util.SubsonicOK(c, gin.H{
			"searchResult3": gin.H{
				"artist": artists,
				"album":  albums,
				"song":   songs,
			},
		})
	}
}
