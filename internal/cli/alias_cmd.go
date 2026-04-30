package cli

import (
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/yourusername/vaultpull/internal/vault"
)

// AliasCmdFlags holds parsed flags for the alias command.
type AliasCmdFlags struct {
	Address string
	Token   string
	AliasID string
}

// ParseAliasFlags parses command-line flags for the alias subcommand.
func ParseAliasFlags(args []string) (*AliasCmdFlags, error) {
	fs := flag.NewFlagSet("alias", flag.ContinueOnError)

	address := fs.String("address", os.Getenv("VAULT_ADDR"), "Vault server address")
	token := fs.String("token", os.Getenv("VAULT_TOKEN"), "Vault token")
	aliasID := fs.String("id", "", "Entity alias ID to look up (required)")

	if err := fs.Parse(args); err != nil {
		return nil, err
	}
	if *aliasID == "" {
		return nil, fmt.Errorf("flag -id is required")
	}
	return &AliasCmdFlags{
		Address: *address,
		Token:   *token,
		AliasID: *aliasID,
	}, nil
}

// RunAlias executes the alias lookup command, writing output to w.
func RunAlias(flags *AliasCmdFlags, w io.Writer) error {
	client, err := vault.NewAliasClient(flags.Address, flags.Token)
	if err != nil {
		return fmt.Errorf("creating alias client: %w", err)
	}
	alias, err := client.GetAlias(flags.AliasID)
	if err != nil {
		return fmt.Errorf("fetching alias: %w", err)
	}
	fmt.Fprintf(w, "ID:             %s\n", alias.ID)
	fmt.Fprintf(w, "Name:           %s\n", alias.Name)
	fmt.Fprintf(w, "Mount Accessor: %s\n", alias.MountAccessor)
	fmt.Fprintf(w, "Mount Type:     %s\n", alias.MountType)
	fmt.Fprintf(w, "Entity ID:      %s\n", alias.EntityID)
	return nil
}
