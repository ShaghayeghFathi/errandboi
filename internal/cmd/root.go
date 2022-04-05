package cmd

import (
	"log"

	"github.com/ShaghayeghFathi/errandboi/internal/cmd/server"
	"github.com/ShaghayeghFathi/errandboi/internal/config"
	"github.com/ShaghayeghFathi/errandboi/internal/logger"

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
