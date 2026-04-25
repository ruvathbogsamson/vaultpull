package cli

import (
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/example/vaultpull/internal/vault"
)

// QuotaFlags holds parsed flags for the quota subcommand.
type QuotaFlags struct {
	Address string
	Token   string
	Name    string
}

// ParseQuotaFlags parses quota subcommand flags from args.
func ParseQuotaFlags(args []string) (*QuotaFlags, error) {
	fs := flag.NewFlagSet("quota", flag.ContinueOnError)

	address := fs.String("address", os.Getenv("VAULT_ADDR"), "Vault address")
	token := fs.String("token", os.Getenv("VAULT_TOKEN"), "Vault token")
	name := fs.String("name", "", "Quota name to look up (required)")

	if err := fs.Parse(args); err != nil {
		return nil, err
	}
	if *name == "" {
		return nil, fmt.Errorf("flag -name is required")
	}
	return &QuotaFlags{
		Address: *address,
		Token:   *token,
		Name:    *name,
	}, nil
}

// RunQuota fetches and prints quota information to w.
func RunQuota(f *QuotaFlags, w io.Writer) error {
	client, err := vault.NewQuotaClient(f.Address, f.Token)
	if err != nil {
		return fmt.Errorf("creating quota client: %w", err)
	}

	info, err := client.GetQuota(f.Name)
	if err != nil {
		return fmt.Errorf("fetching quota: %w", err)
	}

	fmt.Fprintf(w, "Name:          %s\n", info.Name)
	fmt.Fprintf(w, "Path:          %s\n", info.Path)
	fmt.Fprintf(w, "Type:          %s\n", info.Type)
	fmt.Fprintf(w, "Max Requests:  %d\n", info.MaxRequests)
	fmt.Fprintf(w, "Rate:          %.2f req/s\n", info.Rate)
	fmt.Fprintf(w, "Interval:      %.0fs\n", info.Interval)
	if info.BlockInterval > 0 {
		fmt.Fprintf(w, "Block Interval: %.0fs\n", info.BlockInterval)
	}
	return nil
}
