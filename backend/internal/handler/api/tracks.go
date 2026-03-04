package api

import (
	"net/http"
	"strconv"

	"pineapple-music/internal/model"
	"pineapple-music/internal/util"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func ListTracks(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var params util.PaginationParams
		c.ShouldBindQuery(&params)
		params.Normalize()

		var tracks []model.Track
		var total int64

		query := db.Model(&model.Track{})

		// Search
		if q := c.Query("q"); q != "" {
			like := "%" + q + "%"
			query = query.Where("title LIKE ? OR artist LIKE ? OR album LIKE ?", like, like, like)
		}

		// Filter by artist
		if artist := c.Query("artist"); artist != "" {
			query = query.Where("artist = ?", artist)
		}

		// Filter by album
		if album := c.Query("album"); album != "" {
			query = query.Where("album = ?", album)
		}

		query.Count(&total)
		query.Offset(params.Offset).Limit(params.Limit).Order("artist, album, track_number, title").Find(&tracks)

		c.JSON(http.StatusOK, gin.H{
			"total":  total,
			"offset": params.Offset,
			"limit":  params.Limit,
			"items":  tracks,
		})
	}
}

func GetTrack(db *gorm.DB) gin.HandlerFunc {
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

		c.JSON(http.StatusOK, track)
	}
}

func GetArtists(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var artists []struct {
			Artist string `json:"artist"`
			Count  int64  `json:"count"`
		}
		db.Model(&model.Track{}).
			Select("artist, COUNT(*) as count").
			Where("artist != ''").
			Group("artist").
			Order("artist").
			Find(&artists)

		c.JSON(http.StatusOK, artists)
	}
}

func GetAlbums(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var albums []struct {
			Album       string `json:"album"`
			AlbumArtist string `json:"album_artist"`
			Artist      string `json:"artist"`
			Count       int64  `json:"count"`
			Year        int    `json:"year"`
		}

		query := db.Model(&model.Track{}).
			Select("album, album_artist, MIN(artist) as artist, COUNT(*) as count, MAX(year) as year").
			Where("album != ''").
			Group("album, album_artist").
			Order("album")

		if artist := c.Query("artist"); artist != "" {
			query = query.Where("artist = ? OR album_artist = ?", artist, artist)
		}

		query.Find(&albums)

		c.JSON(http.StatusOK, albums)
	}
}
