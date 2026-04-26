package cli

import (
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/yourusername/vaultpull/internal/vault"
)

// StepDownFlags holds parsed flags for the step-down command.
type StepDownFlags struct {
	Address string
	Token   string
	Verbose bool
}

// ParseStepDownFlags parses CLI flags for the step-down command.
func ParseStepDownFlags(args []string) (*StepDownFlags, error) {
	fs := flag.NewFlagSet("stepdown", flag.ContinueOnError)

	address := fs.String("address", os.Getenv("VAULT_ADDR"), "Vault server address")
	token := fs.String("token", os.Getenv("VAULT_TOKEN"), "Vault token")
	verbose := fs.Bool("verbose", false, "Enable verbose output")

	if err := fs.Parse(args); err != nil {
		return nil, err
	}
	if *address == "" {
		return nil, fmt.Errorf("vault address is required (set -address or VAULT_ADDR)")
	}
	if *token == "" {
		return nil, fmt.Errorf("vault token is required (set -token or VAULT_TOKEN)")
	}
	return &StepDownFlags{
		Address: *address,
		Token:   *token,
		Verbose: *verbose,
	}, nil
}

// RunStepDown executes the step-down command.
func RunStepDown(flags *StepDownFlags, out io.Writer) error {
	client, err := vault.NewStepDownClient(flags.Address, flags.Token)
	if err != nil {
		return fmt.Errorf("initialising step-down client: %w", err)
	}

	if flags.Verbose {
		fmt.Fprintf(out, "Sending step-down request to %s\n", flags.Address)
	}

	status, err := client.StepDown()
	if err != nil {
		return fmt.Errorf("step-down failed: %w", err)
	}

	fmt.Fprintf(out, "Step-down: %s\n", status.Message)
	return nil
}
