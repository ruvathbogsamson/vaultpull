package cli

import (
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/user/vaultpull/internal/vault"
)

// AuthTokenFlags holds parsed flags for the auth-token subcommand.
type AuthTokenFlags struct {
	Address string
	Token   string
	Verbose bool
}

// ParseAuthTokenFlags parses CLI flags for the auth-token subcommand.
func ParseAuthTokenFlags(args []string) (*AuthTokenFlags, error) {
	fs := flag.NewFlagSet("auth-token", flag.ContinueOnError)

	address := fs.String("address", os.Getenv("VAULT_ADDR"), "Vault server address")
	token := fs.String("token", os.Getenv("VAULT_TOKEN"), "Vault token to validate")
	verbose := fs.Bool("verbose", false, "Print full token details")

	if err := fs.Parse(args); err != nil {
		return nil, err
	}
	if *token == "" {
		return nil, fmt.Errorf("flag -token or VAULT_TOKEN is required")
	}

	return &AuthTokenFlags{
		Address: *address,
		Token:   *token,
		Verbose: *verbose,
	}, nil
}

// RunAuthToken validates the configured Vault token and prints token metadata.
func RunAuthToken(flags *AuthTokenFlags, out io.Writer) error {
	c, err := vault.NewAuthTokenClient(flags.Address, flags.Token)
	if err != nil {
		return fmt.Errorf("initialising auth token client: %w", err)
	}

	info, err := c.ValidateToken()
	if err != nil {
		return fmt.Errorf("token validation failed: %w", err)
	}

	fmt.Fprintf(out, "Token is valid\n")
	fmt.Fprintf(out, "  Accessor:  %s\n", info.Accessor)
	fmt.Fprintf(out, "  Renewable: %v\n", info.Renewable)
	fmt.Fprintf(out, "  TTL:       %s\n", info.TTL)

	if flags.Verbose {
		fmt.Fprintf(out, "  Policies:  %v\n", info.Policies)
		fmt.Fprintf(out, "  Expires:   %s\n", info.ExpireTime)
	}
	return nil
}
