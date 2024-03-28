package main

import (
	"context"
	"fmt"
	"os"

	"github.com/topology-gg/gram/config"
	"github.com/topology-gg/gram/execution"
	"github.com/topology-gg/gram/internal/app"
	"github.com/topology-gg/gram/network"
	"github.com/topology-gg/gram/storage"
)

func main() {
	app := app.NewApp()

	app.Name = "gram"
	app.Description = "The official Go implementation of the RAM network"
	app.Action = gram

	if err := app.Run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func gram() error {
	ctx := context.Background()

	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load configuration: %v\n", err)
		os.Exit(1)
	}

	storage := storage.NewStorage(ctx, &cfg.Storage)
	execution := execution.NewExecution(ctx, storage, &cfg.Execution)
	network := network.NewNetwork(ctx, execution, storage, &cfg.Network, &cfg.Grpc)

	network.Start()

	return nil
}
