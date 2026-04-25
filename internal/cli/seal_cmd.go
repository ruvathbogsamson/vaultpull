package cli

import (
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/yourusername/vaultpull/internal/vault"
)

// SealFlags holds parsed flags for the seal subcommand.
type SealFlags struct {
	Address string
	Token   string
	Action  string // "status" or "seal"
}

// ParseSealFlags parses CLI flags for the seal command.
func ParseSealFlags(args []string) (*SealFlags, error) {
	fs := flag.NewFlagSet("seal", flag.ContinueOnError)
	address := fs.String("address", os.Getenv("VAULT_ADDR"), "Vault address")
	token := fs.String("token", os.Getenv("VAULT_TOKEN"), "Vault token")
	action := fs.String("action", "status", "Action to perform: status or seal")

	if err := fs.Parse(args); err != nil {
		return nil, err
	}
	if *action != "status" && *action != "seal" {
		return nil, fmt.Errorf("invalid action %q: must be 'status' or 'seal'", *action)
	}
	return &SealFlags{
		Address: *address,
		Token:   *token,
		Action:  *action,
	}, nil
}

// RunSeal executes the seal command with the given flags.
func RunSeal(f *SealFlags, out io.Writer) error {
	client, err := vault.NewSealClient(f.Address, f.Token)
	if err != nil {
		return fmt.Errorf("creating seal client: %w", err)
	}

	switch f.Action {
	case "status":
		status, err := client.GetSealStatus()
		if err != nil {
			return fmt.Errorf("getting seal status: %w", err)
		}
		sealState := "unsealed"
		if status.Sealed {
			sealState = "sealed"
		}
		fmt.Fprintf(out, "Vault is %s\n", sealState)
		fmt.Fprintf(out, "Initialized: %v\n", status.Initialized)
		fmt.Fprintf(out, "Version:     %s\n", status.Version)
		if status.ClusterName != "" {
			fmt.Fprintf(out, "Cluster:     %s\n", status.ClusterName)
		}
	case "seal":
		if err := client.Seal(); err != nil {
			return fmt.Errorf("sealing vault: %w", err)
		}
		fmt.Fprintln(out, "Vault sealed successfully")
	}
	return nil
}
