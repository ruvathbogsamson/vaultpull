package main

import (
	"fmt"
	"os"

	"github.com/example/vaultpull/internal/config"
	"github.com/example/vaultpull/internal/sync"
	"github.com/example/vaultpull/internal/vault"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("loading config: %w", err)
	}

	client, err := vault.New(cfg.VaultAddr, cfg.VaultToken)
	if err != nil {
		return fmt.Errorf("creating vault client: %w", err)
	}

	syncer := sync.New(client, cfg)
	if err := syncer.Run(); err != nil {
		return fmt.Errorf("syncing secrets: %w", err)
	}

	fmt.Println("secrets synced successfully")
	return nil
}
