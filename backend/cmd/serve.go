package cmd

import (
	"fmt"
	"log"
	"os"
	"time"

	"pineapple-music/internal/config"
	"pineapple-music/internal/handler"
	"pineapple-music/internal/model"
	"pineapple-music/internal/service"

	"github.com/spf13/cobra"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the music server",
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.Load()
		if err != nil {
			log.Fatalf("Failed to load config: %v", err)
		}

		// Ensure data and music directories exist
		os.MkdirAll(cfg.DataDir, 0755)
		os.MkdirAll(cfg.MusicDir(), 0755)

		db, err := model.InitDB(cfg.DBPath())
		if err != nil {
			log.Fatalf("Failed to init database: %v", err)
		}

		// Start session cleanup ticker
		authSvc := &service.AuthService{DB: db, Cfg: cfg}
		go func() {
			ticker := time.NewTicker(1 * time.Hour)
			defer ticker.Stop()
			for range ticker.C {
				authSvc.CleanExpiredSessions()
			}
		}()

		r := handler.SetupRouter(db, cfg)

		addr := fmt.Sprintf(":%d", cfg.Port)
		log.Printf("Pineapple Music listening on %s", addr)
		if err := r.Run(addr); err != nil {
			log.Fatalf("Server failed: %v", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
}
