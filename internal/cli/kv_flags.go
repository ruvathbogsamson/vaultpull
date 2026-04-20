package cli

import (
	"flag"
	"fmt"
	"os"
)

// KVFlags holds parsed flags for KV-engine-aware operations.
type KVFlags struct {
	Mount      string
	SecretPath string
	KVVersion  int // 0 = auto-detect, 1 or 2 = explicit
	Verbose    bool
}

// ParseKVFlags parses command-line flags for KV subcommands.
// It exits with status 2 on unknown flags.
func ParseKVFlags(args []string) (*KVFlags, error) {
	fs := flag.NewFlagSet("kv", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)

	var f KVFlags
	fs.StringVar(&f.Mount, "mount", "secret", "KV engine mount path")
	fs.StringVar(&f.SecretPath, "path", "", "Secret path within the mount")
	fs.IntVar(&f.KVVersion, "kv-version", 0, "KV engine version (0=auto, 1, or 2)")
	fs.BoolVar(&f.Verbose, "verbose", false, "Enable verbose output")

	if err := fs.Parse(args); err != nil {
		return nil, fmt.Errorf("kv flags: %w", err)
	}
	if f.SecretPath == "" {
		return nil, fmt.Errorf("kv flags: -path is required")
	}
	if f.KVVersion != 0 && f.KVVersion != 1 && f.KVVersion != 2 {
		return nil, fmt.Errorf("kv flags: -kv-version must be 0, 1, or 2")
	}
	return &f, nil
}
