package cli

import (
	"flag"
	"fmt"
	"os"

	"github.com/example/vaultpull/internal/vault"
)

// SnapshotFlags holds parsed flags for the snapshot sub-command.
type SnapshotFlags struct {
	Address   string
	Token     string
	Path      string
	OutputFile string
}

// ParseSnapshotFlags parses snapshot sub-command flags from args.
func ParseSnapshotFlags(args []string) (*SnapshotFlags, error) {
	fs := flag.NewFlagSet("snapshot", flag.ContinueOnError)

	addr := fs.String("addr", "http://127.0.0.1:8200", "Vault address")
	token := fs.String("token", "", "Vault token")
	path := fs.String("path", "", "Secret path to snapshot (required)")
	out := fs.String("out", "snapshot.json", "Output file for the snapshot")

	if err := fs.Parse(args); err != nil {
		return nil, err
	}
	if *path == "" {
		return nil, fmt.Errorf("flag -path is required")
	}
	return &SnapshotFlags{
		Address:    *addr,
		Token:      *token,
		Path:       *path,
		OutputFile: *out,
	}, nil
}

// RunSnapshot captures a secret snapshot and writes it to disk.
func RunSnapshot(args []string) error {
	flags, err := ParseSnapshotFlags(args)
	if err != nil {
		return err
	}

	c, err := vault.New(flags.Address, flags.Token)
	if err != nil {
		return fmt.Errorf("vault client: %w", err)
	}

	sc := vault.NewSnapshotClient(c)
	snap, err := sc.Capture(flags.Path)
	if err != nil {
		return fmt.Errorf("capture: %w", err)
	}

	b, err := snap.Marshal()
	if err != nil {
		return fmt.Errorf("marshal snapshot: %w", err)
	}

	if err := os.WriteFile(flags.OutputFile, b, 0600); err != nil {
		return fmt.Errorf("write snapshot: %w", err)
	}

	fmt.Printf("snapshot saved to %s (%d keys)\n", flags.OutputFile, len(snap.Data))
	return nil
}
