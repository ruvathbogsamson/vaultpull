package cli

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/example/vaultpull/internal/vault"
)

// AppRoleFlags holds parsed flags for the approle subcommand.
type AppRoleFlags struct {
	VaultAddr  string
	MountPath  string
	RoleID     string
	SecretID   string
	OutputEnv  string
}

// ParseAppRoleFlags parses CLI flags for the approle login command.
func ParseAppRoleFlags(args []string) (*AppRoleFlags, error) {
	fs := flag.NewFlagSet("approle", flag.ContinueOnError)

	addr := fs.String("addr", "http://127.0.0.1:8200", "Vault server address")
	mount := fs.String("mount", "approle", "AppRole auth mount path")
	roleID := fs.String("role-id", "", "AppRole role_id (required)")
	secretID := fs.String("secret-id", "", "AppRole secret_id (required)")
	output := fs.String("output-env", "VAULT_TOKEN", "Environment variable name to export the token to")

	if err := fs.Parse(args); err != nil {
		return nil, err
	}

	if *roleID == "" {
		return nil, fmt.Errorf("approle: -role-id is required")
	}
	if *secretID == "" {
		return nil, fmt.Errorf("approle: -secret-id is required")
	}

	return &AppRoleFlags{
		VaultAddr: *addr,
		MountPath: *mount,
		RoleID:    *roleID,
		SecretID:  *secretID,
		OutputEnv: *output,
	}, nil
}

// RunAppRole performs AppRole login and prints the resulting token.
func RunAppRole(ctx context.Context, flags *AppRoleFlags) error {
	client := vault.NewAppRoleClient(flags.VaultAddr, flags.MountPath)

	result, err := client.Login(ctx, vault.AppRoleCredentials{
		RoleID:   flags.RoleID,
		SecretID: flags.SecretID,
	})
	if err != nil {
		return fmt.Errorf("approle login: %w", err)
	}

	fmt.Fprintf(os.Stdout, "export %s=%s\n", flags.OutputEnv, result.ClientToken)
	fmt.Fprintf(os.Stderr, "token lease: %ds, renewable: %v\n", result.LeaseDuration, result.Renewable)
	return nil
}
