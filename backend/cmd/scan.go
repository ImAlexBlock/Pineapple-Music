package cmd

import (
	"fmt"
	"log"
	"os"
	"time"

	"pineapple-music/internal/config"
	"pineapple-music/internal/model"
	"pineapple-music/internal/scanner"

	"github.com/spf13/cobra"
)

var scanCmd = &cobra.Command{
	Use:   "scan",
	Short: "Scan music directory for audio files",
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.Load()
		if err != nil {
			log.Fatalf("Failed to load config: %v", err)
		}

		os.MkdirAll(cfg.MusicDir(), 0755)

		db, err := model.InitDB(cfg.DBPath())
		if err != nil {
			log.Fatalf("Failed to init database: %v", err)
		}

		job := &model.ScanJob{Status: "pending"}
		db.Create(job)

		fmt.Printf("Scanning %s ...\n", cfg.MusicDir())
		scanner.Scan(db, cfg.MusicDir(), job)

		// Reload to get final state
		db.First(job, job.ID)
		fmt.Printf("Done: %d total, %d added, %d updated, %d errors\n",
			job.Total, job.Added, job.Updated, job.Errors)
		if job.FinishedAt != nil {
			fmt.Printf("Duration: %s\n", job.FinishedAt.Sub(job.StartedAt).Round(time.Millisecond))
		}
	},
}

func init() {
	rootCmd.AddCommand(scanCmd)
}
