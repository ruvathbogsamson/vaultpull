package cli

import (
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/user/vaultpull/internal/vault"
)

// EgressFlags holds parsed flags for the egress subcommand.
type EgressFlags struct {
	Address   string
	Token     string
	Namespace string
}

// ParseEgressFlags parses egress subcommand flags from args.
func ParseEgressFlags(args []string) (*EgressFlags, error) {
	fs := flag.NewFlagSet("egress", flag.ContinueOnError)
	address := fs.String("address", os.Getenv("VAULT_ADDR"), "Vault address")
	token := fs.String("token", os.Getenv("VAULT_TOKEN"), "Vault token")
	namespace := fs.String("namespace", "", "Vault namespace to list egress rules for")

	if err := fs.Parse(args); err != nil {
		return nil, err
	}
	if strings.TrimSpace(*namespace) == "" {
		return nil, fmt.Errorf("egress: --namespace is required")
	}
	return &EgressFlags{
		Address:   *address,
		Token:     *token,
		Namespace: *namespace,
	}, nil
}

// RunEgress executes the egress command, printing rules to out.
func RunEgress(flags *EgressFlags, out io.Writer) error {
	client, err := vault.NewEgressClient(flags.Address, flags.Token)
	if err != nil {
		return fmt.Errorf("egress: %w", err)
	}

	rules, err := client.ListRules(flags.Namespace)
	if err != nil {
		return fmt.Errorf("egress: %w", err)
	}

	if len(rules) == 0 {
		fmt.Fprintln(out, "no egress rules found")
		return nil
	}

	for _, r := range rules {
		fmt.Fprintf(out, "path=%s capabilities=[%s]\n", r.Path, strings.Join(r.Capabilities, ","))
	}
	return nil
}
