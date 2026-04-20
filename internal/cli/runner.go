package cli

import (
	"context"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/your-org/vaultpull/internal/config"
	"github.com/your-org/vaultpull/internal/sync"
	"github.com/your-org/vaultpull/internal/vault"

	vaultapi "github.com/hashicorp/vault/api"
)

// Runner orchestrates config loading, optional token renewal, and secret sync.
type Runner struct {
	flags  *Flags
	stdout io.Writer
	stderr io.Writer
}

// NewRunner returns a Runner configured with the given flags.
func NewRunner(f *Flags, stdout, stderr io.Writer) *Runner {
	if stdout == nil {
		stdout = os.Stdout
	}
	if stderr == nil {
		stderr = os.Stderr
	}
	return &Runner{flags: f, stdout: stdout, stderr: stderr}
}

// Run loads config, optionally starts token renewal, and executes the sync.
func (r *Runner) Run(ctx context.Context) error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("config: %w", err)
	}

	if r.flags.Namespace != "" {
		cfg.Namespace = r.flags.Namespace
	}

	if r.flags.Verbose {
		fmt.Fprintf(r.stdout, "[vaultpull] address=%s path=%s namespace=%q\n",
			cfg.VaultAddr, cfg.SecretPath, cfg.Namespace)
	}

	if r.flags.DryRun {
		fmt.Fprintln(r.stdout, "[vaultpull] dry-run: no files written")
		return nil
	}

	vCfg := vaultapi.DefaultConfig()
	vCfg.Address = cfg.VaultAddr
	vClient, err := vaultapi.NewClient(vCfg)
	if err != nil {
		return fmt.Errorf("vault client: %w", err)
	}
	vClient.SetToken(cfg.VaultToken)

	renewer := vault.NewTokenRenewer(vClient, 5*time.Minute)
	renewer.Start(ctx)
	defer renewer.Stop()

	s := sync.New(vClient, cfg)
	if err := s.Run(ctx); err != nil {
		return fmt.Errorf("sync: %w", err)
	}

	if r.flags.Verbose {
		fmt.Fprintln(r.stdout, "[vaultpull] sync complete")
	}
	return nil
}
