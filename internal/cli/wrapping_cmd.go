package cli

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/fmartingr/vaultpull/internal/vault"
)

// WrappingFlags holds parsed flags for the wrapping subcommand.
type WrappingFlags struct {
	Address       string
	Token         string
	WrappingToken string
	Lookup        bool
}

// ParseWrappingFlags parses CLI flags for the wrapping subcommand.
func ParseWrappingFlags(args []string) (*WrappingFlags, error) {
	fs := flag.NewFlagSet("wrapping", flag.ContinueOnError)

	address := fs.String("address", os.Getenv("VAULT_ADDR"), "Vault server address")
	token := fs.String("token", os.Getenv("VAULT_TOKEN"), "Vault token")
	wrappingToken := fs.String("wrapping-token", "", "Wrapping token to unwrap or inspect")
	lookup := fs.Bool("lookup", false, "Lookup wrapping token metadata without consuming it")

	if err := fs.Parse(args); err != nil {
		return nil, err
	}
	if strings.TrimSpace(*wrappingToken) == "" {
		return nil, fmt.Errorf("--wrapping-token is required")
	}
	return &WrappingFlags{
		Address:       *address,
		Token:         *token,
		WrappingToken: *wrappingToken,
		Lookup:        *lookup,
	}, nil
}

// RunWrapping executes the wrapping subcommand.
func RunWrapping(f *WrappingFlags) error {
	client, err := vault.NewWrappingClient(f.Address, f.Token)
	if err != nil {
		return fmt.Errorf("initialising wrapping client: %w", err)
	}

	if f.Lookup {
		ws, err := client.LookupWrappingToken(f.WrappingToken)
		if err != nil {
			return fmt.Errorf("lookup failed: %w", err)
		}
		fmt.Printf("Token:         %s\n", ws.Token)
		fmt.Printf("Accessor:      %s\n", ws.Accessor)
		fmt.Printf("Creation Time: %s\n", ws.Creation)
		return nil
	}

	data, err := client.Unwrap(f.WrappingToken)
	if err != nil {
		return fmt.Errorf("unwrap failed: %w", err)
	}
	for k, v := range data {
		fmt.Printf("%s=%s\n", k, v)
	}
	return nil
}
