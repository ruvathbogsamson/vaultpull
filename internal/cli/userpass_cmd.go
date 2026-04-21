package cli

import (
	"flag"
	"fmt"
	"os"

	"github.com/user/vaultpull/internal/vault"
)

// UserpassFlags holds parsed flags for the userpass subcommand.
type UserpassFlags struct {
	Address  string
	Mount    string
	Username string
	Password string
	Verbose  bool
}

// ParseUserpassFlags parses command-line flags for the userpass auth command.
func ParseUserpassFlags(args []string) (*UserpassFlags, error) {
	fs := flag.NewFlagSet("userpass", flag.ContinueOnError)

	address := fs.String("address", os.Getenv("VAULT_ADDR"), "Vault server address")
	mount := fs.String("mount", "userpass", "Userpass auth mount path")
	username := fs.String("username", "", "Vault username")
	password := fs.String("password", "", "Vault password")
	verbose := fs.Bool("verbose", false, "Enable verbose output")

	if err := fs.Parse(args); err != nil {
		return nil, err
	}
	if *username == "" {
		return nil, fmt.Errorf("flag -username is required")
	}
	if *password == "" {
		return nil, fmt.Errorf("flag -password is required")
	}

	return &UserpassFlags{
		Address:  *address,
		Mount:    *mount,
		Username: *username,
		Password: *password,
		Verbose:  *verbose,
	}, nil
}

// RunUserpass executes the userpass login flow and prints the resulting token.
func RunUserpass(f *UserpassFlags) error {
	client, err := vault.NewUserpassClient(f.Address, f.Mount)
	if err != nil {
		return fmt.Errorf("failed to create userpass client: %w", err)
	}

	tok, err := client.Login(vault.UserpassCredentials{
		Username: f.Username,
		Password: f.Password,
	})
	if err != nil {
		return fmt.Errorf("userpass login failed: %w", err)
	}

	if f.Verbose {
		fmt.Printf("lease_duration: %d\n", tok.LeaseDuration)
		fmt.Printf("renewable: %v\n", tok.Renewable)
	}
	fmt.Printf("client_token: %s\n", tok.ClientToken)
	return nil
}
