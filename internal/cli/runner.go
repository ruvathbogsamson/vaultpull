package cli

import (
	"fmt"
	"io"
	"os"

	"github.com/example/vaultpull/internal/config"
	"github.com/example/vaultpull/internal/sync"
	"github.com/example/vaultpull/internal/vault"
)

// Runner executes the vaultpull workflow.
type Runner struct {
	Opts   *Options
	stdout io.Writer
}

// NewRunner creates a Runner with the given options.
func NewRunner(opts *Options) *Runner {
	return &Runner{Opts: opts, stdout: os.Stdout}
}

// Run loads config, connects to Vault, and syncs secrets.
func (r *Runner) Run() error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("config: %w", err)
	}

	if r.Opts.Namespace != "" {
		cfg.Namespace = r.Opts.Namespace
	}
	if r.Opts.OutputFile != "" {
		cfg.OutputFile = r.Opts.OutputFile
	}

	client, err := vault.New(cfg.VaultAddr, cfg.VaultToken)
	if err != nil {
		return fmt.Errorf("vault client: %w", err)
	}

	if r.Opts.DryRun {
		if r.Opts.Verbose {
			fmt.Fprintln(r.stdout, "[dry-run] skipping file write")
		}
		return nil
	}

	syncer := sync.New(client, cfg)
	if err := syncer.Run(); err != nil {
		return fmt.Errorf("sync: %w", err)
	}

	if r.Opts.Verbose {
		fmt.Fprintf(r.stdout, "wrote secrets to %s\n", cfg.OutputFile)
	}
	return nil
}
