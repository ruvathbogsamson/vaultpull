package sync

import (
	"fmt"

	"github.com/example/vaultpull/internal/config"
	"github.com/example/vaultpull/internal/env"
	"github.com/example/vaultpull/internal/vault"
)

// Syncer orchestrates pulling secrets from Vault and writing them to a .env file.
type Syncer struct {
	cfg    *config.Config
	client *vault.Client
	writer *env.Writer
}

// New creates a new Syncer from the provided config.
func New(cfg *config.Config) (*Syncer, error) {
	client, err := vault.New(cfg.VaultAddr, cfg.VaultToken)
	if err != nil {
		return nil, fmt.Errorf("syncer: failed to create vault client: %w", err)
	}

	writer, err := env.NewWriter(cfg.OutputFile)
	if err != nil {
		return nil, fmt.Errorf("syncer: failed to create env writer: %w", err)
	}

	return &Syncer{
		cfg:    cfg,
		client: client,
		writer: writer,
	}, nil
}

// Run fetches secrets from Vault, applies namespace filtering, and writes the result.
func (s *Syncer) Run() (int, error) {
	secrets, err := s.client.GetSecrets(s.cfg.SecretPath)
	if err != nil {
		return 0, fmt.Errorf("syncer: failed to fetch secrets: %w", err)
	}

	filtered := env.Filter(secrets, s.cfg.Namespace)

	if err := s.writer.Write(filtered); err != nil {
		return 0, fmt.Errorf("syncer: failed to write env file: %w", err)
	}

	return len(filtered), nil
}
