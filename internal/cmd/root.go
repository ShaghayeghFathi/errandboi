package cmd

import (
	"log"

	"errandboi/internal/cmd/server"
	"errandboi/internal/config"
	"errandboi/internal/logger"

	"github.com/spf13/cobra"
)

func Execute() {
	cfg := config.New()

	logger := logger.New(cfg.Logger)

	rootCmd := &cobra.Command{
		Use:   "errandboi",
		Short: "Give your errands to the errandboi!",
	}

	server.Register(rootCmd, cfg, logger)

	if err := rootCmd.Execute(); err != nil {
		log.Fatal("error executing command")
	}
}
