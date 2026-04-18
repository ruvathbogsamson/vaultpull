// Package sync orchestrates fetching secrets from Vault and writing them
// to a local .env file via the env package.
package sync

import (
	"context"
	"fmt"

	"github.com/your-org/vaultpull/internal/audit"
	"github.com/your-org/vaultpull/internal/env"
	"github.com/your-org/vaultpull/internal/vault"
)

// Syncer coordinates a single vault-to-env sync operation.
type Syncer struct {
	client  *vault.Client
	writer  *env.Writer
	auditor *audit.Logger
}

// New returns a Syncer ready to run.
func New(client *vault.Client, writer *env.Writer, auditor *audit.Logger) *Syncer {
	return &Syncer{client: client, writer: writer, auditor: auditor}
}

// Run fetches secrets at secretPath, optionally filters by namespace,
// writes them to the configured output file, and emits an audit entry.
func (s *Syncer) Run(ctx context.Context, secretPath, namespace string, dryRun bool) error {
	secrets, err := s.client.GetSecrets(ctx, secretPath)
	if err != nil {
		if s.auditor != nil {
			_ = s.auditor.Log(audit.Entry{
				Operation:  "sync",
				SecretPath: secretPath,
				Namespace:  namespace,
				DryRun:     dryRun,
				Error:      err.Error(),
			})
		}
		return fmt.Errorf("fetching secrets: %w", err)
	}

	filtered := env.Filter(secrets, namespace)

	if !dryRun {
		if err := s.writer.Write(filtered); err != nil {
			return fmt.Errorf("writing env file: %w", err)
		}
	}

	if s.auditor != nil {
		_ = s.auditor.Log(audit.Entry{
			Operation:   "sync",
			SecretPath:  secretPath,
			Namespace:   namespace,
			KeysWritten: len(filtered),
			DryRun:      dryRun,
		})
	}

	return nil
}
