package rest

import (
	"pineapple-music/internal/config"
	"pineapple-music/internal/service"
	"pineapple-music/internal/util"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// subsonicAuth validates Subsonic authentication parameters.
func subsonicAuth(authSvc *service.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		u := c.Query("u")
		p := c.Query("p")
		t := c.Query("t")
		s := c.Query("s")

		if u == "" {
			util.SubsonicErrorResp(c, 10, "Required parameter is missing: u")
			c.Abort()
			return
		}

		// Token-based auth: token = md5(password + salt)
		if t != "" && s != "" {
			role, ok := authSvc.ValidateSubsonicToken(t, s)
			if !ok {
				util.SubsonicErrorResp(c, 40, "Wrong username or password")
				c.Abort()
				return
			}
			c.Set("subsonic_role", role)
			c.Next()
			return
		}

		// Password-based auth (plain or hex-encoded)
		key := p
		if key == "" {
			util.SubsonicErrorResp(c, 10, "Required parameter is missing: p or t+s")
			c.Abort()
			return
		}

		// Strip "enc:" prefix for hex-encoded passwords
		if len(key) > 4 && key[:4] == "enc:" {
			key = hexDecode(key[4:])
		}

		// Try admin key first, then guest
		if authSvc.ValidateAdminKey(key) {
			c.Set("subsonic_role", "admin")
		} else if authSvc.ValidateGuestKey(key) {
			c.Set("subsonic_role", "guest")
		} else {
			util.SubsonicErrorResp(c, 40, "Wrong username or password")
			c.Abort()
			return
		}

		c.Next()
	}
}

func hexDecode(s string) string {
	result := make([]byte, len(s)/2)
	for i := 0; i < len(s)-1; i += 2 {
		result[i/2] = hexByte(s[i])<<4 | hexByte(s[i+1])
	}
	return string(result)
}

func hexByte(b byte) byte {
	switch {
	case b >= '0' && b <= '9':
		return b - '0'
	case b >= 'a' && b <= 'f':
		return b - 'a' + 10
	case b >= 'A' && b <= 'F':
		return b - 'A' + 10
	}
	return 0
}

func RegisterSubsonicRoutes(r *gin.Engine, db *gorm.DB, authSvc *service.AuthService, settingSvc *service.SettingsService, cfg *config.Config) {
	rest := r.Group("/rest")
	rest.Use(func(c *gin.Context) {
		// Check if Subsonic protocol is enabled
		if !settingSvc.GetBool("subsonic_enabled", true) {
			util.SubsonicErrorResp(c, 0, "Subsonic protocol is disabled")
			c.Abort()
			return
		}
		c.Next()
	})
	rest.Use(subsonicAuth(authSvc))

	// Register each endpoint with both /xxx and /xxx.view paths
	endpoints := map[string]gin.HandlerFunc{
		"ping":            Ping(),
		"getLicense":      GetLicense(),
		"getMusicFolders": GetMusicFolders(),
		"getIndexes":      GetIndexes(db),
		"getArtists":      GetArtists(db),
		"getArtist":       GetArtist(db),
		"getAlbum":        GetAlbum(db),
		"getSong":         GetSong(db),
		"search3":         Search3(db),
		"getPlaylists":    GetPlaylists(db),
		"getPlaylist":     GetPlaylistSub(db),
		"stream":          Stream(db),
		"getCoverArt":     GetCoverArt(db),
	}

	for name, handler := range endpoints {
		rest.GET("/"+name, handler)
		rest.GET("/"+name+".view", handler)
	}
}
