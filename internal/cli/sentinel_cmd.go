package cli

import (
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/example/vaultpull/internal/vault"
)

// SentinelFlags holds parsed flags for the sentinel subcommand.
type SentinelFlags struct {
	Address    string
	Token      string
	PolicyName string
}

// ParseSentinelFlags parses CLI flags for the sentinel subcommand.
func ParseSentinelFlags(args []string) (*SentinelFlags, error) {
	fs := flag.NewFlagSet("sentinel", flag.ContinueOnError)
	address := fs.String("address", "", "Vault server address")
	token := fs.String("token", os.Getenv("VAULT_TOKEN"), "Vault token")
	policy := fs.String("policy", "", "Sentinel policy name to fetch")

	if err := fs.Parse(args); err != nil {
		return nil, err
	}
	if *policy == "" {
		return nil, fmt.Errorf("flag -policy is required")
	}
	return &SentinelFlags{
		Address:    *address,
		Token:      *token,
		PolicyName: *policy,
	}, nil
}

// RunSentinel executes the sentinel policy fetch command.
func RunSentinel(flags *SentinelFlags, out io.Writer) error {
	client, err := vault.NewSentinelClient(flags.Address, flags.Token)
	if err != nil {
		return fmt.Errorf("init sentinel client: %w", err)
	}

	policy, err := client.GetPolicy(flags.PolicyName)
	if err != nil {
		return fmt.Errorf("get policy: %w", err)
	}

	fmt.Fprintf(out, "Name: %s\n", policy.Name)
	fmt.Fprintf(out, "Type: %s\n", policy.Type)
	fmt.Fprintf(out, "Body:\n%s\n", policy.Body)
	return nil
}
