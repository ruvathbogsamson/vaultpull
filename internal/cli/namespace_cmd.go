package cli

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/example/vaultpull/internal/vault"
)

// NamespaceFlags holds parsed flags for the namespace subcommand.
type NamespaceFlags struct {
	Address   string
	Token     string
	Namespace string
}

// ParseNamespaceFlags parses args for the namespace subcommand.
func ParseNamespaceFlags(args []string) (*NamespaceFlags, error) {
	fs := flag.NewFlagSet("namespace", flag.ContinueOnError)

	addr := fs.String("addr", "http://127.0.0.1:8200", "Vault address")
	token := fs.String("token", "", "Vault token (or set VAULT_TOKEN)")
	ns := fs.String("namespace", "", "Vault namespace to inspect")

	if err := fs.Parse(args); err != nil {
		return nil, err
	}

	resolved := *token
	if resolved == "" {
		resolved = os.Getenv("VAULT_TOKEN")
	}
	if *ns == "" {
		return nil, fmt.Errorf("flag -namespace is required")
	}

	return &NamespaceFlags{
		Address:   *addr,
		Token:     resolved,
		Namespace: *ns,
	}, nil
}

// RunNamespace lists child namespaces and prints them to stdout.
func RunNamespace(f *NamespaceFlags) error {
	nc, err := vault.NewNamespaceClient(f.Address, f.Token, f.Namespace)
	if err != nil {
		return fmt.Errorf("initialising namespace client: %w", err)
	}

	children, err := nc.ListNamespaces()
	if err != nil {
		return fmt.Errorf("listing namespaces: %w", err)
	}

	if len(children) == 0 {
		fmt.Printf("No child namespaces found under %q\n", f.Namespace)
		return nil
	}

	fmt.Printf("Child namespaces under %q:\n", f.Namespace)
	for _, ns := range children {
		fmt.Printf("  %s/%s\n", strings.TrimSuffix(f.Namespace, "/"), ns)
	}
	return nil
}
