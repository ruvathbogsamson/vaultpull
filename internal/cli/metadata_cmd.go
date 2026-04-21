package cli

import (
	"flag"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/yourusername/vaultpull/internal/vault"
)

// MetadataFlags holds parsed flags for the metadata sub-command.
type MetadataFlags struct {
	Address string
	Token   string
	Mount   string
	Path    string
}

// ParseMetadataFlags parses flags for the metadata sub-command from args.
func ParseMetadataFlags(args []string) (*MetadataFlags, error) {
	fs := flag.NewFlagSet("metadata", flag.ContinueOnError)

	address := fs.String("address", "http://127.0.0.1:8200", "Vault server address")
	token := fs.String("token", "", "Vault token (or set VAULT_TOKEN)")
	mount := fs.String("mount", "secret", "KV v2 mount path")
	path := fs.String("path", "", "Secret path to inspect (required)")

	if err := fs.Parse(args); err != nil {
		return nil, err
	}

	if *path == "" {
		return nil, fmt.Errorf("metadata: -path is required")
	}

	if *token == "" {
		*token = os.Getenv("VAULT_TOKEN")
	}

	return &MetadataFlags{
		Address: *address,
		Token:   *token,
		Mount:   *mount,
		Path:    *path,
	}, nil
}

// RunMetadata fetches and prints secret metadata to stdout.
func RunMetadata(f *MetadataFlags) error {
	client := vault.NewMetadataClient(f.Address, f.Token, f.Mount)

	meta, err := client.FetchMetadata(f.Path)
	if err != nil {
		return fmt.Errorf("metadata: %w", err)
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintf(w, "Path:\t%s\n", f.Path)
	fmt.Fprintf(w, "Current Version:\t%d\n", meta.CurrentVersion)
	fmt.Fprintf(w, "Oldest Version:\t%d\n", meta.OldestVersion)
	fmt.Fprintf(w, "Created:\t%s\n", meta.CreatedTime.Format("2006-01-02T15:04:05Z"))
	fmt.Fprintf(w, "Updated:\t%s\n", meta.UpdatedTime.Format("2006-01-02T15:04:05Z"))
	for k, v := range meta.CustomMetadata {
		fmt.Fprintf(w, "Meta/%s:\t%s\n", k, v)
	}
	return w.Flush()
}
