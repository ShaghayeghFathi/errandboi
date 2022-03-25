package cmd

import (
	"errandboi/internal/cmd/server"
	"errandboi/internal/config"
	"log"

	"github.com/spf13/cobra"
)

func Execute() {

	cfg := config.New()

	// initiate logger 

	rootCmd := &cobra.Command{
		Use:   "errandboi",
		Short: "Give your errands to the errandboi!",
	}

	server.Register(rootCmd, cfg)

	if err := rootCmd.Execute(); err != nil {
		log.Fatal("error executing command")
	}

}



