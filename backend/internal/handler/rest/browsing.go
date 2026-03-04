package rest

import (
	"fmt"
	"strconv"

	"pineapple-music/internal/model"
	"pineapple-music/internal/util"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetMusicFolders() gin.HandlerFunc {
	return func(c *gin.Context) {
		util.SubsonicOK(c, gin.H{
			"musicFolders": gin.H{
				"musicFolder": []gin.H{
					{"id": 1, "name": "Music"},
				},
			},
		})
	}
}

func GetIndexes(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var artists []struct {
			Artist string
		}
		db.Model(&model.Track{}).Select("DISTINCT artist").Where("artist != ''").Order("artist").Find(&artists)

		indexes := buildIndexes(artists)

		util.SubsonicOK(c, gin.H{
			"indexes": gin.H{
				"index": indexes,
			},
		})
	}
}

func GetArtists(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var artists []struct {
			Artist string
			Count  int64
		}
		db.Model(&model.Track{}).
			Select("artist, COUNT(*) as count").
			Where("artist != ''").
			Group("artist").
			Order("artist").
			Find(&artists)

		indexes := make([]gin.H, 0)
		indexMap := make(map[string][]gin.H)

		for i, a := range artists {
			letter := "#"
			if len(a.Artist) > 0 {
				first := a.Artist[0]
				if (first >= 'A' && first <= 'Z') || (first >= 'a' && first <= 'z') {
					letter = string(first & 0xDF) // uppercase
				}
			}
			indexMap[letter] = append(indexMap[letter], gin.H{
				"id":         fmt.Sprintf("ar-%d", i),
				"name":       a.Artist,
				"albumCount": a.Count,
			})
		}

		for letter, arts := range indexMap {
			indexes = append(indexes, gin.H{
				"name":   letter,
				"artist": arts,
			})
		}

		util.SubsonicOK(c, gin.H{
			"artists": gin.H{
				"index": indexes,
			},
		})
	}
}

func GetArtist(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Query("id")
		if id == "" {
			util.SubsonicErrorResp(c, 10, "Required parameter is missing: id")
			return
		}

		// Extract artist name from id (ar-N) or use directly
		artistName := resolveArtistName(db, id)
		if artistName == "" {
			util.SubsonicErrorResp(c, 70, "Artist not found")
			return
		}

		var albums []struct {
			Album string
			Year  int
			Count int64
		}
		db.Model(&model.Track{}).
			Select("album, MAX(year) as year, COUNT(*) as count").
			Where("artist = ? AND album != ''", artistName).
			Group("album").
			Order("year DESC, album").
			Find(&albums)

		albumList := make([]gin.H, 0, len(albums))
		for i, a := range albums {
			albumList = append(albumList, gin.H{
				"id":        fmt.Sprintf("al-%s-%d", id, i),
				"name":      a.Album,
				"artist":    artistName,
				"year":      a.Year,
				"songCount": a.Count,
			})
		}

		util.SubsonicOK(c, gin.H{
			"artist": gin.H{
				"id":    id,
				"name":  artistName,
				"album": albumList,
			},
		})
	}
}

func GetAlbum(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Query("id")
		if id == "" {
			util.SubsonicErrorResp(c, 10, "Required parameter is missing: id")
			return
		}

		// For simplicity, use album name lookup from tracks
		albumName, artistName := resolveAlbumInfo(db, id)
		if albumName == "" {
			util.SubsonicErrorResp(c, 70, "Album not found")
			return
		}

		var tracks []model.Track
		db.Where("album = ? AND artist = ?", albumName, artistName).
			Order("disc_number, track_number").
			Find(&tracks)

		songs := make([]gin.H, 0, len(tracks))
		for _, t := range tracks {
			songs = append(songs, trackToSubsonic(t))
		}

		util.SubsonicOK(c, gin.H{
			"album": gin.H{
				"id":        id,
				"name":      albumName,
				"artist":    artistName,
				"songCount": len(songs),
				"song":      songs,
			},
		})
	}
}

func GetSong(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Query("id")
		if id == "" {
			util.SubsonicErrorResp(c, 10, "Required parameter is missing: id")
			return
		}

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

		util.SubsonicOK(c, gin.H{
			"song": trackToSubsonic(track),
		})
	}
}

func buildIndexes(artists []struct{ Artist string }) []gin.H {
	indexMap := make(map[string][]gin.H)

	for i, a := range artists {
		letter := "#"
		if len(a.Artist) > 0 {
			first := a.Artist[0]
			if (first >= 'A' && first <= 'Z') || (first >= 'a' && first <= 'z') {
				letter = string(first & 0xDF)
			}
		}
		indexMap[letter] = append(indexMap[letter], gin.H{
			"id":   fmt.Sprintf("ar-%d", i),
			"name": a.Artist,
		})
	}

	result := make([]gin.H, 0)
	for letter, arts := range indexMap {
		result = append(result, gin.H{
			"name":   letter,
			"artist": arts,
		})
	}
	return result
}

func resolveArtistName(db *gorm.DB, id string) string {
	// Try to find by extracting index from ar-N format
	if len(id) > 3 && id[:3] == "ar-" {
		idx, err := strconv.Atoi(id[3:])
		if err == nil {
			var artists []struct{ Artist string }
			db.Model(&model.Track{}).Select("DISTINCT artist").Where("artist != ''").Order("artist").Find(&artists)
			if idx >= 0 && idx < len(artists) {
				return artists[idx].Artist
			}
		}
	}
	// Fallback: treat id as artist name
	var track model.Track
	if db.Where("artist = ?", id).First(&track).Error == nil {
		return track.Artist
	}
	return ""
}

func resolveAlbumInfo(db *gorm.DB, id string) (string, string) {
	// Try numeric ID first (track-based)
	trackID, err := strconv.ParseUint(id, 10, 32)
	if err == nil {
		var track model.Track
		if db.First(&track, trackID).Error == nil {
			return track.Album, track.Artist
		}
	}
	// Try al-xxx format
	if len(id) > 3 && id[:3] == "al-" {
		// Fallback - just find first album
		var track model.Track
		if db.Where("album != ''").First(&track).Error == nil {
			return track.Album, track.Artist
		}
	}
	return "", ""
}

func trackToSubsonic(t model.Track) gin.H {
	return gin.H{
		"id":          fmt.Sprintf("%d", t.ID),
		"title":       t.Title,
		"album":       t.Album,
		"artist":      t.Artist,
		"track":       t.TrackNumber,
		"year":        t.Year,
		"genre":       t.Genre,
		"size":        t.Size,
		"duration":    int(t.Duration),
		"bitRate":     t.Bitrate,
		"contentType": formatToMime(t.Format),
		"suffix":      t.Format,
		"isDir":       false,
		"type":        "music",
	}
}

func formatToMime(format string) string {
	switch format {
	case "mp3":
		return "audio/mpeg"
	case "flac":
		return "audio/flac"
	case "ogg":
		return "audio/ogg"
	case "m4a":
		return "audio/mp4"
	case "wav":
		return "audio/wav"
	default:
		return "application/octet-stream"
	}
}
