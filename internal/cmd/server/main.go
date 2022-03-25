package server

import (
	"context"
	"errandboi/internal/config"
	"errandboi/internal/db/rdb"
	"fmt"

	"github.com/spf13/cobra"
)

func main(cfg config.Config){
	
	println("ran serve command")
	ctx := context.Background()

	_, err := rdb.New(ctx, cfg.Redis)
	if err != nil {
		fmt.Errorf("redis initiation failed")
	}

	// mongodb, err := mongodb.New(cfg.Mongo)
	// if err != nil {
	// 	fmt.Errorf("mongo initiation failed")
	// }



}

func Register(root *cobra.Command, cfg config.Config) {
	root.AddCommand(
		&cobra.Command{
			Use:   "serve",
			Short: "Run server",
			Run: func(cmd *cobra.Command, args []string) {
				main(cfg)
			},
		},
	)
}