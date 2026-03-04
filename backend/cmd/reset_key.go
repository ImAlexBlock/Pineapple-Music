package cmd

import (
	"fmt"
	"log"

	"pineapple-music/internal/config"
	"pineapple-music/internal/model"
	"pineapple-music/internal/service"

	"github.com/spf13/cobra"
)

var resetKeyCmd = &cobra.Command{
	Use:   "reset-key [admin|guest]",
	Short: "Reset admin or guest key",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.Load()
		if err != nil {
			log.Fatalf("Failed to load config: %v", err)
		}

		db, err := model.InitDB(cfg.DBPath())
		if err != nil {
			log.Fatalf("Failed to init database: %v", err)
		}

		authSvc := &service.AuthService{DB: db, Cfg: cfg}

		switch args[0] {
		case "admin":
			key, err := authSvc.RotateAdminKey()
			if err != nil {
				log.Fatalf("Failed to reset admin key: %v", err)
			}
			fmt.Printf("New admin key: %s\n", key)
			fmt.Println("Save this key! It will not be shown again.")
		case "guest":
			key, err := authSvc.RotateGuestKey()
			if err != nil {
				log.Fatalf("Failed to reset guest key: %v", err)
			}
			fmt.Printf("New guest key: %s\n", key)
		default:
			log.Fatalf("Unknown key type: %s (use 'admin' or 'guest')", args[0])
		}
	},
}

func init() {
	rootCmd.AddCommand(resetKeyCmd)
}
