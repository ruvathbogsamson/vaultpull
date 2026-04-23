package cli

import (
	"flag"
	"fmt"
	"os"

	"github.com/yourusername/vaultpull/internal/vault"
)

// SSHFlags holds parsed flags for the ssh sub-command.
type SSHFlags struct {
	Address   string
	Token     string
	Mount     string
	Role      string
	PublicKey string
}

// ParseSSHFlags parses CLI flags for the ssh sign sub-command.
func ParseSSHFlags(args []string) (*SSHFlags, error) {
	fs := flag.NewFlagSet("ssh", flag.ContinueOnError)

	address := fs.String("address", os.Getenv("VAULT_ADDR"), "Vault address")
	token := fs.String("token", os.Getenv("VAULT_TOKEN"), "Vault token")
	mount := fs.String("mount", "ssh", "SSH secrets engine mount path")
	role := fs.String("role", "", "SSH role name (required)")
	publicKey := fs.String("public-key", "", "SSH public key to sign (required)")

	if err := fs.Parse(args); err != nil {
		return nil, err
	}
	if *role == "" {
		return nil, fmt.Errorf("flag -role is required")
	}
	if *publicKey == "" {
		return nil, fmt.Errorf("flag -public-key is required")
	}
	return &SSHFlags{
		Address:   *address,
		Token:     *token,
		Mount:     *mount,
		Role:      *role,
		PublicKey: *publicKey,
	}, nil
}

// RunSSH executes the SSH key-signing command.
func RunSSH(args []string, stdout *os.File) error {
	flags, err := ParseSSHFlags(args)
	if err != nil {
		return err
	}

	c, err := vault.NewSSHClient(flags.Address, flags.Token, flags.Mount)
	if err != nil {
		return fmt.Errorf("creating SSH client: %w", err)
	}

	cred, err := c.SignKey(flags.Role, flags.PublicKey)
	if err != nil {
		return fmt.Errorf("signing key: %w", err)
	}

	fmt.Fprintln(stdout, cred.SignedKey)
	return nil
}
