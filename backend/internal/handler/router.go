package handler

import (
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"pineapple-music/internal/config"
	"pineapple-music/internal/handler/api"
	"pineapple-music/internal/handler/rest"
	"pineapple-music/internal/middleware"
	"pineapple-music/internal/service"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRouter(db *gorm.DB, cfg *config.Config) *gin.Engine {
	r := gin.Default()

	// Configure trusted proxies
	if proxies := cfg.TrustedProxiesList(); len(proxies) > 0 {
		r.SetTrustedProxies(proxies)
	} else {
		r.SetTrustedProxies(nil)
	}

	// Services
	authSvc := &service.AuthService{DB: db, Cfg: cfg}
	auditSvc := &service.AuditService{DB: db}
	scanSvc := &service.ScanService{DB: db, MusicDir: cfg.MusicDir()}
	plSvc := &service.PlaylistService{DB: db}
	settingSvc := &service.SettingsService{DB: db}

	// Rate limiter
	rl := middleware.NewRateLimiter(cfg.RateLimitRPS, cfg.RateLimitBurst)
	r.Use(rl.Middleware())

	// Health check
	r.GET("/api/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"version": "0.1.0",
		})
	})

	// Public API routes
	v1 := r.Group("/api/v1")
	{
		v1.Use(middleware.CSRF())

		// Bootstrap — localhost only, no key returned over HTTP
		v1.GET("/setup/status", api.BootstrapStatus(authSvc))
		setupGroup := v1.Group("/setup")
		setupGroup.Use(middleware.LocalOnly())
		setupGroup.Use(middleware.Turnstile(cfg.TurnstileSecret))
		{
			setupGroup.POST("/bootstrap", api.Bootstrap(authSvc, auditSvc))
		}

		// Auth
		authGroup := v1.Group("/auth")
		authGroup.Use(middleware.Turnstile(cfg.TurnstileSecret))
		{
			authGroup.POST("/login", api.Login(authSvc, cfg, auditSvc))
		}
		v1.POST("/auth/logout", api.Logout(authSvc))

		// Public browsing — protected by AccessMode (private mode requires login)
		public := v1.Group("")
		public.Use(middleware.OptionalAuth(authSvc))
		public.Use(middleware.AccessMode(authSvc, settingSvc))
		{
			public.GET("/tracks", api.ListTracks(db))
			public.GET("/tracks/:id", api.GetTrack(db))
			public.GET("/tracks/:id/stream", api.StreamTrack(db))
			public.GET("/tracks/:id/cover", api.TrackCover(db))
			public.GET("/tracks/:id/lyrics", api.GetLyrics(db))
			public.GET("/artists", api.GetArtists(db))
			public.GET("/albums", api.GetAlbums(db))
			public.GET("/playlists", api.ListPlaylists(plSvc))
			public.GET("/playlists/:id", api.GetPlaylist(plSvc))
		}

		// Authenticated routes
		authed := v1.Group("")
		authed.Use(middleware.Auth(authSvc))
		{
			authed.GET("/auth/me", api.Me(authSvc))
			authed.POST("/play-events", api.RecordPlayEvent(db))
		}

		// Admin routes
		admin := authed.Group("")
		admin.Use(middleware.RequireAdmin())
		{
			admin.POST("/upload", api.Upload(cfg, scanSvc, auditSvc, db))
			admin.POST("/scan", api.StartScan(scanSvc, auditSvc))
			admin.GET("/scan/status", api.ScanStatus(scanSvc))

			admin.POST("/playlists", api.CreatePlaylist(plSvc, auditSvc))
			admin.DELETE("/playlists/:id", api.DeletePlaylist(plSvc, auditSvc))
			admin.POST("/playlists/:id/tracks", api.AddTrackToPlaylist(plSvc))
			admin.DELETE("/playlists/:id/tracks/:trackId", api.RemoveTrackFromPlaylist(plSvc))
			admin.PUT("/playlists/:id/reorder", api.ReorderPlaylistTracks(plSvc))

			admin.GET("/admin/dashboard", api.AdminDashboard(db))
			admin.GET("/admin/settings", api.GetSettings(settingSvc))
			admin.PUT("/admin/settings", api.UpdateSettings(settingSvc, auditSvc))
			admin.POST("/admin/rotate-admin-key", api.RotateAdminKey(authSvc, auditSvc))
			admin.POST("/admin/rotate-guest-key", api.RotateGuestKey(authSvc, auditSvc))
			admin.GET("/admin/audit-logs", api.ListAuditLogs(db))
		}
	}

	// Subsonic REST API
	rest.RegisterSubsonicRoutes(r, db, authSvc, settingSvc, cfg)

	// Serve frontend SPA from ../frontend/dist if it exists
	distPath := filepath.Join("..", "frontend", "dist")
	if info, err := os.Stat(distPath); err == nil && info.IsDir() {
		distFS := os.DirFS(distPath)
		fileServer := http.FileServer(http.FS(distFS))
		r.NoRoute(func(c *gin.Context) {
			path := c.Request.URL.Path
			if strings.HasPrefix(path, "/api") || strings.HasPrefix(path, "/rest") {
				return
			}
			cleanPath := strings.TrimPrefix(path, "/")
			if _, err := fs.Stat(distFS, cleanPath); err == nil {
				fileServer.ServeHTTP(c.Writer, c.Request)
				return
			}
			c.Request.URL.Path = "/"
			fileServer.ServeHTTP(c.Writer, c.Request)
		})
	}

	return r
}
