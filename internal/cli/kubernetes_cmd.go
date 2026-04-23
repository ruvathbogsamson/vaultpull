package cli

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/yourusername/vaultpull/internal/vault"
)

// KubernetesFlags holds parsed flags for the kubernetes auth command.
type KubernetesFlags struct {
	Address string
	Mount   string
	Role    string
	JWT     string
}

// ParseKubernetesFlags parses CLI flags for the kubernetes auth subcommand.
func ParseKubernetesFlags(args []string) (*KubernetesFlags, error) {
	fs := flag.NewFlagSet("kubernetes", flag.ContinueOnError)

	address := fs.String("address", os.Getenv("VAULT_ADDR"), "Vault server address")
	mount := fs.String("mount", "kubernetes", "Kubernetes auth mount path")
	role := fs.String("role", "", "Kubernetes role name (required)")
	jwt := fs.String("jwt", "", "Service account JWT token (required)")

	if err := fs.Parse(args); err != nil {
		return nil, err
	}

	if strings.TrimSpace(*role) == "" {
		return nil, fmt.Errorf("kubernetes: -role is required")
	}
	if strings.TrimSpace(*jwt) == "" {
		return nil, fmt.Errorf("kubernetes: -jwt is required")
	}

	return &KubernetesFlags{
		Address: *address,
		Mount:   *mount,
		Role:    *role,
		JWT:     *jwt,
	}, nil
}

// RunKubernetes executes the kubernetes auth login flow and prints the client token.
func RunKubernetes(args []string) error {
	flags, err := ParseKubernetesFlags(args)
	if err != nil {
		return err
	}

	client, err := vault.NewKubernetesClient(flags.Address, flags.Mount)
	if err != nil {
		return fmt.Errorf("kubernetes: failed to create client: %w", err)
	}

	resp, err := client.Login(flags.Role, flags.JWT)
	if err != nil {
		return fmt.Errorf("kubernetes: login failed: %w", err)
	}

	fmt.Printf("client_token: %s\n", resp.ClientToken)
	fmt.Printf("accessor:     %s\n", resp.Accessor)
	fmt.Printf("policies:     %s\n", strings.Join(resp.Policies, ", "))
	return nil
}
